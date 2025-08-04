import type { Component } from 'vue';

export interface ProtectedRoute {
  name: string;
  path: string;
  icon: Component;
  text: string;
}

export interface UseProtectedRoutesState {
  protecredRoutes: Array<ProtectedRoute>;
}
