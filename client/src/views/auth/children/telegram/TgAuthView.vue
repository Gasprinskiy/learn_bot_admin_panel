<script lang="ts" setup>
import { computed, onBeforeUnmount, onBeforeMount, onMounted, shallowRef } from 'vue';
import { NButton, NDivider, NQrCode, NSpace } from 'naive-ui';
import { useRouter } from 'vue-router';

import { useAuth } from '@/composables/use_auth';
import { useConfig } from '@/composables/use_config';

const router = useRouter();
const { redirectWindow, tempData, closeRedirectWindow, closeEventSource } = useAuth();
const { SSETTL } = useConfig();

const countDownLeftMinutes = shallowRef<number>(SSETTL);
const countDownLeftSeconds = shallowRef<number>(0);
const countDownIntervalId = shallowRef<number | undefined>(undefined);

const countDownView = computed<string>(() => {
  const decimal = countDownLeftSeconds.value >= 10 ? `${countDownLeftSeconds.value}` : `0${countDownLeftSeconds.value}`;
  return `${countDownLeftMinutes.value}:${decimal}`;
});

function startCountDown() {
  countDownIntervalId.value = setInterval(() => {
    if (countDownLeftMinutes.value < 0) {
      clearInterval(countDownIntervalId.value);
    }

    if (countDownLeftSeconds.value === 0) {
      countDownLeftSeconds.value = 60;
      countDownLeftMinutes.value -= 1;
    }

    countDownLeftSeconds.value -= 1;
  }, 1000);
}

onBeforeMount(async () => {
  if (!tempData.value) {
    await router.replace({
      name: 'auth',
    });
  }
});

onMounted(() => {
  startCountDown();
});

onBeforeUnmount(() => {
  closeRedirectWindow();
  closeEventSource();
  clearInterval(countDownIntervalId.value);
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
