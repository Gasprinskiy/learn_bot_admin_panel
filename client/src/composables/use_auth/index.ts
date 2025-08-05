import { computed, reactive, shallowRef, toRefs } from 'vue';
import { useMessage, useNotification } from 'naive-ui';
import { useRouter } from 'vue-router';
import { useStorage } from '@vueuse/core';
import type { NavigationGuardNext, RouteLocationNormalizedLoadedGeneric } from 'vue-router';

import $api from '@/packages/api/client';
import { useConfig } from '@/composables/use_config';
import { useRedirectWindow } from '@/composables/use_redirect_window';
import { useProtectedRoutes } from '@/composables/use_protected_routes';
import type { ResponseWithBoolStatus } from '@/shared/types/common';
import type {
  AuthTempData,
  UserFirstLoginAnswer,
  UserShortInfo,
  AccessRight,
  User,
  PasswordLoginParams,
  PasswordLoginResponse,
} from '@/shared/types/profile';

import type { ListenTgAuthSourceParams, TgAuthParams, UseAuthState } from './types';
import { ErrorMessagesByCode, PasswordLoginErrorMessagesByCode } from './constants';

const hasToken = useStorage('has_token', false);

const state = reactive<UseAuthState>({
  redirectWindow: null,
  eventSource: null,
  tempData: null,
});

