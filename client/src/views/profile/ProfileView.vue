<script setup lang="ts">
import { useAuth } from '@/composables/use_auth';
import { useModal } from '@/composables/use_modal';
import type { ConvertEmitType } from '@/packages/types';
import type { User } from '@/shared/types/profile';
import { LogOutOutlined, PasswordOutlined, TelegramOutlined } from '@vicons/material';
import { NAlert, NButton, NCard, NIcon, NTag } from 'naive-ui';
import { computed, defineAsyncComponent, onBeforeMount, ref } from 'vue';
import type { SetPasswordEmits } from './components/set_password/types';

const { getUserInfo, createPassword } = useAuth();
const { showModal, closeModal } = useModal();

const userInfo = ref<User | null>(null);

const userAvatarView = computed<string>(() => {
  if (userInfo.value === null) {
    return 'U';
  }

  const { first_name, last_name } = userInfo.value;

  return `${first_name[0]}${last_name[0]}`.toUpperCase();
});

async function fetchUserInfo() {
  userInfo.value = await getUserInfo();
}

async function onPasswordEmit(password: string) {
  if (userInfo.value === null) {
    return;
  }

  const done = await createPassword(password);
  userInfo.value.is_password_set = done;
  closeModal();
}

function onPasswordSet() {
  const component = defineAsyncComponent(() => {
    return import('./components/set_password/SetPassword.vue');
  });

  const emits: ConvertEmitType<SetPasswordEmits> = {
    onSubmit: onPasswordEmit,
  };

  showModal({
    component,
    width: 400,
    emits,
  });
}

onBeforeMount(fetchUserInfo);
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
        <NButton
          v-if="!userInfo.is_password_set"
          type="warning"
          @click="onPasswordSet"
        >
          <template #icon>
            <NIcon :component="PasswordOutlined" />
          </template>

          <template #default>
            Создать пароль
          </template>
        </NButton>

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

    <NAlert
      v-if="!userInfo.is_password_set"
      type="info"
    >
      <template #header>
        Создайте пароль, чтобы в будущем входить с его помощью.
      </template>

      <template #default>
        В данный момент авторизация доступна только через Telegram-бота. Чтобы использовать стандартный вход — создайте пароль.
      </template>
    </NAlert>
  </div>
</template>

<style lang="scss" src="./style.scss" />
