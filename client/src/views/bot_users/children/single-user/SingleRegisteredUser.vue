<script lang="ts" setup>
import { computed, defineAsyncComponent, onBeforeMount, shallowRef } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { NButton, NTable, NTag } from 'naive-ui';

import $api from '@/packages/api/client';
import { SubscriptionStatus } from '@/shared/types/bot_users';
import type { BotSubscriptionType, BotUserProfile } from '@/shared/types/bot_users';
import { dateToRuLocaleString } from '@/packages/chronos';
import { useModal } from '@/composables/use_modal';

import { UserListTab } from '../../types';
import { SubscriptionStatusTitleMap, SubscriptionUIStatus } from '../registered/constants';
import type { CreatePurchaseEmits, CreatePurchaseProps } from './components/create_purchase/types';
import type { ConvertEmitType } from '@/packages/types';

const route = useRoute();
const router = useRouter();
const { showModal } = useModal();

const user = shallowRef<BotUserProfile | null>(null);

const hasActiveSub = computed<boolean>(() => user.value?.subscription_status === SubscriptionStatus.ACTIVE);

async function getUserByID() {
  try {
    const response = await $api(`/bot_users/${route.params.id}`);
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

  const response = await $api(`/bot_users/purchase/${user.value?.u_id}`, {
    method: 'POST',
    body: formData,
  });

  console.log('response: ', response);
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

onBeforeMount(getUserByID);
</script>

<template>
  <div
    v-if="user"
    class="single-user"
  >
    <div class="single-user__head">
      <h2>{{ user.first_name }} {{ user.last_name }}</h2>

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

        <tr>
          <td>Подписка</td>
          <td>
            <NTag
              :type="SubscriptionUIStatus[user.subscription_status]"
              :bordered="false"
            >
              {{ SubscriptionStatusTitleMap[user.subscription_status] }}
            </NTag>
          </td>
        </tr>
      </tbody>
    </NTable>
  </div>
</template>

<style lang="scss" src="./style.scss" />
