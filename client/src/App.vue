<script setup lang="ts">
import { exclusiveRoutesByAccessRight } from '@/router/exclusive_routes';
import { AccessRight } from '@/shared/types/profile';
import { useRouter } from 'vue-router';

const router = useRouter();

async function addFullAcc() {
  const routesToAdd = exclusiveRoutesByAccessRight[AccessRight.AccessRightManager];
  if (routesToAdd.length <= 0) {
    console.error('no routes found by acces right');
    return;
  }

  for (const route of routesToAdd) {
    router.addRoute(route);
  }

  await router.replace(router.currentRoute.value.fullPath);

  if (!router.hasRoute('users_list')) {
    return;
  }

  await router.push({
    name: 'users_list',
  });
}
</script>

<template>
  <RouterView />
  <button @click="addFullAcc">
    FUCK
  </button>
</template>

<style scoped>

</style>
