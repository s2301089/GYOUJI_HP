import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import react from '@vitejs/plugin-react';

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), vue()],
  css: {
    postcss: './postcss.config.js'
  },
  optimizeDeps: {
    include: ['tailwindcss']
  }
});
