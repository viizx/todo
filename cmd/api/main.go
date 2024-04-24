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
		"+447491163443",
		"https://vv8dxv.api.infobip.com/sms/2/text/advanced",
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
