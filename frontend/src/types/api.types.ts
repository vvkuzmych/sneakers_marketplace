export interface ApiError {
  error: string;
  message?: string;
  statusCode?: number;
}

export interface PaginationParams {
  page?: number;
  pageSize?: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: string;
  page: number;
  pageSize: number;
}
