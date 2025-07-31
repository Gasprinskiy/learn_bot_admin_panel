import { useApiRequestEventBus } from '@/composables/use_api_requests_event_bus';
import { useConfig } from '@/composables/use_config';
import { ofetch } from 'ofetch';

const apiRequestEventBus = useApiRequestEventBus();
const { ApiURL } = useConfig();

const $api = ofetch.create({
  baseURL: ApiURL,
  credentials: 'include',
  headers: {
    Accept: 'application/json',
  },
  onRequest: () => apiRequestEventBus.dispatch('on_request', null),
  onResponse: () => apiRequestEventBus.dispatch('on_response', null),
  onRequestError: () => apiRequestEventBus.dispatch('on_error', null),
  onResponseError: (response) => {
    apiRequestEventBus.dispatch('on_error', { message: response.error?.message.toString() });
  },
});

export default $api;
