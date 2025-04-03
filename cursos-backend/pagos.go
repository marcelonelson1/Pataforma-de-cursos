package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PagoRequest struct {
	CursoID        uint            `json:"curso_id" binding:"required"`
	Monto          float64         `json:"monto" binding:"required"`
	Metodo         string          `json:"metodo" binding:"required"`
	DetallesTarjeta *DetallesTarjeta `json:"detalles_tarjeta,omitempty"`
}

type DetallesTarjeta struct {
	Numero     string `json:"numero"`
	Expiracion string `json:"expiracion"`
	CVV        string `json:"cvv"`
}

func crearPago(c *gin.Context) {
	// Asegurarse de que la respuesta sea JSON
	c.Header("Content-Type", "application/json")

	userValue, exists := c.Get("user")
	if !exists {
		SendErrorResponse(c, ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	user, ok := userValue.(Usuario)
	if !ok {
		SendErrorResponse(c, errors.New("error al obtener información del usuario"), http.StatusInternalServerError)
		return
	}

	var req PagoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error al parsear solicitud de pago: %v", err)
		SendValidationErrorResponse(c, ErrInvalidRequest, gin.H{"details": err.Error()})
		return
	}

	log.Printf("Solicitud de pago recibida: %+v", req)

	if req.Metodo != "tarjeta" && req.Metodo != "paypal" && req.Metodo != "transferencia" && req.Metodo != "dev" {
		SendErrorResponse(c, errors.New("método de pago no válido"), http.StatusBadRequest)
		return
	}

	if req.Metodo == "tarjeta" && req.DetallesTarjeta == nil && req.Metodo != "dev" {
		SendErrorResponse(c, errors.New("se requieren detalles de tarjeta para este método de pago"), http.StatusBadRequest)
		return
	}

	var curso Curso
	if result := db.First(&curso, req.CursoID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			SendErrorResponse(c, ErrResourceNotFound, http.StatusNotFound)
		} else {
			log.Printf("Error de base de datos al buscar curso: %v", result.Error)
			SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		}
		return
	}

	var pagoExistente Pago
	if result := db.Where("usuario_id = ? AND curso_id = ? AND estado = ?", 
		user.ID, req.CursoID, "aprobado").First(&pagoExistente); result.Error == nil {
		SendSuccessResponse(c, gin.H{
			"message": "Ya tienes acceso a este curso",
			"estado":  "aprobado",
			"pago_id": pagoExistente.ID,
		})
		return
	}

	pago := Pago{
		UsuarioID:     user.ID,
		CursoID:       req.CursoID,
		Monto:         req.Monto,
		Metodo:        req.Metodo,
		Estado:        "pendiente",
		TransaccionID: "",
	}

	if result := db.Create(&pago); result.Error != nil {
		log.Printf("Error al guardar pago en base de datos: %v", result.Error)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	go simularPasarelaPago(pago.ID, req.Metodo)

	SendSuccessResponse(c, gin.H{
		"message": "Pago en proceso",
		"pago_id": pago.ID,
		"estado":  pago.Estado,
	})
}

func verificarPagoPorCurso(c *gin.Context) {
	// Asegurarse de que la respuesta sea JSON
	c.Header("Content-Type", "application/json")
	
	userValue, exists := c.Get("user")
	if !exists {
		SendErrorResponse(c, ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	user, ok := userValue.(Usuario)
	if !ok {
		SendErrorResponse(c, errors.New("error al obtener información del usuario"), http.StatusInternalServerError)
		return
	}

	cursoID := c.Param("id")
	if cursoID == "" {
		SendErrorResponse(c, errors.New("se requiere el ID del curso"), http.StatusBadRequest)
		return
	}

	var cursoUint uint
	if _, err := fmt.Sscanf(cursoID, "%d", &cursoUint); err != nil {
		SendErrorResponse(c, errors.New("ID de curso inválido"), http.StatusBadRequest)
		return
	}

	log.Printf("Verificando pago para curso ID: %s, usuario ID: %d", cursoID, user.ID)

	var pago Pago
	result := db.Where("usuario_id = ? AND curso_id = ?", user.ID, cursoUint).
		Order("created_at desc").
		First(&pago)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Enviar respuesta JSON directamente sin anidamiento
			c.JSON(http.StatusOK, gin.H{
				"estado": "no_pagado",
				"message": "No se encontró pago para este curso",
			})
			return
		}
		log.Printf("Error de base de datos al verificar pago: %v", result.Error)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	// Enviar respuesta JSON directamente sin anidamiento
	c.JSON(http.StatusOK, gin.H{
		"estado": pago.Estado,
		"pago":   pago,
	})
}

func webhookPago(c *gin.Context) {
	// Asegurarse de que la respuesta sea JSON
	c.Header("Content-Type", "application/json")
	
	var payload struct {
		PagoID         uint   `json:"pago_id"`
		Estado         string `json:"estado"`
		TransaccionID  string `json:"transaccion_id"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		SendErrorResponse(c, ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	var pago Pago
	if result := db.First(&pago, payload.PagoID); result.Error != nil {
		SendErrorResponse(c, ErrPaymentNotFound, http.StatusNotFound)
		return
	}

	pago.Estado = payload.Estado
	pago.TransaccionID = payload.TransaccionID

	if result := db.Save(&pago); result.Error != nil {
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	SendSuccessResponse(c, gin.H{
		"message": "Estado de pago actualizado",
		"pago_id": pago.ID,
		"estado":  pago.Estado,
	})
}

func simularPasarelaPago(pagoID uint, metodo string) {
	time.Sleep(3 * time.Second)

	var pago Pago
	if result := db.First(&pago, pagoID); result.Error != nil {
		log.Printf("Error al recuperar pago ID %d para simulación: %v", pagoID, result.Error)
		return
	}

	if getEnv("ENV", "development") == "development" {
		pago.Estado = "aprobado"
		pago.TransaccionID = generarIDTransaccion(metodo)
	} else {
		if rand.Intn(100) < 80 {
			pago.Estado = "aprobado"
			pago.TransaccionID = generarIDTransaccion(metodo)
		} else {
			pago.Estado = "rechazado"
		}
	}

	if result := db.Save(&pago); result.Error != nil {
		log.Printf("Error al actualizar estado de pago ID %d: %v", pagoID, result.Error)
	} else {
		log.Printf("Pago ID %d actualizado a estado: %s", pagoID, pago.Estado)
	}
}

func generarIDTransaccion(metodo string) string {
	prefix := "txn"
	switch metodo {
	case "tarjeta":
		prefix = "card"
	case "paypal":
		prefix = "pp"
	case "transferencia":
		prefix = "trf"
	case "dev":
		prefix = "dev"
	}
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().Unix(), rand.Intn(10000))
}