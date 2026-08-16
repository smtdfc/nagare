import ReactDOM from "react-dom/client";
import { createRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";
import "./global.css";
import App from "./App";
import {
  IsServerRunning,
  ShowErrorDialog,
} from "@nagare-agent/service-bindings";

const router = createRouter({
  routeTree,
  defaultPreload: "intent",
  scrollRestoration: true,
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

export async function idleTask() {
  if (!(await IsServerRunning()))
    ShowErrorDialog(
      "Error",
      "Failed to connect to server. Please check your connection and try again.",
    );
}

export function startApp() {
  idleTask().catch((e) => {
    console.log(e);
    return;
  });

  const rootElement = document.getElementById("app")!;

  if (!rootElement.innerHTML) {
    const root = ReactDOM.createRoot(rootElement);
    root.render(
      <>
        <App router={router} />
      </>,
    );
  }
}
