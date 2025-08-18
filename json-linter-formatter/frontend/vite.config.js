import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
    plugins: [
        vue({
            // Tell the Vue plugin to compile to a custom element
            customElement: true,
        }),
    ],
    build: {
        // We are not building a library, but a single JS file to be executed.
        // The output is a single JS bundle.
        outDir: 'dist',
        assetsDir: '',
        rollupOptions: {
            input: 'src/main.js', // Entry point is our new main.js
            output: {
                entryFileNames: 'component.js', // Output a single, predictably named file
            },
        },
        emptyOutDir: true,
    },
})
