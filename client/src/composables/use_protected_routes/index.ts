import { shallowReactive, toRefs } from 'vue';

import type { AccessRight } from '@/shared/types/profile';

import type { UseProtectedRoutesState } from './types';
import { ProtectedRoutes } from './constants';

const state = shallowReactive<UseProtectedRoutesState>({
  protecredRoutes: [],
});

export function useProtectedRoutes() {
  function setRoutesByAccesRight(accRight: AccessRight) {
    const routes = ProtectedRoutes.map((route) => {
      if (!route.access_rights.includes(accRight)) {
        route.disabled = true;
      }

      return route;
    });

    state.protecredRoutes = routes;
  }

  return {
    ...toRefs(state),
    setRoutesByAccesRight,
  };
}
