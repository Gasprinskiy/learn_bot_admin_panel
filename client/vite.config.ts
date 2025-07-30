import { defineConfig } from 'vite';
import type { AliasOptions } from 'vite';
import vue from '@vitejs/plugin-vue';
import tsconfigPaths from 'vite-tsconfig-paths';
// @ts-expect-error freakint description
import path from 'node:path';

// @ts-expect-error freakint description
const root = path.resolve(__dirname, 'src');

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tsconfigPaths()],
  server: {
    port: 3000,
  },
  resolve: {
    alias: {
      '@': root,
    } as AliasOptions,
  },
});
