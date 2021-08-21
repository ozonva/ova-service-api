package saver

import (
	"sync"
	"time"

	"github.com/ozonva/ova-service-api/internal/models"

	flusher_ "github.com/ozonva/ova-service-api/internal/flusher"
)

const sleepTime = 100 * time.Millisecond

type Saver interface {
	Save(service models.Service)
	Init()
	Close()
}

func New(capacity uint, flushTimeout time.Duration, flusher flusher_.Flusher) Saver {
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
	flusher       flusher_.Flusher
}

func (s *saver) Save(service models.Service) {
	// If our local storage is full we are force flush our data, free up the localStorage and add new service to it.
	if len(s.localStorage) == cap(s.localStorage) {
		s.flush()
	}
	s.localStorage = append(s.localStorage, service)
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
					s.flush()
					return
				}
			default:
				time.Sleep(sleepTime)
			}
		}
	}(s.signalChannel)
}

func (s *saver) Close() {
	close(s.signalChannel)
}

func (s *saver) flush() {
	// We need lock here because it is possible situation when timeout and close events occur in the same time.
	// In this case we are possibly could flush the same slice (localStorage) twice without the lock.
	s.Lock()
	defer s.Unlock()

	// Flush some data to storage could be potentially long operation.
	// We will run it asynchronously and free up localStorage immediately.
	go func(services []models.Service) {
		_ = s.flusher.Flush(services)
	}(s.localStorage)

	s.localStorage = make([]models.Service, 0, cap(s.localStorage))
}
