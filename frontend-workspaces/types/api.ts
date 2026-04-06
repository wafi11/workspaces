import { ReactNode } from "react";

export interface BackendResponse<T> {
  message: string;
  data?: T;
  error?: string;
}

export type CursorResponse<T> = {
  items: T;
  nextCursor: string;
  hasNextPage: boolean;
};

export interface ApiError {
  message: string;
  code: number;
  error?: string;
  success: boolean;
}

export interface RequestConfig {
  isMultipart?: boolean;
  headers?: Record<string, string>;
}

export interface WithChildren {
  children: ReactNode;
}
export interface ErrorResponse {
  code: number;
  error: string;
  message?: string;
}

export type ApiResponse<T> = {
  message: string;
  status: number;
  data: T;
};

export type ApiPagination<T> = {
  data: {
    data: T;
    meta: PaginationResponse;
  };
  message: string;
  statusCode: number;
};

export type PaginationResponse = {
  currentPage: number;
  totalPages: number;
  totalItems: number;
  itemsPerPage: number;
  hasNextPage: boolean;
  hasPrevPage: boolean;
};

export type CountResponse = {
  count: number;
};

export interface PaginationParams {
  limit?: string;
  page?: string;
}

export type CursorPagination<T> = {
  data: T;
  metadata: PaginationCursor;
};

export interface PaginationCursor {
  nextCursor?: string | null;
  hasMore: boolean;
  count: number;
}
export type CursorRequest = {
  search?: string;
  cursor?: string;
};
