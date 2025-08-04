import { AccessRight } from '@/shared/types/profile';
import type { ProtectedRoute } from './types';
import { BookRound, PlayLessonRound, SupervisedUserCircleRound, WorkRound } from '@vicons/material';

const createProfileRoute: ProtectedRoute = {
  name: 'staff',
  path: '/staff',
  icon: WorkRound,
  text: 'Работники',
};

const teacherTasksRoute: ProtectedRoute = {
  name: 'teacher_tasks',
  path: '/teacher_tasks',
  icon: BookRound,
  text: 'Домашние задания',
};

const usersListRoute: ProtectedRoute = {
  name: 'users_list',
  path: '/users_list',
  icon: SupervisedUserCircleRound,
  text: 'Пользователи',
};

const videoContentRoute: ProtectedRoute = {
  name: 'video_content',
  path: '/video_content',
  icon: PlayLessonRound,
  text: 'Видео уроки',
};

export const RoutesByAccessRight: Record<AccessRight, Array<ProtectedRoute>> = {
  [AccessRight.AccessRightFull]: [createProfileRoute, teacherTasksRoute, usersListRoute, videoContentRoute],
  [AccessRight.AccessRightManager]: [usersListRoute, teacherTasksRoute, videoContentRoute],
  [AccessRight.AccessRightTeacher]: [teacherTasksRoute, videoContentRoute],
};
