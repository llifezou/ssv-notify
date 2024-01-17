package lark

import (
	"bytes"
	"encoding/json"
	"golang.org/x/xerrors"
	"net/http"
)

type Client struct {
	webhook string
}

func NewLarkClient(webhook string) *Client {
	return &Client{
		webhook: webhook,
	}
}

func (*Client) Platform() string {
	return "lark"
}

type message struct {
	MsgType string  `json:"msg_type"`
	Content Content `json:"content"`
}

type Content struct {
	Text string `json:"text"`
}

func newMessage(msg string) *message {
	return &message{
		"text",
		Content{
			Text: msg,
		},
	}
}

func (c *Client) Send(msg string) error {
	msgByte, err := json.Marshal(newMessage(msg))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.webhook, bytes.NewBuffer(msgByte))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	return xerrors.Errorf("%s", resp.Status)
}
