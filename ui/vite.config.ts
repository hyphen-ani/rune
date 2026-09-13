import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  base: '/ui/',
  build: {
    outDir: '../internal/ui/static',
    emptyOutDir: true,
  },
  resolve: {
    alias: { '@': `${import.meta.dirname}/src` },
  },
  server: {
    proxy: {
      '/secret':    'http://localhost:8080',
      '/token':     'http://localhost:8080',
      '/namespace': 'http://localhost:8080',
      '/unseal':    'http://localhost:8080',
      '/status':    'http://localhost:8080',
      '/seal':      'http://localhost:8080',
    },
  },
})
