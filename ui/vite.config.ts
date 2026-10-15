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
      '/secret':    'http://localhost:2710',
      '/token':     'http://localhost:2710',
      '/namespace': 'http://localhost:2710',
      '/unseal':    'http://localhost:2710',
      '/status':    'http://localhost:2710',
      '/seal':      'http://localhost:2710',
    },
  },
})
