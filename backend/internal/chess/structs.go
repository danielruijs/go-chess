package chess

type Color string

const (
	White Color = "white"
	Black Color = "black"
)

type Piece struct {
	Type  string `json:"type"`
	Color Color `json:"color"`
}

type Position struct {
	Board [8][8]Piece `json:"board"`
}