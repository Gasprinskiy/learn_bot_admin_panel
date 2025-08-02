<script setup lang="ts">
import { useAuth } from '@/composables/use_auth';
import $api from '@/packages/api/client';
import type { User } from '@/shared/types/profile';
import { onBeforeMount, shallowRef } from 'vue';

const { getUserInfo } = useAuth();

const userInfo = shallowRef<User | null>(null);
const dataLoading = shallowRef<boolean>(false);

async function setPass() {
  await $api('/auth/create_password', {
    method: 'POST',
    body: {
      password: '{Somepass12}',
    },
  });
}

onBeforeMount(async () => {
  dataLoading.value = true;
  userInfo.value = await getUserInfo();
  dataLoading.value = false;
});
</script>

<template>
  <div>{{ userInfo }}</div>
  <button @click="setPass">
    set pass
  </button>
</template>

<style scoped>

</style>
