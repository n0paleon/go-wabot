package turuapi

import (
	"TuruBot/configs"
	"TuruBot/pkg/util"
)

func SsWeb(url string, ua string) ([]byte, string, error) {
	bytes, mimetype, err := util.GetImageBytes(configs.GetEnv("API_URL") + "/screenshot-web?type=image&ua=" + ua + "&url=" + url)
	if err != nil {
		return nil, "", err
	}

	return bytes, mimetype, nil
}