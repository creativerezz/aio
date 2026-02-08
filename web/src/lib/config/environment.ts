/**
 * Environment configuration for the Aio web app
 * Centralizes all environment variable handling
 */

// Default values
const DEFAULT_AIO_BASE_URL = 'http://localhost:8080';

/**
 * Get the Aio base URL from environment variable or default
 * This function works in both server and client contexts
 */
export function getAioBaseUrl(): string {
  // In server context (Node.js), use process.env
  if (typeof process !== 'undefined' && process.env) {
    return process.env.AIO_BASE_URL || DEFAULT_AIO_BASE_URL;
  }

  // In client context, check if the environment was injected via Vite
  if (typeof window !== 'undefined' && (window as any).__AIO_CONFIG__) {
    return (window as any).__AIO_CONFIG__.AIO_BASE_URL || DEFAULT_AIO_BASE_URL;
  }

  // Fallback to default
  return DEFAULT_AIO_BASE_URL;
}

/**
 * Get the Aio API base URL (adds /api if not present)
 */
export function getAioApiUrl(): string {
  const baseUrl = getAioBaseUrl();

  // Remove trailing slash if present
  const cleanBaseUrl = baseUrl.replace(/\/$/, '');

  // Check if it already ends with /api
  if (cleanBaseUrl.endsWith('/api')) {
    return cleanBaseUrl;
  }

  return `${cleanBaseUrl}/api`;
}

/**
 * Configuration object for easy access to all environment settings
 */
export const config = {
  aioBaseUrl: getAioBaseUrl(),
  aioApiUrl: getAioApiUrl(),
} as const;

// Type definitions
export interface AioConfig {
  AIO_BASE_URL: string;
}

declare global {
  interface Window {
    __AIO_CONFIG__?: AioConfig;
  }
}
