<script setup lang="ts">
import { useAuth } from '@/composables/use_auth';
import $api from '@/packages/api/client';
import type { User } from '@/shared/types/profile';
import { LogOutOutlined, TelegramOutlined } from '@vicons/material';
import { NButton, NCard, NIcon, NTag } from 'naive-ui';
import { computed, onBeforeMount, shallowRef } from 'vue';

const { getUserInfo } = useAuth();

const userInfo = shallowRef<User | null>(null);
const dataLoading = shallowRef<boolean>(false);

const userAvatarView = computed<string>(() => {
  if (userInfo.value === null) {
    return 'U';
  }

  const { first_name, last_name } = userInfo.value;

  return `${first_name[0]}${last_name[0]}`.toUpperCase();
});

onBeforeMount(async () => {
  dataLoading.value = true;
  userInfo.value = await getUserInfo();
  dataLoading.value = false;
});
</script>

<template>
  <div v-if="userInfo" class="profile-view">
    <NCard class="profile-view__card">
      <div class="profile-view__card-left">
        <NTag
          type="success"
          class="profile-view__avatar"
          :bordered="false"
        >
          {{ userAvatarView }}
        </NTag>

        <div class="profile-view__card-info">
          <h4>{{ userInfo.first_name }} {{ userInfo.last_name }}</h4>

          <NTag
            type="info"
            round
          >
            <template #icon>
              <NIcon :component="TelegramOutlined" />
            </template>

            <template #default>
              {{ userInfo.tg_user_name }}
            </template>
          </NTag>
        </div>
      </div>

      <div class="profile-view__card-right">
        <NButton type="error">
          <template #icon>
            <NIcon :component="LogOutOutlined" />
          </template>

          <template #default>
            Выйти
          </template>
        </NButton>
      </div>
    </NCard>
  </div>
</template>

<style lang="scss" src="./style.scss" />
