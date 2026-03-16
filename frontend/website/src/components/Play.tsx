import { useEffect, useRef, useState, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
    Box,
    CircularProgress,
    Alert,
    Typography,
    Avatar,
    Tooltip,
    Paper,
    Button,
    TextField,
    List,
    ListItem,
    ListItemText,
    Divider,
    IconButton,
    InputAdornment,
    Snackbar,
    useTheme
} from "@mui/material";
import SendIcon from "@mui/icons-material/Send";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";

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
    words_tried: number;
    words_guessed: number;
    name?: string;
}

interface GameConfig {
    language: string;
    "rude-words": boolean;
    "additional-vocabulary": string[];
    clock: number;
}

const GameStatus = {
    RoundOver: 0,
    Explaining: 1,
    Finished: 2
} as const;

type GameStatus = typeof GameStatus[keyof typeof GameStatus];

interface GameState {
    admin: string;
    config: GameConfig;
    game_state: GameStatus;
    players: Record<string, Player>;
    remaining_time: number;
    current_player: string;
}

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || "";

// Client Message Types
const CLIENT_GET_STATE = 0;
const CLIENT_START_ROUND = 1;
const CLIENT_GET_WORD = 2;
const CLIENT_TRY_GUESS = 3;
const CLIENT_FINISH_GAME = 4;

// Server Message Types
const SERVER_NEW_UPDATE = 0;
const SERVER_CURRENT_STATE = 1;
const SERVER_YOUR_WORD = 2;
const SERVER_WORD_GUESSED = 3;
const SERVER_RIGHT_GUESS = 4;
const SERVER_WRONG_GUESS = 5;

