package main

import (
	"bytes"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContactRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone"`
	Message string `json:"message" binding:"required"`
}

type ContactMessage struct {
	gorm.Model
	Name    string `gorm:"size:100;not null" json:"name"`
	Email   string `gorm:"size:100;not null" json:"email"`
	Phone   string `gorm:"size:20" json:"phone"`
	Message string `gorm:"type:text;not null" json:"message"`
	Read    bool   `gorm:"default:false" json:"read"`
	Starred bool   `gorm:"default:false" json:"starred"`
}

func contactHandler(c *gin.Context) {
	var req ContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	// Guardar en la base de datos
	message := ContactMessage{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Message: req.Message,
	}

	if err := db.Create(&message).Error; err != nil {
		log.Printf("Error guardando mensaje: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	// Enviar por email
	if err := sendContactEmail(req); err != nil {
		log.Printf("Error enviando email: %v", err)
		SendErrorResponse(c, ErrEmailSendError, http.StatusInternalServerError)
		return
	}

	SendSuccessResponse(c, gin.H{"message": "Mensaje enviado exitosamente"})
}

func sendContactEmail(contact ContactRequest) error {
	from := getEnv("EMAIL_FROM", "noreply@example.com")
	to := getEnv("CONTACT_EMAIL", "marcelinho.nelson@gmail.com")
	pass := getEnv("EMAIL_PASSWORD", "")
	smtpHost := getEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := getEnv("SMTP_PORT", "587")

	if pass == "" {
		return ErrEmailSendError
	}

	auth := smtp.PlainAuth("", from, pass, smtpHost)

	// Cargar template desde archivo
	tmpl, err := template.ParseFiles(filepath.Join("templates", "contact_email.html"))
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, contact); err != nil {
		return err
	}

	headers := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg := "To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: Nuevo mensaje de contacto - " + contact.Name + "\r\n" +
		headers + "\r\n" +
		body.String()

	// Modo desarrollo
	if getEnv("APP_ENV", "development") == "development" {
		return saveEmailToFile(body.Bytes())
	}

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}

func saveEmailToFile(content []byte) error {
	return os.WriteFile("last_contact_email.html", content, 0644)
}