import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',
    port: 5173,
    allowedHosts: true,
    proxy: {
      '/api': {
        target: 'https://i2u6hsbhf1.execute-api.us-east-2.amazonaws.com',
        changeOrigin: true,
      },
      '/yolo': {
        target: 'https://i2u6hsbhf1.execute-api.us-east-2.amazonaws.com',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/yolo/, ''),
      },
      '/fincas-api': {
        target: 'http://3.23.193.162:8082',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/fincas-api/, ''),
      },
    },
  },
})