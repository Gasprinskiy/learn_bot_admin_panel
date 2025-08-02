import { useApiRequestEventBus } from '@/composables/use_api_requests_event_bus';
import { useConfig } from '@/composables/use_config';
import { generateDeviceID } from '@/packages/device_info';
import { ofetch } from 'ofetch';

const apiRequestEventBus = useApiRequestEventBus();
const { ApiURL } = useConfig();

const $api = ofetch.create({
  baseURL: ApiURL,
  credentials: 'include',
  cache: 'no-store',
  headers: {
    'Accept': 'application/json',
    'Device-ID': generateDeviceID(),
  },
  onRequest: (request) => {
    const { options } = request;
    apiRequestEventBus.dispatch('on_request', { is_blocking: (options.body as any)?._is_blocking });
  },
  onResponse: () => apiRequestEventBus.dispatch('on_response', null),
  onRequestError: () => apiRequestEventBus.dispatch('on_error', null),
  onResponseError: (response) => {
    apiRequestEventBus.dispatch('on_error', { message: response.error?.message.toString() });
  },
});

export default $api;
