import { defineConfig } from 'vite';

export default defineConfig({
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
