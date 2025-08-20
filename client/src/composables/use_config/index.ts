import type { Config } from './types';

export function useConfig(): Config {
  const apiUrl = import.meta.env.VITE_API_URL || '';
  const sseTtl = +import.meta.env.VITE_SSE_TTL || 0;
  const uploadsURL = import.meta.env.VITE_UPLOADS_URL || '';

  return {
    ApiURL: apiUrl,
    SSETTL: sseTtl,
    UploadsURL: uploadsURL,
  };
}
