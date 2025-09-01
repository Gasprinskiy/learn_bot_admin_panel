<script setup lang="ts">
import { computed, shallowRef } from 'vue';
import { required } from '@vuelidate/validators';
import useVuelidate from '@vuelidate/core';
import { NButton, NCard, NSelect } from 'naive-ui';
import type { FormValidationStatus } from 'naive-ui';

import type { SubscriptionCancelReason } from '@/shared/types/bot_users';

import type { CancelPurchaseEmits } from './types';
import { CacnelReasonOptions } from './constants';

const emit = defineEmits<CancelPurchaseEmits>();

const reason = shallowRef<SubscriptionCancelReason | null>(null);

const validators = {
  reason: { required },
};

const v$ = useVuelidate(validators, { reason });

const isReasonInvalid = computed<boolean>(() => v$.value.reason?.$dirty && v$.value.reason?.$invalid);
const reasonValidationStatus = computed<FormValidationStatus>(() => isReasonInvalid.value ? 'error' : 'success');

async function onSubmit() {
  const valid = await v$.value.$validate();
  if (!valid) {
    return;
  }

  emit('onSubmit', reason.value!);
}
</script>

<template>
  <NCard class="cancel-purchase">
    <template #header>
      Отменить подписку
    </template>

    <template #default>
      <form
        class="cancel-purchase__form"
        @submit.prevent="onSubmit"
      >
        <div class="cancel-purchase__input">
          <NSelect
            v-model:value="reason"
            :options="CacnelReasonOptions"
            :status="reasonValidationStatus"
            placeholder="Причина"
          />

          <span
            v-if="isReasonInvalid"
            class="create-purchase__input_message"
          >
            Не выбрана причина
          </span>
        </div>

        <NButton
          type="primary"
          attr-type="submit"
        >
          Отменить
        </NButton>
      </form>
    </template>
  </NCard>
</template>

<style lang="scss" src="./style.scss" />
