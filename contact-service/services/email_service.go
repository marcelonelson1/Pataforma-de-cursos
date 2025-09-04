package services

import (
	"bytes"
	"contact-service/config"
	"contact-service/models"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

type EmailService struct{}

// NewEmailService crea una nueva instancia de EmailService
func NewEmailService() *EmailService {
	return &EmailService{}
}

// SendContactNotification envia notificacion de nuevo mensaje de contacto
func (es *EmailService) SendContactNotification(contact *models.ContactRequest) error {
	cfg := config.AppConfig

	// En desarrollo, solo guardar en archivo
	if cfg.AppEnv == "development" {
		return es.saveContactEmailToFile(contact)
	}

	// Validar configuracion
	if cfg.EmailPassword == "" {
		return fmt.Errorf("EMAIL_PASSWORD no configurado")
	}

	// Preparar autenticacion
	auth := smtp.PlainAuth("", cfg.EmailFrom, cfg.EmailPassword, cfg.SMTPHost)

	// Crear contenido del email
	emailBody := es.createContactEmailBody(contact)

	// Preparar mensaje
	headers := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	subject := fmt.Sprintf("Nuevo mensaje de contacto - %s", contact.Name)
	
	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n%s\r\n%s",
		cfg.ContactEmail, cfg.EmailFrom, subject, headers, emailBody)

	// Enviar email
	serverAddr := cfg.SMTPHost + ":" + cfg.SMTPPort
	return smtp.SendMail(serverAddr, auth, cfg.EmailFrom, []string{cfg.ContactEmail}, []byte(msg))
}

// SendReplyEmail envia respuesta a un mensaje de contacto
func (es *EmailService) SendReplyEmail(message *models.ContactMessage, replyText string) error {
	cfg := config.AppConfig

	// En desarrollo, guardar en archivo
	if cfg.AppEnv == "development" {
		return es.saveReplyEmailToFile(message, replyText)
	}

	// Validar configuracion
	if cfg.EmailPassword == "" {
		return fmt.Errorf("EMAIL_PASSWORD no configurado")
	}

	// Preparar autenticacion
	auth := smtp.PlainAuth("", cfg.EmailFrom, cfg.EmailPassword, cfg.SMTPHost)

	// Crear contenido del email
	emailBody := es.createReplyEmailBody(message, replyText)

	// Preparar mensaje
	headers := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	subject := fmt.Sprintf("Respuesta a tu mensaje - %s", message.Name)
	
	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n%s\r\n%s",
		message.Email, cfg.EmailFrom, subject, headers, emailBody)

	// Enviar email
	serverAddr := cfg.SMTPHost + ":" + cfg.SMTPPort
	return smtp.SendMail(serverAddr, auth, cfg.EmailFrom, []string{message.Email}, []byte(msg))
}

// createContactEmailBody crea el cuerpo del email de contacto usando template externo
func (es *EmailService) createContactEmailBody(contact *models.ContactRequest) string {
	// Cargar la plantilla desde el archivo
	tmpl, err := template.ParseFiles("templates/contact_email.html")
	if err != nil {
		// Fallback al template simple si no se encuentra el archivo
		return es.createSimpleContactEmailBody(contact)
	}

	// Ejecutar la plantilla con los datos
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, contact); err != nil {
		return fmt.Sprintf("Error ejecutando template: %v", err)
	}

	return buf.String()
}

// createSimpleContactEmailBody template de fallback
func (es *EmailService) createSimpleContactEmailBody(contact *models.ContactRequest) string {
	contactTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Nuevo mensaje de contacto</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; }
        .header { background-color: #4a86e8; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .info { background-color: #f5f5f5; padding: 15px; border-left: 4px solid #4a86e8; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; font-size: 0.8em; color: #777; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Nuevo mensaje de contacto</h1>
    </div>
    <div class="content">
        <h2>Detalles del contacto:</h2>
        <div class="info">
            <p><strong>Nombre:</strong> {{.Name}}</p>
            <p><strong>Email:</strong> {{.Email}}</p>
            {{if .Phone}}<p><strong>Telefono:</strong> {{.Phone}}</p>{{end}}
        </div>
        
        <h3>Mensaje:</h3>
        <p>{{.Message}}</p>
    </div>
    <div class="footer">
        <p>Este es un mensaje automatico del sistema de contacto.</p>
    </div>
</body>
</html>`

	tmpl, err := template.New("contact").Parse(contactTemplate)
	if err != nil {
		return fmt.Sprintf("Error en template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, contact); err != nil {
		return fmt.Sprintf("Error ejecutando template: %v", err)
	}

	return buf.String()
}

// createReplyEmailBody crea el cuerpo del email de respuesta usando template externo
func (es *EmailService) createReplyEmailBody(message *models.ContactMessage, replyText string) string {
	// Cargar la plantilla desde el archivo
	tmpl, err := template.ParseFiles("templates/reply_email.html")
	if err != nil {
		// Fallback al template simple si no se encuentra el archivo
		return es.createSimpleReplyEmailBody(message, replyText)
	}

	// Preparar datos para el template moderno
	data := struct {
		Subject      string
		Name         string
		Reply        string
		BusinessName string
		CurrentYear  string
	}{
		Subject:      "Respuesta a tu mensaje",
		Name:         message.Name,
		Reply:        replyText,
		BusinessName: "RM Renders",
		CurrentYear:  "2025",
	}

	// Ejecutar la plantilla con los datos
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Sprintf("Error ejecutando template: %v", err)
	}

	return buf.String()
}

// createSimpleReplyEmailBody template de fallback para respuestas
func (es *EmailService) createSimpleReplyEmailBody(message *models.ContactMessage, replyText string) string {
	replyTemplate := `
<!DOCTYPE html>
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
        <p>Gracias por contactarnos. A continuacion encontraras nuestra respuesta a tu mensaje:</p>
        
        <div class="original-message">
            <p><strong>Tu mensaje original:</strong></p>
            <p>{{.Message}}</p>
        </div>
        
        <p><strong>Nuestra respuesta:</strong></p>
        <p>{{.ReplyMsg}}</p>
        
        <p>Saludos cordiales,<br>El equipo de soporte</p>
    </div>
    <div class="footer">
        <p>Este es un mensaje automatico, por favor no respondas a este correo.</p>
    </div>
</body>
</html>`

	data := models.EmailTemplateData{
		Name:        message.Name,
		Email:       message.Email,
		OriginalMsg: message.Message,
		ReplyMsg:    replyText,
	}

	tmpl, err := template.New("reply").Parse(replyTemplate)
	if err != nil {
		return fmt.Sprintf("Error en template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Sprintf("Error ejecutando template: %v", err)
	}

	return buf.String()
}

// saveContactEmailToFile guarda email de contacto en archivo (desarrollo)
func (es *EmailService) saveContactEmailToFile(contact *models.ContactRequest) error {
	content := es.createContactEmailBody(contact)
	return os.WriteFile("last_contact_email.html", []byte(content), 0644)
}

// saveReplyEmailToFile guarda email de respuesta en archivo (desarrollo)
func (es *EmailService) saveReplyEmailToFile(message *models.ContactMessage, replyText string) error {
	content := es.createReplyEmailBody(message, replyText)
	return os.WriteFile("last_reply_email.html", []byte(content), 0644)
}