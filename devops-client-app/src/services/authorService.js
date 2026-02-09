import { request } from './api';

export const createAuthor = (name) => request('/author/create', {
  method: 'PUT',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ name }),
});

export const getAuthor = (id) => request(`/author/get?id=${id}`, {
  method: 'GET',
});

export const deleteAuthor = (id) => request(`/author/delete?id=${id}`, {
  method: 'DELETE',
});