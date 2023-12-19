package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type loggingService struct {
	next OtpAuthenticator
}

func NewLoggingService(next OtpAuthenticator) OtpAuthenticator {
	return &loggingService{
		next: next,
	}
}

func (s *loggingService) GeneratePrivateKey(ctx context.Context, id string) (privateKey string, err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":  time.Since(begin),
			"error": err,
			"id":    id,
		})
	}(time.Now())

	return s.next.GeneratePrivateKey(ctx, id)
}

func (s *loggingService) GenerateOtp(ctx context.Context, id string) (otp string, err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":  time.Since(begin),
			"error": err,
			"id":    id,
		})
	}(time.Now())

	return s.next.GenerateOtp(ctx, id)
}
