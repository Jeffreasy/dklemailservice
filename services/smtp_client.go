package services

import (
	"fmt"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

// SMTPDialer is een interface voor het testen van SMTP verbindingen
type SMTPDialer interface {
	Dial() error
}

// SMTPConfig bevat de configuratie voor een SMTP verbinding
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// RealSMTPClient implementeert de SMTP client interface
type RealSMTPClient struct {
	defaultConf *SMTPConfig
	regConf     *SMTPConfig
	dialer      SMTPDialer
}

// NewRealSMTPClient creates a new SMTP client with the given configuration
func NewRealSMTPClient(host, port, user, password, from, regUser, regPassword string) *RealSMTPClient {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		portNum = 587 // Default SMTP port
	}

	defaultConf := &SMTPConfig{
		Host:     host,
		Port:     portNum,
		Username: user,
		Password: password,
		From:     from,
	}

	regConf := &SMTPConfig{
		Host:     host,
		Port:     portNum,
		Username: regUser,
		Password: regPassword,
		From:     from,
	}

	return &RealSMTPClient{
		defaultConf: defaultConf,
		regConf:     regConf,
	}
}

// SetDialer stelt een custom dialer in voor tests
func (c *RealSMTPClient) SetDialer(d SMTPDialer) {
	c.dialer = d
}

// Send verzendt een email met de standaard configuratie
func (c *RealSMTPClient) Send(msg *EmailMessage) error {
	return c.sendWithConfig(msg, c.defaultConf)
}

// SendRegistration verzendt een email met de registratie configuratie
func (c *RealSMTPClient) SendRegistration(msg *EmailMessage) error {
	return c.sendWithConfig(msg, c.regConf)
}

// sendWithConfig verzendt een email met de opgegeven configuratie
func (c *RealSMTPClient) sendWithConfig(msg *EmailMessage, conf *SMTPConfig) error {
	if conf == nil {
		return fmt.Errorf("smtp configuration is nil")
	}

	if msg.To == "" {
		return fmt.Errorf("invalid recipient")
	}

	// Als we een test dialer hebben, gebruik die
	if c.dialer != nil {
		if err := c.dialer.Dial(); err != nil {
			return err
		}
		if mockSMTP, ok := c.dialer.(SMTPClient); ok {
			return mockSMTP.Send(msg)
		}
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", conf.From)
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Body)

	d := gomail.NewDialer(conf.Host, conf.Port, conf.Username, conf.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	return nil
}

// SendEmail is een helper functie voor backwards compatibility
func (c *RealSMTPClient) SendEmail(to, subject, body string) error {
	msg := &EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}
	return c.Send(msg)
}
