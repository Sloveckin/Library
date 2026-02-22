const API_BASE = "http://localhost:8080";

export const request = async (url, options) => {
  const res = await fetch(`${API_BASE}${url}`, options);
  if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
  return res.json();
};