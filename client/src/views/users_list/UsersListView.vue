<script setup lang="ts">
import $api from '@/packages/api/client';
import type { PaginationParams } from '@/shared/types/common';
import type { BotUserProfileQueryParam, BotUserProfileListResponse } from '@/shared/types/profile';
import { shallowReactive, shallowRef } from 'vue';

const searchParams = shallowReactive<BotUserProfileQueryParam>({});
const pagidationParams = shallowReactive<PaginationParams>({
  limit: 100,
  page: 1,
});
const count = shallowRef<number>(1);

async function findUsers() {
  try {
    const response = await $api<BotUserProfileListResponse>('/bot_users', {
      params: {
        ...searchParams,
        ...pagidationParams,
      },
    });
    pagidationParams.next_cursor_date = response.data[response.data.length - 1].join_date;
    pagidationParams.next_cursor_id = response.data[response.data.length - 1].u_id;
    count.value += 1;
  } catch (e) {
    console.error(e);
  }
}
</script>

<template>
  <div>
    USERS
    <button @click="findUsers">
      FIND
    </button>
  </div>
</template>

<style scoped>

</style>
