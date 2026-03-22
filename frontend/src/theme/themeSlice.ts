import { createSlice } from "@reduxjs/toolkit";
import type { PaletteMode } from "@mui/material";

interface ThemeState {
  mode: PaletteMode;
}

function getInitialMode(): PaletteMode {
  const savedMode = localStorage.getItem("themeMode");
  if (savedMode === "dark" || savedMode === "light") {
    return savedMode;
  }
  localStorage.setItem("themeMode", "light");
  return "light";
}

const initialState: ThemeState = {
  mode: getInitialMode(),
};

export const themeSlice = createSlice({
  name: "theme",
  initialState,
  reducers: {
    toggleTheme: (state) => {
      state.mode = state.mode === "light" ? "dark" : "light";
      localStorage.setItem("themeMode", state.mode);
    },
  },
});

export const { toggleTheme } = themeSlice.actions;

export default themeSlice.reducer;
