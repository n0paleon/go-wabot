package handler

import (
	"TuruBot/pkg/util"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type helper struct {
	Client *whatsmeow.Client
	Msg *events.Message
}

func NewHelper(Client *whatsmeow.Client, msg *events.Message) *helper {
	return &helper {
	    Client: Client,
	    Msg: msg,
	}
}

func (bot *helper) GetExpiration () uint32 {
	if bot.Msg.Info.IsGroup {
		return 30 * 86400
	} else {
		return 30 * 86400
	}
}

func (bot *helper) Sender() string {
	return bot.Msg.Info.Sender.ToNonAD().String()
}

func (bot *helper) Reply (text string) {
	bot.Client.SendMessage(context.Background(), bot.Msg.Info.Chat, &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waProto.ContextInfo{
				Expiration: 	proto.Uint32(bot.GetExpiration()),
				StanzaId:		proto.String(bot.Msg.Info.ID),
				Participant:	proto.String(bot.Sender()),
				QuotedMessage: bot.Msg.Message,
			},
		},
	})
}

func (bot *helper) Upload (image []byte, mediaType whatsmeow.MediaType) (resp whatsmeow.UploadResponse, err error) {
	upload, err := bot.Client.Upload(context.Background(), image, mediaType)
	if err != nil {
		return whatsmeow.UploadResponse{}, err
	}

	return upload, nil
}

func (bot *helper) SendImageAsSticker(image []byte) error {
	mimetype := util.GetFileExtFromBytes(image)

	convert, isAnimated, err := util.ImageToWEBP(image, mimetype, time.Now().Unix())
	if err != nil {
		fmt.Println(err)
		return err
	}
	thumbnail, err := util.GenerateStickerThumbnail(image, mimetype, "png")
	if err != nil {
		fmt.Println(err)
		return err
	}

	if (uint64(len(convert)) + uint64(len(thumbnail))) > 1000000 {
		return errors.New("gambar lu kegedean kek ktl gue")
	}

	upload, err := bot.Upload(convert, whatsmeow.MediaImage)
	if err != nil {
		return err
	}

	message := &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:                &upload.URL,
			FileSha256:         upload.FileSHA256,
			FileEncSha256:      upload.FileEncSHA256,
			MediaKey:           upload.MediaKey,
			Mimetype:           proto.String("image/webp"),
			DirectPath:         &upload.DirectPath,
			FileLength:         &upload.FileLength,
			FirstFrameSidecar:  thumbnail,
			PngThumbnail:       thumbnail,
			IsAnimated: 		proto.Bool(isAnimated),
			ContextInfo: 		&waProto.ContextInfo{
				Expiration: 	proto.Uint32(bot.GetExpiration()),
				StanzaId:		&bot.Msg.Info.ID,
				Participant:	proto.String(bot.Sender()),
				QuotedMessage: bot.Msg.Message,
			},
		},
	}

	_, haveErr := bot.Client.SendMessage(context.Background(), bot.Msg.Info.Chat, message)
	if haveErr != nil {
		return err
	}

	return nil
}

func (bot *helper) SendVideoAsSticker(image []byte) error {
	mimetype := util.GetFileExtFromBytes(image)

	convert, err := util.VideoToWEBP(image, time.Now().Unix())
	if err != nil {
		fmt.Println(err)
		return err
	}
	thumbnail, err := util.GenerateStickerThumbnail(image, mimetype, "png")
	if err != nil {
		fmt.Println(err)
		return err
	}

	if (uint64(len(convert)) + uint64(len(thumbnail))) > 1000000 {
		return errors.New("video lu kegedean kek ktl gue")
	}

	upload, err := bot.Upload(convert, whatsmeow.MediaImage)
	if err != nil {
		return err
	}

	message := &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:                &upload.URL,
			FileSha256:         upload.FileSHA256,
			FileEncSha256:      upload.FileEncSHA256,
			MediaKey:           upload.MediaKey,
			Mimetype:           proto.String("image/webp"),
			DirectPath:         &upload.DirectPath,
			FileLength:         &upload.FileLength,
			FirstFrameSidecar:  thumbnail,
			PngThumbnail:       thumbnail,
			IsAnimated: 		proto.Bool(true),
			ContextInfo: 		&waProto.ContextInfo{
				Expiration: 	proto.Uint32(bot.GetExpiration()),
				StanzaId:		&bot.Msg.Info.ID,
				Participant:	proto.String(bot.Sender()),
				QuotedMessage: bot.Msg.Message,
			},
		},
	}

	_, haveErr := bot.Client.SendMessage(context.Background(), bot.Msg.Info.Chat, message)
	if haveErr != nil {
		return err
	}

	return nil
}

func (bot *helper) ReplyWithImage (image []byte, mimetype string, caption string) error {
	fileExt:= util.GetFileExtFromBytes(image)
	thumbnail, err := util.GenerateMediaThumbnail(image, fileExt, mimetype)
	if err != nil {
		fmt.Println(err)
		return err
	}

	upload, errUpload := bot.Upload(image, whatsmeow.MediaImage)
	if errUpload != nil {
		return errUpload
	}

	var message *waProto.Message = &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Url:           &upload.URL,
			Mimetype:      proto.String("image/" + mimetype),
			Caption:       proto.String(caption),
			FileSha256:    upload.FileSHA256,
			FileEncSha256: upload.FileEncSHA256,
			FileLength:    &upload.FileLength,
			MediaKey:      upload.MediaKey,
			DirectPath:    &upload.DirectPath,
			JpegThumbnail: thumbnail,
			ContextInfo: 		&waProto.ContextInfo{
				Expiration: 	proto.Uint32(bot.GetExpiration()),
				StanzaId:		&bot.Msg.Info.ID,
				Participant:	proto.String(bot.Sender()),
				QuotedMessage: bot.Msg.Message,
			},
		},
	}

	bot.Client.SendMessage(context.Background(), bot.Msg.Info.Chat, message)

	return nil
}