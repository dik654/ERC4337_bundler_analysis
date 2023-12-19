package main

import "context"

type metricService struct {
	next OtpAuthenticator
}

func NewMetricService(next OtpAuthenticator) OtpAuthenticator {
	return &metricService{
		next: next,
	}
}

func (s *metricService) GeneratePrivateKey(ctx context.Context, id string) (string, error) {
	return s.next.GeneratePrivateKey(ctx, id)
}

func (s *metricService) GenerateOtp(ctx context.Context, id string) (string, error) {
	return s.next.GenerateOtp(ctx, id)
}
