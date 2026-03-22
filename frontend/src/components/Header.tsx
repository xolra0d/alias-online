import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Container from "@mui/material/Container";
import Stack from "@mui/material/Stack";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import AppBar from "@mui/material/AppBar";
import { alpha, useTheme } from "@mui/material/styles";
import { Link, useLocation } from "react-router";
import ThemeSwitcher from "./ThemeSwitcher";

export default function Header() {
  const { pathname } = useLocation();
  const theme = useTheme();

  const isLiveRoom = pathname.startsWith("/play");
  const statusLabel = isLiveRoom ? "Live room" : "Setup table";

  return (
    <AppBar
      position="fixed"
      color="transparent"
      elevation={0}
      sx={{
        displayPrint: "none",
        backgroundColor: alpha(
          theme.palette.background.paper,
          theme.palette.mode === "light" ? 0.88 : 0.82,
        ),
        backdropFilter: "blur(20px)",
        borderBottom: 1,
        borderColor: "divider",
      }}
    >
      <Toolbar disableGutters sx={{ minHeight: 80 }}>
        <Container
          maxWidth="xl"
          sx={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            gap: 2,
            px: { xs: 2, sm: 3 },
          }}
        >
          <Box
            component={Link}
            to="/"
            sx={{
              display: "flex",
              alignItems: "center",
              gap: 1.5,
              textDecoration: "none",
              color: "text.primary",
              minWidth: 0,
            }}
          >
            <Box
              aria-hidden
              sx={{
                width: 40,
                height: 40,
                borderRadius: 2,
                background: `linear-gradient(135deg, ${theme.palette.primary.main} 0 52%, ${theme.palette.secondary.main} 52% 100%)`,
                boxShadow: `inset 0 0 0 1px ${alpha(theme.palette.common.white, 0.24)}`,
                flexShrink: 0,
              }}
            />
            <Box sx={{ minWidth: 0 }}>
              <Typography
                variant="h6"
                sx={{
                  lineHeight: 1,
                  fontWeight: 700,
                  whiteSpace: "nowrap",
                }}
              >
                Alias Online
              </Typography>
              <Typography
                variant="caption"
                sx={{
                  display: "block",
                  color: "text.secondary",
                  letterSpacing: "0.18em",
                  textTransform: "uppercase",
                }}
              >
                Word room
              </Typography>
            </Box>
          </Box>

          <Stack direction="row" spacing={1.5} alignItems="center">
            <Box
              sx={{
                display: { xs: "none", md: "flex" },
                alignItems: "center",
                gap: 1,
                px: 1.5,
                py: 0.85,
                borderRadius: 999,
                border: "1px solid",
                borderColor: "divider",
                bgcolor: alpha(theme.palette.background.paper, 0.72),
              }}
            >
              <Box
                sx={{
                  width: 8,
                  height: 8,
                  borderRadius: "50%",
                  bgcolor: isLiveRoom ? "secondary.main" : "primary.main",
                  boxShadow: `0 0 0 4px ${alpha(
                    isLiveRoom
                      ? theme.palette.secondary.main
                      : theme.palette.primary.main,
                    0.12,
                  )}`,
                }}
              />
              <Typography
                variant="caption"
                sx={{ letterSpacing: "0.14em", textTransform: "uppercase" }}
              >
                {statusLabel}
              </Typography>
            </Box>

            <Button
              component={Link}
              to="/"
              variant="text"
              sx={{
                minWidth: 0,
                px: 1.25,
                color: pathname === "/" ? "primary.main" : "text.secondary",
              }}
            >
              Lobby
            </Button>

            <ThemeSwitcher />
          </Stack>
        </Container>
      </Toolbar>
    </AppBar>
  );
}
