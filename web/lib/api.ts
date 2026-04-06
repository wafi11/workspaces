import { API_URL } from "@/constants";
import storage from "@/hooks/Storage";

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (token: string) => void;
  reject: (err: unknown) => void;
}> = [];

function processQueue(error: unknown, token: string | null = null) {
  failedQueue.forEach((promise) => {
    if (error) {
      promise.reject(error);
    } else {
      promise.resolve(token!);
    }
  });
  failedQueue = [];
}

function buildError(status: number, data: any, message?: string) {
  if (status) {
    return {
      code: status,
      message: data?.message || message || "Something went wrong",
      error: data?.error,
      success: false,
    };
  }
  return {
    code: 0,
    message: message || "Network error - no response received",
    error: message,
    success: false,
  };
}

type RequestConfig = {
  method?: string;
  headers?: Record<string, string>;
  body?: any;
  _retry?: boolean;
};

async function request<T = any>(
  endpoint: string,
  config: RequestConfig = {}
): Promise<{ data: T; status: number; message: string }> {
  const token = storage.getAccessToken();

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...config.headers,
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const url = endpoint.startsWith("http") ? endpoint : `${API_URL}${endpoint}`;

  let response: Response;

  try {
    response = await fetch(url, {
      method: config.method || "GET",
      headers,
      body: config.body ? JSON.stringify(config.body) : undefined,
    });
  } catch (err) {
    return Promise.reject({
      code: 0,
      message: "Network error - no response received",
      error: String(err),
      success: false,
    });
  }

  // Handle 401 + auto refresh
  if (response.status === 401 && !config._retry) {
    if (isRefreshing) {
      return new Promise<string>((resolve, reject) => {
        failedQueue.push({ resolve, reject });
      }).then((newToken) => {
        return request<T>(endpoint, {
          ...config,
          headers: { ...config.headers, Authorization: `Bearer ${newToken}` },
          _retry: true,
        });
      });
    }

    isRefreshing = true;
    const refreshToken = storage.getRefreshToken();

    if (!refreshToken) {
      storage.clear();
      window.location.href = "/login";
      return Promise.reject(buildError(401, null, "No refresh token"));
    }

    try {
      const refreshRes = await fetch(`${API_URL}auth/refresh`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      const refreshData = await refreshRes.json();

      if (!refreshRes.ok) {
        throw refreshData;
      }

      const newAccessToken = refreshData.data.access_token;
      const newRefreshToken = refreshData.data.refresh_token;

      storage.setAccessToken(newAccessToken);
      storage.setRefreshToken(newRefreshToken);

      processQueue(null, newAccessToken);

      return request<T>(endpoint, {
        ...config,
        headers: {
          ...config.headers,
          Authorization: `Bearer ${newAccessToken}`,
        },
        _retry: true,
      });
    } catch (refreshError) {
      processQueue(refreshError, null);
      storage.clear();
      window.location.href = "/login";
      return Promise.reject(refreshError);
    } finally {
      isRefreshing = false;
    }
  }

  const responseData = await response.json();

  if (!response.ok) {
    return Promise.reject(buildError(response.status, responseData));
  }

  return {
    data: responseData.data,
    status: response.status,
    message: responseData.message,
  };
}

// Method helpers — sama kayak axios instance sebelumnya
export const api = {
  get: <T = any>(
    url: string,
    config?: Omit<RequestConfig, "method" | "body">
  ) => request<T>(url, { ...config, method: "GET" }),

  post: <T = any>(
    url: string,
    body?: any,
    config?: Omit<RequestConfig, "method" | "body">
  ) => request<T>(url, { ...config, method: "POST", body }),

  put: <T = any>(
    url: string,
    body?: any,
    config?: Omit<RequestConfig, "method" | "body">
  ) => request<T>(url, { ...config, method: "PUT", body }),

  patch: <T = any>(
    url: string,
    body?: any,
    config?: Omit<RequestConfig, "method" | "body">
  ) => request<T>(url, { ...config, method: "PATCH", body }),

  delete: <T = any>(
    url: string,
    config?: Omit<RequestConfig, "method" | "body">
  ) => request<T>(url, { ...config, method: "DELETE" }),
};
