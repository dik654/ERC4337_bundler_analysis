package main

import "context"

func main() {
	service := NewLoggingService(NewMetricService(&otpAuthenticator{}))

	privKey, err := service.GeneratePrivateKey(context.Background(), "test")

}
