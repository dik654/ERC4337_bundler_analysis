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

func (s *loggingService) GeneratePrivateKey(ctx context.Context, id string) (url string, err error) {
	defer func(begin time.Time) {
		fields := logrus.Fields{
			"took": time.Since(begin),
			"id":   id,
		}
		if err != nil {
			fields["error"] = err
			logrus.WithFields(fields).Error("GeneratePrivateKey failed")
		} else {
			logrus.WithFields(fields).Info("GeneratePrivateKey success")
		}
	}(time.Now())

	return s.next.GeneratePrivateKey(ctx, id)
}

func (s *loggingService) GenerateOtp(ctx context.Context, id string) (otp string, err error) {
	defer func(begin time.Time) {
		fields := logrus.Fields{
			"took": time.Since(begin),
			"id":   id,
		}
		if err != nil {
			fields["error"] = err
			logrus.WithFields(fields).Error("GenerateOtp failed")
		} else {
			logrus.WithFields(fields).Info("GenerateOtp success")
		}
	}(time.Now())

	return s.next.GenerateOtp(ctx, id)
}

func (s *loggingService) VerifyOtp(ctx context.Context, id string, otp string) (verification bool, err error) {
	defer func(begin time.Time) {
		fields := logrus.Fields{
			"took": time.Since(begin),
			"id":   id,
		}
		if err != nil {
			fields["error"] = err
			logrus.WithFields(fields).Error("VeirifyOtp failed")
		} else {
			logrus.WithFields(fields).Info("VeirifyOtp success")
		}
	}(time.Now())

	return s.next.VerifyOtp(ctx, id, otp)
}
