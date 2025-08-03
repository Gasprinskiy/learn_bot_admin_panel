<script lang="ts" setup>
import { onBeforeUnmount, onBeforeMount } from 'vue';
import { NButton, NDivider, NQrCode, NSpace } from 'naive-ui';
import { useRouter } from 'vue-router';

import { useAuth } from '@/composables/use_auth';
import { useSSETimeOutCountDown } from '@/views/auth/composables/use_sse_count_down';

const router = useRouter();
const { redirectWindow, tempData, closeRedirectWindow, closeEventSource } = useAuth();
const { countDownView } = useSSETimeOutCountDown();

onBeforeMount(async () => {
  if (!tempData.value) {
    await router.replace({
      name: 'auth',
    });
  }
});

onBeforeUnmount(() => {
  closeRedirectWindow();
  closeEventSource();
});
</script>

<template>
  <div
    v-if="tempData"
    class="tg_auth"
  >
    <NButton
      type="primary"
      @click="redirectWindow!.open(tempData.auth_url)"
    >
      Перейдите по ссылке
    </NButton>

    <NDivider>
      или
    </NDivider>

    <NSpace class="tg_auth__qr-code" vertical>
      <p>Отсканируйте QR-код</p>
      <div class="tg_auth__qr-code__itself">
        <NQrCode
          style="padding: 0;"
          :size="200"
          :value="tempData.auth_url"
          error-correction-level="H"
        />
      </div>
    </NSpace>

    <div class="tg_auth__countdown">
      <h4>Осталось времени: <span class="tg_auth__countdown_time">{{ countDownView }}</span></h4>
    </div>
  </div>
</template>

<style lang="scss" src="./style.scss" />
