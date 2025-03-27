import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  base: "/static/",
  plugins: [react()],
  build: {
    outDir: "dist",  // Output folder
    sourcemap: true, // Generate source maps
    rollupOptions: {
      input: "game.html",  // Entry file
      output: {
        entryFileNames: "bundle.js",  // Output filename
        format: "esm",
      },
    },
  },
})
