import { createRouter, createWebHistory } from 'vue-router';

import AuthView from '@/views/auth/AuthView.vue';
import SetPassView from '@/views/set_pass/SetPassView.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/auth',
      component: AuthView,
      children: [
        { path: '/set_pass', component: SetPassView },
      ],
    },
  ],
});

export default router;
