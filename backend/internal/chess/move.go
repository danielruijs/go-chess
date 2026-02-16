package chess

type Move struct {
	From Square `json:"from"`
	To   Square `json:"to"`
}
