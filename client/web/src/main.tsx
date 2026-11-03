import { createRoot } from "react-dom/client";
import { App } from "@nagare-app/ui";
import { applyBootstrapFromURL, refresh } from "@nagare-app/services";
import "./styles.css";
import "./bindings";

window.addEventListener("load", async () => {
  await window.bindings.loadConfig();
  // The launcher may pass the gateway URL and token in the page URL so the
  // dashboard signs in without asking for them.
  await applyBootstrapFromURL();
  refresh();
  createRoot(document.getElementById("root")!).render(<App />);
});
