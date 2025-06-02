import { defineConfig } from 'vite'

export default defineConfig({
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',  // Your Go backend
        changeOrigin: true,
        ws: true  // Required for WebSocket proxy
      }
    }
  }
})