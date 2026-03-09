import type { PaletteMode } from "@mui/material";
import { amber, deepOrange, grey, brown } from "@mui/material/colors";

export const getDesignTokens = (mode: PaletteMode) => ({
  palette: {
    mode,
    ...(mode === "light"
      ? {
          // palette values for light mode
          primary: {
            main: brown[700],
          },
          secondary: {
            main: amber[900],
          },
          divider: brown[100],
          background: {
            default: "#fffcf5",
            paper: "#fffcf5",
          },
          text: {
            primary: brown[900],
            secondary: brown[700],
          },
        }
      : {
          // palette values for dark mode
          primary: {
            main: amber[200],
          },
          secondary: {
            main: deepOrange[300],
          },
          divider: brown[800],
          background: {
            default: "#1b120f",
            paper: "#261b18",
          },
          text: {
            primary: "#f5ebe0",
            secondary: grey[400],
          },
        }),
  },
});
