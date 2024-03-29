package handler

import (
	"TuruBot/configs"
	msgtemplate "TuruBot/handler/msg_template"
	cmdhelper "TuruBot/pkg/cmd_helper"
	"TuruBot/pkg/turuapi"
	"fmt"
	"log"
	"strings"
	"unicode"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

/**
*	made with hand, and keyboard of course
*	powered by Nopaleon Bonaparte
 */



func EventHandler(sock *whatsmeow.Client, msg *events.Message) {
	bot := NewHelper(sock, msg)

	v := msg

	if !v.IsEdit && v.Info.Edit == "" && !v.Info.IsFromMe {
		// go func for enhanced performance!
		func () {
			var (
				msgType, messageContent, command string
			)
			var (
				isQuotedImage, isImage, isQuotedVideo, isVideo bool
			)
			var _ bool = msg.Info.IsGroup // original => isGroup
			var _ = msg.Info.Chat
			var sender = msg.Info.Sender.ToNonAD() // original => sender
			msgType = msg.Info.Type
			extended := msg.Message.GetExtendedTextMessage()
			isOwner := sender.User == configs.GetEnv("OWNER_NUMBER")

			var quotedMsg *waProto.Message = extended.GetContextInfo().GetQuotedMessage().GetEphemeralMessage().GetMessage()
			if quotedMsg == nil {
				quotedMsg = extended.GetContextInfo().GetQuotedMessage()
			}

			quotedImage := quotedMsg.GetImageMessage()
			msgImage := msg.Message.GetImageMessage()
			quotedVideo := quotedMsg.GetVideoMessage()
			msgVideo := msg.Message.GetVideoMessage()

			
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

			if extendedMessageText := msg.Message.GetExtendedTextMessage().GetText(); extendedMessageText != "" {
				messageContent = extendedMessageText
			} else if textConversation := msg.Message.GetConversation(); textConversation != "" {
				messageContent = textConversation
			} else if imageConversation := msg.Message.ImageMessage.GetCaption(); imageConversation != "" {
				messageContent = imageConversation
			} else if videoConversation := msg.Message.VideoMessage.GetCaption(); videoConversation != "" {
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

			// argument parameter parsing
			var args = strings.Split(strings.TrimSpace(messageContent), " ")[1:]

			if len(messageContent) > 0 {
				input := strings.Split(strings.TrimSpace(messageContent), " ")
				if len(input) > 0 {
					command = input[0][1:]
				} else {
					command = ""
				}
			} else {
				command = ""
			}

			reason, isAllowed := ValidateCmd(sock, msg, command)
			if !isAllowed {
				if reason == "" {
					return
				}
				bot.Reply(reason)
				return
			}

			switch (command) {
				case "setgroupcmd":
					if isOwner {
						if len(args) > 1 {
							switch (strings.ToLower(args[1])) {
								case "on", "enable", "enabled":
									updateCmd, _ := cmdhelper.UpdateGroupPermission(args[0], true)
									bot.Reply(updateCmd)
								case "off", "disable", "disabled":
									updateCmd, _ := cmdhelper.UpdateGroupPermission(args[0], false)
									bot.Reply(updateCmd)
								default:
									bot.Reply("opsi yang tersedia:\n\n1. ON/ENABLE\n\n2. OFF/DISABLE")
							}
						} else {
							bot.Reply("contoh: /setgroupcmd menu on")
						}
					} else {
						bot.Reply("perintah ini hanya bisa diakses oleh owner!")
					}
				case "setdmcmd", "setprivatecmd":
					if isOwner {
						if len(args) > 1 {
							switch (strings.ToLower(args[1])) {
								case "on", "enable", "enabled":
									updateCmd, err := cmdhelper.UpdateDMPermission(args[0], true)
									if err != nil {
										fmt.Println(err)
									}
									bot.Reply(updateCmd)
								case "off", "disable", "disabled":
									updateCmd, err := cmdhelper.UpdateDMPermission(args[0], false)
									if err != nil {
										fmt.Println(err)
									}
									bot.Reply(updateCmd)
								default:
									bot.Reply("opsi yang tersedia:\n\n1. ON/ENABLE\n\n2. OFF/DISABLE")
							}
						} else {
							bot.Reply("contoh: /setprivatecmd menu on")
						}
					} else {
						bot.Reply("perintah ini hanya bisa diakses oleh owner!")
					}
				case "ssweb":
					if len(args) > 0 {
						tipe := "desktop"
						if len(args) > 1 {
							switch (args[1]) {
								case "mobile", "hp":
									tipe = "mobile"
								case "desktop", "pc":
									tipe = "desktop"
								default:
									bot.Reply("tipe yang tersedia:\n\n1. mobile/hp\n2. desktop/pc")
									return
							}
						}
						sswebBytes, _, err := turuapi.SsWeb(args[0], tipe)
						if err != nil {
							bot.Reply("ngasih link yg bner anjg\ncontoh: https://detik.com")
							return
						}

						bot.Reply("nungguin yaaa....")

						if err := bot.ReplyWithImage(sswebBytes, "jpeg", "nih boskuh!"); err != nil {
							bot.Reply("error nih banh!")
							return
						}
					} else {
						bot.Reply("ketik !ssweb <spasi> <link website>.\ncontoh !ssweb https://google.com")
						return
					}
				case "sfw":
					sfwContent, _, err := turuapi.Sfw()
					if err != nil {
						bot.Reply("sori boskuh lagi error nih!")
						return
					}

					if err := bot.ReplyWithImage(sfwContent, "jpeg", "nih bos!"); err != nil {
						bot.Reply("sori boskuh lagi error nih!")
						return
					}
				case "nsfw":
					NsfwContent, _, err := turuapi.Nsfw()
					if err != nil {
						bot.Reply("sori boskuh lagi error nih!")
						return
					}

					if err := bot.ReplyWithImage(NsfwContent, "jpeg", "nih bos!"); err != nil {
						bot.Reply("sori boskuh lagi error nih!")
						return
					}
				case "rporn", "randomporn", "rbokep", "randombokep":
					PornContent, _, err := turuapi.RandomPornImage()
					if err != nil {
						bot.Reply("sori boskuh lagi error nih!")
						return
					}

					if err := bot.ReplyWithImage(PornContent, "jpeg", "lu orangnya sange banget bang"); err != nil {
						bot.Reply("sori boskuh lagi error nih!")
						return
					}
				case "menu":
					fmt.Println(msgtemplate.Menu())
					bot.Reply(msgtemplate.Menu())
				case "stiker", "sticker", "s":
					if ImageData == nil {
						bot.Reply("gambarnya mana kocak, kirim/reply gambar lu pake caption /stiker")
						return
					}

					image, err := sock.Download(ImageData)
					if err != nil {
						bot.Reply("error bang!")
						return
					}

					errr := bot.SendImageAsSticker(image)
					if errr != nil {
						fmt.Println(errr)
					}
				case "stikergif", "stickergif", "sgif":
					if VideoData == nil {
						bot.Reply("videonya mana kocak, kirim/reply video lu pake caption /sgif")
					}
					if isVideo && int(*VideoData.Seconds) > 15 {
						bot.Reply("maksimal 15 detik anjinh")
						return
					}

					video, err := sock.Download(VideoData)
					if err != nil {
						bot.Reply("error bang")
						return
					}

					errr := bot.SendVideoAsSticker(video)
					if errr != nil {
						bot.Reply("error ktl")
						return
					}
				default:
					return
			}
		}()
	}
}

func isSpecialSymbol(s string) bool {
	r := []rune(s)[0]
	return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}