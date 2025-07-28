import { createRouter, createWebHistory } from 'vue-router';

import AuthView from '@/views/auth/AuthView.vue';
import SetPassView from '@/views/set_pass/SetPassView.vue';
import HomeView from '@/views/home/HomeView.vue';
import ProfileView from '@/views/profile/ProfileView.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: HomeView,
    },
    {
      path: '/auth',
      component: AuthView,
      children: [
        { path: '/set_pass', component: SetPassView },
      ],
    },
    {
      path: '/profile',
      component: ProfileView,
    },
  ],
});

export default router;
