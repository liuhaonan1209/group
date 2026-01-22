import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  server: {
    port: 3000,
    proxy: {
      '/api/route': {
        target: 'http://localhost:8890',
        changeOrigin: true
      },
      '/api/schedule': {
        target: 'http://localhost:8890',
        changeOrigin: true
      },
      '/api/finance': {
        target: 'http://localhost:8895',
        changeOrigin: true
      }
    }
  }
})