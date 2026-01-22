package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct {
	auth smtp.Auth
	host string
	port string
	from string
}

func NewEmailService() *EmailService {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		host = "localhost" // Mailhog default (DEV)
	}

	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "1025" // Mailhog default (DEV)
	}

	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = "noreply@sneakersmarketplace.com"
	}

	// SMTP Authentication for Production
	var auth smtp.Auth
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")

	if username != "" && password != "" {
		// Production: Use SMTP authentication (Gmail, SendGrid, AWS SES, etc.)
		auth = smtp.PlainAuth("", username, password, host)
	} else {
		// Development: Mailhog doesn't require authentication
		auth = nil
	}

	return &EmailService{
		host: host,
		port: port,
		from: from,
		auth: auth,
	}
}

// SendEmail sends a plain text email
func (s *EmailService) SendEmail(to, subject, body string) error {
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", s.from, to, subject, body)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	return smtp.SendMail(addr, s.auth, s.from, []string{to}, []byte(msg))
}

// SendMatchCreatedEmail sends notification when match is created
func (s *EmailService) SendMatchCreatedEmail(to, role, productName, size string, price float64, orderNumber string) error {
	subject := "🎯 Your " + role + " has been matched!"

	var body string
	if role == "bid" {
		body = fmt.Sprintf(`Congratulations! 

Your bid for %s (%s) has been matched at $%.2f.

Order Number: %s

Please complete payment within 24 hours to secure your purchase.

View your order: http://localhost:8080/orders

---
Sneakers Marketplace Team
`, productName, size, price, orderNumber)
	} else {
		body = fmt.Sprintf(`Good news!

Your ask for %s (%s) has been matched at $%.2f.

Order Number: %s

Please prepare the item for shipment. You'll receive payout after delivery confirmation.

View your order: http://localhost:8080/orders

---
Sneakers Marketplace Team
`, productName, size, price, orderNumber)
	}

	return s.SendEmail(to, subject, body)
}

// SendOrderShippedEmail sends notification when order is shipped
func (s *EmailService) SendOrderShippedEmail(to, orderNumber, trackingNumber, carrier string) error {
	subject := "📦 Your order has shipped!"
	body := fmt.Sprintf(`Your order %s has been shipped!

Tracking Number: %s
Carrier: %s

Track your package: http://localhost:8080/orders/%s

Expected delivery: 3-5 business days

---
Sneakers Marketplace Team
`, orderNumber, trackingNumber, carrier, orderNumber)

	return s.SendEmail(to, subject, body)
}

// SendPaymentReceivedEmail sends notification when payment is successful
func (s *EmailService) SendPaymentReceivedEmail(to, orderNumber string, amount float64) error {
	subject := "✅ Payment confirmed"
	body := fmt.Sprintf(`Thank you! Your payment has been received.

Order Number: %s
Amount: $%.2f

Your order is now being processed.

View order: http://localhost:8080/orders

---
Sneakers Marketplace Team
`, orderNumber, amount)

	return s.SendEmail(to, subject, body)
}

// SendPayoutCompletedEmail sends notification when payout is completed
func (s *EmailService) SendPayoutCompletedEmail(to, orderNumber string, amount float64) error {
	subject := "💰 Payout completed"
	body := fmt.Sprintf(`Your payout has been completed!

Order Number: %s
Amount: $%.2f

Funds have been transferred to your account.

View details: http://localhost:8080/orders

---
Sneakers Marketplace Team
`, orderNumber, amount)

	return s.SendEmail(to, subject, body)
}
