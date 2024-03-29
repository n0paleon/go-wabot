package msgtemplate

import (
	"TuruBot/configs"
	cmdhelper "TuruBot/pkg/cmd_helper"
	"fmt"
	"strings"
)

func Menu() (string) {
	commands, _ := cmdhelper.GetAllCmd(configs.GetEnv("CMD_FILE"))

	var output strings.Builder
	for i, cmd := range commands {
		output.WriteString(fmt.Sprintf("%d. *%s*\n", i+1, cmd.Name))
		output.WriteString(fmt.Sprintf("   Command = %s\n", strings.Join(cmd.Alias, ", ")))
		output.WriteString(fmt.Sprintf("   Group Allowed? %t\n", cmd.AllowGroup))
		output.WriteString(fmt.Sprintf("   Private Allowed? %t\n\n", cmd.AllowPrivate))
	}

	return fmt.Sprintf("*TuruBot v2024 Menu 🚀*\n\nBot Prefix: [All Unique Symbols]\n\n" + output.String() + "_Powered by n0paleon_ [https://instagram.com/nopaleon.real]")
}