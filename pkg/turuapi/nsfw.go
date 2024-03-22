package turuapi

import (
	"TuruBot/configs"
	"TuruBot/pkg/util"
)

// return nsfw collection data from API

func Nsfw() ([]byte, string, error) {
	bytes, mimetype, err := util.GetImageBytes(configs.GetEnv("API_URL") + "/nsfw?type=image")
	if err != nil {
		return nil, "", err
	}

	return bytes, mimetype, nil
}