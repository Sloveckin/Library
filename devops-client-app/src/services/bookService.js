import { request } from './api';

export const createBook = (name, authors) => request('/book/create', {
  method: 'PUT',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ name, authors }),
});

export const getBook = (id) => request(`/book/get?id=${id}`, {
  method: 'GET',
});

export const deleteBook = (id) => request(`/book/delete?id=${id}`, {
  method: 'DELETE',
});

export const updateBook = (id, name, authors) => request('/book/update', {
  method: 'PUT',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ id, name, authors }),
});