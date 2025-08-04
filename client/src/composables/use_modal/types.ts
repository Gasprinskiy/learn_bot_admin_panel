import type { Component } from 'vue';

export interface UseModalState {
  component: Component | null;
  props?: any;
  emits?: any;
  width: number;
}
