export interface UseApiRequestEventBusEvents {
  on_request: null;
  on_response: null;
  on_error: { message: string | undefined } | null;
}
