import React, { useState } from 'react';
import { createAuthor, getAuthor, deleteAuthor } from '../services/authorService';

function AuthorSection() {
  const [authorGetId, setGetAuthorId] = useState('');
  const [authorDeleteId, setDeleteAuthorId] = useState('');
  const [authorName, setAuthorName] = useState('');

  const handleCreateAuthor = async () => {
    try {
      await createAuthor(authorName);
      alert('Author created!');
      setAuthorName('');
    } catch (err) {
      alert(`Create author failed: ${err.message}`);
    }
  };

  const handleGetAuthor = async () => {
    try {
      const data = await getAuthor(authorGetId);
      alert(`Author: ${data.name}`);
    } catch (err) {
      alert(`Get author failed: ${err.message}`);
    }
  };

  const handleDeleteAuthor = async () => {
    try {
      await deleteAuthor(authorDeleteId);
      alert('Author deleted!');
      setDeleteAuthorId('');
    } catch (err) {
      alert(`Delete author failed: ${err.message}`);
    }
  };

  return (
    <section style={sectionStyle}>
      <div style={cardStyle('#ffeaa7', '#fd7e14')}>
        <h3>Create author</h3>
        <input
          type="text"
          placeholder="Name"
          value={authorName}
          onChange={(e) => setAuthorName(e.target.value)}
          style={inputStyle}
        />
        <button onClick={handleCreateAuthor} style={btnStyle('#fd7e14')}>Create</button>
      </div>

      <div style={cardStyle('#ffeaa7', '#fd7e14')}>
        <h3>Get author</h3>
        <input
          type="text"
          placeholder="Id"
          value={authorGetId}
          onChange={(e) => setGetAuthorId(e.target.value)}
          style={inputStyle}
        />
        <button onClick={handleGetAuthor} style={btnStyle('#fd7e14')}>Get</button>
      </div>

      <div style={cardStyle('#ffeaa7', '#fd7e14')}>
        <h3>Delete author</h3>
        <input
          type="text"
          placeholder="Id"
          value={authorDeleteId}
          onChange={(e) => setDeleteAuthorId(e.target.value)}
          style={inputStyle}
        />
        <button onClick={handleDeleteAuthor} style={btnStyle('#fd7e14')}>Delete</button>
      </div>
    </section>
  );
}

const sectionStyle = {
  display: 'flex',
  gap: '24px',
  flexWrap: 'wrap',
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

export default AuthorSection;