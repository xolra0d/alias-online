import { useCallback, useEffect, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import {
  Alert,
  Avatar,
  Box,
  Button,
  Chip,
  CircularProgress,
  Divider,
  IconButton,
  InputAdornment,
  List,
  ListItem,
  ListItemText,
  Paper,
  Snackbar,
  Stack,
  TextField,
  Tooltip,
  Typography,
  useMediaQuery,
} from "@mui/material";
import { alpha, useTheme } from "@mui/material/styles";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import DoneAllIcon from "@mui/icons-material/DoneAll";
import LogoutIcon from "@mui/icons-material/Logout";
import PlayArrowIcon from "@mui/icons-material/PlayArrow";
import SendIcon from "@mui/icons-material/Send";
import SkipNextIcon from "@mui/icons-material/SkipNext";
import { ensureAuthenticated } from "../auth";

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
  Finished: 2,
} as const;

type GameStatus = (typeof GameStatus)[keyof typeof GameStatus];

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
const CLIENT_GET_NEW_WORD = 5;
const CLIENT_CREATE_ROOM = 6;
const CLIENT_LOAD_ROOM = 7;

// Server Message Types
const SERVER_NEW_UPDATE = 0;
const SERVER_CURRENT_STATE = 1;
const SERVER_YOUR_WORD = 2;
const SERVER_WORD_GUESSED = 3;
const SERVER_RIGHT_GUESS = 4;
const SERVER_WRONG_GUESS = 5;
const SERVER_REDIRECT = 6;

interface PlayWorkerResponse {
  err?: string;
  worker?: string;
}

interface CreatorBootstrap {
  creatorLogin: string;
  config: {
    language: string;
    "rude-words": boolean;
    "additional-vocabulary": string[];
    clock: number;
  };
}

const CREATOR_ROOM_STORAGE_PREFIX = "creator-room:";

const workerToWsBase = (worker: string): string => {
  if (worker.startsWith("http://")) {
    return worker.replace(/^http:/, "ws:");
  }
  if (worker.startsWith("https://")) {
    return worker.replace(/^https:/, "wss:");
  }
  if (worker.startsWith("ws://") || worker.startsWith("wss://")) {
    return worker;
  }
  const secure = window.location.protocol === "https:";
  return `${secure ? "wss" : "ws"}://${worker}`;
};

const getPlayerLabel = (player: Player, currentUserId: string | null) =>
  player.id === currentUserId
    ? `${player.name ?? player.id.slice(0, 8)} (You)`
    : player.name ?? player.id.slice(0, 8);

