package service

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"time"

	"github.com/romanchechyotkin/avito_test_task/internal/repo"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"
)

type SendEmailInput struct {
	Recipient string
	Message   string
}

// todo close chan
type SenderService struct {
	log *slog.Logger

	notifyChan chan uint
	houseRepo  repo.House
}

func NewSenderService(log *slog.Logger, houseRepo repo.House) *SenderService {
	service := &SenderService{
		log:        log,
		notifyChan: make(chan uint),
		houseRepo:  houseRepo,
	}

	go service.Run()

	return service
}

func (s *SenderService) Run() {
	for {
		select {
		case houseID, ok := <-s.notifyChan:
			if !ok {
				s.log.Debug("notify channel is closed")
				return
			}

			s.log.Debug("got request to send", slog.Any("house id", houseID))

			emailsToSend, err := s.houseRepo.GetHouseSubscriptions(context.Background(), houseID)
			if err != nil {
				s.log.Error("failed to get emails to send", logger.Error(err))
				s.notifyChan <- houseID
				continue
			}

			for _, email := range emailsToSend {
				s.log.Debug("sent email", slog.String("email", email))
				err = s.sendEmail(context.Background(), email, "MESSAGE")
				if err != nil {
					// todo append in slice to repeat infinity times
					continue
				}
			}
		}
	}
}

func (s *SenderService) sendEmail(ctx context.Context, recipient string, message string) error {
	// Имитация отправки сообщения
	duration := time.Duration(rand.Int63n(3000)) * time.Millisecond
	time.Sleep(duration)

	// Имитация неуспешной отправки сообщения
	errorProbability := 0.1
	if rand.Float64() < errorProbability {
		return errors.New("internal error")
	}

	s.log.Debug("send message", slog.String("message", message), slog.String("recipient", recipient))

	return nil
}

func (s *SenderService) Notify() chan<- uint {
	return s.notifyChan
}
