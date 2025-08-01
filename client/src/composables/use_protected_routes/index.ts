import { reactive, toRefs } from 'vue';

import { RoutesByAccessRightList } from '@/router/protected_routes';
import type { AccessRight } from '@/shared/types/profile';

import type { UseProtectedRoutesState } from './types';

const state = reactive<UseProtectedRoutesState>({
  AccessRightRoutes: [],
});

export function useProtectedRoutes() {
  async function buildRoutesByAccessRight(accRight: AccessRight) {
    const routesToAdd = RoutesByAccessRightList[accRight] || [];
    if (!routesToAdd) {
      console.error('no routes found by acces right');
    }

    for (const route of [...routesToAdd]) {
      state.AccessRightRoutes = [
        ...state.AccessRightRoutes,
        {
          name: route.name?.toString(),
          path: route.path,
        },
      ];
    }
  }

  return {
    ...toRefs(state),
    buildRoutesByAccessRight,
  };
}
