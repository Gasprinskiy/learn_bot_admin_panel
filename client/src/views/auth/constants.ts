import { defineAsyncComponent } from 'vue';
import type { Component } from 'vue';
import { AuthMethod } from './types';

export const AuthComponentMap: Record<AuthMethod, () => Component> = {
  [AuthMethod.STANDART]: () => {
    return defineAsyncComponent(() => {
      return import('./components/standart_auth/StandartAuth.vue');
    });
  },

  [AuthMethod.TELEGRAM]: () => {
    return defineAsyncComponent(() => {
      return import('./components/tg_auth/TgAuth.vue');
    });
  },
};
