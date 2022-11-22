package mail

import (
	"crypto/tls"
	"log"

	"gopkg.in/gomail.v2"
)

// Client ..
type Client struct {
	Host       string
	Port       int
	User       string
	Password   string
	SenderName string
	SSLEnable  bool
}

// NewClient ..
func NewClient(host, user, password, senderName string, port int, sslEnable bool) *Client {
	return &Client{
		Host:       host,
		Port:       port,
		User:       user,
		Password:   password,
		SenderName: senderName,
		SSLEnable:  sslEnable,
	}
}

// SendMail ..
func (c *Client) SendMail(receiver []string, title, content string, filePath string) bool {
	if len(receiver) == 0 {
		return false
	}
	m := gomail.NewMessage()

	sendUser := m.FormatAddress(c.User, c.SenderName)
	// newReceiver := make([]string, 0, len(receiver))
	// for _, r := range receiver {
	// 	new := m.FormatAddress(r, "刘园")
	// 	newReceiver = append(newReceiver, new)
	// }
	m.SetHeader("From", sendUser)
	m.SetHeader("To", receiver...)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)
	m.Attach(filePath)
	d := gomail.NewDialer(c.Host, c.Port, c.User, c.Password)
	if c.SSLEnable {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("send mail error: %s", err.Error())
		return false
	}
	log.Printf("send success")
	return true
}
