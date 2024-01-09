package email

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/newtoallofthis123/noob_social/templates"
	"github.com/newtoallofthis123/noob_social/utils"
)

func SendOtp(otp string, name string, to string) error {
	env := utils.ReadEnv()

	e := email.NewEmail()
	e.From = fmt.Sprintf("Noob Social <%s>", env.Mail)
	e.To = []string{to}
	e.Subject = "Noob Social OTP"

	var htmlBuffer bytes.Buffer

	templates.SendOTP(otp).Render(context.Background(), &htmlBuffer)

	e.HTML = htmlBuffer.Bytes()

	done := make(chan error)

	go func() {
		done <- e.Send("smtp.gmail.com:587", smtp.PlainAuth("", env.Mail, env.MailPassword, "smtp.gmail.com"))
	}()

	err := <-done

	return err
}
