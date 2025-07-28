import { reactive, toRefs } from 'vue';
import { useRouter } from 'vue-router';

import { DefaultRoutes, RoutesByAccessRightList } from '@/router/protected_routes';
import type { AccessRight } from '@/shared/types/profile';

import type { UseProtectedRoutesState } from './types';

const state = reactive<UseProtectedRoutesState>({
  AccessRightRoutes: [],
});

export function useProtectedRoutes() {
  const router = useRouter();

  async function buildRoutesByAccessRight(accRight: AccessRight) {
    const routesToAdd = RoutesByAccessRightList[accRight];
    if (!routesToAdd) {
      console.error('no routes found by acces right');
    }

    for (const route of [...routesToAdd, ...DefaultRoutes]) {
      router.addRoute(route);

      if (route.meta?.relatedToAccessRights) {
        state.AccessRightRoutes = [
          ...state.AccessRightRoutes,
          {
            name: route.name?.toString(),
            path: route.path,
          },
        ];
      }
    }

    await router.replace(router.currentRoute.value.fullPath);
  }

  return {
    ...toRefs(state),
    buildRoutesByAccessRight,
  };
}
