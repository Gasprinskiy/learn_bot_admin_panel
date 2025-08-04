<script lang="ts" setup>
import { NButton, NCard, NInput } from 'naive-ui';
import type { FormValidationStatus } from 'naive-ui';
import { computed, shallowReactive } from 'vue';
import { minLength, required, sameAs } from '@vuelidate/validators';
import useVuelidate from '@vuelidate/core';

import type { SetPasswordState, SetPasswordEmits } from './types';

const emit = defineEmits<SetPasswordEmits>();

const state = shallowReactive<SetPasswordState>({
  passwordSource: '',
  passwordСonfirm: '',
});

const validators: Record<keyof SetPasswordState, Record<string, any>> = {
  passwordSource: { required, minLength: minLength(7) },
  passwordСonfirm: { required, sameAs: sameAs(computed<string>(() => state.passwordSource)) },
};

const v$ = useVuelidate(validators, state);

const passwordSourceInvalidMessage = computed<string | null>(() => {
  if (isFieldInvalid('passwordSource')) {
    if (v$.value.passwordSource.required.$invalid) {
      return 'Поле не заполнено';
    }

    return 'Минимальная длина 7 символов';
  }

  return null;
});
const passwordСonfirmMessage = computed<string | null>(() => {
  if (isFieldInvalid('passwordСonfirm')) {
    if (v$.value.passwordСonfirm.required.$invalid) {
      return 'Поле не заполнено';
    }

    return 'Пароли не совпадают';
  }

  return null;
});
const passwordSourceValidationStatus = computed<FormValidationStatus>(() => getFieldValidationStatus('passwordSource'));
const passwordСonfirmValidationStatus = computed<FormValidationStatus>(() => getFieldValidationStatus('passwordСonfirm'));

function isFieldInvalid(fieldName: keyof SetPasswordState): boolean {
  return v$.value[fieldName]?.$dirty && v$.value[fieldName]?.$invalid;
}

function getFieldValidationStatus(fieldName: keyof SetPasswordState): FormValidationStatus {
  return isFieldInvalid(fieldName) ? 'error' : 'success';
}

async function onSubmit() {
  const valid = await v$.value.$validate();
  if (!valid) {
    return;
  }

  v$.value.$reset();

  emit('onSubmit', state.passwordSource);
}
</script>

<template>
  <NCard class="set-password">
    <template #header>
      Создать пароль
    </template>

    <template #default>
      <form
        class="set-password__form"
        @submit.prevent="onSubmit"
      >
        <div class="set-password__input">
          <NInput
            v-model:value="state.passwordSource"
            type="password"
            placeholder="Пароль"
            show-password-on="click"
            :status="passwordSourceValidationStatus"
          />

          <span
            v-if="passwordSourceInvalidMessage"
            class="set-password__input_message"
          >
            {{ passwordSourceInvalidMessage }}
          </span>
        </div>

        <div class="set-password__input">
          <NInput
            v-model:value="state.passwordСonfirm"
            type="password"
            placeholder="Повторите пароль"
            show-password-on="click"
            :status="passwordСonfirmValidationStatus"
          />

          <span
            v-if="passwordСonfirmMessage"
            class="set-password__input_message"
          >
            {{ passwordСonfirmMessage }}
          </span>
        </div>

        <NButton
          type="primary"
          attr-type="submit"
        >
          Создать
        </NButton>
      </form>
    </template>
  </NCard>
</template>

<style lang="scss" src="./style.scss" />
