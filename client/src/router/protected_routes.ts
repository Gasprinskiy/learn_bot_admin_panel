import type { RouteRecordRaw } from 'vue-router';
import { AccessRight } from '@/shared/types/profile';
import { UserListTab } from '@/views/bot_users/types';

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
  name: 'staff',
  path: '/staff',
  meta: {
    relatedToAccessRights: true,
    accessRights: { [AccessRight.AccessRightFull]: true },
  },
  component: () => {
    return import('@/views/staff/StaffView.vue');
  },
};

const teacherTasksRoute: ReadonlyRecordRaw = {
  name: 'teacher_tasks',
  path: '/teacher_tasks',
  meta: {
    relatedToAccessRights: true,
    accessRights: {
      [AccessRight.AccessRightFull]: true,
      [AccessRight.AccessRightManager]: true,
      [AccessRight.AccessRightTeacher]: true,
    },
  },
  component: () => {
    return import('@/views/teacher_tasks/TeacherTaskView.vue');
  },
};

const usersListRoute: ReadonlyRecordRaw = {
  name: 'bot_users',
  path: '/bot_users',
  meta: {
    relatedToAccessRights: true,
    accessRights: {
      [AccessRight.AccessRightFull]: true,
      [AccessRight.AccessRightManager]: true,
    },
  },
  component: () => {
    return import('@/views/bot_users/UsersListView.vue');
  },
  children: [
    {
      path: 'registered',
      name: UserListTab.REGISTERED,
      component: () => {
        return import('@/views/bot_users/children/registered/RegisteredUsersView.vue');
      },
    },
    {
      path: 'un_registered',
      name: UserListTab.UNREGISTERED,
      component: () => {
        return import('@/views/bot_users/children/unregistered/UnregisteredUsersView.vue');
      },
    },
  ],
};

const videoContentRoute: ReadonlyRecordRaw = {
  name: 'video_content',
  path: '/video_content',
  meta: {
    relatedToAccessRights: true,
    accessRights: {
      [AccessRight.AccessRightFull]: true,
      [AccessRight.AccessRightManager]: true,
      [AccessRight.AccessRightTeacher]: true,
    },
  },
  component: () => {
    return import('@/views/video_content/VideoContent.vue');
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
    videoContentRoute,
  ],

  [AccessRight.AccessRightManager]: [
    teacherTasksRoute,
    usersListRoute,
    videoContentRoute,
  ],

  [AccessRight.AccessRightTeacher]: [
    videoContentRoute,
    teacherTasksRoute,
  ],
};
