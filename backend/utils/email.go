// utils/email.go
package utils

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"text/template"
	"time"
)

// EmailConfig contiene la configuración para enviar correos electrónicos
type EmailConfig struct {
	From     string
	Password string
	SMTPHost string
	SMTPPort string
	To       string
}

// GetEmailConfig obtiene la configuración de correo electrónico desde variables de entorno
func GetEmailConfig() EmailConfig {
	return EmailConfig{
		From:     GetEnv("EMAIL_FROM", "noreply@example.com"),
		Password: GetEnv("EMAIL_PASSWORD", ""),
		SMTPHost: GetEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort: GetEnv("SMTP_PORT", "587"),
	}
}

// SendEmail envía un correo electrónico con el contenido especificado
func SendEmail(to, subject, body, contentType string) error {
	config := GetEmailConfig()
	
	log.Println("===== INICIO LOG DE ENVÍO DE EMAIL =====")
	log.Printf("Intentando enviar email a: %s", to)
	log.Printf("Usando configuración:")
	log.Printf("- HOST: %s", config.SMTPHost)
	log.Printf("- PORT: %s", config.SMTPPort)
	log.Printf("- FROM: %s", config.From)
	log.Printf("- PASSWORD CONFIGURADA: %v", config.Password != "")
	log.Printf("- APP_ENV: %s", GetEnv("APP_ENV", "development"))
	
	if config.Password == "" {
		log.Println("ERROR: No hay contraseña de email configurada")
		log.Println("===== FIN LOG DE ENVÍO DE EMAIL (ERROR) =====")
		return ErrEmailSendError
	}

	// Intentar crear autenticación SMTP
	auth := smtp.PlainAuth("", config.From, config.Password, config.SMTPHost)

	// Construir el mensaje de correo
	headers := "MIME-version: 1.0;\r\nContent-Type: " + contentType + ";\r\n"
	msg := "To: " + to + "\r\n" +
		"From: " + config.From + "\r\n" +
		"Subject: " + subject + "\r\n" +
		headers + "\r\n" +
		body

	// Para diagnóstico, guardar una copia del correo en un archivo
	// independientemente del entorno
	filename := fmt.Sprintf("email_to_%s_%d.html", 
		strings.Replace(to, "@", "_at_", -1), 
		time.Now().Unix())
	if err := os.WriteFile(filename, []byte(body), 0644); err != nil {
		log.Printf("No se pudo guardar una copia del correo: %v", err)
	} else {
		log.Printf("Copia del correo guardada en: %s", filename)
	}

	// Preparación para envío de correo real
	log.Println("Intentando conectar y enviar email real...")
	log.Printf("Dirección SMTP: %s:%s", config.SMTPHost, config.SMTPPort)
	
	// Intentar enviar el correo
	err := smtp.SendMail(
		config.SMTPHost+":"+config.SMTPPort, 
		auth, 
		config.From, 
		[]string{to}, 
		[]byte(msg),
	)
	
	if err != nil {
		log.Printf("ERROR AL ENVIAR EMAIL: %v", err)
		log.Printf("Error de tipo: %T", err)
		
		// Detectar errores comunes de SMTP
		errMsg := err.Error()
		if strings.Contains(errMsg, "535") || strings.Contains(errMsg, "534") {
			log.Println("ERROR DE AUTENTICACIÓN: Credenciales rechazadas por el servidor.")
			log.Println("Para Gmail, necesitas:")
			log.Println("1. Habilitar 'Acceso a aplicaciones menos seguras' en tu cuenta de Google, o")
			log.Println("2. Usar una 'Contraseña de aplicaciones' si tienes la verificación en dos pasos habilitada.")
			log.Println("Visita: https://myaccount.google.com/security")
		} else if strings.Contains(errMsg, "530") {
			log.Println("ERROR: Autenticación requerida. El servidor SMTP requiere autenticación.")
		} else if strings.Contains(errMsg, "553") {
			log.Println("ERROR: El remitente no fue aceptado. Verifica que la dirección de correo remitente sea válida.")
		} else if strings.Contains(errMsg, "554") {
			log.Println("ERROR: Transacción fallida. El servidor rechazó el mensaje (posible spam o problemas de entrega).")
		} else if strings.Contains(errMsg, "connection refused") {
			log.Printf("ERROR: Conexión rechazada al host %s en puerto %s", config.SMTPHost, config.SMTPPort)
			log.Println("Posibles causas:")
			log.Println("1. El servidor SMTP no está ejecutándose o no está accesible.")
			log.Println("2. Un firewall está bloqueando la conexión.")
			log.Println("3. El puerto SMTP está mal configurado.")
		}
		
		log.Println("===== FIN LOG DE ENVÍO DE EMAIL (ERROR) =====")
		return fmt.Errorf("error al enviar email: %w", err)
	}

	log.Println("Correo enviado correctamente")
	log.Println("===== FIN LOG DE ENVÍO DE EMAIL (ÉXITO) =====")
	return nil
}

// RenderEmailTemplate renderiza una plantilla de correo electrónico con los datos proporcionados
func RenderEmailTemplate(templateName string, data interface{}) (string, error) {
	log.Printf("Renderizando plantilla de email: %s", templateName)
	
	var templateContent string

	// En una aplicación real, cargaríamos la plantilla desde un archivo
	// Por simplicidad, definimos algunas plantillas aquí
	switch templateName {
	case "reset_password":
		templateContent = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Recuperación de contraseña</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; }
        .header { background-color: #4a86e8; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .button { display: inline-block; background-color: #4a86e8; color: white; text-decoration: none; padding: 10px 20px; border-radius: 5px; }
        .footer { text-align: center; padding: 20px; font-size: 0.8em; color: #777; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Recuperación de contraseña</h1>
    </div>
    <div class="content">
        <p>Hola {{.Name}},</p>
        <p>Has solicitado restablecer tu contraseña. Haz clic en el siguiente enlace para crear una nueva contraseña:</p>
        <p><a href="{{.ResetLink}}" class="button">Restablecer contraseña</a></p>
        <p>O copia y pega este enlace en tu navegador:</p>
        <p>{{.ResetLink}}</p>
        <p>Si no solicitaste este cambio, puedes ignorar este correo.</p>
        <p>El enlace expirará en 24 horas.</p>
        <p>Saludos cordiales,<br>El equipo de soporte</p>
    </div>
    <div class="footer">
        <p>Este es un mensaje automático, por favor no respondas a este correo.</p>
    </div>
</body>
</html>`
	case "contact_reply":
		templateContent = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Respuesta a tu mensaje</title>
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
        <h1>Respuesta a tu mensaje</h1>
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
	default:
		return "", errors.New("plantilla no encontrada")
	}

	// Crear el template desde el string
	tmpl, err := template.New("email_template").Parse(templateContent)
	if err != nil {
		log.Printf("Error al parsear plantilla: %v", err)
		return "", err
	}

	var htmlBuffer bytes.Buffer
	if err := tmpl.Execute(&htmlBuffer, data); err != nil {
		log.Printf("Error al ejecutar plantilla: %v", err)
		return "", err
	}

	log.Println("Plantilla renderizada correctamente")
	return htmlBuffer.String(), nil
}