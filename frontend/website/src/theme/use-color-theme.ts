import { createTheme } from "@mui/material";
import { useMemo } from "react";
import { getDesignTokens } from "./theme";
import { useSelector, useDispatch } from "react-redux";
import type { RootState } from "../store";
import { toggleTheme } from "./themeSlice";

export default function useColorTheme() {
  const mode = useSelector((state: RootState) => state.theme.mode);
  const dispatch = useDispatch();

  const toggleColorMode = () => {
    dispatch(toggleTheme());
  };

  const modifiedTheme = useMemo(
    () => createTheme(getDesignTokens(mode)),
    [mode],
  );

  return {
    theme: modifiedTheme,
    mode,
    toggleColorMode,
  };
}
