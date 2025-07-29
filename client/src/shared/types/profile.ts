export enum AccessRight {
  AccessRightFull = 'full_acces',
  AccessRightManager = 'manager_access',
  AccessRightTeacher = 'teacher_access',
}

export interface User {
  u_id: number;
  first_name: string;
  last_name: string;
  tg_user_name: string;
  tg_id: number | null;
  access_right: AccessRight;
}

export interface AuthTempData {
  uu_id: string;
  auth_url: string;
}

export interface UserFirstLoginAnswer {
  is_password_set: boolean;
}
