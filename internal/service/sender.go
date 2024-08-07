package service

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/romanchechyotkin/avito_test_task/internal/repo"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"
)

type SendEmailInput struct {
	Recipient string
	Message   string
}

type SenderService struct {
	log *slog.Logger

	sendChan   chan uint
	notifyChan chan struct{}
	houseRepo  repo.House
}

func NewSenderService(log *slog.Logger, houseRepo repo.House) *SenderService {
	log = log.With(slog.String("component", "sender service"))

	service := &SenderService{
		log:        log,
		notifyChan: make(chan struct{}),
		sendChan:   make(chan uint),
		houseRepo:  houseRepo,
	}

	go service.Run()

	return service
}

func (s *SenderService) Run() {
	for {
		select {
		case _, ok := <-s.notifyChan:
			if !ok {
				s.log.Debug("notify channel is closed")
				return
			}

			s.log.Debug("got request to stop sending")
			close(s.sendChan)
			close(s.notifyChan)
			return

		case houseID, ok := <-s.sendChan:
			if !ok {
				s.log.Debug("send channel is closed")
				return
			}
			s.log.Debug("got request to send", slog.Any("house id", houseID))

			emailsToSend, err := s.houseRepo.GetHouseSubscriptions(context.Background(), houseID)
			if err != nil {
				s.log.Error("failed to get emails to send", logger.Error(err))
				s.sendChan <- houseID
				continue
			}

			var wg sync.WaitGroup

			for _, email := range emailsToSend {
				wg.Add(1)

				go func(email string) {
					defer wg.Done()

					s.log.Debug("sent email", slog.String("email", email))

					err = s.sendEmail(context.Background(), email, "MESSAGE")
					if err != nil {
						return
					}
				}(email)
			}

			wg.Wait()
		}
	}
}

func (s *SenderService) sendEmail(ctx context.Context, recipient string, message string) error {
	duration := time.Duration(rand.Int63n(3000)) * time.Millisecond
	time.Sleep(duration)

	errorProbability := 0.1
	if rand.Float64() < errorProbability {
		return errors.New("internal error")
	}

	s.log.Info("send message", slog.String("message", message), slog.String("recipient", recipient))

	return nil
}

func (s *SenderService) Send() chan<- uint {
	return s.sendChan
}

func (s *SenderService) Notify() chan<- struct{} {
	return s.notifyChan
}
