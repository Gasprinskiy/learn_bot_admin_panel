import { defineConfig } from 'vite';
import type { AliasOptions } from 'vite';
import vue from '@vitejs/plugin-vue';
import tsconfigPaths from 'vite-tsconfig-paths';
// import basicSsl from '@vitejs/plugin-basic-ssl';
// @ts-expect-error freakint description
import path from 'node:path';

// import fs from 'node:fs';

// @ts-expect-error freakint description
const root = path.resolve(__dirname, 'src');

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    tsconfigPaths(),
  ],
  server: {
    host: 'admin-panel.local',
    port: 3000,
    // https: {
    //   key: fs.readFileSync('./.cert/key.pem'),
    //   cert: fs.readFileSync('./.cert/cert.pem'),
    // },
  },
  resolve: {
    alias: {
      '@': root,
    } as AliasOptions,
  },
});
