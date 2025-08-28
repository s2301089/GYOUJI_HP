import React from 'react';
import { NavLink } from 'react-router-dom';
import './Header.css';

const Header = () => {
  return (
    <header className="main-header">
      <div className="logo">
        <NavLink to="/">行事委員会</NavLink>
      </div>
      <nav>
        <ul>
          <li><NavLink to="/">ホーム</NavLink></li>
          <li><NavLink to="/events">主なイベント</NavLink></li>
          <li><NavLink to="/work">仕事内容</NavLink></li>
          <li><NavLink to="/roles">役職紹介</NavLink></li>
          <li><NavLink to="/join">参加方法</NavLink></li>
        </ul>
      </nav>
    </header>
  );
};

export default Header;
