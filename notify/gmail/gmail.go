package gmail

import (
	"net/smtp"
)

type Client struct {
	from     string
	password string
	to       string
}

func NewGmailClient(from, password, to string) *Client {
	return &Client{
		from:     from,
		password: password,
		to:       to,
	}
}

func (*Client) Platform() string {
	return "gmail"
}

func (c *Client) Send(msg string) error {
	text := "From: " + c.from + "\n" +
		"To: " + c.to + "\n" +
		"Subject: " + "Monitoring alarm" + "\n\n" + msg

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", c.from, c.password, "smtp.gmail.com"),
		c.from, []string{c.to}, []byte(text))

	if err != nil {
		return err
	}

	return nil
}
