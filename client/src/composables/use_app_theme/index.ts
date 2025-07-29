import { darkTheme } from 'naive-ui';
import type { BuiltInGlobalTheme } from 'naive-ui/es/themes/interface';
import { shallowRef } from 'vue';

const currentTheme = shallowRef<BuiltInGlobalTheme>(darkTheme);

export function useAppTheme() {
  function setCurrentTheme(theme: BuiltInGlobalTheme): void {
    currentTheme.value = theme;
  }

  return {
    currentTheme,
    setCurrentTheme,
  };
}
