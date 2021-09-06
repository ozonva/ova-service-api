package saver

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ozonva/ova-service-api/internal/models"
)

type Flusher interface {
	Flush(ctx context.Context, services []models.Service) []models.Service
}

type Saver interface {
	Save(service models.Service) error
	Init()
	Close()
}

func New(capacity uint, flushTimeout time.Duration, flusher Flusher) Saver {
	return &saver{
		localStorage: make([]models.Service, 0, capacity),
		flushTimeout: flushTimeout,
		flusher:      flusher,
	}
}

type saver struct {
	sync.Mutex
	signalChannel chan struct{}
	localStorage  []models.Service
	flushTimeout  time.Duration
	flusher       Flusher
}

func (s *saver) Save(service models.Service) error {
	s.Lock()
	defer s.Unlock()

	if len(s.localStorage) == cap(s.localStorage) {
		return fmt.Errorf("local storage is full, please wait for the next flash operation")
	}

	s.localStorage = append(s.localStorage, service)
	return nil
}

func (s *saver) Init() {
	// Initialize channel in init instead of constructor to keep parity with close method.
	// For example, you may call Close and thus close the channel, and then re-create it by calling Init
	// keeping the single saver object.
	s.signalChannel = make(chan struct{})

	go func(ch <-chan struct{}) {
		ticker := time.NewTicker(s.flushTimeout)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.flush()
			case _, ok := <-ch:
				if !ok {
					return
				}
			}
		}
	}(s.signalChannel)
}

func (s *saver) Close() {
	s.flush()
	close(s.signalChannel)
}

func (s *saver) flush() {
	// We need lock here because it is possible situation when timeout and close events occur in the same time.
	// In this case we are possibly could flush the same slice (localStorage) twice without the lock.
	s.Lock()
	defer s.Unlock()

	if len(s.localStorage) == 0 {
		return
	}

	unsaved := s.flusher.Flush(context.Background(), s.localStorage)
	if len(unsaved) > 0 {
		log.Printf("warning: some entities can't be saved to database and will be discraded: \n%v\n", unsaved)
	}

	s.localStorage = make([]models.Service, 0, cap(s.localStorage))
}
