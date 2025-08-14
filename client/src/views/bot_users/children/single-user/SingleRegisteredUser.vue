<script lang="ts" setup>
import $api from '@/packages/api/client';
import type { BotUserProfile } from '@/shared/types/profile';
import { onBeforeMount, shallowRef } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { UserListTab } from '../../types';
import { NTable } from 'naive-ui';

const route = useRoute();
const router = useRouter();

const user = shallowRef<BotUserProfile | null>(null);

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

onBeforeMount(getUserByID);
</script>

<template>
  <div
    v-if="user"
    class="single-user"
  >
    <h2>{{ user.first_name }} {{ user.last_name }}</h2>
    <NTable>
      <thead>
        <tr>
          <th />
          <th />
        </tr>
      </thead>

      <tbody>
        <tr>
          <!-- <td>Имя Фамилия</td> -->
          <!-- <td>{{ user.first_name }} {{ user.last_name }}</td> -->
        </tr>
      </tbody>
    </NTable>
  </div>
</template>
