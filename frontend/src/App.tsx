import { RouterProvider, createBrowserRouter } from "react-router";
import { ThemeProvider } from "@mui/material/styles";
import { CssBaseline } from "@mui/material";
import useColorTheme from "./theme/use-color-theme";
import Layout from "./components/Layout";
import Home from "./components/Home.tsx";
import Play from "./components/Play.tsx";

const router = createBrowserRouter([
  {
    Component: Layout,
    children: [
      {
        path: "/",
        Component: Home,
      },
      {
        path: "/play/:room_id",
        Component: Play,
      },
      {
        path: "*",
        Component: Home,
      },
    ],
  },
]);

function App() {
  const { theme } = useColorTheme();
  return (
    <>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <RouterProvider router={router} />
      </ThemeProvider>
    </>
  );
}

export default App;
