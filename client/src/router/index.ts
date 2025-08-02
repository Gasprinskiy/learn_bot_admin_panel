import { createRouter, createWebHistory } from 'vue-router';

import AuthView from '@/views/auth/AuthView.vue';

import TgAuthView from '@/views/auth/children/telegram/TgAuthView.vue';
import PasswordAuthView from '@/views/auth/children/password_auth/PasswordAuthView.vue';
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
          name: 'password-auth',
          path: 'password',
          component: PasswordAuthView,
          meta: {
            hasBackAction: true,
          },
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
    ...RoutesByAccessRightList.teacher_access,
  ],
});

export default router;
