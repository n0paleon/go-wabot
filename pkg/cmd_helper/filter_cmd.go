package cmdhelper

import (
	"TuruBot/configs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func GetCmdByAlias(alias string) (*Commands, error) {
	data, err := GetAllCmd(configs.GetEnv("CMD_FILE"))
	if err != nil {
	    	return nil, err
	}
 
	for i := range data {
		for _, a := range data[i].Alias {
			if a == alias {
				return &data[i], nil
			}
		}
	}
 
	return nil, nil
} 

func UpdateGroupPermission(alias string, status bool) (string, error) {
	var cmds []Commands
	var aliasFound bool
	var existingStatus bool
 
	data, err := ioutil.ReadFile(configs.GetEnv("CMD_FILE"))
	if err != nil {
	    	return "error", err
	}
 
	err = json.Unmarshal(data, &cmds)
	if err != nil {
	    	return "error", err
	}
 
	for i := range cmds {
		for _, a := range cmds[i].Alias {
			if a == alias {
				aliasFound = true
				if cmds[i].AllowGroup != status {
					cmds[i].AllowGroup = status
				} else {
					existingStatus = true
				}
				break
			}
		}
		if aliasFound {
			break
		}
	}
 
	if !aliasFound {
	    	return "command yang anda minta tidak ditemukan!", nil
	}
 
	updatedData, err := json.MarshalIndent(cmds, "", "     ")
	if err != nil {
	    	return "error ketika memproses permintaan!", err
	}
 
	err = ioutil.WriteFile(configs.GetEnv("CMD_FILE"), updatedData, os.ModePerm)
	if err != nil {
	    	return "gagal menyimpan data!", err
	}
 
	if existingStatus {
	    	return fmt.Sprintf("err, command ini sudah disetting *%t* untuk semua group sebelumnya!", existingStatus), nil
	}
 
	return fmt.Sprintf("sukses, cmd *%s* sudah berhasil diset ke *%t* untuk semua group!", alias, status), nil
} 

func UpdateDMPermission(alias string, status bool) (string, error) {
	var cmds []Commands
	var aliasFound bool
	var existingStatus bool
 
	data, err := ioutil.ReadFile(configs.GetEnv("CMD_FILE"))
	if err != nil {
	    	return "error", err
	}
 
	err = json.Unmarshal(data, &cmds)
	if err != nil {
	    	return "error", err
	}
 
	for i := range cmds {
		for _, a := range cmds[i].Alias {
			if a == alias {
				aliasFound = true
				if cmds[i].AllowPrivate != status {
					cmds[i].AllowPrivate = status
				} else {
					existingStatus = true
				}
				break
			}
		}
		if aliasFound {
			break
		}
	}
 
	if !aliasFound {
	    	return "command yang anda minta tidak ditemukan!", nil
	}
 
	updatedData, err := json.MarshalIndent(cmds, "", "     ")
	if err != nil {
	    	return "error ketika memproses permintaan!", err
	}
 
	err = ioutil.WriteFile(configs.GetEnv("CMD_FILE"), updatedData, os.ModePerm)
	if err != nil {
	    	return "gagal menyimpan data!", err
	}
 
	if existingStatus {
	    	return fmt.Sprintf("err, command ini sudah disetting *%t* untuk semua private chat sebelumnya!", existingStatus), nil
	}
 
	return fmt.Sprintf("sukses, cmd *%s* sudah berhasil diset ke *%t* untuk semua private chat!", alias, status), nil
}