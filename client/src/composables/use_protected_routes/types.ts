import type { Component } from 'vue';

export interface ProtectedRoute {
  name: string;
  path: string;
  icon: Component;
  text: string;
  disabled?: boolean;
}

export interface UseProtectedRoutesState {
  protecredRoutes: Array<ProtectedRoute>;
}
