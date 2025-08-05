import { darkTheme, lightTheme } from 'naive-ui';
import { Theme } from './types';
import type { BuiltInGlobalTheme } from 'naive-ui/es/themes/interface';

export const ThemeMap: Record<Theme, BuiltInGlobalTheme> = {
  [Theme.DARK]: darkTheme,
  [Theme.LIGHT]: lightTheme,
};
