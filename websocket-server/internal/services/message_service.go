package services

import (
	"database/sql"
	"errors"
	"time"
	"websocket-server/internal/models"
	"websocket-server/internal/repository"
)

type MessageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) SendMessage(senderID, recipientID, content string) (models.Message, error) {
	if senderID == "" || recipientID == "" || content == "" {
		return models.Message{}, errors.New("sender, recipient, and content must not be empty")
	}

	message := models.Message{
		SenderID:    senderID,
		RecipientID: recipientID,
		Content:     content,
		CreatedAt:   time.Now(),
	}

	err := s.repo.StoreMessage(message)
	if err != nil {
		return models.Message{}, err
	}

	return message, nil
}

func (s *MessageService) GetMessagesForUser(userID string, limit int) ([]models.Message, error) {
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	messages, err := s.repo.FetchMessagesForUser(userID, limit)
	if err != nil {
		return nil, err
	}

	return messages, nil
}