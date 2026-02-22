import React, { useState } from 'react';
import BookSection from './components/BookSection';
import AuthorSection from './components/AuthorSection';

function App() {
  const [sharedBookId, setSharedBookId] = useState('');

  return (
    <div className="App" style={{
      padding: '40px',
      fontFamily: 'system-ui, sans-serif',
      maxWidth: '1200px',
      margin: '0 auto',
      backgroundColor: '#f9f9f9',
    }}>
      <h1 style={{ textAlign: 'center', marginBottom: '40px', color: '#333' }}>
        Library Management Client
      </h1>

       <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '40px' }}>
      <BookSection
        sharedBookId={sharedBookId}
        onSharedBookIdChange={setSharedBookId}
      />

      <AuthorSection />
      </div>
    </div>
  );
}

export default App;