import type { BusEvent, BusEventCallBack, BusEventDispatch, UnsubscribeCallBack } from './types';

export class EventBus {
  private eventsMap = new Map<string, Map<number, BusEventCallBack>>();

  subscribe(event: BusEvent): UnsubscribeCallBack {
    const eventsByKey = this.eventsMap.get(event.key);
    const count = eventsByKey ? eventsByKey.size + 1 : 1;

    const newEvents = eventsByKey || new Map<number, BusEventCallBack>();
    newEvents.set(count, event.callBack);

    this.eventsMap.set(event.key, newEvents);

    return () => {
      const unCount = count - 1;
      if (unCount <= 0) {
        this.eventsMap.delete(event.key);
        return;
      }

      newEvents.delete(count);
    };
  };

  dispatch(dispatchEvent: BusEventDispatch) {
    const events = this.eventsMap.get(dispatchEvent.key);
    if (!events) {
      return;
    }

    events.forEach((callBack: BusEventCallBack) => {
      callBack(dispatchEvent.arg);
    });
  };
}
