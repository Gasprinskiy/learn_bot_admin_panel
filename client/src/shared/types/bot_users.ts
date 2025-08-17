import type { ListResponse, PaginationParams } from './common';

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

export interface BotSubscriptionType {
  sub_id: number;
  term_in_month: number;
  price: number;
}
