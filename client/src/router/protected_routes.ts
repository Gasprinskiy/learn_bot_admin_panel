import type { RouteRecordRaw } from 'vue-router';
import { AccessRight } from '@/shared/types/profile';

//
type RoutesByAccessRight = {
  [key in AccessRight]: Readonly<RouteRecordRaw[]>
};

type ReadonlyRecordRaw = Readonly<RouteRecordRaw>;
type ReadonlyRoutesByAccessRight = Readonly<RoutesByAccessRight>;
//
const homeRoute: ReadonlyRecordRaw = {
  name: 'home',
  path: '/',
  component: () => {
    return import('@/views/home/HomeView.vue');
  },
};

const profileRoute: ReadonlyRecordRaw = {
  name: 'profile',
  path: '/profile',
  component: () => {
    return import('@/views/profile/ProfileView.vue');
  },
};

const createProfileRoute: ReadonlyRecordRaw = {
  name: 'create_profile',
  path: '/create_profle',
  meta: {
    relatedToAccessRights: true,
  },
  component: () => {
    return import('@/views/create_profile/CreateProfileView.vue');
  },
};

const teacherTasksRoute: ReadonlyRecordRaw = {
  name: 'teacher_tasks',
  path: '/teacher_tasks',
  meta: {
    relatedToAccessRights: true,
  },
  component: () => {
    return import('@/views/create_profile/CreateProfileView.vue');
  },
};

const usersListRoute: ReadonlyRecordRaw = {
  name: 'users_list',
  path: '/users_list',
  meta: {
    relatedToAccessRights: true,
  },
  component: () => {
    return import('@/views/users_list/UsersListView.vue');
  },
};
//
export const DefaultRoutes: ReadonlyRecordRaw[] = [
  homeRoute,
  profileRoute,
];

export const RoutesByAccessRightList: ReadonlyRoutesByAccessRight = {
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
