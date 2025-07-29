import type { AuthTempData, User, UserFirstLoginAnswer } from '@/shared/types/profile';

export interface UseAuthState {
  user: User | null;
}

export interface ListenTgAuthSourceParams {
  authId: string;
  authAnswer: (firstLoginAnswer: UserFirstLoginAnswer) => void;
  authErrorAnswer: () => void;
}

export interface TgAuthParams extends Omit<ListenTgAuthSourceParams, 'authId'> {
  onTempDataCreate: (tempData: AuthTempData) => void;
}
