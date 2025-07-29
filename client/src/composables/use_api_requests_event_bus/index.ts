import { EventBus } from '@/packages/event_bus';
import type { UseApiRequestEventBusEvents } from './types';
import type { UnsubscribeCallBack } from '@/packages/event_bus/types';

const eventBus = new EventBus();

export function useApiRequestEventBus() {
  function subscribe<K extends keyof UseApiRequestEventBusEvents>(key: K, callBack: (args: UseApiRequestEventBusEvents[K]) => void): UnsubscribeCallBack {
    return eventBus.subscribe({
      key,
      callBack,
    });
  }

  function dispatch<K extends keyof UseApiRequestEventBusEvents>(key: K, arg: UseApiRequestEventBusEvents[K]) {
    eventBus.dispatch({
      key,
      arg,
    });
  }

  return {
    subscribe,
    dispatch,
  };
}
