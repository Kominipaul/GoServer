package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"path/filepath"
)

// ContactHandler is an api that gets user form date and sends an email to admin
func ContactHandler(w http.ResponseWriter, r *http.Request) {
	//Check if the request method is GET using tagger switch on r.Method
	if r.Method == http.MethodGet {
		tmplPath := filepath.Join("web", "templates", "contact.html")
		renderTemplate(w, tmplPath, nil)
	} else if r.Method == http.MethodPost {
		name := r.FormValue("username")
		email := r.FormValue("email")
		message := r.FormValue("message")

		// Send email to admin
		println(name, email, message)
		err := sendEmail(name, email, message)
		if err != nil {
			log.Printf("Error sending email: %v", err)
			http.Error(w, "Error sending email", http.StatusInternalServerError)
			return
		}

		// Redirect to the home page after successful submission
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// sendEmail function sends an email to the admin
func sendEmail(name, email, message string) error {
	from := "kominipaul@gmail.com"
	password := "drpxlebrabjwtgby"

	// SMTP server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	msg := "From: " + from + "\n" +
		"To: " + "kominipaulo@gmail.com" + "\n" +
		"Subject: Contact Form Submission\n\n" +
		"Name: " + name + "\n" +
		"Email: " + email + "\n" +
		"Message: " + message

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{"kominipaulo@gmail.com"}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Print("Email sent successfully")
	return nil
}