export default function Play() {
  const { room_id } = useParams<{ room_id: string }>();
  const navigate = useNavigate();
  const theme = useTheme();
  const isSmall = useMediaQuery(theme.breakpoints.down("sm"));
  const isMedium = useMediaQuery(theme.breakpoints.between("sm", "lg"));

  const [error, setError] = useState<string | null>(null);
  const [connecting, setConnecting] = useState(true);
  const [gameState, setGameState] = useState<GameState | null>(null);
  const gameStateRef = useRef<GameState | null>(null);
  const [timeLeft, setTimeLeft] = useState<number>(0);
  const [currentWord, setCurrentWord] = useState<string | null>(null);
  const [guess, setGuess] = useState("");
  const [lastGuessResult, setLastGuessResult] = useState<{
    msg: string;
    type: "success" | "error" | "info";
  } | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const [userId, setUserId] = useState<string | null>(
    localStorage.getItem("login"),
  );

  const boardSize = isSmall ? 320 : isMedium ? 440 : 540;
  const centerSize = isSmall ? 112 : isMedium ? 128 : 148;
  const avatarSize = isSmall ? 56 : isMedium ? 64 : 72;
  const radius = boardSize * 0.34;

  const sortedPlayers = gameState
    ? Object.values(gameState.players).sort(
        (a, b) => b.words_guessed - a.words_guessed,
      )
    : [];

  const readyPlayers = gameState
    ? Object.values(gameState.players).filter((player) => player.ready)
    : [];

  const currentPlayerName = gameState
    ? gameState.players[gameState.current_player]?.name ??
      gameState.current_player.slice(0, 8)
    : "";

  const isOurTurn = gameState?.current_player === userId;
  const isExplaining = gameState?.game_state === GameStatus.Explaining;
  const isFinished = gameState?.game_state === GameStatus.Finished;
  const isRoundOver = gameState?.game_state === GameStatus.RoundOver;

  const roomStatus =
    isFinished ? "Finished" : isExplaining ? "Explaining" : "Lobby";
  const roomStatusColor = isFinished
    ? "success.main"
    : isExplaining
      ? "secondary.main"
      : "primary.main";
  const roomStatusBackground = isFinished
    ? alpha(theme.palette.success.main, 0.12)
    : isExplaining
      ? alpha(theme.palette.secondary.main, 0.12)
      : alpha(theme.palette.primary.main, 0.12);

  const sendMessage = useCallback(
    (type: number, data: Record<string, unknown> = {}) => {
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.send(
          JSON.stringify({
            user_id: userId,
            type,
            data,
          }),
        );
      }
    },
    [userId],
  );

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

    if (isExplaining) {
      timer = window.setInterval(() => {
        setTimeLeft((prev) => (prev > 0 ? prev - 1 : 0));
      }, 1000);
    }

    return () => {
      if (timer) clearInterval(timer);
    };
  }, [isExplaining]);

  useEffect(() => {
    let cancelled = false;

    const connect = async () => {
      const httpBase = BACKEND_URL || window.location.origin;
      if (!room_id) {
        setError("Room id is missing.");
        return;
      }
      const creds = await ensureAuthenticated(httpBase);
      const login = creds.login;
      const name = creds.name;
      setUserId(creds.login);

      if (cancelled) return;

      const creatorBootstrapRaw = sessionStorage.getItem(
        `${CREATOR_ROOM_STORAGE_PREFIX}${room_id}`,
      );
      let firstMessageType = CLIENT_LOAD_ROOM;
      let firstMessageData: Record<string, unknown> = {};

      if (creatorBootstrapRaw) {
        try {
          const parsed = JSON.parse(creatorBootstrapRaw) as CreatorBootstrap;
          if (
            parsed &&
            parsed.creatorLogin === login &&
            parsed.config &&
            typeof parsed.config.language === "string" &&
            typeof parsed.config["rude-words"] === "boolean" &&
            Array.isArray(parsed.config["additional-vocabulary"]) &&
            typeof parsed.config.clock === "number"
          ) {
            firstMessageType = CLIENT_CREATE_ROOM;
            firstMessageData = parsed.config;
            sessionStorage.removeItem(`${CREATOR_ROOM_STORAGE_PREFIX}${room_id}`);
          }
        } catch {
          sessionStorage.removeItem(`${CREATOR_ROOM_STORAGE_PREFIX}${room_id}`);
        }
      }

      const playResp = await fetch(`${httpBase}/api/protected/play/${room_id}`, {
        method: "GET",
        credentials: "include",
      });
      const playData: PlayWorkerResponse = await playResp.json();
      if (!playResp.ok || !playData.worker) {
        setError(playData.err || "Failed to resolve room worker.");
        setConnecting(false);
        return;
      }

      const connectToWorker = (
        worker: string,
        initialType: number,
        initialData: Record<string, unknown>,
      ) => {
        const wsUrl = new URL(`/api/play/${room_id}`, workerToWsBase(worker));
        wsUrl.searchParams.set("name", name!);

        const ws = new WebSocket(wsUrl.toString());
        wsRef.current = ws;

        ws.onopen = () => {
          if (!cancelled) {
            ws.send(
              JSON.stringify({
                user_id: login,
                type: initialType,
                data: initialData,
              }),
            );
            setConnecting(false);
            setError(null);
            fetchState();
          }
        };

        ws.onmessage = async (event) => {
          const text =
            event.data instanceof Blob ? await event.data.text() : event.data;

          try {
            const json = JSON.parse(text);

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
                setLastGuessResult({
                  msg: `Word guessed by ${json.msg_data.guesser.slice(0, 8)}!`,
                  type: "info",
                });
                if (gameStateRef.current?.current_player === userId) {
                  setCurrentWord(null);
                }
                break;
              case SERVER_RIGHT_GUESS:
                setLastGuessResult({ msg: "Correct guess!", type: "success" });
                break;
              case SERVER_WRONG_GUESS:
                setLastGuessResult({
                  msg: "Wrong guess, try again!",
                  type: "error",
                });
                break;
              case SERVER_REDIRECT: {
                const nextWorker = json.msg_data?.worker;
                if (typeof nextWorker === "string" && nextWorker.length > 0) {
                  ws.close();
                  connectToWorker(nextWorker, initialType, initialData);
                }
                break;
              }
            }
          } catch (parseError) {
            console.error("[WebSocket] Error parsing message:", parseError, text);
          }
        };

        ws.onerror = (event) => {
          console.error("[WebSocket] Error:", event);
          if (!cancelled) {
            setError("WebSocket connection error.");
            setConnecting(false);
          }
        };

        ws.onclose = (event) => {
          console.log("[WebSocket] Closed:", event.code, event.reason);
          if (!cancelled) setConnecting(false);
        };
      };

      connectToWorker(playData.worker, firstMessageType, firstMessageData);
    };

    connect().catch((e) => {
      console.error("[Play] connect failed", e);
      if (!cancelled) {
        setError("Failed to connect to room.");
        setConnecting(false);
      }
    });

    return () => {
      cancelled = true;
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [room_id, fetchState, userId]);

  // Handle "Get Word" automatically if it's our turn to explain
  useEffect(() => {
    if (
      gameState &&
      gameState.game_state === GameStatus.Explaining &&
      gameState.current_player === userId
    ) {
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

  const handleSkipWord = () => {
    sendMessage(CLIENT_GET_NEW_WORD);
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

  const handleShare = async () => {
    try {
      await navigator.clipboard.writeText(window.location.href);
      setLastGuessResult({
        msg: "Room link copied to clipboard!",
        type: "success",
      });
    } catch {
      setError("Failed to copy the room link.");
    }
  };

  const renderPlayerCircle = () => {
    if (!gameState) return null;

    return (
      <Box
        sx={{
          position: "relative",
          width: `min(100%, ${boardSize}px)`,
          height: "auto",
          aspectRatio: "1 / 1",
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          mx: "auto",
        }}
      >
        <Paper
          elevation={0}
          sx={{
            width: centerSize,
            height: centerSize,
            borderRadius: "50%",
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            alignItems: "center",
            bgcolor: isFinished
              ? "success.main"
              : isExplaining
                ? "secondary.main"
                : "primary.main",
            color: isFinished
              ? "success.contrastText"
              : isExplaining
                ? "secondary.contrastText"
                : "primary.contrastText",
            borderColor: "transparent",
            zIndex: 1,
            textAlign: "center",
            p: 2,
          }}
        >
          <Typography
            variant="h4"
            sx={{
              fontSize: `clamp(${isSmall ? "1.8rem" : "2rem"}, 4vw, 3rem)`,
            }}
          >
            {isFinished ? "GG" : `${timeLeft}s`}
          </Typography>
          <Typography variant="caption" sx={{ letterSpacing: "0.2em" }}>
            {roomStatus.toUpperCase()}
          </Typography>
        </Paper>

        {readyPlayers.map((player, index) => {
          const safeCount = Math.max(readyPlayers.length, 1);
          const angle = (index * 2 * Math.PI) / safeCount - Math.PI / 2;
          const x = radius * Math.cos(angle);
          const y = radius * Math.sin(angle);

          const isAdmin = player.id === gameState.admin;
          const isCurrent = player.id === gameState.current_player;
          const isMe = player.id === userId;

          return (
            <Box
              key={player.id}
              sx={{
                position: "absolute",
                transform: `translate(${x}px, ${y}px)`,
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                transition: "all 0.5s ease-in-out",
              }}
            >
              <Tooltip
                title={`${player.ready ? "Ready" : "Not Ready"} | Guessed: ${player.words_guessed}/${player.words_tried}`}
              >
                <Avatar
                  sx={{
                    width: avatarSize,
                    height: avatarSize,
                    border: `4px solid ${
                      isCurrent
                        ? theme.palette.secondary.main
                        : player.ready
                          ? theme.palette.success.main
                          : theme.palette.warning.main
                    }`,
                    boxShadow: isMe
                      ? `0 0 0 6px ${alpha(theme.palette.primary.main, 0.1)}`
                      : "none",
                    bgcolor: isCurrent
                      ? theme.palette.secondary.light
                      : alpha(theme.palette.background.paper, 0.9),
                    color: isCurrent
                      ? theme.palette.secondary.contrastText
                      : theme.palette.text.primary,
                    transition: "all 0.3s ease",
                    transform: isCurrent ? "scale(1.08)" : "none",
                    fontFamily: '"Josefin Sans", Georgia, serif',
                    fontWeight: 700,
                  }}
                >
                  {player.id.slice(0, 2).toUpperCase()}
                </Avatar>
              </Tooltip>
              <Typography
                variant="body2"
                sx={{
                  mt: 1,
                  fontWeight: isMe || isCurrent ? 700 : 500,
                  color: isCurrent ? "secondary.main" : "inherit",
                }}
              >
                {player.name}
                {isAdmin && " 👑"}
                {isCurrent && " 🎙️"}
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

  const renderPrimaryPanel = () => {
    if (!gameState) return null;

    if (isFinished) {
      return (
        <Stack spacing={2.5}>
          <Box>
            <Typography variant="overline" sx={{ color: "success.main" }}>
              Game complete
            </Typography>
            <Typography variant="h5">The table has wrapped up</Typography>
          </Box>
          <Typography variant="body1" color="text.secondary">
            Final scores live in the sidebar. You can jump back to the lobby and
            start a new room whenever you&apos;re ready.
          </Typography>
          <Button variant="contained" onClick={() => navigate("/")}>
            Back to home
          </Button>
        </Stack>
      );
    }

    if (isRoundOver) {
      return (
        <Stack spacing={2.5}>
          <Box>
            <Typography variant="overline" sx={{ color: "primary.main" }}>
              Round paused
            </Typography>
            <Typography variant="h5">Waiting for the next clue</Typography>
          </Box>

          {isOurTurn ? (
            <Stack spacing={1.5}>
              <Typography variant="body1" color="text.secondary">
                It&apos;s your turn to explain. Start the round when you&apos;re
                ready.
              </Typography>
              <Button
                variant="contained"
                color="secondary"
                size="large"
                startIcon={<PlayArrowIcon />}
                onClick={handleStartRound}
              >
                Start my round
              </Button>
            </Stack>
          ) : (
            <Typography variant="body1" color="text.secondary">
              Waiting for <b>{currentPlayerName}</b> to start the round.
            </Typography>
          )}
        </Stack>
      );
    }

    return isOurTurn ? (
      <Stack spacing={2.5}>
        <Box
          sx={{ display: "flex", justifyContent: "space-between", gap: 2, alignItems: "center" }}
        >
          <Box>
            <Typography variant="overline" sx={{ color: "secondary.main" }}>
              Your turn
            </Typography>
            <Typography variant="h5">Your secret word is ready</Typography>
          </Box>
          <Chip
            label={`${timeLeft}s`}
            sx={{
              bgcolor: roomStatusBackground,
              color: roomStatusColor,
              border: "1px solid",
              borderColor: roomStatusColor,
            }}
          />
        </Box>

        <Paper
          elevation={0}
          sx={{
            p: 3,
            bgcolor: "secondary.main",
            color: "secondary.contrastText",
            borderColor: "secondary.main",
          }}
        >
          <Stack spacing={1}>
            <Typography variant="overline" sx={{ opacity: 0.85 }}>
              Keep it hidden
            </Typography>
            <Typography
              variant="h3"
              sx={{ fontSize: "clamp(2.4rem, 6vw, 4.4rem)" }}
            >
              {currentWord || "..."}
            </Typography>
            <Typography variant="body2" sx={{ opacity: 0.9 }}>
              Explain the word with clues, but do not say the word itself.
            </Typography>
          </Stack>
        </Paper>

        <Stack direction="row" spacing={1.5} flexWrap="wrap">
          <Button
            variant="outlined"
            color="secondary"
            startIcon={<SkipNextIcon />}
            onClick={handleSkipWord}
            sx={{ borderStyle: "dashed" }}
          >
            Skip word
          </Button>
        </Stack>
      </Stack>
    ) : (
      <Stack spacing={2.5}>
        <Box>
          <Typography variant="overline" sx={{ color: "primary.main" }}>
            Guessing round
          </Typography>
          <Typography variant="h5">
            {currentPlayerName} is explaining
          </Typography>
        </Box>

        <Typography variant="body1" color="text.secondary">
          Type your guess and send it before the clock runs out.
        </Typography>

        <Box component="form" onSubmit={handleGuess}>
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
              ),
            }}
          />
        </Box>
      </Stack>
    );
  };

  if (connecting && !error) {
    return (
      <Box
        display="flex"
        alignItems="center"
        justifyContent="center"
        gap={2}
        minHeight="50vh"
      >
        <CircularProgress size={20} />
        <Typography>Connecting to room...</Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ maxWidth: 1240, mx: "auto" }}>
      <Stack spacing={3}>
        <Box
          sx={{
            display: "flex",
            flexWrap: "wrap",
            justifyContent: "space-between",
            gap: 2.5,
            alignItems: "flex-start",
          }}
        >
          <Box>
            <Typography
              variant="overline"
              sx={{ color: roomStatusColor, display: "inline-flex" }}
            >
              Live room
            </Typography>
            <Typography
              variant="h3"
              component="h1"
              sx={{ fontSize: { xs: "clamp(2.4rem, 8vw, 3.4rem)", md: "clamp(3rem, 4.4vw, 4.4rem)" } }}
            >
              Room {room_id?.slice(0, 8)}
            </Typography>
            <Stack direction="row" flexWrap="wrap" gap={1} sx={{ mt: 1.5 }}>
              <Chip label={roomStatus} sx={{ bgcolor: roomStatusBackground, color: roomStatusColor }} />
              <Chip
                label={`Language: ${gameState?.config.language ?? "Loading"}`}
                variant="outlined"
              />
              <Chip
                label={`Players: ${gameState ? Object.keys(gameState.players).length : 0}`}
                variant="outlined"
              />
              {gameState && (
                <Chip
                  label={`Clock: ${gameState.config.clock}s`}
                  variant="outlined"
                />
              )}
            </Stack>
          </Box>

          <Stack direction="row" spacing={1.25} alignItems="center">
            <Button
              variant="outlined"
              startIcon={<LogoutIcon />}
              onClick={handleLeaveRoom}
            >
              Leave room
            </Button>
            <Button
              variant="contained"
              startIcon={<ContentCopyIcon />}
              onClick={handleShare}
            >
              Share room
            </Button>
          </Stack>
        </Box>

        {error && (
          <Alert severity="error" sx={{ borderRadius: 2 }}>
            {error}
          </Alert>
        )}

        {!connecting && gameState && (
          <Box
            sx={{
              display: "grid",
              gridTemplateColumns: {
                xs: "1fr",
                lg: "minmax(0, 1.35fr) minmax(320px, 0.65fr)",
              },
              gap: { xs: 3, lg: 4 },
              alignItems: "start",
            }}
          >
            <Stack spacing={3}>
              <Paper
                elevation={0}
                sx={{
                  p: { xs: 3, md: 4 },
                  position: "relative",
                  overflow: "hidden",
                  bgcolor: "background.paper",
                }}
              >
                <Box
                  sx={{
                    position: "absolute",
                    inset: 0,
                    pointerEvents: "none",
                    background:
                      "linear-gradient(135deg, rgba(185, 119, 30, 0.08) 0%, transparent 38%, transparent 62%, rgba(44, 139, 146, 0.08) 100%)",
                  }}
                />
                <Box sx={{ position: "relative" }}>{renderPrimaryPanel()}</Box>
              </Paper>

              <Paper
                elevation={0}
                sx={{
                  p: { xs: 2.5, md: 3 },
                  position: "relative",
                  overflow: "hidden",
                  bgcolor: "background.paper",
                }}
              >
                <Box
                  sx={{
                    position: "absolute",
                    inset: 0,
                    pointerEvents: "none",
                    borderTop: "4px solid",
                    borderColor: "secondary.main",
                    opacity: 0.45,
                  }}
                />
                <Box sx={{ position: "relative" }}>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="h5">Player circle</Typography>
                    <Typography variant="body2" color="text.secondary">
                      Ready players gather around the table.
                    </Typography>
                  </Box>
                  {renderPlayerCircle()}
                </Box>
              </Paper>
            </Stack>

            <Stack
              spacing={3}
              sx={{
                position: { lg: "sticky" },
                top: { lg: 104 },
              }}
            >
              <Paper
                elevation={0}
                sx={{
                  p: { xs: 3, md: 4 },
                  bgcolor: "background.paper",
                }}
              >
                <Stack spacing={2}>
                  <Box>
                    <Typography variant="h5">Scoreboard</Typography>
                    <Typography variant="body2" color="text.secondary">
                      Players are sorted by guessed words.
                    </Typography>
                  </Box>

                  <List disablePadding>
                    {sortedPlayers.map((player, index) => (
                      <Box key={player.id}>
                        <ListItem
                          sx={{
                            px: 0,
                            py: 1.25,
                            alignItems: "flex-start",
                            gap: 1.5,
                          }}
                        >
                          <ListItemText
                            primary={getPlayerLabel(player, userId)}
                            secondary={`Tried: ${player.words_tried} | Guessed: ${player.words_guessed}`}
                          />
                          <Stack alignItems="flex-end" spacing={0.5}>
                            <Typography variant="h6" color="primary">
                              {player.words_guessed}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                              points
                            </Typography>
                          </Stack>
                        </ListItem>
                        {index < sortedPlayers.length - 1 && <Divider />}
                      </Box>
                    ))}
                  </List>
                </Stack>
              </Paper>

              <Paper
                elevation={0}
                sx={{
                  p: { xs: 3, md: 4 },
                  bgcolor: "background.paper",
                }}
              >
                <Stack spacing={2}>
                  <Box>
                    <Typography variant="h5">Room details</Typography>
                    <Typography variant="body2" color="text.secondary">
                      The room settings are locked in for this game.
                    </Typography>
                  </Box>

                  <Stack spacing={1}>
                    <Typography variant="body2" color="text.secondary">
                      Language: <b>{gameState.config.language}</b>
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Clock: <b>{gameState.config.clock}s</b>
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Rude words:{" "}
                      <b>{gameState.config["rude-words"] ? "Allowed" : "Off"}</b>
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Additional words:{" "}
                      <b>{gameState.config["additional-vocabulary"].length}</b>
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Current player: <b>{currentPlayerName}</b>
                    </Typography>
                  </Stack>

                  {gameState.admin === userId && (
                    <Button
                      variant="outlined"
                      color="error"
                      startIcon={<DoneAllIcon />}
                      onClick={handleFinishGame}
                    >
                      Finish game
                    </Button>
                  )}
                </Stack>
              </Paper>
            </Stack>
          </Box>
        )}
      </Stack>

      <Snackbar
        open={!!lastGuessResult}
        autoHideDuration={3000}
        onClose={() => setLastGuessResult(null)}
        anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
      >
        {lastGuessResult ? (
          <Alert
            onClose={() => setLastGuessResult(null)}
            severity={lastGuessResult.type}
            sx={{ width: "100%" }}
          >
            {lastGuessResult.msg}
          </Alert>
        ) : undefined}
      </Snackbar>
    </Box>
  );
}
