import { styled, useTheme } from "@mui/material/styles";
import MuiAppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";
import Stack from "@mui/material/Stack";
import ThemeSwitcher from "./ThemeSwitcher";
import { useLocation, Link } from "react-router";

const AppBar = styled(MuiAppBar)(({ theme }) => ({
  border: "none",
  backgroundColor: theme.palette.mode === "light" 
    ? "rgba(255, 255, 255, 0.72)" 
    : "rgba(15, 23, 42, 0.72)",
  backdropFilter: "blur(12px)",
  position: "fixed",
  top: 0,
  width: "100%",
}));

export default function Header() {
  const { pathname } = useLocation();
  const theme = useTheme();

  const isPlayActive = pathname === "/" || pathname.startsWith("/play");
  const isAboutActive = pathname === "/about";

  return (
    <AppBar
      color="inherit"
      variant="outlined"
      sx={{ displayPrint: "none" }}
    >
      <Toolbar
        sx={{
          maxWidth: 960,
          width: "100%",
          mx: "auto",
          px: { xs: 2, sm: 3 },
          display: "flex",
          justifyContent: "space-between",
        }}
      >
        <Stack direction="row" spacing={2} alignItems="center">
          <Typography 
            variant="h6" 
            component={Link} 
            to="/" 
            sx={{ 
              textDecoration: "none", 
              color: theme.palette.primary.main,
              fontWeight: 800,
              fontSize: "1.25rem",
              mr: 2
            }}
          >
            Alias Online
          </Typography>
          <Box component="nav" sx={{ display: "flex", gap: 1 }}>
            <Button
              component={Link}
              to="/"
              sx={{
                borderRadius: 2,
                px: 2,
                color: isPlayActive ? "primary.main" : "text.secondary",
                backgroundColor: isPlayActive ? theme.palette.primary.light + "1a" : "transparent",
                fontWeight: isPlayActive ? 700 : 500,
                "&:hover": {
                  backgroundColor: isPlayActive 
                    ? theme.palette.primary.light + "26" 
                    : "action.hover",
                },
              }}
            >
              Play
            </Button>
            <Button
              component={Link}
              to="/about"
              sx={{
                borderRadius: 2,
                px: 2,
                color: isAboutActive ? "primary.main" : "text.secondary",
                backgroundColor: isAboutActive ? theme.palette.primary.light + "1a" : "transparent",
                fontWeight: isAboutActive ? 700 : 500,
                "&:hover": {
                  backgroundColor: isAboutActive 
                    ? theme.palette.primary.light + "26" 
                    : "action.hover",
                },
              }}
            >
              About
            </Button>
          </Box>
        </Stack>
        <ThemeSwitcher />
      </Toolbar>
    </AppBar>
  );
}