export default function Play() {
    const { room_id } = useParams<{ room_id: string }>();
    const navigate = useNavigate();
    const theme = useTheme();
    const [error, setError] = useState<string | null>(null);
    const [connecting, setConnecting] = useState(true);
    const [gameState, setGameState] = useState<GameState | null>(null);
    const gameStateRef = useRef<GameState | null>(null);
    const [timeLeft, setTimeLeft] = useState<number>(0);
    const [currentWord, setCurrentWord] = useState<string | null>(null);
    const [guess, setGuess] = useState("");
    const [lastGuessResult, setLastGuessResult] = useState<{ msg: string, type: 'success' | 'error' | 'info' } | null>(null);
    const wsRef = useRef<WebSocket | null>(null);
    const [userId, setUserId] = useState<string | null>(localStorage.getItem("login"));

    const sendMessage = useCallback((type: number, data: Record<string, unknown> = {}) => {
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({
                user_id: userId,
                type,
                data,
            }));
        }
    }, [userId]);

    const fetchState = useCallback(() => {
        sendMessage(CLIENT_GET_STATE);
    }, [sendMessage]);

    // Local countdown timer
    useEffect(() => {
        if (gameState) {
            setTimeLeft(gameState.remaining_time);
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [gameState?.remaining_time, gameState?.game_state]);

    useEffect(() => {
        let timer: number | null = null;
        const isExplaining = gameState?.game_state === GameStatus.Explaining;
        if (isExplaining) {
            timer = window.setInterval(() => {
                setTimeLeft((prev) => (prev > 0 ? prev - 1 : 0));
            }, 1000);
        }
        return () => {
            if (timer) clearInterval(timer);
        };
    }, [gameState?.game_state]);

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
                        setUserId(login);
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
            } else {
                setUserId(login);
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
                    fetchState();
                }
            };

            ws.onmessage = async (event) => {
                const text = event.data instanceof Blob
                    ? await event.data.text()
                    : event.data;

                try {
                    const json = JSON.parse(text);
                    console.log("[WebSocket] Received:", json);

                    switch (json.msg_type) {
                        case SERVER_NEW_UPDATE:
                            fetchState();
                            break;
                        case SERVER_CURRENT_STATE: {
                            const newState = json.msg_data as GameState;
                            setGameState(newState);
                            gameStateRef.current = newState;
                            break;
                        }
                        case SERVER_YOUR_WORD:
                            setCurrentWord(json.msg_data.word);
                            break;
                        case SERVER_WORD_GUESSED:
                            setLastGuessResult({ msg: `Word guessed by ${json.msg_data.guesser.slice(0, 8)}!`, type: 'info' });
                            if (gameStateRef.current?.current_player === userId) {
                                setCurrentWord(null);
                            }
                            break;
                        case SERVER_RIGHT_GUESS:
                            setLastGuessResult({ msg: "Correct guess!", type: 'success' });
                            break;
                        case SERVER_WRONG_GUESS:
                            setLastGuessResult({ msg: "Wrong guess, try again!", type: 'error' });
                            break;
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
    }, [room_id, fetchState, userId]);

    // Handle "Get Word" automatically if it's our turn to explain
    useEffect(() => {
        if (gameState && gameState.game_state === GameStatus.Explaining && gameState.current_player === userId) {
            if (!currentWord) {
                sendMessage(CLIENT_GET_WORD);
            }
        } else {
            setCurrentWord(null);
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [gameState, currentWord, userId]);

    const handleStartRound = () => {
        sendMessage(CLIENT_START_ROUND);
    };

    const handleFinishGame = () => {
        sendMessage(CLIENT_FINISH_GAME);
    };

    const handleGuess = (e?: React.FormEvent) => {
        if (e) e.preventDefault();
        if (!guess.trim()) return;
        sendMessage(CLIENT_TRY_GUESS, { guess: guess.trim() });
        setGuess("");
    };

    const handleLeaveRoom = () => {
        if (wsRef.current) {
            wsRef.current.close();
        }
        navigate("/");
    };

    const handleShare = () => {
        const url = window.location.href;
        navigator.clipboard.writeText(url).then(() => {
            setLastGuessResult({ msg: "Room link copied to clipboard!", type: 'success' });
        });
    };

    const renderPlayerCircle = () => {
        if (!gameState) return null;

        const players = Object.values(gameState.players).filter(p => p.ready);
        const playerCount = players.length;
        const radius = 180;

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
                <Paper elevation={3} sx={{
                    width: 140,
                    height: 140,
                    borderRadius: "50%",
                    display: "flex",
                    flexDirection: "column",
                    justifyContent: "center",
                    alignItems: "center",
                    bgcolor: gameState.game_state === GameStatus.Explaining ? "secondary.main" : 
                             gameState.game_state === GameStatus.Finished ? "success.main" : "primary.main",
                    color: "primary.contrastText",
                    zIndex: 1,
                    textAlign: "center",
                    p: 2
                }}>
                    <Typography variant="h4">{gameState.game_state === GameStatus.Finished ? "GG" : `${timeLeft}s`}</Typography>
                    <Typography variant="caption">
                        {gameState.game_state === GameStatus.Explaining ? "EXPLAINING" : 
                         gameState.game_state === GameStatus.Finished ? "FINISHED" : "LOBBY"}
                    </Typography>
                </Paper>

                {players.map((player, index) => {
                    const angle = (index * 2 * Math.PI) / playerCount - Math.PI / 2;
                    const x = radius * Math.cos(angle);
                    const y = radius * Math.sin(angle);

                    const is_admin = player.id === gameState.admin;
                    const is_current = player.id === gameState.current_player;

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
                            <Tooltip title={`${player.ready ? "Ready" : "Not Ready"} | Guessed: ${player.words_guessed}/${player.words_tried}`}>
                                <Avatar
                                    sx={{
                                        width: 72,
                                        height: 72,
                                        border: `4px solid ${is_current ? theme.palette.secondary.main : (player.ready ? theme.palette.success.main : theme.palette.warning.main)}`,
                                        boxShadow: is_admin ? `0 0 16px ${theme.palette.primary.main}` : "none",
                                        bgcolor: is_current ? theme.palette.secondary.light : theme.palette.grey[400],
                                        color: is_current ? theme.palette.secondary.contrastText : "white",
                                        transition: "all 0.3s ease",
                                        transform: is_current ? "scale(1.1)" : "none",
                                    }}
                                >
                                    {player.id.slice(0, 2).toUpperCase()}
                                </Avatar>
                            </Tooltip>
                            <Typography variant="body2" sx={{ mt: 1, fontWeight: (is_admin || is_current) ? "bold" : "normal", color: is_current ? "secondary.main" : "inherit" }}>
                                {player.id.slice(0, 8)}
                                {is_admin && " 👑"}
                                {is_current && " 🎙️"}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                                Score: {player.words_guessed}
                            </Typography>
                        </Box>
                    );
                })}
            </Box>
        );
    };

    const isOurTurn = gameState?.current_player === userId;

    return (
        <Box sx={{ maxWidth: 900, mx: "auto", mt: 4, p: 2, textAlign: "center" }}>
            <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2 }}>
                <Button variant="outlined" color="inherit" onClick={handleLeaveRoom}>
                    Leave Room
                </Button>
                <Typography variant="h4">
                    Alias Online - Room: {room_id?.slice(0, 8)}
                </Typography>
                <Box>
                    <Button 
                        variant="contained" 
                        size="small" 
                        startIcon={<ContentCopyIcon />} 
                        onClick={handleShare}
                        sx={{ borderRadius: 3, px: 2 }}
                    >
                        Share
                    </Button>
                </Box>
            </Box>

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
                    <Box sx={{ mb: 4, display: "flex", justifyContent: "center", gap: 4 }}>
                        <Typography variant="subtitle1" color="text.secondary">
                            Language: {gameState.config.language}
                        </Typography>
                        <Typography variant="subtitle1" color="text.secondary">
                            Players: {Object.keys(gameState.players).length}
                        </Typography>
                    </Box>

                    {renderPlayerCircle()}

                    <Box sx={{ mt: 6, p: 3, bgcolor: "background.paper", borderRadius: 2, boxShadow: 1 }}>
                        {gameState.game_state === GameStatus.Finished ? (
                            <Box>
                                <Typography variant="h4" gutterBottom color="primary">Game Finished!</Typography>
                                <Typography variant="body1" sx={{ mb: 4 }}>The final scores are shown above.</Typography>
                                <Button variant="contained" color="primary" onClick={() => navigate("/")}>
                                    Back to Home
                                </Button>
                            </Box>
                        ) : gameState.game_state === GameStatus.RoundOver ? (
                            <Box>
                                <Typography variant="h5" gutterBottom>Round Over</Typography>
                                {isOurTurn ? (
                                    <Box>
                                        <Typography variant="body1" sx={{ mb: 2 }}>It's your turn to explain!</Typography>
                                        <Button variant="contained" color="secondary" size="large" onClick={handleStartRound}>
                                            Start My Round
                                        </Button>
                                    </Box>
                                ) : (
                                    <Typography variant="body1">
                                        Waiting for <b>{gameState.current_player.slice(0, 8)}</b> to start the round...
                                    </Typography>
                                )}

                                {gameState.admin === userId && (
                                    <Box sx={{ mt: 3, pt: 2, borderTop: "1px solid", borderColor: "divider" }}>
                                        <Typography variant="body2" sx={{ mb: 1 }} color="text.secondary">Admin Controls:</Typography>
                                        <Button variant="outlined" color="error" onClick={handleFinishGame}>
                                            Finish Game
                                        </Button>
                                    </Box>
                                )}
                            </Box>
                        ) : (
                            <Box>
                                {isOurTurn ? (
                                    <Box>
                                        <Typography variant="h6" color="secondary" gutterBottom>YOU ARE EXPLAINING</Typography>
                                        <Paper elevation={0} sx={{ p: 3, bgcolor: "secondary.light", color: "secondary.contrastText", mb: 2 }}>
                                            <Typography variant="h3">{currentWord || "..."}</Typography>
                                        </Paper>
                                        <Typography variant="body2">Explain this word to others without using its root!</Typography>
                                    </Box>
                                ) : (
                                    <Box>
                                        <Typography variant="h6" color="primary" gutterBottom>GUESS THE WORD!</Typography>
                                        <form onSubmit={handleGuess}>
                                            <TextField
                                                fullWidth
                                                variant="outlined"
                                                placeholder="Type your guess here..."
                                                value={guess}
                                                onChange={(e) => setGuess(e.target.value)}
                                                autoFocus
                                                autoComplete="off"
                                                InputProps={{
                                                    endAdornment: (
                                                        <InputAdornment position="end">
                                                            <IconButton color="primary" onClick={() => handleGuess()}>
                                                                <SendIcon />
                                                            </IconButton>
                                                        </InputAdornment>
                                                    )
                                                }}
                                                sx={{ mb: 2 }}
                                            />
                                        </form>
                                        <Typography variant="body2">
                                            <b>{gameState.current_player.slice(0, 8)}</b> is explaining...
                                        </Typography>
                                    </Box>
                                )}
                            </Box>
                        )}
                    </Box>

                    <Box sx={{ mt: 4 }}>
                        <Typography variant="h6" gutterBottom align="left">Players Scoreboard</Typography>
                        <Paper>
                            <List>
                                {Object.values(gameState.players).sort((a, b) => b.words_guessed - a.words_guessed).map((p, i) => (
                                    <Box key={p.id}>
                                        <ListItem>
                                            <ListItemText
                                                primary={p.id === userId ? `${p.id.slice(0, 8)} (You)` : p.id.slice(0, 8)}
                                                secondary={`Tried: ${p.words_tried} | Guessed: ${p.words_guessed}`}
                                            />
                                            <Typography variant="h6" color="primary">{p.words_guessed}</Typography>
                                        </ListItem>
                                        {i < Object.keys(gameState.players).length - 1 && <Divider />}
                                    </Box>
                                ))}
                            </List>
                        </Paper>
                    </Box>
                </Box>
            )}

            <Snackbar
                open={!!lastGuessResult}
                autoHideDuration={3000}
                onClose={() => setLastGuessResult(null)}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
            >
                {lastGuessResult ? (
                    <Alert onClose={() => setLastGuessResult(null)} severity={lastGuessResult.type} sx={{ width: '100%' }}>
                        {lastGuessResult.msg}
                    </Alert>
                ) : undefined}
            </Snackbar>
        </Box>
    );
}