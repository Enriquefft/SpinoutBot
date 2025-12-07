package main

import (
	"context"
	"log"
	"fmt"
    "os"
    "errors"

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
    waProto "go.mau.fi/whatsmeow/proto/waE2E"
)

const GROUP_ID = "E1kpXZL0veLFeMs6dCL9ly"

type Bot struct {
    client *whatsmeow.Client
}

// NewBot initializes a new Bot instance.
// (Initialization details like session loading or QR scanning are omitted for brevity.)
func NewBot() (*Bot, error) {
	// Replace the following with actual client setup code:
	// e.g., load session, create a database-backed store, etc.
    	db_host := os.Getenv("WHATSAPP_DATABASE_HOST")
	db_password := os.Getenv("WHATSAPP_DATABASE_PASSWORD")

	if db_host == "" || db_password == "" {
		log.Println("WHATSAPP_DATABASE_HOST: " + db_host)
		log.Println("WHATSAPP_DATABASE_PASSWORD: " + db_password)
		return nil, errors.New("Missing environment variables")
	}

	// Database connection
	db_port := "5432"
	db_user := "neondb_owner"
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, container_err := sqlstore.New("postgresql", "postgresql://"+db_user+":"+db_password+"@"+db_host+":"+db_port, dbLog)

	if container_err != nil {
		log.Println("Error connecting to database")
		log.Println(container_err)
		return nil, container_err
	}

	deviceStore, device_err := container.GetFirstDevice()
	if device_err != nil {
		log.Println("Error getting first device")
		log.Println(device_err)
		return nil, device_err
	}

	// WhatsApp client
	clientLog := waLog.Stdout("Client", "DEBUG", true)
    client := whatsmeow.NewClient(deviceStore, clientLog)


    	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
        err := client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
        err := client.Connect()
		if err != nil {
			panic(err)
		}
	}


    bot := &Bot{client: client}
	// Register the event handler to process group events.
	client.AddEventHandler(bot.eventHandler)

	return bot, nil
}

// Start connects the WhatsApp client and starts the event loop.
func (b Bot)  Start() error {
	// Connect the WhatsApp client.
	if err := b.client.Connect(); err != nil {
		return err
	}

	// Block forever (or run an HTTP server for health checks, etc.).
	select {}
	// return nil  // unreachable, unless you change the design.
}

// eventHandler processes incoming WhatsApp events.
// Note: the exact event types and data depend on the whatsmeow API.
func (b Bot) eventHandler(evt interface{}) {
	// This example shows a simple (and somewhat pseudocode) check for a
	// new group participant. Adapt to the actual event types provided by your library.
	switch event := evt.(type) {
	case *waProto.Message:
		// Here we “simulate” detecting a new participant join.
		// In practice, inspect the message's context or protocol message to identify a group participant event.
		if event.GetExtendedTextMessage() != nil &&
			event.GetExtendedTextMessage().GetText() == "group-participant-added" {
			// Extract participant info (this is illustrative; adjust to your event structure)
			participantJID := "newuser@example.com" // You’d extract this from the event.
			welcomeMsg := GenerateWelcomeMessage(participantJID)
			if err := b.sendMessage(welcomeMsg); err != nil {
				log.Printf("Failed to send welcome message: %v", err)
			}
		}
		// Add additional event processing as needed.
	}
}

// sendMessage sends a text message to the specified target (group or user).
func  (b Bot) sendMessage(message string) error {
	// Create the message payload.
	msg := &waProto.Message{
		Conversation: &message,
	}

    var recipient types.JID
    recipient = types.NewJID(GROUP_ID, types.GroupServer)
	// The SendMessage method and its parameters depend on the whatsmeow API.
	// Adjust the call accordingly.
	_, err := b.client.SendMessage(context.Background(), recipient, msg)
	return err
}
