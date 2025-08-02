<script setup lang="ts">
import {
  NConfigProvider,
  NGlobalStyle,
  NMessageProvider,
  NLoadingBarProvider,
  NNotificationProvider,
} from 'naive-ui';

import AppRouterView from '@/AppRouterView.vue';
import { useAppTheme } from '@/composables/use_app_theme';

const { currentTheme } = useAppTheme();
</script>

<template>
  <NConfigProvider :theme="currentTheme">
    <NLoadingBarProvider>
      <NMessageProvider
        placement="bottom"
        :duration="3000"
      >
        <NNotificationProvider>
          <NGlobalStyle />
          <AppRouterView />
        </NNotificationProvider>
      </NMessageProvider>
    </NLoadingBarProvider>
  </NConfigProvider>
</template>
