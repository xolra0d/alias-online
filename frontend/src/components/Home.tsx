import { useCallback, useEffect, useState } from "react";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Checkbox from "@mui/material/Checkbox";
import CircularProgress from "@mui/material/CircularProgress";
import Alert from "@mui/material/Alert";
import Divider from "@mui/material/Divider";
import FormControlLabel from "@mui/material/FormControlLabel";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import TextField from "@mui/material/TextField";
import ToggleButton from "@mui/material/ToggleButton";
import ToggleButtonGroup from "@mui/material/ToggleButtonGroup";
import Typography from "@mui/material/Typography";
import { alpha, useTheme } from "@mui/material/styles";
import { useNavigate } from "react-router";

interface UserCredentials {
  id: string;
  secret: string;
  name: string;
}

interface CreateUserResponse {
  err?: string;
  credentials?: UserCredentials;
}

interface CreateRoomResponse {
  err?: string;
  room?: string;
}

interface AvailableVocabsResponse {
  err?: string;
  languages?: string[];
}

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || "";

const roomSteps = [
  {
    title: "Create the room",
    body: "Pick a language, set the round clock, and add any house vocabulary you want to keep in play.",
  },
  {
    title: "Share the link",
    body: "Send the room URL to friends so they can join the same table instantly.",
  },
  {
    title: "Take turns explaining",
    body: "One player sees the secret word while the others race to guess it before the timer runs out.",
  },
  {
    title: "Score every round",
    body: "Correct guesses award points and keep the room moving until the game is finished.",
  },
];

