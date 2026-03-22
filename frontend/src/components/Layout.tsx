import Box from "@mui/material/Box";
import Container from "@mui/material/Container";
import Toolbar from "@mui/material/Toolbar";
import { Outlet } from "react-router";
import Header from "./Header";
import Footer from "./Footer";

export default function Layout() {
  return (
    <Box
      sx={{
        minHeight: "100vh",
        display: "flex",
        flexDirection: "column",
        position: "relative",
        overflow: "clip",
      }}
    >
      <Box
        aria-hidden
        sx={{
          position: "fixed",
          inset: 0,
          pointerEvents: "none",
          overflow: "hidden",
          zIndex: 0,
        }}
      >
        <Box
          sx={{
            position: "absolute",
            top: -120,
            left: -120,
            width: 360,
            height: 360,
            borderRadius: "50%",
            bgcolor: "primary.main",
            opacity: 0.12,
            filter: "blur(72px)",
          }}
        />
        <Box
          sx={{
            position: "absolute",
            top: 140,
            right: -160,
            width: 420,
            height: 420,
            borderRadius: "50%",
            bgcolor: "secondary.main",
            opacity: 0.1,
            filter: "blur(92px)",
          }}
        />
        <Box
          sx={{
            position: "absolute",
            bottom: -180,
            left: "32%",
            width: 520,
            height: 520,
            borderRadius: "50%",
            bgcolor: "warning.main",
            opacity: 0.07,
            filter: "blur(100px)",
          }}
        />
      </Box>

      <Header />
      <Toolbar />

      <Container
        component="main"
        maxWidth="xl"
        sx={{
          position: "relative",
          zIndex: 1,
          flex: 1,
          width: "100%",
          py: { xs: 3, md: 5 },
        }}
      >
        <Outlet />
      </Container>

      <Footer />
    </Box>
  );
}
