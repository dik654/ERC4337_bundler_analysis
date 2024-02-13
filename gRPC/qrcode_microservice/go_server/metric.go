package main

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	GeneratePrivateKeyCounter prometheus.Counter
	GenerateOtpCounter        prometheus.Counter
	VerifyOtpCounter          prometheus.Counter
)

func init() {
	GeneratePrivateKeyCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "generate_private_key_count",
		Help: "generate private key count",
	})

	GenerateOtpCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "generate_otp_count",
		Help: "generate otp count",
	})

	VerifyOtpCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "verify_otp_count",
		Help: "verify count",
	})
}

type metricService struct {
	next OtpAuthenticator
}

func NewMetricService(next OtpAuthenticator) OtpAuthenticator {
	return &metricService{
		next: next,
	}
}

func (s *metricService) GeneratePrivateKey(ctx context.Context, id string) (string, error) {
	GeneratePrivateKeyCounter.Inc()

	return s.next.GeneratePrivateKey(ctx, id)
}

func (s *metricService) GenerateOtp(ctx context.Context, id string) (string, error) {
	GenerateOtpCounter.Inc()
	return s.next.GenerateOtp(ctx, id)
}

func (s *metricService) VerifyOtp(ctx context.Context, id string, otp string) (bool, error) {
	VerifyOtpCounter.Inc()
	return s.next.VerifyOtp(ctx, id, otp)
}