export default function Home() {
  const [languages, setLanguages] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const theme = useTheme();
  const navigate = useNavigate();

  const [formData, setFormData] = useState({
    language: "",
    "rude-words": false,
    additionalVocabulary: "",
    clock: 60,
  });

  const loadLanguages = useCallback(() => {
    fetch(`${BACKEND_URL}/api/available-vocabs`)
      .then((res) => res.json())
      .then((data: AvailableVocabsResponse) => {
        if (data.languages && data.languages.length > 0) {
          setLanguages(data.languages);
          setFormData((prev) =>
            prev.language && data.languages!.includes(prev.language)
              ? prev
              : { ...prev, language: data.languages![0] },
          );
        } else {
          setLoadError(
            `Failed to load languages: ${data.err ?? "no languages available"}`,
          );
        }
      })
      .catch(() => {
        setLoadError(
          "Failed to load languages. Please check if the backend is running.",
        );
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  useEffect(() => {
    loadLanguages();
  }, [loadLanguages]);

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    const { name, value, type } = e.target;
    const val =
      type === "checkbox" ? (e.target as HTMLInputElement).checked : value;

    setFormData((prev) => ({
      ...prev,
      [name]: val,
    }));
  };

  const handleLanguageChange = (
    _event: React.MouseEvent<HTMLElement>,
    value: string | null,
  ) => {
    if (value) {
      setFormData((prev) => ({
        ...prev,
        language: value,
      }));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError(null);

    let login = localStorage.getItem("login");
    let name = localStorage.getItem("name");
    let secret = localStorage.getItem("secret");

    if (!login || !secret) {
      try {
        const response = await fetch(`${BACKEND_URL}/api/create-user`, {
          method: "POST",
        });
        const data: CreateUserResponse = await response.json();

        if (data.credentials) {
          login = data.credentials.id;
          name = data.credentials.name;
          secret = data.credentials.secret;

          localStorage.setItem("login", login);
          localStorage.setItem("secret", secret);
          localStorage.setItem("name", name);
        } else {
          setFormError(`Failed to create user: ${data.err}`);
          return;
        }
      } catch {
        setFormError("Network error while creating user.");
        return;
      }
    }

    try {
      const params = new URLSearchParams();
      params.append("language", formData.language);
      params.append("rude-words", String(formData["rude-words"]));
      params.append("additional-vocabulary", formData.additionalVocabulary);
      params.append("clock", String(formData.clock));

      const response = await fetch(`${BACKEND_URL}/api/protected/create-room`, {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
          "User-Id": login || "",
          "User-Secret": secret || "",
        },
        body: params,
      });

      const data: CreateRoomResponse = await response.json();
      if (data.room) {
        navigate(`/play/${data.room}`);
      } else {
        setFormError(`Failed to create room: ${data.err}`);
      }
    } catch {
      setFormError("Network error while creating room.");
    }
  };

  if (loading) {
    return (
      <Box
        display="flex"
        justifyContent="center"
        alignItems="center"
        minHeight="50vh"
      >
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ maxWidth: 1200, mx: "auto" }}>
      <Stack spacing={2.5} sx={{ mb: { xs: 4, md: 6 } }}>
        <Typography
          variant="overline"
          sx={{ color: "primary.main", display: "inline-flex" }}
        >
          Word play
        </Typography>

        <Typography
          variant="h2"
          component="h1"
          sx={{
            maxWidth: "12ch",
            fontSize: {
              xs: "clamp(3rem, 12vw, 4.2rem)",
              md: "clamp(4.5rem, 6.8vw, 6rem)",
            },
          }}
        >
          Explain the word without saying it.
        </Typography>

        <Typography
          variant="body1"
          color="text.secondary"
          sx={{ maxWidth: 720, fontSize: { xs: "1rem", md: "1.1rem" } }}
        >
          Create a private room, tune the round clock, and hand the link to your
          friends for a fast, noisy, score-chasing game night.
        </Typography>
      </Stack>

      <Box
        sx={{
          display: "grid",
          gridTemplateColumns: {
            xs: "1fr",
            lg: "minmax(0, 1fr) minmax(320px, 0.82fr)",
          },
          gap: { xs: 3, lg: 4 },
          alignItems: "start",
        }}
      >
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
                "linear-gradient(135deg, rgba(185, 119, 30, 0.08) 0%, transparent 34%, transparent 66%, rgba(44, 139, 146, 0.08) 100%)",
            }}
          />

          <Stack spacing={3} sx={{ position: "relative" }}>
            <Box>
              <Typography variant="h5">Room settings</Typography>
              <Typography variant="body2" color="text.secondary">
                Tune the room before you hand the link to the table.
              </Typography>
            </Box>

            {loadError ? (
              <Stack spacing={2}>
                <Alert severity="error" sx={{ borderRadius: 2 }}>
                  {loadError}
                </Alert>
                <Typography variant="body2" color="text.secondary">
                  The room settings cannot load until the backend responds.
                  Start the server and retry the connection.
                </Typography>
                <Button
                  variant="contained"
                  size="large"
                  onClick={() => {
                    setLoading(true);
                    setLoadError(null);
                    loadLanguages();
                  }}
                  sx={{ alignSelf: "flex-start" }}
                >
                  Retry connection
                </Button>
              </Stack>
            ) : (
              <Box component="form" onSubmit={handleSubmit}>
                <Stack spacing={3}>
                  {formError && (
                    <Alert severity="error" sx={{ borderRadius: 2 }}>
                      {formError}
                    </Alert>
                  )}

                  <Box>
                    <Typography variant="subtitle2" sx={{ mb: 1 }}>
                      Language
                    </Typography>
                    <ToggleButtonGroup
                      exclusive
                      fullWidth
                      value={formData.language}
                      onChange={handleLanguageChange}
                      sx={{
                        display: "flex",
                        flexWrap: "wrap",
                        gap: 1,
                        "& .MuiToggleButton-root": {
                          flex: "1 1 160px",
                        },
                      }}
                    >
                      {languages.map((lang) => (
                        <ToggleButton key={lang} value={lang}>
                          {lang}
                        </ToggleButton>
                      ))}
                    </ToggleButtonGroup>
                  </Box>

                  <Divider />

                  <TextField
                    fullWidth
                    label="Clock in seconds"
                    name="clock"
                    type="number"
                    value={formData.clock}
                    onChange={handleInputChange}
                    inputProps={{ min: 15, step: 15 }}
                    helperText="Longer rounds work well for larger groups."
                  />

                  <FormControlLabel
                    control={
                      <Checkbox
                        checked={formData["rude-words"]}
                        onChange={handleInputChange}
                        name="rude-words"
                        color="secondary"
                      />
                    }
                    label="Allow rude words"
                  />

                  <TextField
                    fullWidth
                    label="Additional vocabulary"
                    name="additionalVocabulary"
                    placeholder="word1, word2, word3"
                    multiline
                    rows={4}
                    value={formData.additionalVocabulary}
                    onChange={handleInputChange}
                    helperText="Comma-separated extras that will be mixed into the room."
                  />

                  <Button
                    type="submit"
                    variant="contained"
                    size="large"
                    fullWidth
                  >
                    Create room
                  </Button>
                </Stack>
              </Box>
            )}
          </Stack>
        </Paper>

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
                borderLeft: "4px solid",
                borderColor: "primary.main",
                opacity: 0.65,
              }}
            />

            <Stack spacing={2.25} sx={{ position: "relative" }}>
              <Box>
                <Typography variant="h5">How the table works</Typography>
                <Typography variant="body2" color="text.secondary">
                  Four quick beats to get the room moving.
                </Typography>
              </Box>

              <Stack spacing={1.5}>
                {roomSteps.map((step, index) => (
                  <Box
                    key={step.title}
                    sx={{
                      display: "flex",
                      gap: 1.5,
                      p: 2,
                      border: "1px solid",
                      borderColor: "divider",
                      borderRadius: 3,
                      bgcolor: alpha(theme.palette.background.default, 0.36),
                    }}
                  >
                    <Box
                      sx={{
                        width: 36,
                        height: 36,
                        borderRadius: 1.5,
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "center",
                        bgcolor: "primary.main",
                        color: "primary.contrastText",
                        flexShrink: 0,
                        fontFamily: '"Josefin Sans", Georgia, serif',
                        fontWeight: 700,
                        letterSpacing: "0.08em",
                      }}
                    >
                      {String(index + 1).padStart(2, "0")}
                    </Box>
                    <Box>
                      <Typography variant="subtitle1" sx={{ display: "block" }}>
                        {step.title}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        {step.body}
                      </Typography>
                    </Box>
                  </Box>
                ))}
              </Stack>

              <Divider />

              <Box>
                <Typography variant="h5">About Alias</Typography>
                <Typography variant="body2" color="text.secondary">
                  The classic word explanation game, tuned for fast rooms and
                  short rounds.
                </Typography>
              </Box>

              <Alert
                severity="warning"
                sx={{
                  borderRadius: 3,
                  "& .MuiAlert-message": { width: "100%" },
                }}
              >
                <Typography
                  variant="subtitle2"
                  sx={{ fontWeight: 700, mb: 0.5 }}
                >
                  Don&apos;t abuse skipping.
                </Typography>
                <Typography variant="body2">
                  Every room has a limited vocabulary for the chosen language.
                  Skipping too aggressively can exhaust the room and end the
                  game early. For example, in English vocabulary, there are only
                  1100 words.
                </Typography>
              </Alert>
            </Stack>
          </Paper>
        </Stack>
      </Box>
    </Box>
  );
}
