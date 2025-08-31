<script lang="ts" setup>
import { NButton, NCard, NInput } from 'naive-ui';
import type { FormValidationStatus } from 'naive-ui';
import { computed, shallowReactive } from 'vue';
import { required } from '@vuelidate/validators';
import useVuelidate from '@vuelidate/core';
import type { CreateRedactProfileEmits, CreateRedactProfileProps, ProfileFormState } from './types';

const props = defineProps<CreateRedactProfileProps>();
const emit = defineEmits<CreateRedactProfileEmits>();

const state = shallowReactive<ProfileFormState>({
  first_name: props.form_data?.first_name || '',
  last_name: props.form_data?.last_name || '',
  tg_user_name: props.form_data?.tg_user_name || '',
});

const validators: Record<keyof ProfileFormState, Record<string, any>> = {
  first_name: { required },
  last_name: { required },
  tg_user_name: { required },
};

const v$ = useVuelidate(validators, state);

const isPropsSet = computed<boolean>(() => !!props.form_data);
const headerTitle = computed<string>(() => isPropsSet.value ? 'Редактировать пользователя' : 'Создать пользователя');
const submitTitle = computed<string>(() => isPropsSet.value ? 'Сохранить' : 'Создать');

const firstNameValidationStatus = computed<FormValidationStatus>(() => getFieldValidationStatus('first_name'));
const lastNameValidationStatus = computed<FormValidationStatus>(() => getFieldValidationStatus('last_name'));
const tgUserNameValidationStatus = computed<FormValidationStatus>(() => getFieldValidationStatus('tg_user_name'));
const firstNameInvalidMessage = computed<string | null>(() => getInvalidMessage('first_name'));
const lastNameInvalidMessage = computed<string | null>(() => getInvalidMessage('last_name'));
const tgUserNameInvalidMessage = computed<string | null>(() => getInvalidMessage('tg_user_name'));

function isFieldInvalid(fieldName: keyof ProfileFormState): boolean {
  return v$.value[fieldName]?.$dirty && v$.value[fieldName]?.$invalid;
}

function getFieldValidationStatus(fieldName: keyof ProfileFormState): FormValidationStatus {
  return isFieldInvalid(fieldName) ? 'error' : 'success';
}

function getInvalidMessage(fieldName: keyof ProfileFormState): string | null {
  if (isFieldInvalid(fieldName)) {
    return 'Поле не заполнено';
  }

  return null;
}

async function onSubmit() {
  const valid = await v$.value.$validate();
  if (!valid) {
    return;
  }

  v$.value.$reset();

  emit('onSubmit', { ...state });
}
</script>

<template>
  <NCard class="create-redact-profile">
    <template #header>
      {{ headerTitle }}
    </template>

    <template #default>
      <form
        class="create-redact-profile__form"
        @submit.prevent="onSubmit"
      >
        <div class="create-redact-profile__input">
          <NInput
            v-model:value="state.first_name"
            placeholder="Имя"
            :status="firstNameValidationStatus"
          />

          <span
            v-if="firstNameInvalidMessage"
            class="create-redact-profile__input_message"
          >
            {{ firstNameInvalidMessage }}
          </span>
        </div>

        <div class="create-redact-profile__input">
          <NInput
            v-model:value="state.last_name"
            placeholder="Фамилия"
            :status="lastNameValidationStatus"
          />

          <span
            v-if="lastNameInvalidMessage"
            class="create-redact-profile__input_message"
          >
            {{ lastNameInvalidMessage }}
          </span>
        </div>

        <div class="create-redact-profile__input">
          <NInput
            v-model:value="state.tg_user_name"
            placeholder="Юзернейм в телеграм (без @)"
            :status="tgUserNameValidationStatus"
          />

          <span
            v-if="tgUserNameInvalidMessage"
            class="create-redact-profile__input_message"
          >
            {{ tgUserNameInvalidMessage }}
          </span>
        </div>

        <NButton
          type="primary"
          attr-type="submit"
        >
          {{ submitTitle }}
        </NButton>
      </form>
    </template>
  </NCard>
</template>

<style lang="scss" src="./style.scss" />
