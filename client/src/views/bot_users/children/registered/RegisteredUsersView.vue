<script setup lang="ts">
import { NButton, NDatePicker, NIcon, NInput, NInputNumber, NSelect, NTable, NTag, NTooltip } from 'naive-ui';
import { RefreshOutlined, SearchOutlined } from '@vicons/material';
import { DocumentTextOutline } from '@vicons/ionicons5';
import { computed, onBeforeMount } from 'vue';

import type { SubscriptionStatus } from '@/shared/types/profile';

import { SubscriptionStatusSelectOptions, SubscriptionStatusTitleMap, SubscriptionUIStatus } from './constants';
import { useUsersList } from '../../composables/use_users_list';

const {
  data,
  searchParams,
  showLoadMoreButton,
  isDataLeft,
  isLoadingMore,
  fetchRegisteredUsers,
  loadMoreRegisteredUsers,
  resetSearchParams,
  printRegisteredUsers,
} = useUsersList();

const searchQuery = computed<string>({
  get(): string {
    return searchParams.query || '';
  },

  set(value: string) {
    searchParams.query = value;
  },
});
const subStatus = computed<SubscriptionStatus | null>({
  get(): SubscriptionStatus | null {
    return searchParams.subscription_status || null;
  },

  set(value: SubscriptionStatus | null) {
    if (!value) {
      searchParams.subscription_status = undefined;
      return;
    }
    searchParams.subscription_status = value;
  },
});
const joinDateFrom = computed<number | null>({
  get(): number | null {
    if (!searchParams.join_date_from) {
      return null;
    }
    const date = new Date(searchParams.join_date_from);
    return date.getTime();
  },

  set(value: number | null) {
    if (!value) {
      searchParams.join_date_from = undefined;
      return;
    }
    searchParams.join_date_from = new Date(value).toISOString();
  },
});
const joinDateTill = computed<number | null>({
  get(): number | null {
    if (!searchParams.join_date_till) {
      return null;
    }
    const date = new Date(searchParams.join_date_till);
    return date.getTime();
  },

  set(value: number | null) {
    if (!value) {
      searchParams.join_date_till = undefined;
      return;
    }
    searchParams.join_date_till = new Date(value).toISOString();
  },
});
const ageFrom = computed<number | null>({
  get(): number | null {
    return searchParams.age_from || null;
  },

  set(value: number | null) {
    if (!value) {
      searchParams.age_from = undefined;
      return;
    }
    searchParams.age_from = value;
  },
});
const ageTill = computed<number | null>({
  get(): number | null {
    return searchParams.age_till ? searchParams.age_till : null;
  },

  set(value: number | null) {
    if (!value) {
      searchParams.age_till = undefined;
      return;
    }
    searchParams.age_till = value;
  },
});

function dateToRuLocaleString(date: string): string {
  return new Date(date).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
}

async function onReset() {
  resetSearchParams();
  await fetchRegisteredUsers();
}

onBeforeMount(fetchRegisteredUsers);
</script>

<template>
  <div class="registered-users-view">
    <form
      class="registered-users-view__filters"
      @submit.prevent="fetchRegisteredUsers(true)"
    >
      <div class="registered-users-view__filters_fileds">
        <div
          class="registered-users-view__filters_search-input"
        >
          <NInput
            v-model:value="searchQuery"
            placeholder="Юзернейм, имя или фамилия"
          />
        </div>

        <div class="registered-users-view__filters_sub-select">
          <NSelect
            v-model:value="subStatus"
            :options="SubscriptionStatusSelectOptions"
            placeholder="Статус подписки"
          />
        </div>

        <div class="registered-users-view__filters_join-date-range">
          <NDatePicker
            v-model:value="joinDateFrom"
            placeholder="Дата вступления: От"
          />
          <NDatePicker
            v-model:value="joinDateTill"
            placeholder="Дата вступления: До"
          />
        </div>

        <div class="registered-users-view__filters_age-range">
          <NInputNumber
            v-model:value="ageFrom"
            placeholder="Возраст: От"
          />
          <NInputNumber
            v-model:value="ageTill"
            placeholder="Возраст: До"
          />
        </div>
      </div>

      <div class="registered-users-view__filters-buttons">
        <NTooltip>
          <template #default>
            Поиск
          </template>

          <template #trigger>
            <NButton
              attr-type="submit"
              type="primary"
            >
              <template #icon>
                <NIcon
                  :component="SearchOutlined"
                />
              </template>
            </NButton>
          </template>
        </NTooltip>

        <NTooltip>
          <template #default>
            Сброс фильтров
          </template>

          <template #trigger>
            <NButton
              attr-type="reset"
              type="warning"
              @click="onReset"
            >
              <template #icon>
                <NIcon
                  :component="RefreshOutlined"
                />
              </template>
            </NButton>
          </template>
        </NTooltip>

        <NTooltip>
          <template #default>
            Excel выгрузка
          </template>

          <template #trigger>
            <NButton
              type="info"
              @click="printRegisteredUsers"
            >
              <template #icon>
                <NIcon
                  :component="DocumentTextOutline"
                />
              </template>
            </NButton>
          </template>
        </NTooltip>
      </div>
    </form>
    <div class="registered-users-view__list">
      <NTable :single-line="false">
        <thead>
          <tr>
            <th>Имя</th>
            <th>Фамилия</th>
            <th>Юзернейм</th>
            <th>Дата рождения</th>
            <th>Дата вступления</th>
            <th>Подписка</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="user in data"
            :key="user.u_id"
          >
            <td>{{ user.first_name }}</td>
            <td>{{ user.last_name }}</td>
            <td>{{ user.tg_user_name }}</td>
            <td>
              {{ dateToRuLocaleString(user.birth_date) }}
            </td>
            <td>
              {{ dateToRuLocaleString(user.join_date) }}
            </td>
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
      <div
        class="registered-users-view__list_bottom"
      >
        <NButton
          v-if="showLoadMoreButton"
          :disabled="!isDataLeft || isLoadingMore"
          :loading="isLoadingMore"
          type="primary"
          @click="loadMoreRegisteredUsers"
        >
          Загрузить еще
        </NButton>
      </div>
    </div>
  </div>
</template>

<style lang="scss" src="./style.scss" />
