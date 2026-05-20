import { useWebSocket } from "../contexts/WebSocketContext";

function ConnectionStatus() {
  const { isConnected } = useWebSocket();

  return (
    <div className="flex items-center gap-2 px-4 py-2 rounded-lg bg-white shadow-md border border-gray-200">
      {/* Status indicator dot */}
      <div className={`w-3 h-3 rounded-full ${isConnected ? "bg-green-500" : "bg-red-500"}`} />
      {/* Status text */}
      <span className="text-sm font-medium text-gray-700">
        {isConnected ? "Connected" : "Disconnected"}
      </span>
    </div>
  );
}

export default ConnectionStatus;
