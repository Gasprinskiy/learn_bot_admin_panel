export type BusEventCallBack = (arg: any | null) => void;
export type UnsubscribeCallBack = () => void;

export interface BusEvent {
  key: string;
  callBack: any | null;
}

export interface BusEventDispatch {
  key: string;
  arg: any | null;
}
