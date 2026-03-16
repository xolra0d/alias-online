import type { PaletteMode } from "@mui/material";
import { grey, indigo, pink } from "@mui/material/colors";

export const getDesignTokens = (mode: PaletteMode) => ({
  palette: {
    mode,
    ...(mode === "light"
      ? {
          // Standard light mode palette
          primary: {
            main: "#1976d2", // Standard MUI Blue
          },
          secondary: {
            main: "#9c27b0", // Standard MUI Purple
          },
          divider: "rgba(0, 0, 0, 0.12)",
          background: {
            default: "#ffffff",
            paper: "#ffffff",
          },
          text: {
            primary: "rgba(0, 0, 0, 0.87)",
            secondary: "rgba(0, 0, 0, 0.6)",
          },
        }
      : {
          // Standard dark mode palette
          primary: {
            main: "#90caf9",
          },
          secondary: {
            main: "#ce93d8",
          },
          divider: "rgba(255, 255, 255, 0.12)",
          background: {
            default: "#121212",
            paper: "#1e1e1e",
          },
          text: {
            primary: "#ffffff",
            secondary: "rgba(255, 255, 255, 0.7)",
          },
        }),
  },
  shape: {
    borderRadius: 12,
  },
  typography: {
    fontFamily: [
      '"Inter"',
      '"system-ui"',
      '"-apple-system"',
      '"BlinkMacSystemFont"',
      '"Segoe UI"',
      '"Roboto"',
      '"Helvetica Neue"',
      '"Arial"',
      '"sans-serif"',
      '"Apple Color Emoji"',
      '"Segoe UI Emoji"',
      '"Segoe UI Symbol"',
    ].join(","),
    h4: {
      fontWeight: 700,
    },
    button: {
      textTransform: "none" as const,
      fontWeight: 600,
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          padding: "8px 20px",
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundImage: "none",
          boxShadow: mode === "light" 
            ? "0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)"
            : "0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)",
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          boxShadow: "none",
          backgroundColor: mode === "light" ? "rgba(255, 255, 255, 0.9)" : "rgba(18, 18, 18, 0.9)",
          backdropFilter: "blur(8px)",
          borderBottom: `1px solid ${mode === "light" ? "rgba(0, 0, 0, 0.12)" : "rgba(255, 255, 255, 0.12)"}`,
        },
      },
    },
  },
});
