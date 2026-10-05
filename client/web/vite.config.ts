import { defineConfig } from "vite";
import path from "node:path";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import viteReact from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

const webPort = Number.parseInt(process.env.NAGARE_WEB_UI_PORT ?? "3005", 10);

const config = defineConfig({
  server: {
    host: "127.0.0.1",
    port: webPort,
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
      routesDirectory: path.resolve(import.meta.dirname, "../ui/src/routes"),
      generatedRouteTree: path.resolve(
        import.meta.dirname,
        "../ui/src/routeTree.gen.ts",
      ),
    }),
    viteReact(),
    tailwindcss(),
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

export default config;
