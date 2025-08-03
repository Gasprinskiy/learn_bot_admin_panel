<script lang="ts" setup>
import { computed } from 'vue';
import type { Component } from 'vue';
import { darkTheme, lightTheme, NButton, NIcon } from 'naive-ui';
import { MoonOutline, SunnyOutline } from '@vicons/ionicons5';

import { useAppTheme } from '@/composables/use_app_theme';
import { useAuth } from '@/composables/use_auth';

const { currentTheme, setCurrentTheme } = useAppTheme();
const { isAuthorized } = useAuth();

const isDarkTheme = computed<boolean>(() => currentTheme.value === darkTheme);
const themeIcon = computed<Component>(() => isDarkTheme.value ? SunnyOutline : MoonOutline);

function toggleTheme() {
  const theme = isDarkTheme.value ? lightTheme : darkTheme;
  setCurrentTheme(theme);
}
</script>

<template>
  <div class="app-header">
    <div class="app-header__inner">
      <div class="app-header__logo">
        <img
          src="@/assets/vue.svg"
          alt="Logo"
        >
      </div>

      <div class="app-header__right">
        <RouterLink
          v-if="isAuthorized"
          class="app-header__link"
          to="/profile"
        >
          Профиль
        </RouterLink>

        <NButton
          quaternary
          circle
          @click="toggleTheme"
        >
          <template #icon>
            <NIcon
              size="20"
              :component="themeIcon"
            />
          </template>
        </NButton>
      </div>
    </div>
  </div>
</template>

<style lang="scss" src="./style.scss"  />
