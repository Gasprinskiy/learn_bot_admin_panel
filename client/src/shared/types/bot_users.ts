import type { ListResponse, PaginationParams } from './common';

export enum SubscriptionStatus {
  ACTIVE = 'active',
  EXPIRED = 'expired',
  NOTEXISTS = 'not_exists',
}

export enum PaymentTypeID {
  P2P = 1,
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

export interface PurchaseData {
  p_id: number;
  sub_id: number;
  p_time: string;
  payment_type_id: PaymentTypeID;
  discount: number | null;
  manager_id: number | null;
  subscription_term: number;
  price: number;
  receipt_file_name: number | null;
  subscription_status: SubscriptionStatus;
  manager_first_name: string | null;
  manager_last_name: string | null;
}

export interface BotUserDetailData extends Omit<BotUserProfile, 'subscription_status'> {
  purchase_data: PurchaseData[] | null;
}
