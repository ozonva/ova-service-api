package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID             uuid.UUID
	UserID         uint64
	Description    string
	ServiceName    string
	ServiceAddress string
	WhenLocal      *time.Time
	WhenUTC        *time.Time
}

func NewService(userID uint64, description string, serviceName string, serviceAddress string, when *time.Time) (*Service, error) {
	if userID == 0 {
		return nil, fmt.Errorf("can't create service entry for non-existing user")
	}

	service := &Service{
		ID:             uuid.New(),
		UserID:         userID,
		Description:    description,
		ServiceName:    serviceName,
		ServiceAddress: serviceAddress,
	}

	err := service.UpdateCalendar(when)

	if err != nil {
		return nil, err
	}

	return service, nil
}

func (service *Service) UpdateDescription(description string) {
	service.Description = description
}

func (service *Service) UpdateService(serviceName string, serviceAddress string) {
	service.ServiceName = serviceName
	service.ServiceAddress = serviceAddress
}

func (service *Service) UpdateCalendar(when *time.Time) error {
	if when == nil {
		service.WhenLocal = nil
		service.WhenUTC = nil
		return nil
	}

	if !when.After(time.Now()) {
		return fmt.Errorf("can't update calendar to the date in the past")
	}

	service.WhenLocal = when
	utc := when.UTC()
	service.WhenUTC = &utc

	return nil
}

func (service *Service) String() string {
	var (
		local time.Time
		utc   time.Time
	)

	if service.WhenLocal != nil {
		local = *service.WhenLocal
	}

	if service.WhenUTC != nil {
		utc = *service.WhenUTC
	}

	return fmt.Sprintf(`Service entry:
	ID: 			%v
	UserID: 		%d
	Description: 	%s
	ServiceName: 	%s
	ServiceAddress: %s
	WhenLocal: 		%v
	WhenUTC: 		%v
`, service.ID, service.UserID, service.Description, service.ServiceName, service.ServiceAddress, local, utc)
}
