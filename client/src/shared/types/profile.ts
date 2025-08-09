import type { ListResponse, PaginationParams } from './common';

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
  uu_id: string;
}

export enum SubscriptionStatus {
  ACTIVE = 'active',
  EXPIRED = 'expired',
  NOTEXISTS = 'not_exists',
}

export interface BotUserProfile {
  u_id: number;
  tg_id: number;
  tg_user_name: string;
  first_name: string;
  last_name: string;
  birth_date: string;
  phone_number: string;
  join_date: string;
  register_date: string;
  subscription_status: SubscriptionStatus;
}

export type BotUserProfileListResponse = ListResponse<BotUserProfile>;

export interface BotUserProfileQueryParam {
  query?: string;
  age_from?: number;
  age_till?: number;
  subscription_status?: SubscriptionStatus;
  join_date_from?: string;
  join_date_till?: string;
}

export type BotUserProfileQueryParamCommon = BotUserProfileQueryParam & PaginationParams;
