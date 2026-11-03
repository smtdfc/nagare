import { createRoot } from "react-dom/client";
import { App } from "@nagare-app/ui";
import { applyBootstrapFromURL, refresh } from "@nagare-app/services";
import "./styles.css";
import "./bindings";

await window.bindings.loadConfig();
// The desktop shell may pass the gateway URL and token in the window URL so it
// opens signed in, mirroring `nagare web`.
await applyBootstrapFromURL();
refresh();

createRoot(document.getElementById("root")!).render(<App />);
