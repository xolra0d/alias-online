import type { PaletteMode } from "@mui/material";
import { grey, indigo, pink } from "@mui/material/colors";

export const getDesignTokens = (mode: PaletteMode) => ({
  palette: {
    mode,
    ...(mode === "light"
      ? {
          // Modern palette for light mode
          primary: {
            main: indigo[500],
          },
          secondary: {
            main: pink[500],
          },
          divider: grey[200],
          background: {
            default: "#f8fafc", // Slate 50
            paper: "#ffffff",
          },
          text: {
            primary: grey[900],
            secondary: grey[600],
          },
        }
      : {
          // Modern palette for dark mode
          primary: {
            main: indigo[400],
          },
          secondary: {
            main: pink[400],
          },
          divider: grey[800],
          background: {
            default: "#0f172a", // Slate 900
            paper: "#1e293b", // Slate 800
          },
          text: {
            primary: "#f8fafc",
            secondary: grey[400],
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
          backgroundColor: mode === "light" ? "rgba(255, 255, 255, 0.8)" : "rgba(15, 23, 42, 0.8)",
          backdropFilter: "blur(8px)",
          borderBottom: `1px solid ${mode === "light" ? grey[200] : grey[800]}`,
        },
      },
    },
  },
});
