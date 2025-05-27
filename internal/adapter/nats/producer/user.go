package producer

import (
	"context"
	"github.com/sorawaslocked/ap2final_base/pkg/nats"
	"github.com/sorawaslocked/ap2final_user_service/internal/adapter/nats/producer/dto"
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
	"google.golang.org/protobuf/proto"
	"time"
)

const PushTimeout = time.Second * 30

type UserProducer struct {
	natsClient *nats.Client
	subject    string
}

func NewUserProducer(natsClient *nats.Client, subject string) *UserProducer {
	return &UserProducer{
		natsClient: natsClient,
		subject:    subject,
	}
}

func (p *UserProducer) Push(ctx context.Context, user model.User) error {
	ctx, cancel := context.WithTimeout(ctx, PushTimeout)
	defer cancel()

	pbUser := dto.FromUserToRegisterEvent(user)
	data, err := proto.Marshal(pbUser)
	if err != nil {
		return err
	}

	err = p.natsClient.Conn.Publish(p.subject, data)
	if err != nil {
		return err
	}

	return nil
}
