import Box from "@mui/material/Box";
import Container from "@mui/material/Container";
import Link from "@mui/material/Link";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { alpha, useTheme } from "@mui/material/styles";

export default function Footer() {
  const theme = useTheme();

  return (
    <Box
      component="footer"
      sx={{
        mt: "auto",
        borderTop: "1px solid",
        borderColor: "divider",
        backgroundColor: alpha(
          theme.palette.background.paper,
          theme.palette.mode === "light" ? 0.72 : 0.54,
        ),
        backdropFilter: "blur(12px)",
      }}
    >
      <Container maxWidth="xl" sx={{ py: 3, px: { xs: 2, sm: 3 } }}>
        <Stack
          direction={{ xs: "column", sm: "row" }}
          spacing={2}
          justifyContent="space-between"
          alignItems={{ xs: "flex-start", sm: "center" }}
        >
          <Box>
            <Typography variant="body2" sx={{ fontWeight: 600 }}>
              Alias Online
            </Typography>
            <Typography variant="body2" color="text.secondary">
              A word game built for lively rooms, sharp clues, and quick rounds.
            </Typography>
          </Box>

          <Stack direction="row" spacing={2} alignItems="center">
            <Typography variant="caption" color="text.secondary">
              © {new Date().getFullYear()}
            </Typography>
            <Link
              href="mailto:a@xolra0d.com"
              underline="hover"
              sx={{ fontWeight: 600, letterSpacing: "0.08em" }}
            >
              Contact support
            </Link>
          </Stack>
        </Stack>
      </Container>
    </Box>
  );
}
