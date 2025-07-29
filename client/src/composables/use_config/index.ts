import type { Config } from './types';

export function useConfig(): Config {
  const apiUrl = import.meta.env.VITE_API_URL || '';
  const sseTtl = +import.meta.env.VITE_SSE_TTL || 0;

  return {
    ApiURL: apiUrl,
    SSETTL: sseTtl,
  };
}
