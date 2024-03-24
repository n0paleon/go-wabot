package turuapi

import (
	"TuruBot/configs"
	"TuruBot/pkg/util"
)

func RandomPornImage() ([]byte, string, error) {
	bytes, mimetype, err := util.GetImageBytes(configs.GetEnv("API_URL") + "/random-porn?type=image")
	if err != nil {
		return nil, "", err
	}

	return bytes, mimetype, nil
}