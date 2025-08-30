<script lang="ts" setup>
import { computed, h } from 'vue';
import type { Component, VNodeChild } from 'vue';
import { NButton, NDropdown, NIcon } from 'naive-ui';
import type { DropdownMixedOption } from 'naive-ui/es/dropdown/src/interface';

import { DoneFilled, PersonSharp, SettingsSharp } from '@vicons/material';
import { MoonOutline, SunnyOutline } from '@vicons/ionicons5';

import { useAppTheme } from '@/composables/use_app_theme';
import { useAuth } from '@/composables/use_auth';
import { usePerformanceSettings } from '@/composables/use_performance_settings';
import type { VisualMode } from '@/composables/use_performance_settings/types';

const { isAuthorized } = useAuth();
const { isDarkTheme, toggleTheme } = useAppTheme();
const { VisualModeOptions, currentMode, setMode } = usePerformanceSettings();

const themeIcon = computed<Component>(() => isDarkTheme.value ? SunnyOutline : MoonOutline);
const visualModeOptionsView = computed<Array<DropdownMixedOption>>(() => {
  return VisualModeOptions.map(({ value, label }) => {
    return {
      label,
      key: value,
      icon: () => renderVisualModeOptionIcon(value),
    };
  });
});

function renderVisualModeOptionIcon(value: VisualMode): VNodeChild | undefined {
  if (value !== currentMode.value) {
    return;
  }

  return h(NIcon, null, {
    default: () => h(DoneFilled),
  });
}
</script>

<template>
  <div class="app-header">
    <div class="app-header__inner">
      <div class="app-header__logo">
        <RouterLink v-if="isAuthorized" to="/">
          <img
            src="@/assets/logo.png"
            alt="Logo"
          >
        </RouterLink>
        <img
          v-else
          src="@/assets/logo.png"
          alt="Logo"
        >
      </div>

      <div class="app-header__right">
        <RouterLink
          v-if="isAuthorized"
          class="app-header__link"
          to="/profile"
        >
          <NButton
            quaternary
            circle
            class="app-header__link-button"
          >
            <template #icon>
              <NIcon
                size="20"
                :component="PersonSharp"
              />
            </template>
          </NButton>
        </RouterLink>

        <NDropdown
          :options="visualModeOptionsView"
          trigger="click"
          @select="setMode"
        >
          <NButton
            quaternary
            circle
          >
            <template #icon>
              <NIcon
                size="20"
                :component="SettingsSharp"
              />
            </template>
          </NButton>
        </NDropdown>

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
