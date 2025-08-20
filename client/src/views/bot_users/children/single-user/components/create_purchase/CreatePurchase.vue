<script lang="ts" setup>
import type { FormValidationStatus, UploadFileInfo } from 'naive-ui';
import { NButton, NCard, NIcon, NSelect, NText, NUpload, NUploadDragger } from 'naive-ui';
import { ArchiveOutline } from '@vicons/ionicons5';
import { required } from '@vuelidate/validators';
import useVuelidate from '@vuelidate/core';
import { computed, shallowRef } from 'vue';
import type { SelectMixedOption } from 'naive-ui/es/select/src/interface';

import type { CreatePurchaseEmits, CreatePurchaseProps } from './types';

import { pluralize } from '@/packages/words';

const props = defineProps<CreatePurchaseProps>();
const emit = defineEmits<CreatePurchaseEmits>();

const selectedType = shallowRef<number | null>(null);
const uploadedFile = shallowRef<UploadFileInfo | null>(null);

const validators = {
  type: { required },
  file: { required },
};

const v$ = useVuelidate(validators, {
  type: selectedType,
  file: uploadedFile,
});

const singleFileList = computed<UploadFileInfo[]>({
  get(): UploadFileInfo[] {
    return uploadedFile.value ? [uploadedFile.value] : [];
  },

  set(value: UploadFileInfo[]) {
    uploadedFile.value = value[value.length - 1];
  },
});

const subscriptionTypeOptions = computed<SelectMixedOption[]>(() => {
  return props.subscriptionTypes.map(({ sub_id, term_in_month, price }): SelectMixedOption => {
    const plural = pluralize(term_in_month, 'месяц', 'месяца', 'месяцев');
    return {
      label: `${term_in_month} ${plural} за ${price}`,
      value: sub_id,
    };
  });
});

const isTypeInvalid = computed<boolean>(() => isFieldInvalid('type'));
const isFileInvalid = computed<boolean>(() => isFieldInvalid('file'));
const typeValidationStatus = computed<FormValidationStatus>(() => getFieldValidationStatus('type'));
const fileValidationStatus = computed<FormValidationStatus>(() => getFieldValidationStatus('file'));

function isFieldInvalid(fieldName: string): boolean {
  return v$.value[fieldName]?.$dirty && v$.value[fieldName]?.$invalid;
}

function getFieldValidationStatus(fieldName: string): FormValidationStatus {
  return isFieldInvalid(fieldName) ? 'error' : 'success';
}

async function onSubmit() {
  const valid = await v$.value.$validate();
  if (!valid) {
    return;
  }

  emit('onSubmit', selectedType.value!, uploadedFile.value!.file!);
}
</script>

<template>
  <NCard class="create-purchase">
    <template #header>
      Подключить подписку
    </template>

    <template #default>
      <form
        class="create-purchase__form"
        @submit.prevent="onSubmit"
      >
        <div class="create-purchase__input">
          <NSelect
            v-model:value="selectedType"
            :options="subscriptionTypeOptions"
            :status="typeValidationStatus"
          />

          <span
            v-if="isTypeInvalid"
            class="create-purchase__input_message"
          >
            Тип подписки не выбран
          </span>
        </div>

        <div
          class="create-purchase__input"
          :class="fileValidationStatus"
        >
          <NUpload
            v-model:file-list="singleFileList"
            :multiple="false"
            accept=".jpg, .png, .pdf"
          >
            <NUploadDragger>
              <div style="margin-bottom: 12px">
                <NIcon
                  size="48"
                  :depth="3"
                  :component="ArchiveOutline"
                />
              </div>
              <NText style="font-size: 16px">
                Нажмите или перетащите файл квитанции об оплате в эту область для загрузки.
              </NText>
            </NUploadDragger>
          </NUpload>

          <span
            v-if="isFileInvalid"
            class="create-purchase__input_message"
          >
            Файл квитанции не выбран
          </span>
        </div>

        <div>
          <NButton
            type="primary"
            attr-type="submit"
          >
            Подключить
          </NButton>
        </div>
      </form>
    </template>
  </NCard>
</template>

<style lang="scss" src="./style.scss" />
