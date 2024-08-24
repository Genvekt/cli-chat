package user

import (
	"log"
	"net/http"

	"github.com/Genvekt/cli-chat/services/auth_producer/internal/service"
)

// Service is api of this application
type Service struct {
	userCreator service.UserCreatorService
}

// NewService creates user api instance
func NewService(userCreator service.UserCreatorService) *Service {
	return &Service{
		userCreator: userCreator,
	}
}

// HandleCreate triggers user creation
func (s *Service) HandleCreate(w http.ResponseWriter, r *http.Request) {
	err := s.userCreator.Create(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr := w.Write([]byte(err.Error()))
		if writeErr != nil {
			log.Printf("can't write response: %v", writeErr)
		}
	}
}
