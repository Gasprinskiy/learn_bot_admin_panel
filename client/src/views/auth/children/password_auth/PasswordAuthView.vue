<script lang="ts" setup>
import { useAuth } from '@/composables/use_auth';
import type { PasswordLoginParams } from '@/shared/types/profile';
import { useVuelidate } from '@vuelidate/core';
import { required, minLength } from '@vuelidate/validators';
import { NButton, NInput } from 'naive-ui';
import type { FormValidationStatus } from 'naive-ui';
import { computed, shallowReactive, shallowRef } from 'vue';

const params = shallowReactive<PasswordLoginParams>({
  user_name: '',
  password: '',
});

const validators: Record<keyof PasswordLoginParams, Record<string, any>> = {
  user_name: { required, minLength: minLength(5) },
  password: { required, minLength: minLength(7) },
};

const { passwordLogin } = useAuth();
const v$ = useVuelidate(validators, params);

const loading = shallowRef<boolean>(false);

const userNameInvalidMessage = computed<string | null>(() => {
  if (isFieldInvalid('user_name')) {
    if (v$.value.user_name.required.$invalid) {
      return 'Поле не заполнено';
    }

    return 'Минимальная длина 5 символов';
  }

  return null;
});
const passwordInvalidMessage = computed<string | null>(() => {
  if (isFieldInvalid('password')) {
    if (v$.value.password.required.$invalid) {
      return 'Поле не заполнено';
    }

    return 'Минимальная длина 7 символов';
  }

  return null;
});
const userNameInputStatus = computed<FormValidationStatus>(() => getFieldValidationStatus('user_name'));
const passwordInputStatus = computed<FormValidationStatus>(() => getFieldValidationStatus('password'));

function isFieldInvalid(fieldName: keyof PasswordLoginParams): boolean {
  return v$.value[fieldName]?.$dirty && v$.value[fieldName]?.$invalid;
}

function getFieldValidationStatus(fieldName: keyof PasswordLoginParams): FormValidationStatus {
  return isFieldInvalid(fieldName) ? 'error' : 'success';
}

async function onLoginSubmit() {
  if (loading.value) {
    return;
  }

  const valid = await v$.value.$validate();

  if (!valid) {
    return;
  }

  v$.value.$reset();

  await passwordLogin({ ...params });
};
</script>

<template>
  <form
    type="submit"
    class="password-auth"
    @submit.prevent="onLoginSubmit"
  >
    <div class="password-auth__input">
      <NInput
        v-model:value="params.user_name"
        type="text"
        placeholder="Юсернейм в Telegram"
        :status="userNameInputStatus"
      />

      <span
        v-if="userNameInvalidMessage"
        class="password-auth__input_message"
      >
        {{ userNameInvalidMessage }}
      </span>
    </div>

    <div class="password-auth__input">
      <NInput
        v-model:value="params.password"
        type="password"
        placeholder="Пароль"
        :status="passwordInputStatus"
        show-password-on="click"
      />

      <span
        v-if="passwordInvalidMessage"
        class="password-auth__input_message"
      >
        {{ passwordInvalidMessage }}
      </span>
    </div>

    <NButton
      class="password-auth__submit-btn"
      type="primary"
      attr-type="submit"
    >
      Войти
    </NButton>
  </form>
</template>

<style lang="scss" src="./style.scss" />
