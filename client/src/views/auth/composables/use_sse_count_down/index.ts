import { computed, onBeforeUnmount, onMounted, shallowRef } from 'vue';

import { useConfig } from '@/composables/use_config';

export function useSSETimeOutCountDown() {
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

  onMounted(() => {
    startCountDown();
  });

  onBeforeUnmount(() => {
    clearInterval(countDownIntervalId.value);
  });

  return {
    countDownView,
  };
}
