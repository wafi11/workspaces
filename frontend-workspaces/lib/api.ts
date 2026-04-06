import { API_URL } from "@/constants";
import storage from "@/hooks/Storage";
import { ApiResponse } from "@/types";
import axios, { AxiosInstance, AxiosResponse, AxiosError } from "axios";

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (token: string) => void;
  reject: (err: unknown) => void;
}> = [];

// Proses antrian request yang gagal — jalankan ulang setelah token baru didapat
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

class Api {
  private instance: AxiosInstance;

  constructor() {
    this.instance = axios.create({
      baseURL: API_URL,
      headers: {
        "Content-Type": "application/json",
      },
    });

    // Request interceptor — inject access token di setiap request
    this.instance.interceptors.request.use((config) => {
      const token = storage.getAccessToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Response interceptor — handle 401 + auto refresh
    this.instance.interceptors.response.use(
      (response: AxiosResponse) => {
        return {
          ...response,
          data: {
            data: response.data.data,
            status: response.status,
            message: response.data.message,
          },
        };
      },
      async (error: AxiosError) => {
        const originalRequest = error.config as any;

        // Kalau bukan 401 atau sudah pernah di-retry → lempar error
        if (error.response?.status !== 401 || originalRequest._retry) {
          return Promise.reject(buildError(error));
        }

        // Kalau sedang refresh — masukkan ke antrian, tunggu token baru
        if (isRefreshing) {
          return new Promise((resolve, reject) => {
            failedQueue.push({ resolve, reject });
          })
            .then((token) => {
              originalRequest.headers.Authorization = `Bearer ${token}`;
              return this.instance(originalRequest);
            })
            .catch((err) => {
              return Promise.reject(err);
            });
        }

        // Mulai refresh
        originalRequest._retry = true;
        isRefreshing = true;

        const refreshToken = storage.getRefreshToken();

        if (!refreshToken) {
          // Tidak ada refresh token → paksa logout
          storage.clear();
          window.location.href = "/login";
          return Promise.reject(buildError(error));
        }

        try {
          const { data } = await api.post<
            ApiResponse<{
              access_token: string;
              refresh_token: string;
            }>
          >(`${API_URL}auth/refresh`, {
            refresh_token: refreshToken,
          });

          const newAccessToken = data.data.access_token;
          const newRefreshToken = data.data.refresh_token;

          // Simpan token baru
          storage.setAccessToken(newAccessToken);
          storage.setRefreshToken(newRefreshToken);

          // Update header default
          this.instance.defaults.headers.common.Authorization = `Bearer ${newAccessToken}`;

          // Jalankan semua request yang tertunda
          processQueue(null, newAccessToken);

          // Retry request original
          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
          return this.instance(originalRequest);
        } catch (refreshError) {
          // Refresh gagal → paksa logout
          processQueue(refreshError, null);
          storage.clear();
          window.location.href = "/login";
          return Promise.reject(refreshError);
        } finally {
          isRefreshing = false;
        }
      }
    );
  }

  get instance_() {
    return this.instance;
  }
}

function buildError(error: AxiosError) {
  if (error.response) {
    const errorData = error.response.data as any;
    return {
      code: error.response.status,
      message: errorData.message || error.message,
      error: errorData.error,
      success: false,
    };
  } else if (error.request) {
    return {
      code: 0,
      message: "Network error - no response received",
      error: error.message,
      success: false,
    };
  }
  return {
    code: -1,
    message: error.message,
    error: error.message,
    success: false,
  };
}

export const api = new Api().instance_;
