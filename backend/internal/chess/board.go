package chess

type Color string

const (
	White Color = "white"
	Black Color = "black"
)

type PieceType string

const (
	Pawn   PieceType = "pawn"
	Knight PieceType = "knight"
	Bishop PieceType = "bishop"
	Rook   PieceType = "rook"
	Queen  PieceType = "queen"
	King   PieceType = "king"
)

var PieceTypeToSAN = map[PieceType]string{
	Pawn:   "",
	Knight: "N",
	Bishop: "B",
	Rook:   "R",
	Queen:  "Q",
	King:   "K",
}

type Piece struct {
	Type  PieceType `json:"type"`
	Color Color     `json:"color"`
}

type Board [8][8]*Piece // Files(columns) a-h, ranks(rows) 1-8
