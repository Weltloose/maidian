package model

type EmitPvData struct {
	Url     string `json:"url"`
	Referer string `json:"referer"`
	Title   string `json:"title"`
	Time    string `json:"time"`
}

type EmitPVStruct struct {
	Data    EmitPvData `json:"data"`
	Uv      string     `json:"uv"`
	AppId   string     `json:"app_id"`
	AppName string     `json:"app_name"`
}

type Point struct {
	X int64   `json:"x"`
	Y float64 `json:"y"`
}
