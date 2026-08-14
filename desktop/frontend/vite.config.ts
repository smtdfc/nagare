import { defineConfig } from "vite";
import path from "node:path";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import viteReact from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

const config = defineConfig({
  server: {
    host: "127.0.0.1",
    port: 9245,
    strictPort: true,
  },
  resolve: {
    tsconfigPaths: true,
    alias: {
      "@wails": path.resolve(import.meta.dirname, "./bindings"),
    },
  },
  plugins: [
    tailwindcss(),
    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
      routesDirectory: path.resolve(
        import.meta.dirname,
        "../../ui/main/src/routes",
      ),
    }),
    viteReact(),
  ],
});

export default config;
