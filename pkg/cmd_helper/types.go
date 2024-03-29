package cmdhelper

type Commands struct {
	Name         string   `json:"name"`
	Cmd          string   `json:"cmd"`
	Alias        []string `json:"alias"`
	AllowGroup   bool     `json:"allowGroup"`
	AllowPrivate bool     `json:"allowPrivate"`
}