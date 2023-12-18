package main

import "context"

type loggingService struct {
	next JwtGenerator
}

func (s *loggingService) GenerateJWT(ctx context.Context, session string) (string, error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(begin),
			"error": err,
			"session": session,
		})
	}
}