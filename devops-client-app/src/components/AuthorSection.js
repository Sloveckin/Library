import React, { useState } from 'react';
import { createAuthor, getAuthor, deleteAuthor, updateAuthor } from '../services/authorService';
import { StatusModal, sectionStyle, cardStyle, inputStyle, btnStyle } from './SectionUI';

function AuthorSection() {
  const [authorGetId, setGetAuthorId] = useState('');
  const [authorDeleteId, setDeleteAuthorId] = useState('');
  const [authorName, setAuthorName] = useState('');
  const [authorUpdateId, setAuthorUpdateId] = useState('');
  const [authorUpdateName, setAuthorUpdateName] = useState('');
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
    setModalTitle('Ошибка');
  };

  const handleCreateAuthor = async () => {
    try {
      const data = await createAuthor(authorName);
      showModal(`Автор создан! Id: ${data.id}`, 'Успех');
      setAuthorName('');
    } catch (err) {
      showModal(`Не удалось создать автора: ${err.message}`);
    }
  };

  const handleUpdateAuthor = async () => {
    try {
      await updateAuthor(authorUpdateId, authorUpdateName);
      showModal('Автор обновлён!', 'Успех');
      setAuthorUpdateId('');
      setAuthorUpdateName('');
    } catch (err) {
      showModal(`Не удалось обновить автора: ${err.message}`);
    }
  };

  const handleGetAuthor = async () => {
    try {
      const data = await getAuthor(authorGetId);
      showModal(`Автор: ${data.name}`, 'Успех');
    } catch (err) {
      if (err.status === 404) {
        showModal('Автор не найден. Проверьте id и попробуйте снова.');
      } else {
        showModal(`Не удалось получить автора: ${err.message}`);
      }
    }
  };

  const handleDeleteAuthor = async () => {
    try {
      await deleteAuthor(authorDeleteId);
      showModal('Автор удалён!', 'Успех');
      setDeleteAuthorId('');
    } catch (err) {
      showModal(`Не удалось удалить автора: ${err.message}`);
    }
  };

  return (
    <>
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

      <div style={cardStyle('#ffeaa7', '#fd7e14')}>
        <h3>Update author</h3>
        <input
          type="text"
          placeholder="Id"
          value={authorUpdateId}
          onChange={(e) => setAuthorUpdateId(e.target.value)}
          style={inputStyle}
        />
        <input
          type="text"
          placeholder="New name"
          value={authorUpdateName}
          onChange={(e) => setAuthorUpdateName(e.target.value)}
          style={inputStyle}
        />
        <button onClick={handleUpdateAuthor} style={btnStyle('#fd7e14')}>Update</button>
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

export default AuthorSection;
