import { computed, reactive, shallowRef } from 'vue';

import $api from '@/packages/api/client';
import type { AuthTempData, UserFirstLoginAnswer } from '@/shared/types/profile';
import { useConfig } from '@/composables/use_config';
import { useRedirectWindow } from '@/composables/use_redirect_window';

import type { ListenTgAuthSourceParams, TgAuthParams, UseAuthState } from './types';
import { useMessage } from 'naive-ui';
import { ErrorMessagesByCode } from './constants';

const state = reactive<UseAuthState>({
  user: null,
});

export function useAuth() {
  const redirectWindow = useRedirectWindow({
    name: 'Авторизация через телеграм',
  });

  const { ApiURL } = useConfig();
  const message = useMessage();

  const tempDataLoading = shallowRef<boolean>(false);
  const eventSource = shallowRef<EventSource | null>(null);

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
    } catch (e) {
      console.error(e);
    } finally {
      tempDataLoading.value = false;
    }
  }

  function listenTgAuthSource(params: ListenTgAuthSourceParams) {
    const { authId, authAnswer, authErrorAnswer } = params;

    const url = `${ApiURL}/auth/listen/${authId}`;
    eventSource.value = new EventSource(url);

    eventSource.value?.addEventListener('user_data', (event: MessageEvent<string>) => {
      const data: UserFirstLoginAnswer = JSON.parse(event.data);
      authAnswer(data);
      redirectWindow.close();
      eventSource.value?.close();
    });

    eventSource.value?.addEventListener('error', (event: MessageEvent<string>) => {
      message.error(ErrorMessagesByCode[+event.data], {
        duration: 3000,
      });
      authErrorAnswer();
      redirectWindow.close();
      eventSource.value?.close();
    });
  }

  function closeEventSource() {
    eventSource.value?.close();
  }

  return {
    isAuthorized,
    redirectWindow,
    tempDataLoading,
    tgAuth,
    closeEventSource,
  };
}
