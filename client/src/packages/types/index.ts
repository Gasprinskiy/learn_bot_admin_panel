export type ConvertEmitType<T> = {
  [K in keyof T]: T[K] extends unknown[] ? (...args: T[K]) => void : never;
};

export type PartialBy<T, K extends keyof T> = Omit<T, K> & Partial<Pick<T, K>>;
