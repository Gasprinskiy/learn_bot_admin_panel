export enum AccessRight {
  AccessRightFull = 'full_access',
  AccessRightManager = 'manager_access',
  AccessRightTeacher = 'teacher_access',
}

export interface User {
  u_id: number;
  first_name: string;
  last_name: string;
  tg_user_name: string;
  is_password_set: boolean;
  access_right: AccessRight;
}

export interface UserShortInfo {
  u_id: number;
  access_right: AccessRight;
}

export interface AuthTempData {
  uu_id: string;
  auth_url: string;
}

export interface UserFirstLoginAnswer {
  is_password_set: boolean;
}

export interface PasswordLoginParams {
  user_name: string;
  password: string;
}

export interface PasswordLoginResponse {
  need_two_step_auth: boolean;
  u_id: number;
  access_right: AccessRight;
}
