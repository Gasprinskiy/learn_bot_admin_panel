<script setup lang="ts">
import { useLoadingBar, useMessage } from 'naive-ui';
import { useApiRequestEventBus } from '@/composables/use_api_requests_event_bus';
import { onMounted, onBeforeMount } from 'vue';

import { useAuth } from './composables/use_auth';
import { useRouter } from 'vue-router';

const router = useRouter()
const loadingBar = useLoadingBar();
const apiEventBus = useApiRequestEventBus();
const message = useMessage();
const { checkAuthOnFirstRun, checkAuthOnRouteChange } = useAuth();

onBeforeMount(() => {
  // checkAuthOnFirstRun();
});

onMounted(() => {
  apiEventBus.subscribe('on_request', () => loadingBar.start());
  apiEventBus.subscribe('on_response', () => loadingBar.finish());
  apiEventBus.subscribe('on_error', (arg: { message: string | undefined } | null) => {
    loadingBar.error();
    if (arg && arg.message) {
      message.error(arg.message, { duration: 3000 });
    }
  });
});

router.beforeEach(checkAuthOnRouteChange)
</script>

<template>
  <RouterView />
</template>
