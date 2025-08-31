import type { Component } from 'vue';
import type { AccessRight } from '@/shared/types/profile';

export interface ProtectedRoute {
  name: string;
  path: string;
  icon: Component;
  text: string;
  access_rights: Array<AccessRight>;
  disabled: boolean;
}

export interface UseProtectedRoutesState {
  protecredRoutes: Array<ProtectedRoute>;
}
