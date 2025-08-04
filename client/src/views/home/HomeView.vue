<script setup lang="ts">
import { useProtectedRoutes } from '@/composables/use_protected_routes';
import { NButton, NIcon } from 'naive-ui';

const { protecredRoutes } = useProtectedRoutes();
</script>

<template>
  <div class="home-view">
    <div class="home-view__wrapper">
      <RouterLink
        v-for="route in protecredRoutes"
        :key="route.name"
        :to="route.path"
        class="home-view__route"
      >
        <NButton
          class="home-view__route-button"
          type="primary"
          size="large"
          tertiary
        >
          <template #icon>
            <NIcon :component="route.icon" />
          </template>

          <template #default>
            {{ route.text }}
          </template>
        </NButton>
      </RouterLink>
    </div>
  </div>
</template>

<style lang="scss" src="./style.scss" />
