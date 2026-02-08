import { purgeCss } from 'vite-plugin-tailwind-purgecss';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

// Get the Aio base URL from environment variable with fallback
const AIO_BASE_URL = process.env.AIO_BASE_URL || 'http://localhost:8080';

export default defineConfig({
  plugins: [sveltekit(), purgeCss()],
  optimizeDeps: {
    include: ['pdfjs-dist'],
    esbuildOptions: {
      target: 'esnext',
      supported: {
        'top-level-await': true
      }
    }
  },
  define: {
    'process.env': {
      NODE_ENV: JSON.stringify(process.env.NODE_ENV)
    },
    'process.platform': JSON.stringify(process.platform),
    'process.cwd': JSON.stringify('/'),
    'process.browser': true,
    'process': {
      cwd: () => ('/')
    },
    // Inject Aio configuration for client-side access
    '__AIO_CONFIG__': {
      AIO_BASE_URL: JSON.stringify(AIO_BASE_URL)
    }
  },
  resolve: {
    alias: {
      process: 'process/browser'
    }
  },
  server: {
    fs: {
      allow: ['..']  // allows importing from the parent directory
    },
    proxy: {
      '/api': {
        target: AIO_BASE_URL,
        changeOrigin: true,
        timeout: 900000,
        rewrite: (path) => path.replace(/^\/api/, ''),
        configure: (proxy, _options) => {
          proxy.on('error', (err, req, res) => {
            console.log('proxy error', err);
            res.writeHead(500, {
              'Content-Type': 'text/plain',
            });
            res.end('Something went wrong. The backend server may not be running.');
          });
        }
      },
      '^/(patterns|models|sessions)/names': {
        target: AIO_BASE_URL,
        changeOrigin: true,
        timeout: 900000,
        configure: (proxy, _options) => {
          proxy.on('error', (err, req, res) => {
            console.log('proxy error', err);
            res.writeHead(500, {
              'Content-Type': 'application/json',
            });
            res.end(JSON.stringify({ error: 'Backend server not running', names: [] }));
          });
        }
      }
    },
    watch: {
      usePolling: true,
      interval: 100,
      ignored: ['**/node_modules/**', '**/dist/**', '**/.git/**', '**/.svelte-kit/**']
    }
  },
  build: {
    commonjsOptions: {
      transformMixedEsModules: true
    },
    target: 'esnext',
    minify: true,
    rollupOptions: {
      output: {
        format: 'es'
      }
    }
  }
});
