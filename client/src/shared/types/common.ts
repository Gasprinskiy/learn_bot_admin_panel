export interface ResponseWithBoolStatus {
  success: boolean;
}

export interface PaginationParams {
  limit: number;
  page: number;
  next_cursor_id?: number;
  next_cursor_date?: string;
}

export interface ListResponse<T> {
  data: T[];
  total: number;
  left: number;
}
