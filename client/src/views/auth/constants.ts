import { AuthMethod } from './types';

export const AuthMethodPathMap: Record<AuthMethod, string> = {
  [AuthMethod.STANDART]: 'standart-auth',
  [AuthMethod.TELEGRAM]: 'telegram-auth',
} as const;

export const AuthTempDataInjectKey = 'auth_temp_data' as const;
export const ChildRouteNameMap: Record<string, boolean> = {
  'telegram-auth': true,
  'standart-auth': true,
  'set-pass': true,
};
