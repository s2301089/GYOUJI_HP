import React from 'react';
import './App.css';
import Header from './components/Header';
import Home from './pages/Home';
import ScrollToTopButton from './components/ScrollToTopButton';

function App() {
  return (
    <div className="App">
      <Header />
      <Home />
      <ScrollToTopButton />
      <footer className="App-footer">
        <p>&copy; 2025 仙台高専広瀬キャンパス 行事委員会</p>
      </footer>
    </div>
  );
}

export default App;