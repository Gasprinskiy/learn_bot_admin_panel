import type { AuthTempData, LoginParams } from '@/shared/types/profile';
import type { UseRedirectWindowReturnType } from '@/composables/use_redirect_window/types';

export interface UseAuthState {
  redirectWindow: UseRedirectWindowReturnType | null;
  eventSource: EventSource | null;
  tempData: AuthTempData | null;
}

export interface ListenTgAuthSourceParams {
  authId: string;
  onRequestError: () => void;
}

export interface TgAuthParams extends Omit<ListenTgAuthSourceParams, 'authId'> {
  onTempDataCreate: () => void;
}

export interface LoginCommonParams {
  tempID: string | null;
  loginParams: LoginParams | null;
  isPassowrdSet: boolean;
}
