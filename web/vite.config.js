import { defineConfig } from "vite";
import { resolve } from "path";
import { fileURLToPath } from "url";

const __dirname = fileURLToPath(new URL(".", import.meta.url));

export default defineConfig({
  root: ".",
  publicDir: "images", // Copia a pasta images/ para o build
  build: {
    outDir: resolve(__dirname, "../public"),
    emptyOutDir: true,
    rollupOptions: {
      input: resolve(__dirname, "index.html"),
    },
    // Otimizações de produção
    minify: "esbuild",
    cssMinify: true,
    sourcemap: false,
    // Code splitting e tree shaking automáticos
    chunkSizeWarningLimit: 1000,
  },
  // Configuração para desenvolvimento (watch mode)
  server: {
    watch: {
      usePolling: true, // Necessário para Docker
    },
  },
});
