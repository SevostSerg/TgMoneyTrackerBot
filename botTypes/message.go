package botTypes

import (
	"errors"
	"strings"

	tgAPI "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Message struct {
	Message     string
	MessageType MessageType
	LoopMessage bool
}

type MessageType int

const (
	AddWastes MessageType = 0
	AddIncome MessageType = 1
	AddLimit  MessageType = 2
	Unknown   MessageType = 3
)

func IdentifyMessageType(update *tgAPI.Update) (*Message, error) {
	splittedMessage := strings.Split(update.Message.Text, " ")
	if len(splittedMessage) != 0 {
		switch splittedMessage[0] {
		case "+":
			return &Message{
				Message:     splittedMessage[1],
				MessageType: AddIncome,
				LoopMessage: false,
			}, nil
		case "-":
			return &Message{
				Message:     splittedMessage[1],
				MessageType: AddWastes,
				LoopMessage: false,
			}, nil
		}
	}

	if update.Message.Text == "Add limit" || update.Message.Text == "Добавить предел трат" {
		return &Message{
			Message:     update.Message.Text,
			MessageType: AddLimit,
			LoopMessage: false,
		}, nil
	}

	return nil, errors.New("Unknown message!")
}
