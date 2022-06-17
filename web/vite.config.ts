import {defineConfig} from 'vite';
import react from '@vitejs/plugin-react';
import eslintPlugin from 'vite-plugin-eslint';

// https://vitejs.dev/config/
export default defineConfig({
    build: {
        outDir: 'build/dist',
    },
    plugins: [
        react(),
        eslintPlugin({throwOnWarning: true, throwOnError: true}),
    ],
});
