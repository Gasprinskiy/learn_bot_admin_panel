import type { AuthTempData, User, UserFirstLoginAnswer } from '@/shared/types/profile';
import type { UseRedirectWindowReturnType } from '../use_redirect_window/types';

export interface UseAuthState {
  user: User | null;
  redirectWindow: UseRedirectWindowReturnType | null;
  eventSource: EventSource | null;
}

export interface ListenTgAuthSourceParams {
  authId: string;
  authAnswer: (firstLoginAnswer: UserFirstLoginAnswer) => void;
  authErrorAnswer: () => void;
}

export interface TgAuthParams extends Omit<ListenTgAuthSourceParams, 'authId'> {
  onTempDataCreate: (tempData: AuthTempData) => void;
}
