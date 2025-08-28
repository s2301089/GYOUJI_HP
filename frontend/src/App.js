import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';
import Header from './components/Header';
import Home from './pages/Home';
import Events from './pages/Events';
import Work from './pages/Work';
import Roles from './pages/Roles';
import Join from './pages/Join';

function App() {
  return (
    <Router>
      <div className="App">
        <Header />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/events" element={<Events />} />
          <Route path="/work" element={<Work />} />
          <Route path="/roles" element={<Roles />} />
          <Route path="/join" element={<Join />} />
        </Routes>
        <footer className="App-footer">
          <p>&copy; 2025 仙台高専広瀬キャンパス 行事委員会</p>
        </footer>
      </div>
    </Router>
  );
}

export default App;
