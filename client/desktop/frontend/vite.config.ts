import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import wails from "@wailsio/runtime/plugins/vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import tailwindcss from "@tailwindcss/vite";
import path from "node:path";

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
  resolve: {
    tsconfigPaths: true,
    alias: {},
  },
  plugins: [
    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
      routesDirectory: path.resolve(import.meta.dirname, "../../ui/src/routes"),
      generatedRouteTree: path.resolve(
        import.meta.dirname,
        "../../ui/src/routeTree.gen.ts",
      ),
    }),
    react(),
    tailwindcss(),
    wails("./bindings"),
  ],
  build: {
    rollupOptions: {
      output: {
        chunkFileNames: "assets/chunk-[hash].js",
        entryFileNames: "assets/[hash].js",
        assetFileNames: "assets/[hash].[ext]",
      },
    },
  },
});
