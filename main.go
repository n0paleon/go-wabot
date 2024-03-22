package main

import (
	"TuruBot/configs"
	"TuruBot/handler"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	//_ "github.com/jackc/pgx/v4/stdlib" // <= pgsql driver pake ini
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "modernc.org/sqlite"
)


func main() {
	dbLog := waLog.Stdout("Database", "INFO", true)

	// testing config
	// jangan lupa ubah pake pgsql kalo udah production
	sessionName :=  fmt.Sprintf("%s.sqlite?_pragma=foreign_keys=1&_journal_mode=WAL", configs.GetEnv("SESSION_NAME")) // fmt.Sprintf("postgresql://%s:%s@%s/%s", configs.GetEnv("DB_USER"), configs.GetEnv("DB_PASS"), configs.GetEnv("DB_HOST"), configs.GetEnv("DB_NAME"))
	container, err := sqlstore.New("sqlite", sessionName, dbLog)
	if err != nil {
		panic(err)
	}
	
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	myClient := handler.MyClient{
		WAClient: client,
	}
	myClient.Register()


	if client.Store.ID == nil {
		// login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				fmt.Println("QR code received, please scan the qr code!")
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// reconnect, existing session
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	// channel untuk exit terima sinyal exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}