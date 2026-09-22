import { createRoot } from "react-dom/client";
import { App } from "@nagare-app/ui";
import "./styles.css";
import "./bindings";

createRoot(document.getElementById("root")!).render(<App />);
