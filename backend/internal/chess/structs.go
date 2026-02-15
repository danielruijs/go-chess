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

type Piece struct {
	Type  PieceType `json:"type"`
	Color Color     `json:"color"`
}

type CastlingRights struct {
	WhiteOO  bool `json:"white_oo"`
	WhiteOOO bool `json:"white_ooo"`
	BlackOO  bool `json:"black_oo"`
	BlackOOO bool `json:"black_ooo"`
}

type Board [BoardSize][BoardSize]Piece

type Position struct {
	Board          Board          `json:"board"`           // Ranks(rows) 8-1, files(columns) a-h
	ActiveColor    Color          `json:"active_color"`    // Color to move
	CastlingRights CastlingRights `json:"castling_rights"` // Castling rights
	EnPassant      Square         `json:"en_passant"`      // En passant target square, square over which pawn just moved when moving two squares
	Halfmove       uint           `json:"halfmove"`        // Halfmove clock, number of halfmoves since last capture or pawn move, for fifty-move rule
	Fullmove       uint           `json:"fullmove"`        // Fullmove number
}

type Move struct {
	From Square `json:"from"`
	To   Square `json:"to"`
}
