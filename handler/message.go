package handler

import (
	"TuruBot/pkg/turuapi"
	"context"
	"fmt"
	"log"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

/**
*	made with hand, and keyboard of course
*	powered by Nopaleon Bonaparte
 */

type MyClient struct {
	WAClient       *whatsmeow.Client
	eventHandlerID uint32
}

func (cli *MyClient) Register() {
	cli.eventHandlerID = cli.WAClient.AddEventHandler(cli.myEventHandler)
}

func (cli *MyClient) myEventHandler(evt interface{}) {
	var client *whatsmeow.Client = cli.WAClient

	switch v := evt.(type) {
		case *events.Message:
			if !v.IsEdit && v.Info.Edit == "" {
				// go func for enhanced performance!
				go func () {
					var msgType string = v.Info.Type
					var _ bool = v.Info.IsGroup // original => isGroup
					var messageContent string
					var command string
					var chat = v.Info.Chat

					if extendedMessageText := v.Message.GetExtendedTextMessage().GetText(); extendedMessageText != "" {
						messageContent = extendedMessageText
					} else {
						messageContent = v.Message.GetConversation()
					}

					// log
					log.Printf("[TYPE] => %s", msgType)
					log.Printf("[CHAT] => %s", messageContent)

					// generate message ID
					var messageID string = client.GenerateMessageID()

					if len(messageContent) > 0 {
						input := strings.Split(messageContent, " ")
						if len(input) > 0 {
							command = input[0][1:]
						} else {
							command = ""
						}
					} else {
						command = ""
					}
	
					switch (command) {
						case "menu":
							client.SendMessage(context.Background(), chat, &waProto.Message{
								Conversation: proto.String("gada menunya anjing!"),
							}, whatsmeow.SendRequestExtra{ID: messageID})
						case "sfw":
							sfwContent, mimetype, err := turuapi.Sfw()
							if err != nil {
								client.SendMessage(context.Background(), chat, &waProto.Message{
									Conversation: proto.String("sorry bos error, coba ulangin lagi deh!"),
								})
							}
	
							upload, err := client.Upload(context.Background(), sfwContent, whatsmeow.MediaImage)
							if err != nil {
								client.SendMessage(context.Background(), chat, &waProto.Message{
									Conversation: proto.String("sorry bos error, coba ulangin lagi deh!"),
								})
							}

							fmt.Println(mimetype)
							var message *waProto.Message
							if mimetype == "image/gif" {
								message = &waProto.Message{
									ImageMessage: &waProto.ImageMessage{
										Url:           &upload.URL,
										Mimetype:      proto.String("image/jpeg"),
										Caption:       proto.String("nih bos!"),
										FileSha256:    upload.FileSHA256,
										FileEncSha256: upload.FileEncSHA256,
										FileLength:    &upload.FileLength,
										MediaKey:      upload.MediaKey,
										DirectPath:    &upload.DirectPath,
										//FileName: 	proto.String("image.gif"),
									},
								}
							} else {
								message = &waProto.Message{
									ImageMessage: &waProto.ImageMessage{
										Url:           &upload.URL,
										Mimetype:      proto.String(mimetype),
										Caption:		proto.String("nih bos!"),
										FileSha256:    upload.FileSHA256,
										FileEncSha256: upload.FileEncSHA256,
										FileLength:    &upload.FileLength,
										MediaKey:      upload.MediaKey,
										DirectPath:    &upload.DirectPath,
									},
								}
							}

							client.SendMessage(context.Background(), chat, message)
						case "nsfw":
							nsfwContent, mimetype, err := turuapi.Nsfw()
							if err != nil {
								client.SendMessage(context.Background(), chat, &waProto.Message{
									Conversation: proto.String("sorry bos error, coba ulangin lagi deh!"),
								})
							}
	
							upload, err := client.Upload(context.Background(), nsfwContent, whatsmeow.MediaImage)
							if err != nil {
								client.SendMessage(context.Background(), chat, &waProto.Message{
									Conversation: proto.String("sorry bos error, coba ulangin lagi deh!"),
								})
							}

							fmt.Println(mimetype)
							var message *waProto.Message
							if mimetype == "image/gif" {
								message = &waProto.Message{
									ImageMessage: &waProto.ImageMessage{
										Url:           &upload.URL,
										Mimetype:      proto.String("image/jpeg"),
										Caption:       proto.String("nih bos!"),
										FileSha256:    upload.FileSHA256,
										FileEncSha256: upload.FileEncSHA256,
										FileLength:    &upload.FileLength,
										MediaKey:      upload.MediaKey,
										DirectPath:    &upload.DirectPath,
										//FileName: 	proto.String("image.gif"),
									},
								}
							} else {
								message = &waProto.Message{
									ImageMessage: &waProto.ImageMessage{
										Url:           &upload.URL,
										Mimetype:      proto.String(mimetype),
										Caption:		proto.String("nih bos!"),
										FileSha256:    upload.FileSHA256,
										FileEncSha256: upload.FileEncSHA256,
										FileLength:    &upload.FileLength,
										MediaKey:      upload.MediaKey,
										DirectPath:    &upload.DirectPath,
									},
								}
							}

							client.SendMessage(context.Background(), chat, message)
						default:
							// matikan dlu auto reply biar ga spam
	
							// client.SendMessage(context.Background(), chat, &waProto.Message{
							// 	Conversation: proto.String("ngomong apaan sih lu ora mudeng gua anjing"),
							// })
					}
				}()
			}
		default:

	}
}