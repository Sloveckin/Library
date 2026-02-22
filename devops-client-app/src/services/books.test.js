// src/services/books.test.js
import { createBook, getBook, deleteBook } from './bookService';
import { request } from './api'; 

// Мокаем модуль api
jest.mock('./api');

describe('Book API Functions', () => {
  afterEach(() => {
    jest.clearAllMocks(); // очищаем моки после каждого теста
  });

  describe('createBook', () => {
    test('should call request with correct parameters', async () => {
      const mockResponse = { id: 1, name: 'Test Book', authors: [] };
      request.mockResolvedValue(mockResponse);

      const name = 'Test Book';
      const authors = [{ id: 1, name: 'Author 1' }];
      const result = await createBook(name, authors);

      expect(request).toHaveBeenCalledWith('/book/create', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, authors }),
      });
      expect(result).toEqual(mockResponse);
    });

    test('should handle request rejection', async () => {
      const error = new Error('Network Error');
      request.mockRejectedValue(error);

      await expect(createBook('Test', [])).rejects.toThrow('Network Error');
    });
  });

  describe('getBook', () => {
    test('should call request with correct parameters', async () => {
      const mockResponse = { id: 1, name: 'Test Book', authors: [] };
      request.mockResolvedValue(mockResponse);

      const id = 1;
      const result = await getBook(id);

      expect(request).toHaveBeenCalledWith(`/book/get?id=${id}`, {
        method: 'GET',
      });
      expect(result).toEqual(mockResponse);
    });

    test('should handle request rejection', async () => {
      const error = new Error('Not Found');
      request.mockRejectedValue(error);

      await expect(getBook(1)).rejects.toThrow('Not Found');
    });
  });

  describe('deleteBook', () => {
    test('should call request with correct parameters', async () => {
      const mockResponse = null; // или любой ожидаемый ответ
      request.mockResolvedValue(mockResponse);

      const id = 1;
      const result = await deleteBook(id);

      expect(request).toHaveBeenCalledWith(`/book/delete?id=${id}`, {
        method: 'DELETE',
      });
      expect(result).toBeNull(); 
    });

    test('should handle request rejection', async () => {
      const error = new Error('Conflict');
      request.mockRejectedValue(error);

      await expect(deleteBook(1)).rejects.toThrow('Conflict');
    });
  });
});