import React, { useState } from 'react';
import { createBook, getBook, deleteBook } from '../services/bookService';

function BookSection({ sharedBookId, onSharedBookIdChange }) {
  const [getBookId, setGetBookId] = useState('');
  const [deleteBookId, setDeleteBookId] = useState('');
  const [bookName, setBookName] = useState('');
  const [bookAuthors, setBookAuthors] = useState('');

  const handleGetBook = async () => {
    try {
      const data = await getBook(getBookId);
      alert(`Book: ${data.name}\nAuthors: ${data.authors?.join(', ') || '—'}`);
    } catch (err) {
      alert(`Get book failed: ${err.message}`);
    }
  };

  const handleCreateBook = async () => {
    try {
      const authorsList = bookAuthors.split(',').map(s => s.trim()).filter(Boolean);
      await createBook(bookName, authorsList);
      alert('Book created!');
      setBookName('');
      setBookAuthors('');
    } catch (err) {
      alert(`Create book failed: ${err.message}`);
    }
  };

  const handleDeleteBook = async () => {
    try {
      await deleteBook(deleteBookId);
      alert('Book deleted!');
      setDeleteBookId('');
      onSharedBookIdChange('');
    } catch (err) {
      alert(`Delete book failed: ${err.message}`);
    }
  };

  const handleIdChange = (e) => {
    const value = e.target.value;
    setGetBookId(value);
    onSharedBookIdChange(value);
  };

  return (
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
          onChange={handleIdChange}
          style={inputStyle}
        />
        <button onClick={handleDeleteBook} style={btnStyle('#28a745')}>Delete</button>
      </div>
    </section>
  );
}

const sectionStyle = {
  display: 'flex',
  gap: '24px',
  flexWrap: 'wrap',
  marginBottom: '60px',
};

const cardStyle = (bg, border) => ({
  backgroundColor: bg,
  border: `2px dashed ${border}`,
  borderRadius: '12px',
  padding: '20px',
  minWidth: '200px',
  width: '250px',
  boxShadow: '0 4px 6px rgba(0,0,0,0.05)',
});

const inputStyle = {
  width: '100%',
  padding: '8px 12px',
  margin: '8px 0',
  border: '1px solid #ccc',
  borderRadius: '6px',
  fontSize: '14px',
  boxSizing: 'border-box',
};

const btnStyle = (color) => ({
  width: '100%',
  padding: '8px',
  marginTop: '12px',
  backgroundColor: color,
  color: 'white',
  fontWeight: 'bold',
  border: 'none',
  borderRadius: '6px',
  cursor: 'pointer',
});

export default BookSection;