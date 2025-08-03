import { AuthMethod } from './types';

export const AuthMethodPathMap: Record<AuthMethod, string> = {
  [AuthMethod.STANDART]: 'password-auth',
  [AuthMethod.TELEGRAM]: 'telegram-auth',
} as const;

export const AuthTempDataInjectKey = 'auth_temp_data' as const;
export const ChildRouteNameMap: Record<string, boolean> = {
  'telegram-auth': true,
  'password-auth': true,
  'two-step-verification': true,
};
