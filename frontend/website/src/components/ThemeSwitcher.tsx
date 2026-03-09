import { Fragment } from "react";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import DarkModeIcon from "@mui/icons-material/DarkMode";
import LightModeIcon from "@mui/icons-material/LightMode";
import useColorTheme from "../theme/use-color-theme";

export default function ThemeSwitcher() {
  const { mode, toggleColorMode } = useColorTheme();

  return (
    <Tooltip
      title={`${mode === "dark" ? "Light" : "Dark"} mode`}
      enterDelay={1000}
    >
      <div>
        <IconButton
          size="small"
          aria-label={`Switch to ${mode === "dark" ? "light" : "dark"} mode`}
          onClick={toggleColorMode}
        >
          <Fragment>
            {mode === "dark" ? <LightModeIcon /> : <DarkModeIcon />}
          </Fragment>
        </IconButton>
      </div>
    </Tooltip>
  );
}
