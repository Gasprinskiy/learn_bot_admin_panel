<script setup lang="ts">
import { NConfigProvider, NGlobalStyle, NScrollbar, NMessageProvider } from 'naive-ui';
// import AppPageHeader from '@/components/header/Header.vue';
// import ModalProvider from '@/components/modal-provier/ModalProvider.vue';
</script>

<template>
  <NConfigProvider>
    <NMessageProvider>
      <NGlobalStyle />

      <!-- <ModalProvider /> -->
      <div class="app-wrapper">
        <!-- <AppPageHeader /> -->
        <div class="app-wrapper__scrollbar-container">
          <NScrollbar class="app-wrapper__scrollbar">
            <RouterView />
          </NScrollbar>
        </div>
      </div>
    </NMessageProvider>
  </NConfigProvider>
</template>

<style lang="scss">
.app-wrapper {
  &__scrollbar-container {
    padding: 16px;
  }

  &__scrollbar {
    height: calc(100vh - 92px);
  }
}
</style>
