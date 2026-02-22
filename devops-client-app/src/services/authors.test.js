// src/services/authors.test.js
import { createAuthor, getAuthor, deleteAuthor } from './authorService';
import { request } from './api'; 

// Мокаем модуль api
jest.mock('./api');

describe('Author API Functions', () => {
  afterEach(() => {
    jest.clearAllMocks(); // очищаем моки после каждого теста
  });

  describe('createAuthor', () => {
    test('should call request with correct parameters', async () => {
      const mockResponse = { id: 1, name: 'Test Author' };
      request.mockResolvedValue(mockResponse);

      const name = 'Test Author';
      const result = await createAuthor(name);

      expect(request).toHaveBeenCalledWith('/author/create', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name }),
      });
      expect(result).toEqual(mockResponse);
    });

    test('should handle request rejection', async () => {
      const error = new Error('Validation Error');
      request.mockRejectedValue(error);

      await expect(createAuthor('')).rejects.toThrow('Validation Error');
    });
  });

  describe('getAuthor', () => {
    test('should call request with correct parameters', async () => {
      const mockResponse = { id: 1, name: 'Test Author' };
      request.mockResolvedValue(mockResponse);

      const id = 1;
      const result = await getAuthor(id);

      expect(request).toHaveBeenCalledWith(`/author/get?id=${id}`, {
        method: 'GET',
      });
      expect(result).toEqual(mockResponse);
    });

    test('should handle request rejection', async () => {
      const error = new Error('Not Found');
      request.mockRejectedValue(error);

      await expect(getAuthor(1)).rejects.toThrow('Not Found');
    });
  });

  describe('deleteAuthor', () => {
    test('should call request with correct parameters', async () => {
      const mockResponse = null; // или любой ожидаемый ответ
      request.mockResolvedValue(mockResponse);

      const id = 1;
      const result = await deleteAuthor(id);

      expect(request).toHaveBeenCalledWith(`/author/delete?id=${id}`, {
        method: 'DELETE',
      });
      expect(result).toBeNull(); // или какое-то ожидаемое значение
    });

    test('should handle request rejection', async () => {
      const error = new Error('Forbidden');
      request.mockRejectedValue(error);

      await expect(deleteAuthor(1)).rejects.toThrow('Forbidden');
    });
  });
});