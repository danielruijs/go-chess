import type { Position } from "../interfaces/chess";

function Board({ position }: { position: Position | null }) {
    if (!position) {
        return <div>Loading board...</div>;
    }

    return <div>
        {position.board.map((row, i) => (
            <div key={i} style={{ display: "flex" }}>
                {row.map(piece => {
                    const imgPath = `/pieces/${piece.color}-${piece.type}.png`;
                    return (
                        <div
                            style={{
                                width: 80,
                                height: 80,
                                display: "flex",
                                alignItems: "center",
                                justifyContent: "center",
                                border: "1px solid #000000",
                            }}
                        >
                            {imgPath && <img src={imgPath} alt={`${piece.color} ${piece.type}`} style={{ width: "80px", height: "80px" }} />}
                        </div>
                    );
                })}
            </div>
        ))}
    </div>
}

export default Board;