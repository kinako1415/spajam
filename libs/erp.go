package libs

type errorJson struct {
	Message string `json:"message"`
}

func ErrorResponse(err string) errorJson {
	return errorJson{Message: err}
}