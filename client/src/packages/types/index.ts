export type ConvertEmitType<T> = {
  [K in keyof T]: T[K] extends unknown[] ? (...args: T[K]) => void : never;
};
