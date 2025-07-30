import { AuthMethod } from './types';

export const AuthMethodPathMap: Record<AuthMethod, string> = {
  [AuthMethod.STANDART]: 'elegram-auth',
  [AuthMethod.TELEGRAM]: 'standart-auth',
} as const;

export const AuthTempDataInjectKey = 'auth_temp_data' as const;
export const ChildRouteNameMap: Record<string, boolean> = {
  'telegram-auth': true,
  'standart-auth': true,
  'set-pass': true,
};
