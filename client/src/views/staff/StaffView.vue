<script setup lang="ts">
import { NButton, NDivider, NIcon, NPopconfirm, NTable, NTag, NTooltip, useMessage } from 'naive-ui';
import { computed, defineAsyncComponent, onBeforeMount, shallowRef } from 'vue';
import { dateToRuLocaleString } from '@/packages/chronos';
import { optionalResult } from '@/packages/words';
import { AddOutline, Pencil, TrashBin } from '@vicons/ionicons5';

import $api from '@/packages/api/client';
import { AccessRight } from '@/shared/types/profile';
import { useModal } from '@/composables/use_modal';
import type { CreateUserParam, UserStaffInfo } from '@/shared/types/profile';
import type { ConvertEmitType } from '@/packages/types';

import type { CreateRedactProfileEmits, CreateRedactProfileProps, ProfileFormState } from './components/create-redact-profile/types';
import { AccesRightTitleMap, StaffListErrorsMap } from './constants';

const message = useMessage();
const { showModal, closeModal } = useModal();

const data = shallowRef<UserStaffInfo[]>([]);
const redactingProfile = shallowRef<UserStaffInfo | null>(null);

const isRedactMode = computed<boolean>(() => redactingProfile.value !== null);

function openCreateRedactModal(redactInfo?: UserStaffInfo) {
  redactingProfile.value = redactInfo || null;

  const component = defineAsyncComponent(() => {
    return import('./components/create-redact-profile/CreateRedactProfile.vue');
  });

  const props: CreateRedactProfileProps = {
    form_data: redactInfo,
  };
  const emits: ConvertEmitType<CreateRedactProfileEmits> = {
    onSubmit: onSubmitProfileForm,
  };

  showModal({
    component,
    props,
    emits,
    width: 500,
  });
}

async function onSubmitProfileForm(state: ProfileFormState) {
  const params: CreateUserParam = {
    id: redactingProfile.value?.id || undefined,
    ...state,
    access_right: redactingProfile.value ? undefined : AccessRight.AccessRightManager,
  };

  try {
    const path = isRedactMode.value ? '/panel_users/redact' : '/panel_users/create';
    const method = isRedactMode.value ? 'PATCH' : 'POST';

    await $api(path, { method, params });
    await fetchStaffList();

    closeModal();
    message.success(isRedactMode.value ? 'Пользовтель отредактирован' : 'Пользователь добавлен');
  } catch (e) {
    const stauts = +(e as any).status || 500;
    message.error(StaffListErrorsMap[stauts]);
  }
}

async function deleteProfile(id: number) {
  try {
    await $api(`/panel_users/delete/${id}`, {
      method: 'DELETE',
    });

    message.warning('Пользователь удален');
  } catch (e) {
    const stauts = +(e as any).status || 500;
    message.error(StaffListErrorsMap[stauts]);
  }
}

async function fetchStaffList() {
  try {
    const response = await $api<UserStaffInfo[]>('/panel_users');
    data.value = response;
  } catch (e) {
    const stauts = +(e as any).status || 500;
    message.error(StaffListErrorsMap[stauts]);
  }
}

onBeforeMount(fetchStaffList);
</script>

<template>
  <div class="staff-view">
    <div class="staff-view__top">
      <h2 class="staff-view__head">
        Работники
      </h2>

      <NTooltip>
        <template #default>
          Добавить работника
        </template>

        <template #trigger>
          <NButton
            type="primary"
            @click="openCreateRedactModal()"
          >
            <template #icon>
              <NIcon
                :component="AddOutline"
              />
            </template>
          </NButton>
        </template>
      </NTooltip>
    </div>

    <NDivider class="staff-view__divider" />

    <div class="staff-view__table">
      <NTable :single-line="false">
        <thead>
          <tr>
            <th>Имя</th>
            <th>Фамилия</th>
            <th>Юзернейм</th>
            <th>Доступ</th>
            <th>Дата последнего входа</th>
            <th />
          </tr>
        </thead>

        <tbody>
          <tr
            v-for="user in data"
            :key="user.id"
            class="staff-view__table-item"
          >
            <td>{{ user.first_name }}</td>
            <td>{{ user.last_name }}</td>
            <td>{{ user.tg_user_name }}</td>
            <td>{{ AccesRightTitleMap[user.access_right] }}</td>
            <td>{{ optionalResult(user.last_login!, dateToRuLocaleString) }}</td>
            <td class="staff-view__table-actions">
              <NTag
                v-if="user.is_you"
                type="primary"
                :bordered="false"
              >
                Вы
              </NTag>

              <template v-else>
                <NTooltip>
                  <template #default>
                    Редактировать
                  </template>

                  <template #trigger>
                    <NButton
                      type="primary"
                      @click="openCreateRedactModal(user)"
                    >
                      <template #icon>
                        <NIcon
                          :component="Pencil"
                        />
                      </template>
                    </NButton>
                  </template>
                </NTooltip>

                <NTooltip>
                  <template #default>
                    Удалить
                  </template>

                  <template #trigger>
                    <NPopconfirm
                      positive-text="Да"
                      negative-text="Отмена"
                      @positive-click="deleteProfile(user.id)"
                    >
                      <template #trigger>
                        <NButton
                          type="error"
                        >
                          <template #icon>
                            <NIcon
                              :component="TrashBin"
                            />
                          </template>
                        </NButton>
                      </template>

                      <template #default>
                        Вы уверены что хотите удалить пользователя?
                      </template>
                    </NPopconfirm>
                  </template>
                </NTooltip>
              </template>
            </td>
          </tr>
        </tbody>
      </NTable>
    </div>
  </div>
</template>

<style lang="scss" src="./style.scss" />
