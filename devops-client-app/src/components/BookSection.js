import React, { useState } from 'react';
import { createBook, getBook, deleteBook, updateBook } from '../services/bookService';
import { StatusModal, sectionStyle as baseSectionStyle, cardStyle, inputStyle, btnStyle } from './SectionUI';

function BookSection({ sharedBookId, onSharedBookIdChange }) {
  const [getBookId, setGetBookId] = useState('');
  const [deleteBookId, setDeleteBookId] = useState('');
  const [bookName, setBookName] = useState('');
  const [bookAuthors, setBookAuthors] = useState('');
  const [updateBookId, setUpdateBookId] = useState('');
  const [updateBookName, setUpdateBookName] = useState('');
  const [updateBookAuthors, setUpdateBookAuthors] = useState('');
  const [modalMessage, setModalMessage] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalTitle, setModalTitle] = useState('Ошибка');

  const showModal = (message, title = 'Ошибка') => {
    setModalMessage(message);
    setModalTitle(title);
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setModalMessage('');
  };

  const handleGetBook = async () => {
    try {
      const data = await getBook(getBookId);
      showModal(`Книга: ${data.name}\nАвторы: ${data.authors?.join(', ') || '—'}`, 'Успех');
      // update shared id only after a successful fetch
      if (onSharedBookIdChange) onSharedBookIdChange(data.id || getBookId);
    } catch (err) {
      if (err.status === 404) {
        showModal('Книга не существует. Проверьте id и попробуйте снова.');
      } else {
        showModal(`Не удалось получить книгу: ${err.message}`);
      }
    }
  };

  const handleCreateBook = async () => {
    try {
      const authorsList = bookAuthors.split(',').map(s => s.trim()).filter(Boolean);
      const data = await createBook(bookName, authorsList);
      showModal(`Книга создана! Id: ${data.id}`, 'Успех');
      onSharedBookIdChange(data.id);
      setBookName('');
      setBookAuthors('');
    } catch (err) {
      showModal(`Не удалось создать книгу: ${err.message}`);
    }
  };

  const handleDeleteBook = async () => {
    try {
      await deleteBook(deleteBookId);
      showModal('Книга удалена!', 'Успех');
      setDeleteBookId('');
      onSharedBookIdChange('');
    } catch (err) {
      showModal(`Не удалось удалить книгу: ${err.message}`);
    }
  };

  const handleUpdateBook = async () => {
    try {
      const authorsList = updateBookAuthors.split(',').map(s => s.trim()).filter(Boolean);
      const idToUpdate = updateBookId || sharedBookId;
      if (!idToUpdate) {
        showModal('Требуется id для обновления книги');
        return;
      }
      await updateBook(idToUpdate, updateBookName, authorsList);
      showModal('Книга обновлена!', 'Успех');
      setUpdateBookId('');
      setUpdateBookName('');
      setUpdateBookAuthors('');
      onSharedBookIdChange('');
    } catch (err) {
      showModal(`Не удалось обновить книгу: ${err.message}`);
    }
  };

  const handleIdChange = (e) => {
    const value = e.target.value;
    setGetBookId(value);
    // do not update shared id on every keystroke
  };

  return (
    <>
      <section style={sectionStyle}>
        <div style={cardStyle('#d4edda', '#28a745')}>
          <h3>Get book</h3>
          <input
            type="text"
            placeholder="Id"
            value={getBookId}
            onChange={handleIdChange}
            style={inputStyle}
          />
          <button onClick={handleGetBook} style={btnStyle('#28a745')}>Get</button>
        </div>

        <div style={cardStyle('#d4edda', '#28a745')}>
          <h3>Create book</h3>
          <input
            type="text"
            placeholder="Name"
            value={bookName}
            onChange={(e) => setBookName(e.target.value)}
            style={inputStyle}
          />
          <input
            type="text"
            placeholder="Authors (comma-separated)"
            value={bookAuthors}
            onChange={(e) => setBookAuthors(e.target.value)}
            style={inputStyle}
          />
          <button onClick={handleCreateBook} style={btnStyle('#28a745')}>Create</button>
        </div>

        <div style={cardStyle('#d4edda', '#28a745')}>
          <h3>Delete book</h3>
          <input
            type="text"
            placeholder="Id"
            value={deleteBookId}
            onChange={(e) => setDeleteBookId(e.target.value)}
            style={inputStyle}
          />
          <button onClick={handleDeleteBook} style={btnStyle('#28a745')}>Delete</button>
        </div>

        <div style={cardStyle('#d4edda', '#28a745')}>
          <h3>Update book</h3>
          <input
            type="text"
            placeholder="Id"
            value={updateBookId || sharedBookId}
            onChange={(e) => { setUpdateBookId(e.target.value); }}
            style={inputStyle}
          />
          <input
            type="text"
            placeholder="Name"
            value={updateBookName}
            onChange={(e) => setUpdateBookName(e.target.value)}
            style={inputStyle}
          />
          <input
            type="text"
            placeholder="Authors (comma-separated)"
            value={updateBookAuthors}
            onChange={(e) => setUpdateBookAuthors(e.target.value)}
            style={inputStyle}
          />
          <button onClick={handleUpdateBook} style={btnStyle('#28a745')}>Update</button>
        </div>
      </section>

      <StatusModal
        isOpen={isModalOpen}
        title={modalTitle}
        message={modalMessage}
        onClose={closeModal}
      />
    </>
  );
}

const sectionStyle = {
  ...baseSectionStyle,
  marginBottom: '60px',
};

export default BookSection;
