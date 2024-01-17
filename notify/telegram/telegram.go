package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/xerrors"
	"io"
	"net/http"
)

type Client struct {
	webhook string
	chatId  string
}

func NewTelegramClient(accessToken string, chatId string) *Client {
	webhook := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", accessToken)
	return &Client{
		webhook: webhook,
		chatId:  chatId,
	}
}

func (*Client) Platform() string {
	return "telegram"
}

type message struct {
	ChatId string `json:"chat_id"`
	Msg    string `json:"text"`
}

func newMessage(chatId, msg string) *message {
	return &message{
		chatId,
		msg,
	}
}

type Response struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func (c *Client) Send(msg string) error {
	msgByte, err := json.Marshal(newMessage(c.chatId, msg))
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

	if resp.StatusCode != 200 {
		return xerrors.Errorf("%s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var r *Response
	err = json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	if !r.Ok {
		return xerrors.New(r.Description)
	}

	return nil
}
