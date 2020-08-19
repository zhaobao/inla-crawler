package model

type Tag struct {
	Id      string            `json:"id"`
	Name    map[string]string `json:"name"`
	SortKey int               `json:"sort_key"`
	Key     string            `json:"key"`
	Type    string            `json:"type"`
	Status  string            `json:"status"`
}
