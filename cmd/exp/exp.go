package main

import (
	"fmt"

	"github.com/alorents/lenslocked/models"
)

// Mailtrap credentials
const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 587
	username = "e3d8389ef2588d"
	password = "ccc570e0112cd6"
)

func main() {
	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
	err := es.ForgotPassword("andreasphoenix@gmail.com", "fakeresettoken")
	if err != nil {
		panic(err)
	}
	fmt.Println("Email sent!")
}
