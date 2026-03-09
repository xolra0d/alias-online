import Box from "@mui/material/Box";
import Container from "@mui/material/Container";
import Paper from "@mui/material/Paper";
import Link from "@mui/material/Link";

export default function Footer() {
  return (
    <Paper
      sx={{
        bottom: 0,
      }}
      component="footer"
      variant="outlined"
    >
      <Container maxWidth="lg">
        <Box
          sx={{
            flexGrow: 1,
            justifyContent: "center",
            display: "flex",
            mb: 2,
          }}
        >
          <Link
            href="mailto:a@xolra0d.com"
            sx={{
              textDecoration: "none",
              "&:hover": {
                textDecoration: "underline",
              },
            }}
          >
            Email me
          </Link>
        </Box>
      </Container>
    </Paper>
  );
}
