import { createRouter, createWebHistory } from 'vue-router';

import AuthView from '@/views/auth/AuthView.vue';
import SetPassView from '@/views/set_pass/SetPassView.vue';
import NotFoundView from '@/views/not_found/NotFoundView.vue';
import TgAuthView from '@/views/auth/children/telegram/TgAuthView.vue';

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
          component: TgAuthView,
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
      path: '/:pathMatch(.*)*', // Catch-all route for 404
      name: 'NotFound',
      component: NotFoundView,
    },
  ],
});

export default router;
