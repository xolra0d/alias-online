import { styled } from "@mui/material/styles";
import MuiAppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import Stack from "@mui/material/Stack";
import ThemeSwitcher from "./ThemeSwitcher";
import ToggleButton from "@mui/material/ToggleButton";
import ToggleButtonGroup from "@mui/material/ToggleButtonGroup";
import { useLocation, Link } from "react-router";

const AppBar = styled(MuiAppBar)(({ theme }) => ({
  borderWidth: 0,
  borderBottomWidth: 1,
  borderStyle: "solid",
  borderColor: theme.palette.divider,
  boxShadow: "none",
  zIndex: theme.zIndex.drawer + 1,
}));

export default function Header() {
  const { pathname } = useLocation();

  let alignment = "left";
  if (pathname.startsWith("/about")) {
    alignment = "right";
  }

  return (
    <AppBar
      color="inherit"
      variant="outlined"
      position="static"
      sx={{ displayPrint: "none" }}
    >
      <Toolbar
        sx={{
          maxWidth: 720,
          width: "100%",
          mx: "auto",
          px: { xs: 2, sm: 0 },
        }}
      >
        <Stack
          direction="row"
          justifyContent="space-between"
          alignItems="center"
          width="100%"
          flexWrap="wrap"
        >
          <ToggleButtonGroup
            value={alignment}
            exclusive
            aria-label="text alignment"
          >
            <ToggleButton
              value="left"
              aria-label="left aligned"
              component={Link}
              to="/"
            >
              Play
            </ToggleButton>
            <ToggleButton
              value="right"
              aria-label="right aligned"
              component={Link}
              to="/about"
            >
              About
            </ToggleButton>
          </ToggleButtonGroup>
          <ThemeSwitcher />
        </Stack>
      </Toolbar>
    </AppBar>
  );
}
