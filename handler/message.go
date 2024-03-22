package handler

import (
	"TuruBot/pkg/turuapi"
	"TuruBot/pkg/util"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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
					var (
						msgType, messageContent, command string
					)
					var (
						isQuotedImage, isImage bool
					)
					var _ bool = v.Info.IsGroup // original => isGroup
					var chat = v.Info.Chat
					var _ = v.Info.Sender.String() // original => sender
					msgType = v.Info.Type
					extended := v.Message.GetExtendedTextMessage()
					quotedMsg := extended.GetContextInfo().GetQuotedMessage()
					quotedImage := quotedMsg.GetImageMessage()
					msgImage := v.Message.GetImageMessage()

					if quotedImage != nil {
						isQuotedImage = true
					}
					if msgImage != nil {
						isImage = true
					}

					if extendedMessageText := v.Message.GetExtendedTextMessage().GetText(); extendedMessageText != "" {
						messageContent = extendedMessageText
					} else if textConversation := v.Message.GetConversation(); textConversation != "" {
						messageContent = textConversation
					} else if imageConversation := v.Message.ImageMessage.GetCaption(); imageConversation != "" {
						messageContent = imageConversation
					} else {
						messageContent = ""
					}

					// log
					log.Printf("[TYPE] => %s", msgType)
					log.Printf("[CHAT] => %s", messageContent)

					// generate message ID
					var messageID string = v.Info.ID

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
						case "stiker", "sticker", "s":
							fmt.Println("isImage?", isImage)
							fmt.Println("isQuotedImage?", isQuotedImage)
							var ImageData *waProto.ImageMessage
							if isImage {
								ImageData = msgImage
							} else if isQuotedImage {
								ImageData = quotedImage
							} else {
								ImageData = nil
							}
							
							if ImageData != nil {
								data, err := client.Download(ImageData)
								if err != nil {
									client.SendMessage(context.Background(), chat, &waProto.Message{
										Conversation: proto.String("error bang!"),
									})
								}

								convert,  err := util.WebpWriteExifData(data, time.Now().Unix())
								if err != nil {
									log.Println(err)
								}

								upload, err := client.Upload(context.Background(), convert, whatsmeow.MediaImage)
								if err != nil {
									client.SendMessage(context.Background(), chat, &waProto.Message{
										Conversation: proto.String("error bang!"),
									})
								}

								client.SendMessage(context.Background(), chat, &waProto.Message{
									StickerMessage: &waProto.StickerMessage{
										Url:               &upload.URL,
										FileSha256:        upload.FileSHA256,
										FileEncSha256:     upload.FileEncSHA256,
										MediaKey:          upload.MediaKey,
										Mimetype:          proto.String("image/webp"),
										DirectPath:        &upload.DirectPath,
										FileLength:        proto.Uint64(uint64(len(convert))),
										FirstFrameSidecar: convert,
										PngThumbnail:      convert,
									},
								})
							} else {
								client.SendMessage(context.Background(), chat, &waProto.Message{
									Conversation: proto.String("kirim/reply gambar dengan caption !stiker"),
								})
							}
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