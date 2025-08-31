import { AccessRight } from '@/shared/types/profile';
import type { ProtectedRoute } from './types';
import { SupervisedUserCircleRound, WorkRound } from '@vicons/material';

const createProfileRoute: ProtectedRoute = {
  name: 'staff',
  path: '/staff',
  icon: WorkRound,
  text: 'Работники',
  access_rights: [AccessRight.AccessRightFull],
  disabled: false,
};

// const teacherTasksRoute: ProtectedRoute = {
//   name: 'teacher_tasks',
//   path: '/teacher_tasks',
//   icon: BookRound,
//   text: 'Домашние задания',
//   disabled: true,
// };

const usersListRoute: ProtectedRoute = {
  name: 'bot_users',
  path: '/bot_users/registered',
  icon: SupervisedUserCircleRound,
  text: 'Пользователи',
  access_rights: [AccessRight.AccessRightFull, AccessRight.AccessRightManager],
  disabled: false,
};

// const videoContentRoute: ProtectedRoute = {
//   name: 'video_content',
//   path: '/video_content',
//   icon: PlayLessonRound,
//   text: 'Видео уроки',
//   disabled: true,
// };

export const ProtectedRoutes: Array<ProtectedRoute> = [createProfileRoute, usersListRoute];
