function coordsToString(file: number, rank: number): string {
    const fileChar = String.fromCharCode("a".charCodeAt(0) + file);
    const rankChar = (rank + 1).toString();
    return fileChar + rankChar;
}

function parsePgn(pgn: string): string[] {
    const moves: string[] = [];
    const tokens = pgn.split(/\s+/); // split on whitespace
    for (const token of tokens) {
        if (token.trim() === "") continue; // skip empty tokens
        if (token.includes(".")) {
            const parts = token.split(".");
            if (parts.length === 2) {
                moves.push(parts[1]);
            }
        }
        else {
            moves.push(token);
        }
    }
    return moves;
}

export { coordsToString, parsePgn };