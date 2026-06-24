import { defineConfig } from 'vite';

export default defineConfig({
    base: './',
    build: {
        outDir: 'dist',
        assetsDir: '',
        emptyOutDir: true,
    },
    server: {
        watch: {
            ignored: ['**/wailsjs/**'],
        },
    },
});
