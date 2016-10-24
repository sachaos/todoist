package lib

type Command struct {
	Args   interface{} `json:"args"`
	TempID string      `json:"temp_id"`
	Type   string      `json:"type"`
	UUID   string      `json:"uuid"`
}
