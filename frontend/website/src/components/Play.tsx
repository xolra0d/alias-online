import React, { useEffect, useState } from "react";
import {
  Typography,
  Box,
  Radio,
  RadioGroup,
  FormControlLabel,
  FormControl,
  FormLabel,
  Checkbox,
  TextField,
  Button,
  CircularProgress,
  Alert,
} from "@mui/material";

interface UserCredentials {
  id: string;
  secret: string;
}

interface CreateUserResponse {
  ok: boolean;
  reason?: string;
  credentials?: UserCredentials;
  name?: string;
}

interface CreateRoomResponse {
  ok: boolean;
  reason?: string;
  room_id?: string;
}

interface AvailableVocabsResponse {
  ok: boolean;
  reason?: string;
  languages?: string[];
}

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || "";

export default function Play() {
  const [languages, setLanguages] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    language: "en",
    "rude-words": false,
    additionalVocabulary: "",
    clock: 60,
  });

  useEffect(() => {
    fetch(`${BACKEND_URL}/api/7available-vocabs`)
      .then((res) => res.json())
      .then((data: AvailableVocabsResponse) => {
        if (data.ok && data.languages) {
          setLanguages(data.languages);
          if (data.languages.length > 0 && !data.languages.includes("en")) {
            setFormData((prev) => ({ ...prev, language: data.languages![0] }));
          }
        } else {
          setError(`Failed to load languages: ${data.reason}`);
        }
        setLoading(false);
      })
      .catch(() => {
        setError("Failed to load languages. Please check if the backend is running.");
        setLoading(false);
      });
  }, []);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value, type } = e.target;
    const val = type === "checkbox" ? (e.target as HTMLInputElement).checked : value;

    setFormData((prev) => ({
      ...prev,
      [name]: val,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

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
      if (data.ok) {
        console.log("Room created, ID:", data.room_id);
      } else {
        setError(`Failed to create room: ${data.reason}`);
      }
    } catch {
      setError("Network error while creating room.");
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" mt={4}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ maxWidth: 600, mx: "auto", mt: 4, p: 2 }}>
      <Typography variant="h4" gutterBottom>
        Create a Room
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <form onSubmit={handleSubmit}>
        <FormControl component="fieldset" fullWidth margin="normal">
          <FormLabel component="legend">Language</FormLabel>
          <RadioGroup
            name="language"
            value={formData.language}
            onChange={handleInputChange}
          >
            {languages.map((lang) => (
              <FormControlLabel
                key={lang}
                value={lang}
                control={<Radio />}
                label={lang}
              />
            ))}
          </RadioGroup>
        </FormControl>

        <Box display="flex" flexDirection="column" gap={1} mb={2}>
          <FormControlLabel
            control={
              <Checkbox
                checked={formData["rude-words"]}
                onChange={handleInputChange}
                name="rude-words"
              />
            }
            label="Allow rude words"
          />
        </Box>

        <TextField
          fullWidth
          label="Additional vocabulary"
          name="additionalVocabulary"
          placeholder="word1,word2,word3"
          multiline
          rows={3}
          value={formData.additionalVocabulary}
          onChange={handleInputChange}
          margin="normal"
        />

        <TextField
          fullWidth
          label="Clock in seconds (-1 for no clock)"
          name="clock"
          type="number"
          value={formData.clock}
          onChange={handleInputChange}
          margin="normal"
        />

        <Button
          type="submit"
          variant="contained"
          color="primary"
          fullWidth
          sx={{ mt: 3 }}
        >
          Create Room
        </Button>
      </form>
    </Box>
  );
}
