function MatchmakingComponent({ queueLength, inQueue }: { queueLength: number | null, inQueue: boolean }) {
    return (
        <div className="p-2">
            <div className="text-xl font-semibold mb-2">Matchmaking Queue</div>
            {queueLength !== null ? (
                <p>Current queue length: {queueLength}</p>
            ) : (
                <p>Loading queue length...</p>
            )}
            {inQueue && (
                <p className="mt-2 text-gray-600">You are in the matchmaking queue.</p>
            )}
        </div>
    );
}

export default MatchmakingComponent;