import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import ErrorPage from "./routers/error";
import Home, { action as homeAction } from "./routers/home";
import Playground, {
  loader as playgroundLoader,
  action as playgroundAction,
} from "./routers/playground/router";

import "./index.css";

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

// biome-ignore lint/style/noNonNullAssertion: <explanation>
ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
