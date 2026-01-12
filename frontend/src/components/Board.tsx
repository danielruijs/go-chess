import type { Position } from "../interfaces/chess";

function Board({ position }: { position: Position | null }) {
    if (!position) {
        return <div>Loading board...</div>;
    }

    return <div>
        {position.board.map((row, i) => (
            <div key={i} style={{ display: "flex" }}>
                {row.map((piece, j) => {
                    const imgPath = `/pieces/${piece.type}_${piece.color}.png`;
                    return (
                        <div
                            key={j}
                            style={{
                                width: 60,
                                height: 60,
                                display: "flex",
                                alignItems: "center",
                                justifyContent: "center",
                                border: "1px solid #000000",
                            }}
                        >
                            {imgPath && <img src={imgPath} alt={`${piece.color} ${piece.type}`} style={{ width: "80%" }} />}
                        </div>
                    );
                })}
            </div>
        ))}
    </div>
}

export default Board;