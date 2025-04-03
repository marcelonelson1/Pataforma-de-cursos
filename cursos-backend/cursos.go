package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CursoRequest struct {
	Titulo      string  `json:"titulo" binding:"required"`
	Descripcion string  `json:"descripcion" binding:"required"`
	Contenido   string  `json:"contenido" binding:"required"`
	Precio      float64 `json:"precio"`
}

func getCursos(c *gin.Context) {
	var cursos []Curso
	
	result := db.Find(&cursos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener cursos"})
		return
	}

	c.JSON(http.StatusOK, cursos)
}

func getCursoById(c *gin.Context) {
	id := c.Param("id")
	
	var curso Curso
	if result := db.First(&curso, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
		return
	}

	c.JSON(http.StatusOK, curso)
}

func createCurso(c *gin.Context) {
	var req CursoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	curso := Curso{
		Titulo:      req.Titulo,
		Descripcion: req.Descripcion,
		Contenido:   req.Contenido,
		Precio:      req.Precio,
	}

	if result := db.Create(&curso); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear curso"})
		return
	}

	c.JSON(http.StatusCreated, curso)
}

func updateCurso(c *gin.Context) {
	id := c.Param("id")
	
	var curso Curso
	if result := db.First(&curso, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
		return
	}

	var req CursoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	curso.Titulo = req.Titulo
	curso.Descripcion = req.Descripcion
	curso.Contenido = req.Contenido
	curso.Precio = req.Precio

	if result := db.Save(&curso); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar curso"})
		return
	}

	c.JSON(http.StatusOK, curso)
}

func deleteCurso(c *gin.Context) {
	id := c.Param("id")
	
	var curso Curso
	if result := db.First(&curso, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
		return
	}

	if result := db.Delete(&curso); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar curso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Curso eliminado correctamente"})
}