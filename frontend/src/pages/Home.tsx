import { useEffect, useRef } from "react";

function Home() {
    const socket = useRef<WebSocket | null>(null);

    useEffect(() => {
        // Initialize WebSocket once
        socket.current = new WebSocket("ws://localhost:8080/ws");

        socket.current.addEventListener("open", () => {
            console.log("WebSocket connected");
            socket.current?.send("Client connection established");
        });

        socket.current.addEventListener("message", (event) => {
            console.log("Message from server:", event.data);
        });

        // Cleanup on unmount
        return () => {
            socket.current?.close();
        };
    }, []); // Empty dependency array ensures this runs only once

    return (
        <div>
            <h1>Test</h1>
            <img src="/pieces/pawn_w.png" alt="White pawn" />
            <button
                onClick={() => {
                    socket.current?.send("Hello from Home page");
                }}
            >
                Send Message
            </button>
        </div>
    );
}

export default Home;
