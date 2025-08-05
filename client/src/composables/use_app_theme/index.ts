import { useStorage } from '@vueuse/core';
import type { BuiltInGlobalTheme } from 'naive-ui/es/themes/interface';
import { computed } from 'vue';
import { Theme } from './types';
import { ThemeMap } from './constants';

const theme = useStorage<Theme>('theme', Theme.LIGHT);

export function useAppTheme() {
  const currentTheme = computed<BuiltInGlobalTheme>(() => ThemeMap[theme.value]);
  const isDarkTheme = computed<boolean>(() => theme.value === Theme.DARK);

  function toggleTheme() {
    const theme = isDarkTheme.value ? Theme.LIGHT : Theme.DARK;
    setCurrentTheme(theme);
  }

  function setCurrentTheme(value: Theme): void {
    theme.value = value;
  }

  return {
    currentTheme,
    isDarkTheme,
    setCurrentTheme,
    toggleTheme,
  };
}
