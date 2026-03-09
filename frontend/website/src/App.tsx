import { RouterProvider, createBrowserRouter } from "react-router";
import { ThemeProvider } from "@mui/material/styles";
import { CssBaseline } from "@mui/material";
import useColorTheme from "./theme/use-color-theme";
import Layout from "./components/Layout";
import Play from "./components/Play";
import About from "./components/About";

const router = createBrowserRouter([
  {
    Component: Layout,
    children: [
      {
        path: "/",
        Component: Play,
      },
      {
        path: "/about",
        Component: About,
      },
      {
        path: "*",
        Component: Play,
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
