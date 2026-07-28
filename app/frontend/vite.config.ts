import { defineConfig } from 'vite'
import { devtools } from '@tanstack/devtools-vite'
import path from 'node:path';
import { tanstackRouter } from '@tanstack/router-plugin/vite'
import viteReact from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

const config = defineConfig({
  resolve: {
    tsconfigPaths: true, alias: {
      '@wails': path.resolve(__dirname, './wailsjs'),
    }
  },
  plugins: [
    devtools(),
    tailwindcss(),
    tanstackRouter({ target: 'react', autoCodeSplitting: true }),
    viteReact(),
  ],
})

export default config
