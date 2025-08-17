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
