import { type AliasOptions, defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tsconfigPaths from "vite-tsconfig-paths";
//@ts-ignore
import path from "path";

//@ts-ignore
const root = path.resolve(__dirname, "src");


// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tsconfigPaths()],
  server: {
    port: 5173
  },
  resolve: {
    alias: {
      "@": root,
    } as AliasOptions,
  },
})