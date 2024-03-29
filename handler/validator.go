package handler

import (
	cmdhelper "TuruBot/pkg/cmd_helper"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func ValidateCmd(sock *whatsmeow.Client, msg *events.Message, cmd string)  (string, bool) {
	getCmd, err := cmdhelper.GetCmdByAlias(cmd)
	if err != nil {
		return "error ketika memvalidasi command!", false
	}

	if getCmd == nil {
		return "", true
	}

	if !getCmd.AllowGroup && !getCmd.AllowPrivate {
		return fmt.Sprintf("command *%s* sedang dalam perbaikan!", cmd), false
	}

	if msg.Info.IsGroup {
		if getCmd.AllowGroup {
			return "", true
		} else {
			return fmt.Sprintf("command *%s* hanya bisa digunakan di japri/private chat!", cmd), false
		}
	} else {
		if getCmd.AllowPrivate {
			return "", true
		} else {
			return fmt.Sprintf("comand *%s* hanya bisa digunakan di grup chat!", cmd), false
		}
	}
}