export function useAuth() {
  const message = useMessage();
  const notification = useNotification();
  const router = useRouter();

  const { ApiURL } = useConfig();
  const { setRoutesByAccesRight } = useProtectedRoutes();

  const tempDataLoading = shallowRef<boolean>(false);

  const isAuthorized = computed<boolean>(() => hasToken.value);

  function _authCheck(): Promise<UserShortInfo> {
    return $api<UserShortInfo>('/auth/check');
  }

  function _listenTgAuthSource(params: ListenTgAuthSourceParams) {
    state.redirectWindow = useRedirectWindow({
      name: 'Авторизация через телеграм',
    });

    _listenAuthSource(params);
  }

  function _listenAuthSource(params: ListenTgAuthSourceParams) {
    const { authId, onRequestError } = params;

    const url = `${ApiURL}/auth/listen/${authId}`;
    state.eventSource = new EventSource(url);

    state.eventSource?.addEventListener('user_data', async (event: MessageEvent<string>) => {
      const data: UserFirstLoginAnswer = JSON.parse(event.data);

      closeEventSource();
      closeRedirectWindow();

      await tgLogin(authId, data.is_password_set);

      clearTempData();
    });

    state.eventSource?.addEventListener('error', (event: MessageEvent<string>) => {
      message.error(ErrorMessagesByCode[+event.data]);

      onRequestError();
      closeEventSource();
      closeRedirectWindow();
    });
  };

  async function getTgAuthDataAndListen(params: TgAuthParams) {
    if (tempDataLoading.value) {
      return;
    }

    tempDataLoading.value = true;

    const { onTempDataCreate } = params;

    try {
      const response = await $api<AuthTempData>('/auth/temp_data');
      state.tempData = response;

      onTempDataCreate();
      _listenTgAuthSource({
        authId: response.uu_id,
        ...params,
      });
    } catch (e) {
      const stauts = +(e as any).status || 500;
      message.error(ErrorMessagesByCode[stauts]);
    } finally {
      tempDataLoading.value = false;
    }
  }

  async function checkAuthOnFirstRun() {
    try {
      const response = await _authCheck();
      setRoutesByAccesRight(response.access_right);
    } catch (e) {
      if (hasToken.value) {
        const stauts = +(e as any).status || 500;
        message.error(ErrorMessagesByCode[stauts]);
        hasToken.value = false;
      }

      await router.replace({
        name: 'auth',
      });
    }
  };

  async function checkAuthOnRouteChange(to: RouteLocationNormalizedLoadedGeneric, from: RouteLocationNormalizedLoadedGeneric, next: NavigationGuardNext) {
    if (to.fullPath.includes('auth')) {
      if (hasToken.value) {
        next(from.fullPath);
        return;
      }
      next();
      return;
    }
    const arMap = to.meta.accessRights ? (to.meta.accessRights as Record<AccessRight, boolean>) : null;

    if (!hasToken.value) {
      next('/auth');
      message.error('Требуется авторизация');
      return;
    }

    try {
      const user = await _authCheck();

      if (arMap) {
        const exists = arMap[user.access_right];
        if (!exists) {
          next(from.fullPath);
          message.warning('Отказано в доступе');
          return;
        }
      }

      if (to.name === 'home') {
        setRoutesByAccesRight(user.access_right);
      }

      hasToken.value = true;
      next();
    } catch (e) {
      const stauts = +(e as any).status || 500;
      message.error(ErrorMessagesByCode[stauts]);
      hasToken.value = false;
      next('/auth');
    }
  }

  async function tgLogin(tempID: string, isPassowrdSet: boolean) {
    try {
      const response = await $api<UserShortInfo>(`/auth/tg_login/${tempID}`);

      setRoutesByAccesRight(response.access_right);

      message.success('Авторизация прошла успешно');

      if (!isPassowrdSet) {
        notification.info({
          title: 'Необходимо создать пароль для учетной записи',
          description: 'Необходимо создать пароль для учетной записи, для последующего входа с использованием пароля',
          duration: 5000,
        });
      }

      hasToken.value = true;
      await router.replace({
        name: 'home',
      });
    } catch (e) {
      const stauts = +(e as any).status || 500;
      message.error(ErrorMessagesByCode[stauts]);

      hasToken.value = false;
    }
  };

  async function passwordLogin(params: PasswordLoginParams) {
    try {
      const response = await $api<PasswordLoginResponse>(`/auth/password_login`, {
        method: 'POST',
        body: {
          ...params,
          _is_blocking: true,
        },
      });

      if (response.need_two_step_auth) {
        message.info('Требуется подтверждение');

        await router.replace({
          name: 'two-step-verification',
        });

        _listenAuthSource({
          authId: response.uu_id,
          onRequestError: async () => {
            await router.replace({
              name: 'password-auth',
            });
          },
        });
        return;
      }

      setRoutesByAccesRight(response.access_right);

      message.success('Авторизация прошла успешно');

      hasToken.value = true;
      await router.replace({
        name: 'home',
      });
    } catch (e) {
      const stauts = +(e as any).status || 500;
      message.error(PasswordLoginErrorMessagesByCode[stauts]);

      hasToken.value = false;
    }
  };

  async function getUserInfo(): Promise<User | null> {
    try {
      return $api<User>('/auth/profile');
    } catch (e) {
      const stauts = +(e as any).status || 500;
      message.error(ErrorMessagesByCode[stauts]);
      return null;
    }
  }

  async function createPassword(password: string): Promise<boolean> {
    try {
      const response = await $api<ResponseWithBoolStatus>('/auth/create_password', {
        method: 'POST',
        body: {
          password,
        },
      });
      if (response.success) {
        message.success('Пароль для учетной записи создан');
      }
      return response.success;
    } catch (e) {
      const stauts = +(e as any).status || 500;
      message.error(ErrorMessagesByCode[stauts]);
      return false;
    }
  }

  async function logOut() {
    try {
      const response = await $api<ResponseWithBoolStatus>('/auth/log_out', {
        method: 'POST',
      });

      if (response.success) {
        hasToken.value = false;

        await router.replace({
          name: 'auth',
        });
        message.info('Выход из аккаунта');
      }
      return response.success;
    } catch (e) {
      const stauts = +(e as any).status || 500;
      message.error(ErrorMessagesByCode[stauts]);
      return false;
    }
  }

  function closeEventSource() {
    state.eventSource?.close();
    state.eventSource = null;
  }

  function closeRedirectWindow() {
    state.redirectWindow?.close();
    state.redirectWindow = null;
  }

  function clearTempData() {
    state.tempData = null;
  }

  return {
    ...toRefs(state),
    isAuthorized,
    tempDataLoading,
    tgLogin,
    passwordLogin,
    getTgAuthDataAndListen,
    getUserInfo,
    checkAuthOnFirstRun,
    closeEventSource,
    closeRedirectWindow,
    clearTempData,
    checkAuthOnRouteChange,
    createPassword,
    logOut,
  };
}
