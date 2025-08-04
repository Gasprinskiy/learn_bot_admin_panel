<script setup lang="ts">
import { NDivider, NScrollbar, useLoadingBar, useMessage } from 'naive-ui';
import { useApiRequestEventBus } from '@/composables/use_api_requests_event_bus';
import { onMounted, shallowRef } from 'vue';
import { useRouter } from 'vue-router';

import { useAuth } from './composables/use_auth';
import AppHeader from './components/header/AppHeader.vue';
import type { UseApiRequestEventBusEvents } from './composables/use_api_requests_event_bus/types';
import ModalProvider from './components/modal-provider/ModalProvider.vue';
import { useModal } from './composables/use_modal';

const router = useRouter();
const loadingBar = useLoadingBar();
const apiEventBus = useApiRequestEventBus();
const message = useMessage();
const { isVisible } = useModal();
const { checkAuthOnRouteChange } = useAuth();

const isBlocked = shallowRef<boolean>(false);

onMounted(() => {
  apiEventBus.subscribe('on_request', (arg: UseApiRequestEventBusEvents['on_request']) => {
    isBlocked.value = arg.is_blocking || false;
    loadingBar.start();
  });

  apiEventBus.subscribe('on_response', () => {
    isBlocked.value = false;
    loadingBar.finish();
  });

  apiEventBus.subscribe('on_error', (arg: UseApiRequestEventBusEvents['on_error']) => {
    isBlocked.value = false;
    loadingBar.error();
    if (arg && arg.message) {
      message.error(arg.message, { duration: 3000 });
    }
  });
});

router.beforeEach(checkAuthOnRouteChange);
</script>

<template>
  <ModalProvider />

  <div
    class="app-wrapper"
    :class="{ blocked: isBlocked || isVisible }"
  >
    <AppHeader />
    <NDivider class="app-wrapper__divider" />
    <div class="app-wrapper__scrollbar-container">
      <NScrollbar class="app-wrapper__scrollbar">
        <RouterView />
      </NScrollbar>
    </div>
  </div>
</template>

<style lang="scss">
.app-wrapper {
  position: relative;

  &.blocked {
    pointer-events: none;
  }

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
