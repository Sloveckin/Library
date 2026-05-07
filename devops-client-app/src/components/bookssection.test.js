import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import BookSection from './BookSection';

jest.mock('../services/bookService', () => ({
  createBook: jest.fn(),
  getBook: jest.fn(),
  deleteBook: jest.fn(),
  updateBook: jest.fn(),
}));

import {
  createBook,
  getBook,
  deleteBook,
  updateBook,
} from '../services/bookService';

describe('BookSection', () => {
  const onSharedBookIdChange = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('creates book successfully', async () => {
    createBook.mockResolvedValue({ id: '1' });

    render(<BookSection sharedBookId="" onSharedBookIdChange={onSharedBookIdChange} />);

    const nameInputs = screen.getAllByPlaceholderText('Name');
    const authorsInputs = screen.getAllByPlaceholderText('Authors (comma-separated)');

    fireEvent.change(nameInputs[0], { target: { value: 'Book One' } });
    fireEvent.change(authorsInputs[0], { target: { value: 'John, Mike' } });

    fireEvent.click(screen.getByText('Create'));

    await waitFor(() => {
      expect(screen.getByText(/Книга создана/i)).toBeInTheDocument();
    });

    expect(createBook).toHaveBeenCalledWith('Book One', ['John', 'Mike']);
    expect(onSharedBookIdChange).toHaveBeenCalledWith('1');
  });

  test('gets book successfully', async () => {
    getBook.mockResolvedValue({
      id: '10',
      name: 'Clean Code',
      authors: ['Robert Martin'],
    });

    render(<BookSection sharedBookId="" onSharedBookIdChange={onSharedBookIdChange} />);

    const idInputs = screen.getAllByPlaceholderText('Id');

    fireEvent.change(idInputs[0], { target: { value: '10' } });
    fireEvent.click(screen.getByText('Get'));

    await waitFor(() => {
      expect(screen.getByText(/Книга:/i)).toBeInTheDocument();
      expect(screen.getByText(/Clean Code/i)).toBeInTheDocument();
      expect(screen.getByText(/Robert Martin/i)).toBeInTheDocument();
    });

    expect(getBook).toHaveBeenCalledWith('10');
    expect(onSharedBookIdChange).toHaveBeenCalledWith('10');
  });

  test('shows not found message', async () => {
    getBook.mockRejectedValue({ status: 404 });

    render(<BookSection sharedBookId="" onSharedBookIdChange={onSharedBookIdChange} />);

    const idInputs = screen.getAllByPlaceholderText('Id');

    fireEvent.change(idInputs[0], { target: { value: '404' } });
    fireEvent.click(screen.getByText('Get'));

    await waitFor(() => {
      expect(screen.getByText(/Книга не существует/i)).toBeInTheDocument();
    });
  });

  test('updates book successfully', async () => {
    updateBook.mockResolvedValue({});

    render(<BookSection sharedBookId="55" onSharedBookIdChange={onSharedBookIdChange} />);

    const nameInputs = screen.getAllByPlaceholderText('Name');
    const authorsInputs = screen.getAllByPlaceholderText('Authors (comma-separated)');
    const idInputs = screen.getAllByPlaceholderText('Id');

    fireEvent.change(idInputs[1], { target: { value: '55' } });
    fireEvent.change(nameInputs[1], { target: { value: 'Updated Book' } });
    fireEvent.change(authorsInputs[1], { target: { value: 'Alice, Bob' } });

    fireEvent.click(screen.getByText('Update'));

    await waitFor(() => {
      expect(screen.getByText(/Книга обновлена/i)).toBeInTheDocument();
    });

    expect(updateBook).toHaveBeenCalledWith('55', 'Updated Book', ['Alice', 'Bob']);
    expect(onSharedBookIdChange).toHaveBeenCalledWith('');
  });

  test('delete book successfully', async () => {
    deleteBook.mockResolvedValue({});

    render(<BookSection sharedBookId="1" onSharedBookIdChange={onSharedBookIdChange} />);

    const idInputs = screen.getAllByPlaceholderText('Id');

    fireEvent.change(idInputs[1], { target: { value: '1' } });
    fireEvent.click(screen.getByText('Delete'));

    await waitFor(() => {
      expect(screen.getByText(/Книга удалена/i)).toBeInTheDocument();
    });

    expect(deleteBook).toHaveBeenCalledWith('1');
    expect(onSharedBookIdChange).toHaveBeenCalledWith('');
  });

  test('shows validation when update id missing', async () => {
    render(<BookSection sharedBookId="" onSharedBookIdChange={onSharedBookIdChange} />);

    fireEvent.click(screen.getByText('Update'));

    await waitFor(() => {
      expect(screen.getByText(/Требуется id для обновления книги/i)).toBeInTheDocument();
    });

    expect(updateBook).not.toHaveBeenCalled();
  });
});