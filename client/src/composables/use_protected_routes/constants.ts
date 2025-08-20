import { AccessRight } from '@/shared/types/profile';
import type { ProtectedRoute } from './types';
import { BookRound, PlayLessonRound, SupervisedUserCircleRound, WorkRound } from '@vicons/material';

const createProfileRoute: ProtectedRoute = {
  name: 'staff',
  path: '/staff',
  icon: WorkRound,
  text: 'Работники',
  disabled: true,
};

const teacherTasksRoute: ProtectedRoute = {
  name: 'teacher_tasks',
  path: '/teacher_tasks',
  icon: BookRound,
  text: 'Домашние задания',
  disabled: true,
};

const usersListRoute: ProtectedRoute = {
  name: 'bot_users',
  path: '/bot_users/registered',
  icon: SupervisedUserCircleRound,
  text: 'Пользователи',
};

const videoContentRoute: ProtectedRoute = {
  name: 'video_content',
  path: '/video_content',
  icon: PlayLessonRound,
  text: 'Видео уроки',
  disabled: true,
};

export const RoutesByAccessRight: Record<AccessRight, Array<ProtectedRoute>> = {
  [AccessRight.AccessRightFull]: [createProfileRoute, usersListRoute, videoContentRoute, teacherTasksRoute],
  [AccessRight.AccessRightManager]: [usersListRoute, videoContentRoute, teacherTasksRoute],
  [AccessRight.AccessRightTeacher]: [videoContentRoute, teacherTasksRoute],
};
