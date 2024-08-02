import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    proxy: process.env.NODE_ENV === "production"? undefined: {
      "/rpc": {
        target: `http://localhost${process.env.PORT}`,
        changeOrigin: true,
      },
    }
  },
  plugins: [react()],
})
