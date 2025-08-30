<script setup lang="ts">
import { NDivider, NTabPane, NTabs } from 'naive-ui';
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { UserListTab } from './types';

const route = useRoute();
const router = useRouter();

const currentTab = computed<UserListTab>({
  get(): UserListTab {
    const param = route.matched[1].name as UserListTab;
    return param;
  },

  async set(value: UserListTab) {
    await router.push({
      name: value,
    });
  },
});

const isSingleUser = computed<boolean>(() => route.name === 'registered-user');
</script>

<template>
  <div class="users-list-view">
    <template v-if="!isSingleUser">
      <h2 class="users-list-view__head">
        Пользователи бота
      </h2>

      <div class="users-list-view__tabs">
        <NTabs
          v-model:value="currentTab"
          type="segment"
        >
          <NTabPane tab="Зарегистрированные" :name="UserListTab.REGISTERED" />
          <NTabPane tab="Не зарегистрированные" :name="UserListTab.UNREGISTERED" />
        </NTabs>
      </div>

      <NDivider class="users-list-view__divider" />
    </template>

    <div class="users-list-view__body">
      <RouterView />
    </div>
  </div>
</template>

<style lang="scss" src="./style.scss" />
