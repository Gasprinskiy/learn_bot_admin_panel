<script setup lang="ts">
import {
  NConfigProvider,
  NGlobalStyle,
  NScrollbar,
  NMessageProvider,
  NLoadingBarProvider,
  NDivider,
} from 'naive-ui';

import AppHeader from '@/components/header/AppHeader.vue';
import AppRouterView from '@/AppRouterView.vue';
import { useAppTheme } from '@/composables/use_app_theme';

const { currentTheme } = useAppTheme();
// import ModalProvider from '@/components/modal-provier/ModalProvider.vue';
</script>

<template>
  <NConfigProvider :theme="currentTheme">
    <NLoadingBarProvider>
      <NMessageProvider>
        <NGlobalStyle />

        <!-- <ModalProvider /> -->
        <div class="app-wrapper">
          <AppHeader />
          <NDivider class="app-wrapper__divider" />
          <div class="app-wrapper__scrollbar-container">
            <NScrollbar class="app-wrapper__scrollbar">
              <AppRouterView />
            </NScrollbar>
          </div>
        </div>
      </NMessageProvider>
    </NLoadingBarProvider>
  </NConfigProvider>
</template>

<style lang="scss">
.app-wrapper {
  &__scrollbar-container {
    width: 100%;
    margin: auto;
  }

  &__divider {
    margin: 0 !important;
  }

  &__scrollbar {
    height: calc(100vh - 57px);
    padding: 16px;

    .n-scrollbar-content {
      height: 100%;
    }
  }
}
</style>
