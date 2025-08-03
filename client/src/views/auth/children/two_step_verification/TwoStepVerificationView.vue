<script setup lang="ts">
import { useSSETimeOutCountDown } from '@/views/auth/composables/use_sse_count_down';

const { countDownView } = useSSETimeOutCountDown();
</script>

<template>
  <div class="two-step-verification">
    <h4 class="two-step-verification__title">
      Дополнительный шаг чтобы убедиться, что это действительно вы.
    </h4>

    <div class="two-step-verification__image">
      <img src="@/assets/two_step.png">
    </div>

    <div class="two-step-verification__description">
      Мы отправили запрос на подтверждение в Telegram.
      Нажмите <span>"✅ Да"</span> чтобы продолжить.
    </div>

    <div class="two-step-verification__countdown">
      <h4>Осталось времени: <span class="two-step-verification__countdown_time">{{ countDownView }}</span></h4>
    </div>
  </div>
</template>

<style lang="scss" src="./style.scss" />
