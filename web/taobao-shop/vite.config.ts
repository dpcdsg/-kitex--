import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

const root = path.dirname(fileURLToPath(import.meta.url));

// 开发时把 /seckill 代理到本机 API，避免浏览器跨域；后端默认监听 10001
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: { '@': path.join(root, 'src') },
  },
  server: {
    port: 5173,
    proxy: {
      '/seckill': {
        target: 'http://127.0.0.1:10001',
        changeOrigin: true,
      },
    },
  },
});
