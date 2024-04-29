package controller

import (
	"backup-client/service"
	"context"

	v1 "backup-client/api/helloworld/v1"
)

// GreeterService is a greeter controller.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *service.GreeterUsecase
}

// NewGreeterService new a greeter controller.
func NewGreeterService(uc *service.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &service.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}
