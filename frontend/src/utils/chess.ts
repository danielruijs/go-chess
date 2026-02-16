function coordsToString(file: number, rank: number): string {
    const fileChar = String.fromCharCode("a".charCodeAt(0) + file);
    const rankChar = (rank + 1).toString();
    return fileChar + rankChar;
}

export { coordsToString };