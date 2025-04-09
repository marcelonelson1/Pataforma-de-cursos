package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	
	"text/template"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getContactMessages(c *gin.Context) {
	var messages []ContactMessage

	if err := db.Order("created_at DESC").Find(&messages).Error; err != nil {
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	// Estructura de respuesta estandarizada
	response := gin.H{
		"success": true,
		"data":    messages,
	}

	c.JSON(http.StatusOK, response)
}

func getContactMessage(c *gin.Context) {
	id := c.Param("id")

	var message ContactMessage
	if err := db.First(&message, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			SendErrorResponse(c, ErrResourceNotFound, http.StatusNotFound)
		} else {
			SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		}
		return
	}

	if !message.Read {
		db.Model(&message).Update("read", true)
	}

	SendSuccessResponse(c, gin.H{"data": message})
}

func updateMessageStatus(c *gin.Context) {
	id := c.Param("id")
	action := c.Param("action")

	var message ContactMessage
	if err := db.First(&message, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			SendErrorResponse(c, ErrResourceNotFound, http.StatusNotFound)
		} else {
			SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		}
		return
	}

	switch action {
	case "read":
		db.Model(&message).Update("read", !message.Read)
	case "star":
		db.Model(&message).Update("starred", !message.Starred)
	default:
		SendErrorResponse(c, ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	SendSuccessResponse(c, gin.H{"data": message})
}

func deleteContactMessage(c *gin.Context) {
	id := c.Param("id")

	if err := db.Delete(&ContactMessage{}, id).Error; err != nil {
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	SendSuccessResponse(c, gin.H{"message": "Mensaje eliminado correctamente"})
}

func replyToMessage(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	var contactMsg ContactMessage
	if err := db.First(&contactMsg, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			SendErrorResponse(c, ErrResourceNotFound, http.StatusNotFound)
		} else {
			SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		}
		return
	}

	// Enviar email de respuesta al usuario
	if err := sendReplyEmail(contactMsg, req.Message); err != nil {
		log.Printf("Error al enviar email de respuesta: %v", err)
		SendErrorResponse(c, ErrEmailSendError, http.StatusInternalServerError)
		return
	}

	log.Printf("Respuesta enviada a %s (%s): %s", contactMsg.Name, contactMsg.Email, req.Message)

	SendSuccessResponse(c, gin.H{"message": "Respuesta enviada correctamente"})
}

func sendReplyEmail(contact ContactMessage, replyText string) error {
	from := getEnv("EMAIL_FROM", "noreply@example.com")
	pass := getEnv("EMAIL_PASSWORD", "")
	smtpHost := getEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := getEnv("SMTP_PORT", "587")

	// Log detallado de la configuración (ocultando la contraseña)
	log.Printf("Configuración de correo: FROM=%s, HOST=%s, PORT=%s, PASS=***", 
		from, smtpHost, smtpPort)

	// Si estamos en modo desarrollo, simular el envío del correo
	if getEnv("APP_ENV", "development") == "development" && getEnv("MOCK_EMAIL", "true") == "true" {
		log.Println("Modo desarrollo con MOCK_EMAIL=true: Simulando envío de correo")
		
		// Crear un mensaje simple para modo desarrollo
		msg := fmt.Sprintf(`
De: %s
Para: %s (%s)
Asunto: Respuesta a tu mensaje

Mensaje original:
%s

Respuesta:
%s
`, from, contact.Name, contact.Email, contact.Message, replyText)
		
		return os.WriteFile("last_reply_email.txt", []byte(msg), 0644)
	}
	
	// Verificar la configuración requerida
	if pass == "" {
		log.Println("Error: EMAIL_PASSWORD no está configurado")
		return ErrEmailSendError
	}

	// Autenticación SMTP
	auth := smtp.PlainAuth("", from, pass, smtpHost)

	// Preparar contenido del email
	subject := "Respuesta a tu mensaje"
	
	// Email plano (sin HTML) como fallback si hay problemas con el template
	plainBody := fmt.Sprintf("Hola %s,\n\nGracias por contactarnos. A continuación encontrarás nuestra respuesta a tu mensaje:\n\nTu mensaje original:\n%s\n\nNuestra respuesta:\n%s\n\nSaludos cordiales,\nEl equipo de soporte", 
		contact.Name, contact.Message, replyText)
	
	// Intentar crear un email HTML
	var emailBody string
	
	// Datos para el template
	data := struct {
		Name        string
		Email       string
		Subject     string
		OriginalMsg string
		ReplyMsg    string
	}{
		Name:        contact.Name,
		Email:       contact.Email,
		Subject:     subject,
		OriginalMsg: contact.Message,
		ReplyMsg:    replyText,
	}
	
	// Intentar usar un template HTML
	htmlTemplate := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>{{.Subject}}</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; }
        .header { background-color: #4a86e8; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .original-message { background-color: #f5f5f5; padding: 15px; border-left: 4px solid #ccc; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; font-size: 0.8em; color: #777; }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.Subject}}</h1>
    </div>
    <div class="content">
        <p>Hola {{.Name}},</p>
        <p>Gracias por contactarnos. A continuación encontrarás nuestra respuesta a tu mensaje:</p>
        
        <div class="original-message">
            <p><strong>Tu mensaje original:</strong></p>
            <p>{{.OriginalMsg}}</p>
        </div>
        
        <p><strong>Nuestra respuesta:</strong></p>
        <p>{{.ReplyMsg}}</p>
        
        <p>Saludos cordiales,<br>El equipo de soporte</p>
    </div>
    <div class="footer">
        <p>Este es un mensaje automático, por favor no respondas a este correo.</p>
    </div>
</body>
</html>`
	
	// Crear el template desde el string en lugar de intentar cargar el archivo
	tmpl, err := template.New("email_template").Parse(htmlTemplate)
	if err != nil {
		log.Printf("Error al procesar template HTML: %v", err)
		// Fallback a email de texto plano
		emailBody = plainBody
	} else {
		var htmlBuffer bytes.Buffer
		if err := tmpl.Execute(&htmlBuffer, data); err != nil {
			log.Printf("Error al renderizar template: %v", err)
			// Fallback a email de texto plano
			emailBody = plainBody
		} else {
			emailBody = htmlBuffer.String()
		}
	}
	
	// Determinar si estamos usando HTML o texto plano
	var contentType string
	if emailBody == plainBody {
		contentType = "text/plain; charset=\"UTF-8\""
	} else {
		contentType = "text/html; charset=\"UTF-8\""
	}
	
	// Construir el mensaje de correo
	headers := "MIME-version: 1.0;\r\nContent-Type: " + contentType + ";\r\n"
	msg := "To: " + contact.Email + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + " - " + contact.Name + "\r\n" +
		headers + "\r\n" +
		emailBody
	
	// En modo desarrollo, también guardar una copia del email
	if getEnv("APP_ENV", "development") == "development" {
		log.Println("Guardando copia del email en last_reply_email.html")
		os.WriteFile("last_reply_email.html", []byte(emailBody), 0644)
	}
	
	// Enviar email
	serverAddr := smtpHost + ":" + smtpPort
	log.Printf("Intentando enviar email a %s a través de %s", contact.Email, serverAddr)
	
	err = smtp.SendMail(serverAddr, auth, from, []string{contact.Email}, []byte(msg))
	if err != nil {
		log.Printf("Error SMTP al enviar email: %v", err)
		return err
	}
	
	log.Printf("Email enviado exitosamente a %s", contact.Email)
	return nil
}

func testSmtpConnection(c *gin.Context) {
	from := getEnv("EMAIL_FROM", "noreply@example.com")
	pass := getEnv("EMAIL_PASSWORD", "")
	smtpHost := getEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := getEnv("SMTP_PORT", "587")
	testEmail := c.Query("email") // Email para prueba
	
	if testEmail == "" {
		testEmail = from // Si no se proporciona, usar el mismo email del remitente
	}
	
	// Añadir información de configuración
	result := gin.H{
		"smtp_host": smtpHost,
		"smtp_port": smtpPort,
		"from_email": from,
		"password_configured": pass != "",
		"test_email": testEmail,
	}
	
	// Verificar que las variables de entorno necesarias estén configuradas
	if pass == "" {
		result["error"] = "EMAIL_PASSWORD no está configurado"
		result["success"] = false
		c.JSON(http.StatusBadRequest, result)
		return
	}
	
	// Autenticación SMTP
	auth := smtp.PlainAuth("", from, pass, smtpHost)
	
	// Mensaje simple para prueba
	headers := "MIME-version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n"
	subject := "Prueba de conexión SMTP"
	body := "Este es un mensaje de prueba para verificar la configuración SMTP."
	
	msg := "To: " + testEmail + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		headers + "\r\n" +
		body
	
	// Intentar conectar al servidor SMTP
	serverAddr := smtpHost + ":" + smtpPort
	err := smtp.SendMail(serverAddr, auth, from, []string{testEmail}, []byte(msg))
	
	if err != nil {
		log.Printf("Error al probar conexión SMTP: %v", err)
		result["error"] = err.Error()
		result["success"] = false
		c.JSON(http.StatusInternalServerError, result)
		return
	}
	
	result["success"] = true
	result["message"] = "Conexión SMTP exitosa. Email de prueba enviado a " + testEmail
	c.JSON(http.StatusOK, result)
}