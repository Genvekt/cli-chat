package converter

import (
	"encoding/json"

	"github.com/IBM/sarama"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/auth_producer/internal/model"
)

// UserToKafkaUser converts user service model to kafka message model
func UserToKafkaUser(user *model.User) *userApi.CreateRequest {
	return &userApi.CreateRequest{
		Password:        user.Password,
		PasswordConfirm: user.Password,
		Info: &userApi.UserInfo{
			Name:  user.Name,
			Email: user.Email,
		},
	}
}

// StructToMsg converts struct to kafka json message
func StructToMsg(obj interface{}, topic string) (*sarama.ProducerMessage, error) {
	encodedUser, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(encodedUser),
	}

	return msg, nil
}
