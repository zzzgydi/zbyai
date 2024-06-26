import { vitePlugin as remix } from "@remix-run/dev";
import { installGlobals } from "@remix-run/node";
import { defineConfig } from "vite";
import svgr from "vite-plugin-svgr";
import tsconfigPaths from "vite-tsconfig-paths";

installGlobals();

export default defineConfig({
  server: {
    proxy: {
      "/v1": {
        target: "http://localhost:14090",
        changeOrigin: true,
        secure: false,
      },
    },
  },
  plugins: [remix(), tsconfigPaths(), svgr()],
});
