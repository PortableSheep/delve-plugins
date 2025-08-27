import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { resolve } from "path";

export default defineConfig({
    plugins: [
        vue({
            customElement: true,
        }),
    ],
    build: {
        outDir: "dist",
        emptyOutDir: true,
        lib: {
            entry: resolve(__dirname, "src/main.js"),
            name: "GitHubDashboard",
            fileName: "component",
            formats: ["iife"],
        },
        rollupOptions: {
            output: {
                format: "iife",
                inlineDynamicImports: true,
                manualChunks: undefined,
                assetFileNames: "assets/[name].[ext]",
                chunkFileNames: "[name].js",
                entryFileNames: "component.js",
                exports: "named",
            },
        },
        minify: "terser",
        terserOptions: {
            compress: {
                drop_console: false,
                drop_debugger: true,
            },
            mangle: {
                keep_classnames: true,
                keep_fnames: true,
            },
        },
        target: "es2015",
        cssCodeSplit: false,
    },
    define: {
        __VUE_OPTIONS_API__: true,
        __VUE_PROD_DEVTOOLS__: false,
        __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: false,
    },
    resolve: {
        alias: {
            "@": resolve(__dirname, "src"),
        },
    },
    css: {
        postcss: {
            plugins: [],
        },
    },
    base: "./",
    server: {
        port: 3000,
        cors: true,
        host: true,
        open: false,
    },
    preview: {
        port: 3001,
        cors: true,
        host: true,
    },
});
