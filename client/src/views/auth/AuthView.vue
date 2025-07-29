<script setup lang="ts">
import type { Component } from 'vue';
import { computed, shallowRef } from 'vue';
import { NCard, NButton, NIcon, NDivider, NAlert } from 'naive-ui';
import { TelegramOutlined, PasswordOutlined, ArrowBackIosOutlined } from '@vicons/material';

import type { AuthTempData, UserFirstLoginAnswer } from '@/shared/types/profile';
import { useAuth } from '@/composables/use_auth';

import { AuthMethod } from './types';
import { AuthComponentMap } from './constants';

const { tempDataLoading, tgAuth, closeEventSource, redirectWindow } = useAuth();

const choseMethod = shallowRef<AuthMethod | null>(null);
const authComponent = shallowRef<Component | null>(null);
const authTempData = shallowRef<AuthTempData | null>();

const authMethodChosen = computed<boolean>(() => authComponent.value !== null);

function onAuthMehtodChose(method: AuthMethod) {
  choseMethod.value = method;
  authComponent.value = AuthComponentMap[method]();
}

async function onChoseTgAsAuthMethod() {
  await tgAuth({
    onTempDataCreate: (data: AuthTempData) => {
      authTempData.value = data;
      onAuthMehtodChose(AuthMethod.TELEGRAM);
    },
    authAnswer: (firstLoginAnswer: UserFirstLoginAnswer) => {
      console.log('firstLoginAnswer: ', firstLoginAnswer);
    },
    authErrorAnswer: resetAuthMethod,
  });
}

function resetAuthMethod() {
  choseMethod.value = null;
  authComponent.value = null;
  closeEventSource();
}
</script>

<template>
  <div class="auth-view">
    <NCard class="auth-view__card">
      <template #header>
        <div class="auth-view__card_head">
          <NButton
            class="auth-view__card_back-button"
            :class="{ hidded: !authMethodChosen }"
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
        <div
          v-if="authMethodChosen"
          class="auth-view__card_chosen-component"
        >
          <component
            :is="authComponent"
            v-bind="authTempData"
            v-on="{
              onLinkClick: () => redirectWindow.open(authTempData!.auth_url),
            }"
          />
        </div>

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
            Если ваш аккаунт не активирован вам нужно пройти авторизацию через Telegram
          </NAlert>
        </div>
      </template>
    </NCard>
  </div>
</template>

<style lang="scss" src="./style.scss" />
