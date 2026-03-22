import { Fragment } from "react";
import DarkModeIcon from "@mui/icons-material/DarkMode";
import LightModeIcon from "@mui/icons-material/LightMode";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import useColorTheme from "../theme/use-color-theme";

export default function ThemeSwitcher() {
  const { mode, toggleColorMode } = useColorTheme();

  return (
    <Tooltip
      title={`${mode === "dark" ? "Light" : "Dark"} mode`}
      enterDelay={1000}
    >
      <IconButton
        size="small"
        aria-label={`Switch to ${mode === "dark" ? "light" : "dark"} mode`}
        onClick={toggleColorMode}
        sx={{
          width: 40,
          height: 40,
          border: "1px solid",
          borderColor: "divider",
          backgroundColor: "background.paper",
        }}
      >
        <Fragment>
          {mode === "dark" ? <LightModeIcon /> : <DarkModeIcon />}
        </Fragment>
      </IconButton>
    </Tooltip>
  );
}
