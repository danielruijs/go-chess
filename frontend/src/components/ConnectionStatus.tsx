import { useWebSocket } from "../contexts/WebSocketContext";

function ConnectionStatus() {
  const { isConnected } = useWebSocket();

  return (
    <div
      className="flex items-center gap-2 px-3 py-2 rounded-full bg-slate-100 border border-gray-200 select-none"
      title={isConnected ? "Connected to server" : "Disconnected from server"}
    >
      {/* Status indicator dot */}
      <span
        className={`w-2 h-2 rounded-full ${
          isConnected ? "bg-emerald-500" : "bg-rose-500 animate-pulse"
        }`}
      />
      {/* Status text */}
      <span className="text-xs font-medium text-slate-600">
        {isConnected ? "Connected" : "Disconnected"}
      </span>
    </div>
  );
}

export default ConnectionStatus;
