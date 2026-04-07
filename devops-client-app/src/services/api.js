// Prefer relative `/api` so the frontend uses nginx proxy in production.
// Fallback to env var or localhost for local development.
const API_BASE = process.env.REACT_APP_API_URL || '/api';

export const request = async (url, options) => {
  const res = await fetch(`${API_BASE}${url}`, options);
  if (!res.ok) {
    const error = new Error(`HTTP error! status: ${res.status}`);
    error.status = res.status;
    error.statusText = res.statusText;
    throw error;
  }
  return res.json();
};