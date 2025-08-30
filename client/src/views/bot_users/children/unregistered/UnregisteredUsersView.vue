<script setup lang="ts">
import { NButton, NIcon, NTable, NTooltip } from 'naive-ui';
import { DocumentTextOutline } from '@vicons/ionicons5';
import { onBeforeMount } from 'vue';

import { dateToRuLocaleString } from '@/packages/chronos';

import { useUsersList } from '../../composables/use_users_list';

const {
  data,
  showLoadMoreButton,
  isDataLeft,
  isLoadingMore,
  fetchUsers,
  loadMoreUsers,
  printUsers,
} = useUsersList();

function optionalResult<T extends string | number | Date>(value: T, valueProxy?: (value: T) => string): string {
  if (value === null) {
    return '-';
  }

  return valueProxy ? valueProxy(value) : value.toString();
}

onBeforeMount(fetchUsers);
</script>

<template>
  <div class="unregistered-users-view">
    <form
      class="unregistered-users-view__filters"
    >
      <div class="unregistered-users-view__filters-buttons">
        <NTooltip>
          <template #default>
            Excel выгрузка
          </template>

          <template #trigger>
            <NButton
              type="info"
              @click="printUsers"
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

    <div class="unregistered-users-view__list">
      <NTable :single-line="false">
        <thead>
          <tr>
            <th>Юзернейм</th>
            <th>Имя</th>
            <th>Фамилия</th>
            <th>Дата рождения</th>
            <th>Дата вступления</th>
          </tr>
        </thead>

        <tbody>
          <tr
            v-for="user in data"
            :key="user.u_id"
            class="unregistered-users-view__table-item"
          >
            <td>{{ user.tg_user_name }}</td>
            <td>{{ optionalResult(user.first_name) }}</td>
            <td>{{ optionalResult(user.last_name) }}</td>
            <td>{{ optionalResult(user.birth_date, dateToRuLocaleString) }}</td>
            <td>
              {{ optionalResult(user.join_date, dateToRuLocaleString) }}
            </td>
          </tr>
        </tbody>
      </NTable>

      <div
        class="unregistered-users-view__list_bottom"
      >
        <NButton
          v-if="showLoadMoreButton"
          :disabled="!isDataLeft || isLoadingMore"
          :loading="isLoadingMore"
          type="primary"
          @click="loadMoreUsers"
        >
          Загрузить еще
        </NButton>
      </div>
    </div>
  </div>
</template>

<style lang="scss" src="./style.scss" />
