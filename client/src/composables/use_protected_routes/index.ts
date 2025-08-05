import { shallowReactive, toRefs } from 'vue';

import type { AccessRight } from '@/shared/types/profile';

import type { UseProtectedRoutesState } from './types';
import { RoutesByAccessRight } from './constants';

const state = shallowReactive<UseProtectedRoutesState>({
  protecredRoutes: [],
});

export function useProtectedRoutes() {
  function setRoutesByAccesRight(accRight: AccessRight) {
    state.protecredRoutes = RoutesByAccessRight[accRight] || [];
  }

  return {
    ...toRefs(state),
    setRoutesByAccesRight,
  };
}
