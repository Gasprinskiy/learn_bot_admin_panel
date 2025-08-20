<script lang="ts" setup>
import { computed, defineAsyncComponent, onBeforeMount, shallowRef } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ImageOutline } from '@vicons/ionicons5';
import { NButton, NIcon, NImage, NTable, NTag, useMessage } from 'naive-ui';

import $api from '@/packages/api/client';
import { SubscriptionStatus } from '@/shared/types/bot_users';
import type { BotSubscriptionType, BotUserDetailData } from '@/shared/types/bot_users';
import type { ResponseWithBoolStatus } from '@/shared/types/common';
import { dateToRuLocaleString } from '@/packages/chronos';
import { useModal } from '@/composables/use_modal';
import type { ConvertEmitType } from '@/packages/types';
import { pluralize } from '@/packages/words';
import { useConfig } from '@/composables/use_config';

import type { CreatePurchaseEmits, CreatePurchaseProps } from './components/create_purchase/types';
import { PaymentTypeNameByID, UploadErrorMessagesByCode } from './constants';
import { SubscriptionStatusTitleMap, SubscriptionUIStatus } from '../registered/constants';
import { UserListTab } from '../../types';

const route = useRoute();
const router = useRouter();
const message = useMessage();
const { showModal, closeModal } = useModal();
const { UploadsURL } = useConfig();

const user = shallowRef<BotUserDetailData | null>(null);

const hasActiveSub = computed<boolean>(() => {
  const index = user.value?.purchase_data?.findIndex(({ subscription_status }) => subscription_status === SubscriptionStatus.ACTIVE);
  return index !== undefined && index >= 0;
});

async function getUserByID() {
  try {
    const response = await $api<BotUserDetailData>(`/bot_users/${route.params.id}`);
    user.value = response;
  } catch {
    await router.replace({
      name: UserListTab.REGISTERED,
    });
  }
}

async function getSubscriptionTypes(): Promise<BotSubscriptionType[] | null> {
  try {
    return $api<BotSubscriptionType[]>('/bot_users/subscr_types');
  } catch {
    return null;
  }
}

async function uploadPurchase(sub_id: number, file: File) {
  const formData = new FormData();
  formData.append('sub_id', sub_id.toString());
  formData.append('receipt', file);

  try {
    const response = await $api<ResponseWithBoolStatus>(`/bot_users/purchase/${user.value?.u_id}`, {
      method: 'POST',
      body: formData,
    });

    if (response.success) {
      await getUserByID();
      message.success('Подписка подключена');
      closeModal();
    } else {
      message.error(UploadErrorMessagesByCode[500]);
    }
  } catch (e) {
    const stauts = +(e as any).status || 500;
    message.error(UploadErrorMessagesByCode[stauts]);
  }
}

async function openCreatePurchaseModal() {
  const subscriptionTypes = await getSubscriptionTypes();
  if (subscriptionTypes === null) {
    return;
  }

  const component = defineAsyncComponent(() => {
    return import('./components/create_purchase/CreatePurchase.vue');
  });

  const props: CreatePurchaseProps = {
    subscriptionTypes,
  };

  const emits: ConvertEmitType<CreatePurchaseEmits> = {
    onSubmit: uploadPurchase,
  };

  showModal({
    component,
    props,
    emits,
    width: 500,
  });
}

function getSubscriptionName(term_in_month: number, price: number): string {
  const plural = pluralize(term_in_month, 'месяц', 'месяца', 'месяцев');
  return `${term_in_month} ${plural} за ${price}`;
}

onBeforeMount(getUserByID);
</script>

<template>
  <div
    v-if="user"
    class="single-user"
  >
    <div class="single-user__head">
      <h1>{{ user.first_name }} {{ user.last_name }}</h1>

      <div>
        <NButton
          v-if="!hasActiveSub"
          type="primary"
          @click="openCreatePurchaseModal"
        >
          Подключить подписку
        </NButton>
      </div>
    </div>

    <NTable :single-line="false">
      <thead>
        <tr>
          <th />
          <th />
        </tr>
      </thead>

      <tbody>
        <tr>
          <td>Юзернейм</td>
          <td>{{ user.tg_user_name }}</td>
        </tr>

        <tr>
          <td>Номер телефона</td>
          <td>{{ user.phone_number }}</td>
        </tr>

        <tr>
          <td>Дата рождения</td>
          <td>{{ dateToRuLocaleString(user.birth_date) }}</td>
        </tr>

        <tr>
          <td>Дата вступления</td>
          <td>{{ dateToRuLocaleString(user.join_date) }}</td>
        </tr>
      </tbody>
    </NTable>

    <div
      v-if="user.purchase_data"
      class="single-user__purchases"
    >
      <h2>История покупок</h2>

      <div class="single-user__purchases-info">
        <NTable
          v-for="purchase in user.purchase_data"
          :key="purchase.p_id"
          :single-line="false"
        >
          <thead>
            <tr>
              <th style="width: 49%;" />
              <th style="width: 50%;" />
            </tr>
          </thead>

          <tbody>
            <tr>
              <td>Тип подписки</td>
              <td>{{ getSubscriptionName(purchase.subscription_term, purchase.price) }}</td>
            </tr>

            <tr>
              <td>Дата подключения</td>
              <td>{{ dateToRuLocaleString(purchase.p_time) }}</td>
            </tr>

            <tr>
              <td>Оплачено</td>
              <td>{{ purchase.price - (purchase.discount || 0) }}</td>
            </tr>

            <tr v-if="purchase.discount">
              <td>Скидка</td>
              <td>{{ purchase.discount }}</td>
            </tr>

            <tr>
              <td>Способ оплаты</td>
              <td>{{ PaymentTypeNameByID[purchase.payment_type_id] }}</td>
            </tr>

            <tr v-if="purchase.manager_id !== null">
              <td>Менеджер принимавший оплату</td>
              <td>{{ purchase.manager_first_name }} {{ purchase.manager_last_name }}</td>
            </tr>

            <tr v-if="purchase.receipt_file_name !== null">
              <td>Квитанция об оплате</td>
              <td>
                <a
                  v-if="purchase.receipt_file_name.includes('.pdf')"
                  :href="`${UploadsURL}/${purchase.receipt_file_name}`"
                  target="_blank"
                >
                  pdf файл
                </a>
                <NImage
                  v-else
                  width="150"
                  :src="`${UploadsURL}/${purchase.receipt_file_name}`"
                  :show-toolbar="false"
                >
                  <template #error>
                    <NIcon
                      :size="150"
                      color="lightGrey"
                      :component="ImageOutline"
                    />
                  </template>
                </NImage>
              </td>
            </tr>

            <tr>
              <td>Статус</td>
              <td>
                <NTag
                  :type="SubscriptionUIStatus[purchase.subscription_status]"
                  :bordered="false"
                >
                  {{ SubscriptionStatusTitleMap[purchase.subscription_status] }}
                </NTag>
              </td>
            </tr>
          </tbody>
        </NTable>
      </div>
    </div>
  </div>
</template>

<style lang="scss" src="./style.scss" />
