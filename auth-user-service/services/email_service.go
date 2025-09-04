package services

import (
	"auth-user-service/config"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"text/template"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

// SendPasswordResetEmail envía un email de recuperación de contraseña
func (e *EmailService) SendPasswordResetEmail(to, name, resetLink string) error {
	// Verificar configuración
	if config.AppConfig.EmailPassword == "" {
		return errors.New("la contraseña de email no está configurada")
	}

	// Configurar autenticación
	auth := smtp.PlainAuth("", 
		config.AppConfig.EmailFrom, 
		config.AppConfig.EmailPassword, 
		config.AppConfig.SMTPHost)

	// Renderizar el HTML del correo electrónico
	htmlContent, err := e.renderPasswordResetTemplate(name, resetLink)
	if err != nil {
		log.Printf("Error al renderizar la plantilla de correo: %v", err)
		return err
	}

	// Construir el mensaje
	subject := "Recuperación de contraseña"
	headers := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	message := []byte("To: " + to + "\r\n" +
		"From: " + config.AppConfig.EmailFrom + "\r\n" +
		"Subject: " + subject + "\r\n" +
		headers + "\r\n" +
		htmlContent)

	// En desarrollo, podemos simular el envío
	if config.AppConfig.AppEnv == "development" && config.AppConfig.MockEmail == "true" {
		log.Printf("Simulando envío de email a %s. Contenido: %s", to, htmlContent)
		// Guardar una copia del correo en un archivo para probar
		f, err := os.OpenFile("last_reset_email.html", os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			f.WriteString(htmlContent)
		}
		return nil
	}

	// Enviar email
	smtpAddr := config.AppConfig.SMTPHost + ":" + config.AppConfig.SMTPPort
	err = smtp.SendMail(smtpAddr, auth, config.AppConfig.EmailFrom, []string{to}, message)
	if err != nil {
		log.Printf("Error SMTP: %v", err)
		return err
	}

	return nil
}

// renderPasswordResetTemplate renderiza la plantilla del email de reset usando archivo externo
func (e *EmailService) renderPasswordResetTemplate(name, resetLink string) (string, error) {
	// Cargar la plantilla desde el archivo (como en tu código original)
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

// SendWelcomeEmail envía un email de bienvenida (opcional)
func (e *EmailService) SendWelcomeEmail(to, name string) error {
	// Template simple de bienvenida
	subject := "¡Bienvenido!"
	body := fmt.Sprintf(`
		<h1>¡Bienvenido %s!</h1>
		<p>Tu cuenta ha sido creada exitosamente.</p>
		<p>Gracias por unirte a nosotros.</p>
	`, name)

	return e.sendEmail(to, subject, body)
}

// sendEmail función auxiliar para enviar emails
func (e *EmailService) sendEmail(to, subject, body string) error {
	if config.AppConfig.EmailPassword == "" {
		return errors.New("la contraseña de email no está configurada")
	}

	auth := smtp.PlainAuth("", 
		config.AppConfig.EmailFrom, 
		config.AppConfig.EmailPassword, 
		config.AppConfig.SMTPHost)

	headers := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	message := []byte("To: " + to + "\r\n" +
		"From: " + config.AppConfig.EmailFrom + "\r\n" +
		"Subject: " + subject + "\r\n" +
		headers + "\r\n" +
		body)

	// En desarrollo, simular envío
	if config.AppConfig.AppEnv == "development" && config.AppConfig.MockEmail == "true" {
		log.Printf("Simulando envío de email a %s con asunto: %s", to, subject)
		return nil
	}

	smtpAddr := config.AppConfig.SMTPHost + ":" + config.AppConfig.SMTPPort
	return smtp.SendMail(smtpAddr, auth, config.AppConfig.EmailFrom, []string{to}, message)
}