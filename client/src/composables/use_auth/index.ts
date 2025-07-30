import { computed, reactive, shallowRef, toRefs } from 'vue';

import $api from '@/packages/api/client';
import type { AuthTempData, UserFirstLoginAnswer } from '@/shared/types/profile';
import { useConfig } from '@/composables/use_config';
import { useRedirectWindow } from '@/composables/use_redirect_window';

import type { ListenTgAuthSourceParams, TgAuthParams, UseAuthState } from './types';
import { useMessage } from 'naive-ui';
import { ErrorMessagesByCode } from './constants';

const state = reactive<UseAuthState>({
  user: null,
  redirectWindow: null,
  eventSource: null,
});

export function useAuth() {
  const { ApiURL } = useConfig();
  const message = useMessage();

  const tempDataLoading = shallowRef<boolean>(false);

  const isAuthorized = computed<boolean>(() => state.user !== null);

  async function tgAuth(params: TgAuthParams) {
    tempDataLoading.value = true;

    const { onTempDataCreate } = params;

    try {
      const response = await $api<AuthTempData>('/auth/temp_data');
      onTempDataCreate(response);
      listenTgAuthSource({
        authId: response.uu_id,
        ...params,
      });
    } catch {
      // console.error(e);
    } finally {
      tempDataLoading.value = false;
    }
  }

  function listenTgAuthSource(params: ListenTgAuthSourceParams) {
    state.redirectWindow = useRedirectWindow({
      name: 'Авторизация через телеграм',
    });

    const { authId, authAnswer, authErrorAnswer } = params;

    const url = `${ApiURL}/auth/listen/${authId}`;
    state.eventSource = new EventSource(url);

    state.eventSource?.addEventListener('user_data', (event: MessageEvent<string>) => {
      const data: UserFirstLoginAnswer = JSON.parse(event.data);
      authAnswer(data);
      message.success('Авторизация прошла успешно', {
        duration: 3000,
      });

      closeEventSource();
      closeRedirectWindow();
    });

    state.eventSource?.addEventListener('error', (event: MessageEvent<string>) => {
      message.error(ErrorMessagesByCode[+event.data], {
        duration: 3000,
      });
      authErrorAnswer();

      closeEventSource();
      closeRedirectWindow();
    });
  }

  function closeEventSource() {
    state.eventSource?.close();
    state.eventSource = null;
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
    closeEventSource,
    closeRedirectWindow,
  };
}
