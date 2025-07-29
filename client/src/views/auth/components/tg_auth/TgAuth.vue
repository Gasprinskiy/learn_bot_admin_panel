<script lang="ts" setup>
import { computed, onMounted, shallowRef } from 'vue';
import { NButton, NDivider, NQrCode, NSpace, NNumberAnimation } from 'naive-ui';

import { useConfig } from '@/composables/use_config';
import type { AuthTempData } from '@/shared/types/profile';

import type { TgAuthEmits } from './types';

const props = defineProps<AuthTempData>();
const emit = defineEmits<TgAuthEmits>();

const { SSETTL } = useConfig();

const countDownLeftMinutes = shallowRef<number>(SSETTL);
const countDownLeftSeconds = shallowRef<number>(0);

const countDownCommon = computed<string>(() => {
  const decimal = countDownLeftSeconds.value >= 10 ? `${countDownLeftSeconds.value}` : `0${countDownLeftSeconds.value}`;
  return `${countDownLeftMinutes.value}:${decimal}`;
});

function startCountDown() {
  const interval = setInterval(() => {
    if (countDownLeftMinutes.value < 0) {
      clearInterval(interval);
    }

    if (countDownLeftSeconds.value === 0) {
      countDownLeftSeconds.value = 60;
      countDownLeftMinutes.value -= 1;
    }

    countDownLeftSeconds.value -= 1;
  }, 1000);
}

onMounted(() => {
  startCountDown();
});
</script>

<template>
  <div class="tg_auth">
    <NButton @click="emit('onLinkClick')">
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
          :value="props.auth_url"
          error-correction-level="H"
        />
      </div>
    </NSpace>

    <div class="tg_auth__countdown">
      <h4>Осталось времени: <span class="tg_auth__countdown_time">{{ countDownCommon }}</span></h4>
    </div>
  </div>
</template>

<style lang="scss" src="./style.scss" />
