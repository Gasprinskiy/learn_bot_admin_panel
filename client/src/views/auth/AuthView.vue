<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { computed } from 'vue';
import { NCard, NButton, NIcon, NDivider, NAlert } from 'naive-ui';
import { TelegramOutlined, PasswordOutlined, ArrowBackIosOutlined } from '@vicons/material';

import { useAuth } from '@/composables/use_auth';

import { AuthMethod } from './types';
import { AuthMethodPathMap, ChildRouteNameMap } from './constants';

const route = useRoute();
const router = useRouter();
const { tempDataLoading, getTgAuthDataAndListen } = useAuth();

const childRouteActive = computed<boolean>(() => route.name ? (ChildRouteNameMap[route.name.toString()] || false) : false);
const hasBackAction = computed<boolean>(() => Boolean(route.meta?.hasBackAction));

async function onAuthMehtodChose(method: AuthMethod) {
  await router.push({
    name: AuthMethodPathMap[method],
  });
}

async function resetAuthMethod() {
  await router.replace('/auth');
}

async function onChoseTgAsAuthMethod() {
  await getTgAuthDataAndListen({
    onTempDataCreate: () => onAuthMehtodChose(AuthMethod.TELEGRAM),
    onRequestError: resetAuthMethod,
  });
}
</script>

<template>
  <div class="auth-view">
    <NCard class="auth-view__card">
      <template #header>
        <div class="auth-view__card_head">
          <NButton
            class="auth-view__card_back-button"
            :class="{ hidded: !hasBackAction }"
            type="primary"
            size="small"
            quaternary
            @click="resetAuthMethod"
          >
            <template #icon>
              <NIcon size="15">
                <ArrowBackIosOutlined />
              </NIcon>
            </template>
          </NButton>

          <h4 class="auth-view__card_title">
            Авторизация
          </h4>
        </div>
      </template>

      <template #default>
        <RouterView v-if="childRouteActive" />

        <div
          v-else
          class="auth-view__card_content"
        >
          <NButton
            size="large"
            type="primary"
            class="auth-view__button"
            @click="onAuthMehtodChose(AuthMethod.STANDART)"
          >
            <template #icon>
              <NIcon>
                <PasswordOutlined />
              </NIcon>
            </template>

            <template #default>
              По логину и паролю
            </template>
          </NButton>

          <NDivider class="auth-view__divider">
            или
          </NDivider>

          <NButton
            size="large"
            type="info"
            class="auth-view__button"
            :loading="tempDataLoading"
            @click="onChoseTgAsAuthMethod"
          >
            <template #icon>
              <NIcon>
                <TelegramOutlined />
              </NIcon>
            </template>

            <template #default>
              Через Telegram
            </template>
          </NButton>

          <NDivider class="auth-view__divider" />

          <NAlert type="info">
            <template #header>
              Важно!
            </template>

            <template #default>
              В случае если вы не выставили пароль для аккаунта или это ваш первый вход вам нужно пройти авторизацию через Telegram
            </template>
          </NAlert>
        </div>
      </template>
    </NCard>
  </div>
</template>

<style lang="scss" src="./style.scss" />
