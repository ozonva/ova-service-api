package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Event types
const (
	Create = iota
	Update
	Delete
)

// ServiceCUDEvent should be more complicated in real life, e.g. should contain full created service entity,
// but we will only provide serviceID to keep it simple for now.
type ServiceCUDEvent struct {
	EventID   uuid.UUID
	EventType int
	ServiceID uuid.UUID
	Timestamp time.Time
}

func NewServiceCreateEvent(serviceID uuid.UUID) ServiceCUDEvent {
	return newServiceCUDEvent(Create, serviceID)
}

func NewServiceUpdateEvent(serviceID uuid.UUID) ServiceCUDEvent {
	return newServiceCUDEvent(Update, serviceID)
}

func NewServiceDeleteEvent(serviceID uuid.UUID) ServiceCUDEvent {
	return newServiceCUDEvent(Delete, serviceID)
}

func (event ServiceCUDEvent) String() string {
	res, err := json.Marshal(event)

	// Do not return error to satisfy stringer interface.
	// Actually we should not have errors on marshalling ServiceCUDEvent structure.
	if err != nil {
		log.Err(err).Msg("error occurs during marshaling ServiceCUDEvent")
		return ""
	}

	return string(res)
}

func newServiceCUDEvent(eventType int, serviceID uuid.UUID) ServiceCUDEvent {
	return ServiceCUDEvent{
		EventID:   uuid.New(),
		EventType: eventType,
		ServiceID: serviceID,
		Timestamp: time.Now().UTC(),
	}
}
