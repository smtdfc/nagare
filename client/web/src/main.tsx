import { createRoot } from "react-dom/client";
import { App } from "@nagare-app/ui";
import { refresh } from "@nagare-app/services";
import "./styles.css";
import "./bindings";

window.addEventListener("load", async () => {
  await window.bindings.loadConfig();
  refresh();
  createRoot(document.getElementById("root")!).render(<App />);
});
