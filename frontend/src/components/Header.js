import React from 'react';
import './Header.css';

const Header = () => {
  const handleScroll = (e, targetId) => {
    e.preventDefault();
    const targetElement = document.getElementById(targetId);
    if (targetElement) {
      const headerOffset = 80; // Adjust this value to match your header's height
      const elementPosition = targetElement.getBoundingClientRect().top;
      const offsetPosition = elementPosition + window.pageYOffset - headerOffset;

      window.scrollTo({
        top: offsetPosition,
        behavior: 'smooth'
      });
    }
  };

  return (
    <header className="main-header">
      <div className="logo">
        <a href="#home" onClick={(e) => handleScroll(e, 'home')}>行事委員会</a>
      </div>
      <nav>
        <ul>
          <li><a href="#home" onClick={(e) => handleScroll(e, 'home')}>ホーム</a></li>
          <li><a href="#events" onClick={(e) => handleScroll(e, 'events')}>主なイベント</a></li>
          <li><a href="#roles" onClick={(e) => handleScroll(e, 'roles')}>役職紹介</a></li>
          <li><a href="#join" onClick={(e) => handleScroll(e, 'join')}>参加方法</a></li>
        </ul>
      </nav>
    </header>
  );
};

export default Header;