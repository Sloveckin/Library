import React from 'react';

export const sectionStyle = {
  display: 'flex',
  gap: '24px',
  flexWrap: 'wrap',
};

export const cardStyle = (bg, border) => ({
  backgroundColor: bg,
  border: `2px dashed ${border}`,
  borderRadius: '12px',
  padding: '20px',
  minWidth: '200px',
  width: '250px',
  boxShadow: '0 4px 6px rgba(0,0,0,0.05)',
});

export const inputStyle = {
  width: '100%',
  padding: '8px 12px',
  margin: '8px 0',
  border: '1px solid #ccc',
  borderRadius: '6px',
  fontSize: '14px',
  boxSizing: 'border-box',
};

export const btnStyle = (color) => ({
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

const overlayStyle = {
  position: 'fixed',
  top: 0,
  left: 0,
  right: 0,
  bottom: 0,
  backgroundColor: 'rgba(0, 0, 0, 0.45)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  zIndex: 999,
};

const modalStyle = {
  backgroundColor: 'white',
  borderRadius: '12px',
  padding: '24px',
  width: '100%',
  maxWidth: '420px',
  boxShadow: '0 16px 40px rgba(0, 0, 0, 0.15)',
  textAlign: 'center',
};

export function StatusModal({ isOpen, title, message, onClose }) {
  if (!isOpen) {
    return null;
  }

  return (
    <div style={overlayStyle} onClick={onClose}>
      <div style={modalStyle} onClick={(e) => e.stopPropagation()}>
        <h2 style={{ marginTop: 0 }}>{title}</h2>
        <div style={{ margin: '16px 0', lineHeight: '1.5' }}>
          <div style={{ whiteSpace: 'pre-line' }}>{message}</div>
        </div>
        <button onClick={onClose} style={btnStyle('#007bff')}>Закрыть</button>
      </div>
    </div>
  );
}
