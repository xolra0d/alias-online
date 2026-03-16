import Stack from "@mui/material/Stack";
import Box from "@mui/material/Box";
import { Outlet } from "react-router";
import Header from "./Header";
import Footer from "./Footer";
import Toolbar from "@mui/material/Toolbar";

export default function Layout() {
  return (
    <Box sx={{ minHeight: "100vh", display: "flex", flexDirection: "column" }}>
      <Header />
      <Toolbar sx={{ mb: 2 }} />
      <Stack
        component="main"
        maxWidth={960}
        mx="auto"
        px={{ xs: 2, sm: 3 }}
        sx={{ flex: 1, width: "100%", py: 4 }}
      >
        <Outlet />
      </Stack>
      <Footer />
    </Box>
  );
}
