import React from "react";
import { ChakraProvider, extendTheme } from "@chakra-ui/react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import ErrorPage from "./routers/error";
import Home, { action as homeAction } from "./routers/home";
import Playground, {
  loader as playgroundLoader,
  action as playgroundAction,
} from "./routers/playground/router";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
    errorElement: <ErrorPage />,
    action: homeAction,
  },
  {
    path: "playground/:gameId/:userId",
    element: <Playground />,
    errorElement: <ErrorPage />,
    loader: playgroundLoader,
    action: playgroundAction,
  },
]);

const theme = extendTheme({
  colors: {
    dark: {
      500: "#03111e",
    },
    blue: {
      500: "#1c6369",
      600: "#17383b",
    },
    white: {
      500: "#dcdbc2",
    },
    orange: {
      500: "#ce7b31",
    },
    red: {
      500: "#d15416",
    },
    green: {
      500: "#42b063",
    },
  },
  styles: {
    global: {
      body: {
        bg: "dark.500",
        color: "white.500",
      },
    },
  },
  components: {
    Button: {
      baseStyle: {
        fontWeight: "bold", // 全てのボタンに太字を適用
      },
      // ボタンのデフォルトスタイルを定義
      defaultProps: {
        colorScheme: "blue", // デフォルトのカラーシステムをbrandに設定
      },
      variants: {
        solid: {
          bg: "blue.500", // デフォルトの背景色
          color: "white", // デフォルトのフォントカラー
          _hover: {
            bg: "blue.600", // ホバー時の背景色
          },
        },
      },
    },
  },
});

// biome-ignore lint/style/noNonNullAssertion: <explanation>
ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ChakraProvider theme={theme}>
      <RouterProvider router={router} />
    </ChakraProvider>
  </React.StrictMode>,
);
