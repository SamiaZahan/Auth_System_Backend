package response

var (
	BodyParseFailedErrorMsg = "failed to parse the body"
	ValidationFailedMsg     = "validation failed"
)

type Payload struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  error       `json:"errors"`
}
