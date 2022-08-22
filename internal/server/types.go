package server

type LogEntryStreamRequest struct {
	PositionFrom int `json:"positionFrom,omitempty"`
}

type LogEntryStreamResponse struct {
	Position int    `json:"position"`
	Payload  string `json:"payload"`
}
