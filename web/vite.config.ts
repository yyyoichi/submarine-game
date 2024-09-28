import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    proxy: {
      "/rpc": {
        target: `http://localhost:${process.env.PORT}`,
        changeOrigin: true,
      },
    }
  },
  plugins: [react()],
})
