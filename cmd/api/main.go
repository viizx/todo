package main

import (
	"evolutio_to-do/internal/database"
	"evolutio_to-do/internal/notification"
	"evolutio_to-do/internal/server"
	"fmt"
	"os"
)

func main() {
	smsService := notification.NewInfobipSMSService(
		os.Getenv("INFOBIP_API_KEY"),
		os.Getenv("SMS_SENDER"),
		fmt.Sprintf("https://%s/sms/2/text/advanced", os.Getenv("INFOBIP_API_BASE_URL")),
	)

	go smsService.HandleErrors()

	server := server.NewServer(smsService)
	store := database.New()
	store.Init()
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
