import { createRouter, createWebHistory } from 'vue-router';

import AuthView from '@/views/auth/AuthView.vue';

import TgAuthView from '@/views/auth/children/telegram/TgAuthView.vue';
import StandartAuthView from '@/views/auth/children/standart_auth/StandartAuthView.vue';
import SetPassView from '@/views/auth/children/set_pass/SetPassView.vue';
import { DefaultRoutes, RoutesByAccessRightList } from './protected_routes';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/auth',
      name: 'auth',
      component: AuthView,
      children: [
        {
          name: 'telegram-auth',
          path: 'telegram',
          component: TgAuthView,
          meta: {
            hasBackAction: true,
          },
        },
        {
          name: 'standart-auth',
          path: 'standart',
          component: StandartAuthView,
          meta: {
            hasBackAction: true,
          },
        },
        {
          path: 'set_pass',
          component: SetPassView,
          name: 'set-pass',
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/not_found/NotFoundView.vue'),
    },
    ...DefaultRoutes,
    ...RoutesByAccessRightList.full_access,
    ...RoutesByAccessRightList.manager_access,
    ...RoutesByAccessRightList.teacher_access
  ],
});

export default router;
