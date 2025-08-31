export function pluralize(count: number, one: string, few: string, many: string): string {
  count = Math.abs(count) % 100;
  // const lastDigit = count % 10;

  if (count > 1 && count < 5) {
    return few;
  }

  if (count === 1) {
    return one;
  }

  return many;
}

export function optionalResult<T extends string | number | Date>(value: T, valueProxy?: (value: T) => string): string {
  if (value === null) {
    return '-';
  }

  return valueProxy ? valueProxy(value) : value.toString();
}
