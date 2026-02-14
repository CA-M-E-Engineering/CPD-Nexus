import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // Load env file from root directory
  const env = loadEnv(mode, path.resolve(__dirname, '..'), '')

  const apiPort = env.API_PORT || '3000'
  const frontendPort = parseInt(env.FRONTEND_PORT || '5173')

  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src')
      }
    },
    server: {
      port: frontendPort,
      proxy: {
        '/api': {
          target: `http://localhost:${apiPort}`,
          changeOrigin: true,
        }
      }
    }
  }
})
