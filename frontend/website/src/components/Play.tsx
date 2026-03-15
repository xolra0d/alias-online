import { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import { Box, CircularProgress, Alert, Typography, Avatar, Tooltip, Paper } from "@mui/material";

interface CreateUserResponse {
    ok: boolean;
    reason?: string;
    credentials?: {
        id: string;
        secret: string;
    };
    name?: string;
}

interface Player {
    id: string;
    ready: boolean;
    score: number;
    name?: string;
}

interface GameConfig {
    language: string;
    "rude-words": boolean;
    "additional-vocabulary": string;
    clock: number;
}

interface GameState {
    admin: string;
    config: GameConfig;
    game_state: number;
    players: Record<string, Player>;
    remaining_time: number;
}

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || "";

// Constants for WebSocket message types
const SERVER_NEW_UPDATE = 1;
const GET_STATE = 0;

export default function Play() {
    const { room_id } = useParams<{ room_id: string }>();
    const [error, setError] = useState<string | null>(null);
    const [connecting, setConnecting] = useState(true);
    const [gameState, setGameState] = useState<GameState | null>(null);
    const wsRef = useRef<WebSocket | null>(null);

    useEffect(() => {
        let cancelled = false;

        const connect = async () => {
            let login = localStorage.getItem("login");
            let secret = localStorage.getItem("secret");

            if (!login || !secret) {
                try {
                    const response = await fetch(`${BACKEND_URL}/api/create-user`, {
                        method: "POST",
                    });
                    const data: CreateUserResponse = await response.json();

                    if (data.ok && data.credentials) {
                        login = data.credentials.id;
                        secret = data.credentials.secret;
                        localStorage.setItem("login", login);
                        localStorage.setItem("secret", secret);
                        if (data.name) {
                            localStorage.setItem("name", data.name);
                        }
                    } else {
                        setError(`Failed to create user: ${data.reason}`);
                        return;
                    }
                } catch {
                    setError("Network error while creating user.");
                    return;
                }
            }

            if (cancelled) return;

            const wsBase = BACKEND_URL.replace(/^http/, "ws");
            const wsUrl = new URL(`/api/ws/${room_id}`, wsBase || window.location.origin);
            wsUrl.searchParams.set("user_id", login!);
            wsUrl.searchParams.set("user_secret", secret!);

            const ws = new WebSocket(wsUrl.toString());
            wsRef.current = ws;

            ws.onopen = () => {
                if (!cancelled) {
                    console.log("[WebSocket] Connected to room:", room_id);
                    setConnecting(false);
                    // Initial state request
                    ws.send(JSON.stringify({
                        user_id: login!,
                        type: GET_STATE,
                        data: {},
                    }));
                }
            };

            ws.onmessage = async (event) => {
                const text = event.data instanceof Blob
                    ? await event.data.text()
                    : event.data;

                try {
                    const json = JSON.parse(text);
                    if (json.msg_type === SERVER_NEW_UPDATE) {
                        console.log("[WebSocket] ServerNewUpdate received:", json.msg_data);
                        setGameState(json.msg_data);
                    }
                } catch (e) {
                    console.error("[WebSocket] Error parsing message:", e, text);
                }
            };

            ws.onerror = (event) => {
                console.error("[WebSocket] Error:", event);
                if (!cancelled) setError("WebSocket connection error.");
            };

            ws.onclose = (event) => {
                console.log("[WebSocket] Closed:", event.code, event.reason);
                if (!cancelled) setConnecting(false);
            };
        };

        connect();

        return () => {
            cancelled = true;
            if (wsRef.current) {
                wsRef.current.close();
            }
        };
    }, [room_id]);

    const renderPlayerCircle = () => {
        if (!gameState) return null;

        const players = Object.values(gameState.players);
        const playerCount = players.length;
        const radius = 180; // Distance from center

        return (
            <Box sx={{
                position: "relative",
                width: 500,
                height: 500,
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
                mx: "auto",
                mt: 4
            }}>
                {/* Center point - could be game status or timer */}
                <Paper elevation={3} sx={{
                    width: 120,
                    height: 120,
                    borderRadius: "50%",
                    display: "flex",
                    flexDirection: "column",
                    justifyContent: "center",
                    alignItems: "center",
                    bgcolor: "primary.main",
                    color: "primary.contrastText",
                    zIndex: 1
                }}>
                    <Typography variant="h6">{gameState.remaining_time}s</Typography>
                    <Typography variant="caption">Time Left</Typography>
                </Paper>

                {players.map((player, index) => {
                    const angle = (index * 2 * Math.PI) / playerCount - Math.PI / 2;
                    const x = radius * Math.cos(angle);
                    const y = radius * Math.sin(angle);

                    const is_admin = player.id === gameState.admin;

                    return (
                        <Box
                            key={player.id}
                            sx={{
                                position: "absolute",
                                transform: `translate(${x}px, ${y}px)`,
                                display: "flex",
                                flexDirection: "column",
                                alignItems: "center",
                                transition: "all 0.5s ease-in-out"
                            }}
                        >
                            <Tooltip title={`${player.ready ? "Ready" : "Not Ready"} | Score: ${player.score}`}>
                                <Avatar
                                    sx={{
                                        width: 60,
                                        height: 60,
                                        border: `4px solid ${player.ready ? "#4caf50" : "#ff9800"}`,
                                        boxShadow: is_admin ? "0 0 15px #ffeb3b" : "none"
                                    }}
                                >
                                    {player.id.slice(0, 2).toUpperCase()}
                                </Avatar>
                            </Tooltip>
                            <Typography variant="body2" sx={{ mt: 1, fontWeight: is_admin ? "bold" : "normal" }}>
                                {player.id.slice(0, 8)}
                                {is_admin && " (Admin)"}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                                Score: {player.score}
                            </Typography>
                        </Box>
                    );
                })}
            </Box>
        );
    };

    return (
        <Box sx={{ maxWidth: 800, mx: "auto", mt: 4, p: 2, textAlign: "center" }}>
            <Typography variant="h4" gutterBottom>
                Room: {room_id}
            </Typography>

            {error && (
                <Alert severity="error" sx={{ mb: 2 }}>
                    {error}
                </Alert>
            )}

            {connecting && !error && (
                <Box display="flex" alignItems="center" justifyContent="center" gap={2}>
                    <CircularProgress size={20} />
                    <Typography>Connecting to room...</Typography>
                </Box>
            )}

            {!connecting && !error && gameState && (
                <Box>
                    <Typography variant="subtitle1" color="text.secondary">
                        Game State: {gameState.game_state === 0 ? "Lobby" : "Playing"} | Language: {gameState.config.language}
                    </Typography>
                    {renderPlayerCircle()}
                </Box>
            )}

            {!connecting && !error && !gameState && (
                <Typography color="success.main">Connected! Waiting for game state...</Typography>
            )}
        </Box>
    );
}