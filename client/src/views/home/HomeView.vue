<script setup lang="ts">
import { useProtectedRoutes } from '@/composables/use_protected_routes';
import { NButton, NIcon, NTooltip } from 'naive-ui';

const { protecredRoutes } = useProtectedRoutes();
</script>

<template>
  <div class="home-view">
    <div class="home-view__wrapper">
      <RouterLink
        v-for="route in protecredRoutes"
        :key="route.name"
        :to="route.disabled ? '' : route.path"
        class="home-view__route"
      >
        <NTooltip
          v-if="route.disabled"
        >
          <template #default>
            Временно не доступно
          </template>

          <template #trigger>
            <NButton
              class="home-view__route-button"
              type="primary"
              size="large"
              tertiary
              :disabled="true"
            >
              <template #icon>
                <NIcon :component="route.icon" />
              </template>

              <template #default>
                {{ route.text }}
              </template>
            </NButton>
          </template>
        </NTooltip>

        <NButton
          v-else
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
