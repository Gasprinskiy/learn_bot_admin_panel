import type { RouteRecordRaw } from 'vue-router';
//
import { AccessRight } from '@/shared/types/profile';

type RoutesByAccessRight = {
  [key in AccessRight]: Readonly<RouteRecordRaw[]>
};

const createProfileRoute: Readonly<RouteRecordRaw> = {
  name: 'create_profile',
  path: '/create_profle',
  component: () => {
    return import('@/views/create_profile/CreateProfileView.vue');
  },
};

const teacherTasksRoute: Readonly<RouteRecordRaw> = {
  name: 'teacher_tasks',
  path: '/teacher_tasks',
  component: () => {
    return import('@/views/create_profile/CreateProfileView.vue');
  },
};

const usersListRoute: Readonly<RouteRecordRaw> = {
  name: 'users_list',
  path: '/users_list',
  component: () => {
    return import('@/views/users_list/UsersListView.vue');
  },
};

export const exclusiveRoutesByAccessRight: Readonly<RoutesByAccessRight> = {
  [AccessRight.AccessRightFull]: [
    createProfileRoute,
    teacherTasksRoute,
    usersListRoute,
  ],

  [AccessRight.AccessRightManager]: [
    teacherTasksRoute,
    usersListRoute,
  ],

  [AccessRight.AccessRightTeacher]: [
    {
      ...teacherTasksRoute,
      meta: {
        isTeacher: true,
      },
    },
  ],
};
