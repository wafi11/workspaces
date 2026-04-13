"use client";
export const storage = {
  getAccessToken: () => localStorage.getItem("access_token"),
  getRefreshToken: () => localStorage.getItem("refresh_token"),
  setAccessToken: (token: string) =>
    localStorage.setItem("access_token", token),
  setRefreshToken: (token: string) =>
    localStorage.setItem("refresh_token", token),
  clear: () => {
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
  },
};
