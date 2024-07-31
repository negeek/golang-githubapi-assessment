package utils

type Response struct {
	StatusCode int         `json:"statuscode"`
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
