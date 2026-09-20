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

    },
  },
  plugins: [

    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
      routesDirectory: path.resolve(
          import.meta.dirname,
          "../ui/src/routes",
      ),
      generatedRouteTree: path.resolve(
          import.meta.dirname,
          "../ui/src/routeTree.gen.ts",
      ),
    }),
    viteReact(),
    tailwindcss(),
  ],
});

export default config;