package handler

import (
	"TuruBot/pkg/turuapi"
	"TuruBot/pkg/util"
	"context"
	"log"
	"strings"
	"time"
	"unicode"

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
			if !v.IsEdit && v.Info.Edit == "" && !v.Info.IsFromMe {
				// go func for enhanced performance!
				go func () {
					var (
						msgType, messageContent, command string
					)
					var (
						isQuotedImage, isImage, isQuotedVideo, isVideo bool
					)
					var _ bool = v.Info.IsGroup // original => isGroup
					var chat = v.Info.Chat
					var sender = v.Info.Sender.String() // original => sender
					msgType = v.Info.Type
					extended := v.Message.GetExtendedTextMessage()

					var quotedMsg *waProto.Message = extended.GetContextInfo().GetQuotedMessage().GetEphemeralMessage().GetMessage()
					if quotedMsg == nil {
						quotedMsg = extended.GetContextInfo().GetQuotedMessage()
					}

					quotedImage := quotedMsg.GetImageMessage()
					msgImage := v.Message.GetImageMessage()
					quotedVideo := quotedMsg.GetVideoMessage()
					msgVideo := v.Message.GetVideoMessage()

					
					if quotedImage != nil {
						isQuotedImage = true
					}
					if msgImage != nil {
						isImage = true
					}
					if quotedVideo != nil {
						isQuotedVideo = true
					}
					if msgVideo != nil {
						isVideo = true
					}

					var ImageData *waProto.ImageMessage
					if isImage {
						ImageData = msgImage
					} else if isQuotedImage {
						ImageData = quotedImage
					} else {
						ImageData = nil
					}
					
					var VideoData *waProto.VideoMessage
					if isVideo {
						VideoData = msgVideo
					} else if isQuotedVideo {
						VideoData = quotedVideo
					} else {
						VideoData = nil
					}

					if extendedMessageText := v.Message.GetExtendedTextMessage().GetText(); extendedMessageText != "" {
						messageContent = extendedMessageText
					} else if textConversation := v.Message.GetConversation(); textConversation != "" {
						messageContent = textConversation
					} else if imageConversation := v.Message.ImageMessage.GetCaption(); imageConversation != "" {
						messageContent = imageConversation
					} else if videoConversation := v.Message.VideoMessage.GetCaption(); videoConversation != "" {
						messageContent = videoConversation
					} else {
						messageContent = "nil"
					}

					// log
					log.Printf("[CHAT]	=> %s", sender)
					log.Printf("[TYPE]	=> %s", msgType)
					log.Printf("[MSG]	=> %s", messageContent)

					// return if not special prefix
					if !isSpecialSymbol(messageContent) {
						return
					}

					// generate message ID
					var messageID string = v.Info.ID

					// argument parameter parsing
					var args = strings.Split(messageContent, " ")[1:]

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
						case "stikergif", "stickergif", "sgif":
							log.Println("isVideo?", isVideo)
							log.Println("isQuotedVideo?", isQuotedVideo)
							
							if VideoData != nil {
								var maxDuration int = 15
								if isVideo && int(*VideoData.Seconds) > maxDuration {
									client.SendMessage(context.Background(), chat, &waProto.Message{
										Conversation: proto.String("maksimal 15 detik anjing"),
									})
									return
								}

								data, err := client.Download(VideoData)
								if err != nil {
									client.SendMessage(context.Background(), chat, &waProto.Message{
										Conversation: proto.String("error bang!"),
									})
									return
								}

								convert, thumbnail,  err := util.VideoToWEBP(data, time.Now().Unix(), VideoData.FileLength, uint64(*VideoData.Seconds))
								if err != nil {
									client.SendMessage(context.Background(), chat, &waProto.Message{
										Conversation: proto.String("error bang!"),
									})
									return
								}

								if (uint64(len(convert)) + uint64(len(thumbnail))) > 1000000 {
									client.SendMessage(context.Background(), chat, &waProto.Message{
										Conversation: proto.String("error bang, gambar lu kegedean kek ktl gue"),
									})
									return
								}

								upload, err := client.Upload(context.Background(), convert, whatsmeow.MediaImage)
								if err != nil {
									client.SendMessage(context.Background(), chat, &waProto.Message{
										Conversation: proto.String("error bang!"),
									})
									return
								}
								
								client.SendMessage(context.Background(), chat, &waProto.Message{
									StickerMessage: &waProto.StickerMessage{
										Url:                &upload.URL,
										FileSha256:         upload.FileSHA256,
										FileEncSha256:      upload.FileEncSHA256,
										MediaKey:           upload.MediaKey,
										Height: 			proto.Uint32(512),
										Width: 			proto.Uint32(512),
										Mimetype:           proto.String("image/webp"),
										DirectPath:         &upload.DirectPath,
										FileLength:         proto.Uint64(uint64(len(convert))),
										FirstFrameSidecar:  convert,
										PngThumbnail: 		thumbnail,
										IsAnimated: 		proto.Bool(true),
									},
								})
							} else {
								client.SendMessage(context.Background(), chat, &waProto.Message{
									Conversation: proto.String("kirim/reply video dengan caption !stikergif"),
								})
							}
						case "stiker", "sticker", "s":
							log.Println("isImage?", isImage)
							log.Println("isQuotedImage?", isQuotedImage)
							
							if ImageData != nil {
								data, err := client.Download(ImageData)
								if err != nil {
									client.SendMessage(context.Background(), chat, &waProto.Message{
										Conversation: proto.String("error bang!"),
									})
								}

								convert, thumbnail, isAnimated, err := util.ImageToWEBP(data, time.Now().Unix())
								if err != nil {
									client.SendMessage(context.Background(), chat, &waProto.Message{
										Conversation: proto.String("error bang!"),
									})
									return
								}

								if (uint64(len(convert)) + uint64(len(thumbnail))) > 1000000 {
									client.SendMessage(context.Background(), chat, &waProto.Message{
										Conversation: proto.String("error bang, gambar lu kegedean kek ktl gue"),
									})
									return
								}

								upload, err := client.Upload(context.Background(), convert, whatsmeow.MediaImage)
								if err != nil {
									client.SendMessage(context.Background(), chat, &waProto.Message{
										Conversation: proto.String("error bang!"),
									})
								}

								message := &waProto.Message{
									StickerMessage: &waProto.StickerMessage{
										Url:                &upload.URL,
										FileSha256:         upload.FileSHA256,
										FileEncSha256:      upload.FileEncSHA256,
										MediaKey:           upload.MediaKey,
										Mimetype:           proto.String("image/webp"),
										DirectPath:         &upload.DirectPath,
										FileLength:         proto.Uint64(uint64(len(convert))),
										FirstFrameSidecar:  convert,
										PngThumbnail:       thumbnail,
										IsAnimated: 		proto.Bool(isAnimated),
									},
								}
								
								client.SendMessage(context.Background(), chat, message)
							} else {
								client.SendMessage(context.Background(), chat, &waProto.Message{
									Conversation: proto.String("kirim/reply gambar dengan caption !stiker"),
								})
							}
						case "menu":
							client.SendMessage(context.Background(), chat, &waProto.Message{
								Conversation: proto.String(util.Dedent(
								`
								*Menu Bot*
								
								/menu
								/stiker [/s, /sticker]
								/stikergif [/sgif, /stickergif]
								/nsfw
								/sfw
								/randombokep [/rbokep, /rporn]
								/ssweb <link web> (tipe)
								
								Powered by n0paleon (https://www.instagram.com/nopaleon.real)`)),
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
							
							var message *waProto.Message = &waProto.Message{
								ImageMessage: &waProto.ImageMessage{
									Url:           &upload.URL,
									Mimetype:      proto.String("image/jpeg"),
									Caption:       proto.String("nih bos!"),
									FileSha256:    upload.FileSHA256,
									FileEncSha256: upload.FileEncSHA256,
									FileLength:    &upload.FileLength,
									MediaKey:      upload.MediaKey,
									DirectPath:    &upload.DirectPath,
								},
							}
							if mimetype == "image/gif" {
								message.ImageMessage.Mimetype = proto.String("image/gif")
							}

							client.SendMessage(context.Background(), chat, message)
						case "ssweb", "screenshotweb":
							if len(args) > 0 {
								tipe := "desktop"
								if len(args) > 1 {
									switch (args[1]) {
										case "mobile", "hp":
											tipe = "mobile"
										case "desktop", "pc":
											tipe = "desktop"
										default:
											client.SendMessage(context.Background(), chat, &waProto.Message{
												Conversation: proto.String("tipe yang tersedia:\n\n1. mobile/hp\n2. desktop/pc"),
											})
											return
									}
								}
								sswebBytes, mimetype, err := turuapi.SsWeb(args[0], tipe)
								if err != nil {
									client.SendMessage(context.Background(), chat, &waProto.Message{
										Conversation: proto.String("ngasih link yg bner anjg\ncontoh: https://detik.com"),
									})
									return
								}

								client.SendMessage(context.Background(), chat, &waProto.Message{
									Conversation: proto.String("nungguin yaaa..."),
								})

								upload, err := client.Upload(context.Background(), sswebBytes, whatsmeow.MediaImage)
								if err != nil {
									client.SendMessage(context.Background(), chat, &waProto.Message{
										Conversation: proto.String("sorry bos error, coba ulangin lagi deh!"),
									})
								}

								var message *waProto.Message = &waProto.Message{
									ImageMessage: &waProto.ImageMessage{
										Url:           &upload.URL,
										Mimetype:      proto.String(mimetype),
										Caption:       proto.String("nih bos!"),
										FileSha256:    upload.FileSHA256,
										FileEncSha256: upload.FileEncSHA256,
										FileLength:    &upload.FileLength,
										MediaKey:      upload.MediaKey,
										DirectPath:    &upload.DirectPath,
									},
								}

								client.SendMessage(context.Background(), chat, message)
							} else {
								client.SendMessage(context.Background(), chat, &waProto.Message{
									Conversation: proto.String("ketik !ssweb <spasi> <link website>.\ncontoh !ssweb https://google.com"),
								})
								return
							}
						case "rporn", "randomporn", "rbokep", "randombokep":
							PornContent, mimetype, err := turuapi.RandomPornImage()
							if err != nil {
								client.SendMessage(context.Background(), chat, &waProto.Message{
									Conversation: proto.String("sorry bos error, coba ulangin lagi deh!"),
								})
							}
	
							upload, err := client.Upload(context.Background(), PornContent, whatsmeow.MediaImage)
							if err != nil {
								client.SendMessage(context.Background(), chat, &waProto.Message{
									Conversation: proto.String("sorry bos error, coba ulangin lagi deh!"),
								})
							}

							var message *waProto.Message = &waProto.Message{
								ImageMessage: &waProto.ImageMessage{
									Url:           &upload.URL,
									Mimetype:      proto.String("image/png"),
									Caption:       proto.String("lu orangnya sange banget bang"),
									FileSha256:    upload.FileSHA256,
									FileEncSha256: upload.FileEncSHA256,
									FileLength:    &upload.FileLength,
									MediaKey:      upload.MediaKey,
									DirectPath:    &upload.DirectPath,
								},
							}
							if mimetype == "image/gif" {
								message.ImageMessage.Mimetype = proto.String("image/gif")
							}

							client.SendMessage(context.Background(), chat, message)
						case "nsfw":
							NsfwContent, mimetype, err := turuapi.Nsfw()
							if err != nil {
								client.SendMessage(context.Background(), chat, &waProto.Message{
									Conversation: proto.String("sorry bos error, coba ulangin lagi deh!"),
								})
							}
	
							upload, err := client.Upload(context.Background(), NsfwContent, whatsmeow.MediaImage)
							if err != nil {
								client.SendMessage(context.Background(), chat, &waProto.Message{
									Conversation: proto.String("sorry bos error, coba ulangin lagi deh!"),
								})
							}

							var message *waProto.Message = &waProto.Message{
								ImageMessage: &waProto.ImageMessage{
									Url:           &upload.URL,
									Mimetype:      proto.String("image/jpeg"),
									Caption:       proto.String("nih bos!"),
									FileSha256:    upload.FileSHA256,
									FileEncSha256: upload.FileEncSHA256,
									FileLength:    &upload.FileLength,
									MediaKey:      upload.MediaKey,
									DirectPath:    &upload.DirectPath,
								},
							}
							if mimetype == "image/gif" {
								message.ImageMessage.Mimetype = proto.String("image/gif")
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

func isSpecialSymbol(s string) bool {
	r := []rune(s)[0]
	return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}