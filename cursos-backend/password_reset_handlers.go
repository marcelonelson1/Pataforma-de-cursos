package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"text/template"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ForgotPasswordRequest estructura para la solicitud de recuperación
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest estructura para restablecer contraseña
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// Controlador para solicitar restablecimiento de contraseña
func forgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el email existe
	var user Usuario
	result := db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		// No revelamos si el email existe o no por seguridad
		c.JSON(http.StatusOK, gin.H{"message": "Si el email existe en nuestra base de datos, recibirás un enlace para restablecer tu contraseña."})
		log.Printf("Email no encontrado o error de DB: %v", result.Error)
		return
	}

	// Generar token de restablecimiento
	reset, err := GenerateResetToken(req.Email)
	if err != nil {
		log.Printf("Error al generar token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar token de restablecimiento"})
		return
	}

	// Construir el enlace de restablecimiento
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	resetLink := frontendURL + "/reset-password/" + reset.Token

	// Enviar correo electrónico
	appEnv := getEnv("APP_ENV", "development")

	// Intentar enviar email en todos los entornos, pero manejar fallos de forma diferente
	emailError := sendPasswordResetEmail(user.Email, user.Nombre, resetLink)
	if emailError != nil {
		log.Printf("Error al enviar correo: %v", emailError)
		// En producción, esto podría ser un problema crítico
		if appEnv == "production" {
			// En producción registramos el error pero no lo devolvemos al cliente
			log.Printf("ERROR CRÍTICO: Fallo al enviar email en producción: %v", emailError)
		}
	}

	// Respuesta base
	response := gin.H{
		"message": "Si el email existe en nuestra base de datos, recibirás un enlace para restablecer tu contraseña.",
	}

	// Solo enviar el token y enlace en desarrollo
	if appEnv == "development" {
		response["resetToken"] = reset.Token
		response["resetLink"] = resetLink
		// En desarrollo, podemos informar sobre el error de email si ocurrió
		if emailError != nil {
			response["emailError"] = emailError.Error()
		}
	}

	c.JSON(http.StatusOK, response)
}

// Controlador para validar el token de restablecimiento
func validateResetToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token requerido", "valid": false})
		return
	}

	// Validar el token
	reset, err := ValidateResetToken(token)
	if err != nil {
		log.Printf("Token inválido: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token inválido o expirado", "valid": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "email": reset.Email})
}

// Controlador para restablecer la contraseña
func resetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar el token primero
	reset, err := ValidateResetToken(req.Token)
	if err != nil {
		log.Printf("Error al validar token: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token inválido o expirado"})
		return
	}

	// Buscar el usuario por email
	var user Usuario
	if result := db.Where("email = ?", reset.Email).First(&user); result.Error != nil {
		log.Printf("Usuario no encontrado: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Hash de la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error al generar hash: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la contraseña"})
		return
	}

	// Actualizar la contraseña del usuario
	if result := db.Model(&user).Update("password", string(hashedPassword)); result.Error != nil {
		log.Printf("Error al actualizar contraseña: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la contraseña"})
		return
	}

	// Marcar el token como usado
	if err := MarkTokenAsUsed(req.Token); err != nil {
		log.Printf("Error al marcar token como usado: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el estado del token"})
		return
	}

	// Devolver un JSON válido
	c.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada correctamente"})
}

// Función para enviar email de recuperación con manejo mejorado de errores
func sendPasswordResetEmail(to, name, resetLink string) error {
	from := getEnv("EMAIL_FROM", "noreply@example.com")
	password := getEnv("EMAIL_PASSWORD", "")
	smtpHost := getEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := getEnv("SMTP_PORT", "587")

	// Comprobación de configuración mínima
	if password == "" {
		return errors.New("la contraseña de email no está configurada")
	}

	// Configurar autenticación
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Renderizar el HTML del correo electrónico
	htmlContent, err := renderEmailTemplate(name, resetLink)
	if err != nil {
		log.Printf("Error al renderizar la plantilla de correo: %v", err)
		return err
	}

	// Construir el mensaje
	subject := "Recuperación de contraseña"
	headers := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	message := []byte("To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		headers + "\r\n" +
		htmlContent)

	// En desarrollo, podemos simular el envío
	if getEnv("APP_ENV", "development") == "development" && getEnv("MOCK_EMAIL", "true") == "true" {
		log.Printf("Simulando envío de email a %s. Contenido: %s", to, htmlContent)
		// Guardar una copia del correo en un archivo para probar
		f, err := os.OpenFile("last_reset_email.html", os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			f.WriteString(htmlContent)
		}
		return nil
	}

	// Enviar email con timeout adicional
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		log.Printf("Error SMTP: %v", err)
		return err
	}

	return nil
}

// Función para renderizar el HTML del correo electrónico
func renderEmailTemplate(name, resetLink string) (string, error) {
	// Cargar la plantilla desde el archivo
	tmpl, err := template.ParseFiles("templates/email_template.html")
	if err != nil {
		return "", err
	}

	// Ejecutar la plantilla con los datos
	var htmlContent bytes.Buffer
	data := struct {
		Name      string
		ResetLink string
	}{
		Name:      name,
		ResetLink: resetLink,
	}

	if err := tmpl.Execute(&htmlContent, data); err != nil {
		return "", err
	}

	return htmlContent.String(), nil
}
