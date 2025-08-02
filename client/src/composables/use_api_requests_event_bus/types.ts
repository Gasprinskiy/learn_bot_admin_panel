export interface UseApiRequestEventBusEvents {
  on_request: { is_blocking?: boolean };
  on_response: null;
  on_error: { message: string | undefined } | null;
}
