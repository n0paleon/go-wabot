package errorHandling

import "log"

func LogErr(error error) {
	log.Println(error)
}