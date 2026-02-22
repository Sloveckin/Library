import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import BooksPage from './pages/BooksPage';
import AuthorsPage from './pages/AuthorsPage';

function AppRoutes() {
  return (
    <Router>
      <Routes>
        <Route path="/books" element={<BooksPage />} />
        <Route path="/authors" element={<AuthorsPage />} />
        <Route path="/" element={<h2>Welcome to the Library App</h2>} />
      </Routes>
    </Router>
  );
}

export default AppRoutes;