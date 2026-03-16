// import type { PaletteMode } from "@mui/material";
// import { grey, indigo, pink } from "@mui/material/colors";
// 
// export const getDesignTokens = (mode: PaletteMode) => ({
//   palette: {
//     mode,
//     ...(mode === "light"
//       ? {
//           // Standard light mode palette
//           primary: {
//             main: "#1976d2", // Standard MUI Blue
//           },
//           secondary: {
//             main: "#9c27b0", // Standard MUI Purple
//           },
//           divider: "rgba(0, 0, 0, 0.12)",
//           background: {
//             default: "#ffffff",
//             paper: "#ffffff",
//           },
//           text: {
//             primary: "rgba(0, 0, 0, 0.87)",
//             secondary: "rgba(0, 0, 0, 0.6)",
//           },
//         }
//       : {
//           // Standard dark mode palette
//           primary: {
//             main: "#90caf9",
//           },
//           secondary: {
//             main: "#ce93d8",
//           },
//           divider: "rgba(255, 255, 255, 0.12)",
//           background: {
//             default: "#121212",
//             paper: "#1e1e1e",
//           },
//           text: {
//             primary: "#ffffff",
//             secondary: "rgba(255, 255, 255, 0.7)",
//           },
//         }),
//   },
//   shape: {
//     borderRadius: 12,
//   },
//   typography: {
//     fontFamily: [
//       '"Inter"',
//       '"system-ui"',
//       '"-apple-system"',
//       '"BlinkMacSystemFont"',
//       '"Segoe UI"',
//       '"Roboto"',
//       '"Helvetica Neue"',
//       '"Arial"',
//       '"sans-serif"',
//       '"Apple Color Emoji"',
//       '"Segoe UI Emoji"',
//       '"Segoe UI Symbol"',
//     ].join(","),
//     h4: {
//       fontWeight: 700,
//     },
//     button: {
//       textTransform: "none" as const,
//       fontWeight: 600,
//     },
//   },
//   components: {
//     MuiButton: {
//       styleOverrides: {
//         root: {
//           padding: "8px 20px",
//         },
//       },
//     },
//     MuiPaper: {
//       styleOverrides: {
//         root: {
//           backgroundImage: "none",
//           boxShadow: mode === "light" 
//             ? "0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)"
//             : "0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)",
//         },
//       },
//     },
//     MuiAppBar: {
//       styleOverrides: {
//         root: {
//           boxShadow: "none",
//           backgroundColor: mode === "light" ? "rgba(255, 255, 255, 0.9)" : "rgba(18, 18, 18, 0.9)",
//           backdropFilter: "blur(8px)",
//           borderBottom: `1px solid ${mode === "light" ? "rgba(0, 0, 0, 0.12)" : "rgba(255, 255, 255, 0.12)"}`,
//         },
//       },
//     },
//   },
// });

import type { PaletteMode } from "@mui/material";

export const getDesignTokens = (mode: PaletteMode) => ({
  palette: {
    mode,
    ...(mode === "light"
      ? {
          // Parchment field debrief — warm cream base, amber primary
          primary: {
            main: "#D4942A",
            light: "#F2B84B",
            dark: "#A8731A",
            contrastText: "#ffffff",
          },
          secondary: {
            main: "#2AB8C4",
            light: "#4DD4DF",
            dark: "#1D8F99",
            contrastText: "#ffffff",
          },
          divider: "rgba(26, 22, 18, 0.12)",
          background: {
            default: "#F5F0E8",
            paper: "#FDFAF4",
          },
          text: {
            primary: "#1A1612",
            secondary: "rgba(26, 22, 18, 0.6)",
          },
        }
      : {
          // Obsidian covert ops — near-black base, warm amber + cold teal
          primary: {
            main: "#F2B84B",
            light: "#F7CF7E",
            dark: "#D4942A",
            contrastText: "#0E0F14",
          },
          secondary: {
            main: "#4DD4DF",
            light: "#7DE3EB",
            dark: "#2AB8C4",
            contrastText: "#0E0F14",
          },
          divider: "rgba(212, 148, 42, 0.15)",
          background: {
            default: "#0E0F14",
            paper: "#161820",
          },
          text: {
            primary: "#F0EAD6",
            secondary: "rgba(240, 234, 214, 0.65)",
          },
        }),
  },
  shape: {
    borderRadius: 12,
  },
  typography: {
    fontFamily: [
      '"Rajdhani"',      // sharp, intelligence-agency feel for headings
      '"IBM Plex Sans"', // clean, technical for body
      '"system-ui"',
      '"sans-serif"',
    ].join(","),
    h4: { fontWeight: 700, letterSpacing: "0.04em" },
    button: { textTransform: "none" as const, fontWeight: 600, letterSpacing: "0.04em" },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: { padding: "8px 20px" },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundImage: "none",
          boxShadow: mode === "light"
            ? "0 4px 6px -1px rgba(26, 22, 18, 0.08), 0 2px 4px -2px rgba(26, 22, 18, 0.06)"
            : "0 10px 15px -3px rgba(0, 0, 0, 0.4), 0 4px 6px -4px rgba(0, 0, 0, 0.3)",
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          boxShadow: "none",
          backgroundColor: mode === "light"
            ? "rgba(253, 250, 244, 0.9)"
            : "rgba(14, 15, 20, 0.9)",
          backdropFilter: "blur(8px)",
          borderBottom: `1px solid ${
            mode === "light"
              ? "rgba(26, 22, 18, 0.12)"
              : "rgba(212, 148, 42, 0.15)"
          }`,
        },
      },
    },
  },
});
