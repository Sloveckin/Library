import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import AuthorSection from './AuthorSection';

import {
  createAuthor,
  getAuthor,
  deleteAuthor,
  updateAuthor,
} from '../services/authorService';

jest.mock('../services/authorService');

describe('AuthorSection', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('creates author successfully', async () => {
    createAuthor.mockResolvedValue({ id: '123' });

    render(<AuthorSection />);

    const inputs = screen.getAllByRole('textbox');

    fireEvent.change(inputs[0], {
      target: { value: 'Mark' },
    });

    fireEvent.click(screen.getByText('Create'));

    await waitFor(() => {
      expect(createAuthor).toHaveBeenCalledWith('Mark');
    });

    await waitFor(() => {
      expect(screen.getByText('Успех')).toBeInTheDocument();
      expect(
        screen.getByText((content) =>
          content.includes('Автор создан! Id: 123')
        )
      ).toBeInTheDocument();
    });
  });

  test('shows error when create author fails', async () => {
    createAuthor.mockRejectedValue(new Error('create failed'));

    render(<AuthorSection />);

    const inputs = screen.getAllByRole('textbox');

    fireEvent.change(inputs[0], {
      target: { value: 'Mark' },
    });

    fireEvent.click(screen.getByText('Create'));

    await waitFor(() => {
      expect(
        screen.getByText((content) =>
          content.includes('Не удалось создать автора')
        )
      ).toBeInTheDocument();
    });
  });

  test('gets author successfully', async () => {
    getAuthor.mockResolvedValue({ name: 'John' });

    render(<AuthorSection />);

    const inputs = screen.getAllByRole('textbox');

    fireEvent.change(inputs[1], {
      target: { value: '1' },
    });

    fireEvent.click(screen.getByText('Get'));

    await waitFor(() => {
      expect(getAuthor).toHaveBeenCalledWith('1');
    });

    await waitFor(() => {
      expect(
        screen.getByText((content) =>
          content.includes('Автор: John')
        )
      ).toBeInTheDocument();
    });
  });

  test('shows not found message on 404', async () => {
    getAuthor.mockRejectedValue({
      status: 404,
      message: 'not found',
    });

    render(<AuthorSection />);

    const inputs = screen.getAllByRole('textbox');

    fireEvent.change(inputs[1], {
      target: { value: '999' },
    });

    fireEvent.click(screen.getByText('Get'));

    await waitFor(() => {
      expect(
        screen.getByText((content) =>
          content.includes('Автор не найден')
        )
      ).toBeInTheDocument();
    });
  });

  test('updates author successfully', async () => {
    updateAuthor.mockResolvedValue({});

    render(<AuthorSection />);

    const inputs = screen.getAllByRole('textbox');

    fireEvent.change(inputs[3], {
      target: { value: '1' },
    });

    fireEvent.change(inputs[4], {
      target: { value: 'New Name' },
    });

    fireEvent.click(screen.getByText('Update'));

    await waitFor(() => {
      expect(updateAuthor).toHaveBeenCalledWith('1', 'New Name');
    });

    await waitFor(() => {
      expect(
        screen.getByText((content) =>
          content.includes('Автор обновлён!')
        )
      ).toBeInTheDocument();
    });
  });

  test('deletes author successfully', async () => {
    deleteAuthor.mockResolvedValue({});

    render(<AuthorSection />);

    const inputs = screen.getAllByRole('textbox');

    fireEvent.change(inputs[2], {
      target: { value: '1' },
    });

    fireEvent.click(screen.getByText('Delete'));

    await waitFor(() => {
      expect(deleteAuthor).toHaveBeenCalledWith('1');
    });

    await waitFor(() => {
      expect(
        screen.getByText((content) =>
          content.includes('Автор удалён!')
        )
      ).toBeInTheDocument();
    });
  });

  test('closes modal', async () => {
    createAuthor.mockResolvedValue({ id: '123' });

    render(<AuthorSection />);

    const inputs = screen.getAllByRole('textbox');

    fireEvent.change(inputs[0], {
      target: { value: 'Mark' },
    });

    fireEvent.click(screen.getByText('Create'));

    await waitFor(() => {
      expect(
        screen.getByText((content) =>
          content.includes('Автор создан! Id: 123')
        )
      ).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Закрыть'));

    await waitFor(() => {
      expect(
        screen.queryByText((content) =>
          content.includes('Автор создан! Id: 123')
        )
      ).not.toBeInTheDocument();
    });
  });
});