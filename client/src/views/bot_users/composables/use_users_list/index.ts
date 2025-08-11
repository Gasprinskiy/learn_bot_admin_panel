import { useMessage } from 'naive-ui';
import { computed, shallowReactive, shallowRef } from 'vue';

import $api from '@/packages/api/client';
import type { PaginationParams } from '@/shared/types/common';
import type { BotUserProfile, BotUserProfileListResponse, BotUserProfileQueryParam } from '@/shared/types/profile';

import { ErrorMessagesByCode } from './constants';

export function useUsersList() {
  const message = useMessage();

  const data = shallowRef<BotUserProfile[]>([]);
  const isLoadingMore = shallowRef<boolean>(false);
  const leftDataCount = shallowRef<number>(0);

  const searchParams = shallowReactive<BotUserProfileQueryParam>({});
  const pagidationParams = shallowReactive<PaginationParams>({
    limit: 10,
    page: 1,
  });

  const noData = computed<boolean>(() => data.value.length === 0);
  const isDataLeft = computed<boolean>(() => leftDataCount.value > 0);
  const showLoadMoreButton = computed<boolean>(() => !noData.value && isDataLeft.value);

  function resetSearchParams() {
    for (const key in searchParams) {
      const k = key as keyof typeof searchParams;
      searchParams[k] = undefined;
    }

    pagidationParams.limit = 10;
    pagidationParams.page = 1;
    pagidationParams.next_cursor_date = undefined;
    pagidationParams.next_cursor_id = undefined;
  }

  async function fetchRegisteredUsers(reset?: boolean) {
    const isReset = reset !== undefined && reset === true;

    if (isReset) {
      pagidationParams.next_cursor_date = undefined;
      pagidationParams.next_cursor_id = undefined;
    }

    try {
      const response = await $api<BotUserProfileListResponse>('/bot_users', {
        params: {
          ...searchParams,
          ...pagidationParams,
        },
      });
      if (isReset) {
        data.value = response.data;
      } else {
        data.value = [...data.value, ...response.data];
        pagidationParams.next_cursor_date = response.data[response.data.length - 1].join_date;
        pagidationParams.next_cursor_id = response.data[response.data.length - 1].u_id;
      }
      leftDataCount.value = response.left;
    } catch (e) {
      if (isReset) {
        data.value = [];
      }

      const stauts = +(e as any).status || 500;
      message.error(ErrorMessagesByCode[stauts]);
    }
  }

  function _downloadBlob(file: Blob, fileName: string): void {
    const a = document.createElement('a');
    const url = window.URL.createObjectURL(file);
    document.body.appendChild(a);

    a.href = url;
    a.download = fileName;

    const clickEvent = new MouseEvent('click', {
      bubbles: true,
      cancelable: true,
      view: window,
    });
    a.dispatchEvent(clickEvent);

    a.remove();
    window.URL.revokeObjectURL(url);
  }

  async function printRegisteredUsers() {
    try {
      const response = await $api<Blob>('/bot_users/excel_file', {
        params: {
          ...searchParams,
        },
      });

      _downloadBlob(response, 'bot_users.xlsx');
    } catch (e) {
      const stauts = +(e as any).status || 500;
      message.error(ErrorMessagesByCode[stauts]);
    }
  }

  async function loadMoreRegisteredUsers() {
    if (isLoadingMore.value) {
      return;
    }

    pagidationParams.page += 1;

    isLoadingMore.value = true;

    try {
      await fetchRegisteredUsers();
    } finally {
      isLoadingMore.value = false;
    }
  }

  return {
    data,
    noData,
    searchParams,
    pagidationParams,
    showLoadMoreButton,
    isDataLeft,
    isLoadingMore,
    fetchRegisteredUsers,
    loadMoreRegisteredUsers,
    resetSearchParams,
    printRegisteredUsers,
  };
}
