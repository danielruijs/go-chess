function MatchmakingComponent({ queueLength, inQueue }: { queueLength: number | null, inQueue: boolean }) {
    return (
        <div>
            <h2>Matchmaking Queue</h2>
            {queueLength !== null ? (
                <p>Current queue length: {queueLength}</p>
            ) : (
                <p>Loading queue length...</p>
            )}
            {inQueue && (
                <p>You are in the matchmaking queue.</p>
            )}
        </div>
    );
}

export default MatchmakingComponent;