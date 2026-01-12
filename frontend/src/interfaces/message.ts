import type { Position } from "./chess";

interface WSMessage {
    type: string;
    data: Position;
}

export type { WSMessage };