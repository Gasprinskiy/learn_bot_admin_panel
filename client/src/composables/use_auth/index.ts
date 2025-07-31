import { computed, reactive, shallowRef, toRefs } from 'vue';
import type { FetchOptions } from 'ofetch';
import { useMessage } from 'naive-ui';
import { useRouter } from 'vue-router';
import type { RouteLocationRaw } from 'vue-router';

import $api from '@/packages/api/client';
import type { User, AuthTempData, UserFirstLoginAnswer, UserShortInfo } from '@/shared/types/profile';
import { useConfig } from '@/composables/use_config';
import { useRedirectWindow } from '@/composables/use_redirect_window';
import { useProtectedRoutes } from '@/composables/use_protected_routes';

import type { ListenTgAuthSourceParams, LoginCommonParams, TgAuthParams, UseAuthState } from './types';
import { ErrorMessagesByCode } from './constants';

const state = reactive<UseAuthState>({
  user: null,
  redirectWindow: null,
  eventSource: null,
  tempData: null,
});

export function useAuth() {
  const message = useMessage();
  const router = useRouter();

  const { ApiURL } = useConfig();
  const { buildRoutesByAccessRight } = useProtectedRoutes();

  const tempDataLoading = shallowRef<boolean>(false);

  const isAuthorized = computed<boolean>(() => state.user !== null);

  async function tgAuth(params: TgAuthParams) {
    tempDataLoading.value = true;

    const { onTempDataCreate } = params;

    try {
      const response = await $api<AuthTempData>('/auth/temp_data');
      state.tempData = response;

      onTempDataCreate();
      listenTgAuthSource({
        authId: response.uu_id,
        ...params,
      });
    } finally {
      tempDataLoading.value = false;
    }
  }

  function listenTgAuthSource(params: ListenTgAuthSourceParams) {
    state.redirectWindow = useRedirectWindow({
      name: 'Авторизация через телеграм',
    });

    const { authId, onRequestError } = params;

    const url = `${ApiURL}/auth/listen/${authId}`;
    state.eventSource = new EventSource(url);

    state.eventSource?.addEventListener('user_data', async (event: MessageEvent<string>) => {
      const data: UserFirstLoginAnswer = JSON.parse(event.data);

      closeEventSource();
      closeRedirectWindow();

      await login({
        tempID: authId,
        isPassowrdSet: data.is_password_set,
        loginParams: null,
      });
    });

    state.eventSource?.addEventListener('error', (event: MessageEvent<string>) => {
      message.error(ErrorMessagesByCode[+event.data], {
        duration: 3000,
      });

      onRequestError();
      closeEventSource();
      closeRedirectWindow();
    });
  }

  async function checkAuth() {
    try {
      const response = await $api<UserShortInfo>('/auth/check');
      await buildRoutesByAccessRight(response.access_right);

      state.user = response;
      await router.replace({
        name: 'home',
      });
    } catch (e) {
      console.error(e);
      await router.replace({
        name: 'auth',
      });
    }
  };

  async function login(params: LoginCommonParams) {
    const { tempID, loginParams, isPassowrdSet } = params;
    if (!tempID && !loginParams) {
      return;
    }

    const options: FetchOptions<'json', any> = {};

    if (tempID !== null) {
      options.params = {
        temp_id: tempID,
      };
    }

    if (loginParams !== null) {
      options.body = {
        username: loginParams.username,
        password: loginParams.password,
      };
    }

    try {
      const routeRaplaceParams: RouteLocationRaw = {
        name: 'home',
      };

      const response = await $api<User>('/auth/login', {
        method: 'POST',
        ...options,
      });
      await buildRoutesByAccessRight(response.access_right);

      state.user = response;
      message.success('Авторизация прошла успешно', {
        duration: 3000,
      });

      if (!isPassowrdSet) {
        message.success('Вам необходимо создать пароль для учетной записи', {
          duration: 3000,
        });

        routeRaplaceParams.name = 'set-pass';
      }

      await router.replace(routeRaplaceParams);
    } catch (e) {
      console.error('login error: ', e);
      message.error(ErrorMessagesByCode[500], {
        duration: 3000,
      });
    }
  };

  function closeEventSource() {
    state.eventSource?.close();
    state.eventSource = null;
    state.tempData = null;
  }

  function closeRedirectWindow() {
    state.redirectWindow?.close();
    state.redirectWindow = null;
  }

  return {
    ...toRefs(state),
    isAuthorized,
    tempDataLoading,
    tgAuth,
    checkAuth,
    closeEventSource,
    closeRedirectWindow,
  };
}
