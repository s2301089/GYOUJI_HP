<script>
  let isMobileMenuOpen = false;

  const toggleMobileMenu = () => {
    isMobileMenuOpen = !isMobileMenuOpen;
  };

  const scrollToSection = (id) => {
    const element = document.getElementById(id);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
    }
  };

  const closeMenuAndScroll = (id) => {
    isMobileMenuOpen = false;
    scrollToSection(id);
  };
</script>

<div class="app-container">
  <!-- Header -->
  <header>
    <div class="container header-content">
      <div class="header-title-group">
        <img src="/icon.svg" alt="行事委員会 Icon" class="header-icon" />
        <h1>行事委員会</h1>
      </div>
      <nav class="desktop-nav">
        <button on:click={() => scrollToSection('home')}>ホーム</button>
        <button on:click={() => scrollToSection('about')}>概要</button>
        <button on:click={() => scrollToSection('events')}>イベント</button>
        <button on:click={() => scrollToSection('roles')}>役職</button>
        <button on:click={() => scrollToSection('join')}>参加方法</button>
      </nav>
      <button class="hamburger-menu" on:click={toggleMobileMenu} aria-label="メニューを開閉する">
        <span class="line"></span>
        <span class="line"></span>
        <span class="line"></span>
      </button>
    </div>
  </header>

  {#if isMobileMenuOpen}
    <div class="mobile-nav-container" on:click={toggleMobileMenu}>
      <nav on:click|stopPropagation>
        <button on:click={() => closeMenuAndScroll('home')}>ホーム</button>
        <button on:click={() => closeMenuAndScroll('about')}>概要</button>
        <button on:click={() => closeMenuAndScroll('events')}>イベント</button>
        <button on:click={() => closeMenuAndScroll('roles')}>役職</button>
        <button on:click={() => closeMenuAndScroll('join')}>参加方法</button>
      </nav>
    </div>
  {/if}

  <!-- Hero Section -->
  <section id="home" class="hero">
    <h1>仙台高専広瀬キャンパス 行事委員会</h1>
    <p>学校を、もっと面白く。君のアイデアで、最高のイベントを創ろう！</p>
  </section>

  <!-- Main Content -->
  <main class="container">
    <!-- About Section -->
    <section id="about">
      <h2>行事委員会へようこそ！</h2>
      <p class="about-text">
        私たち行事委員会は、学校生活を彩る様々なイベントを企画・運営しています。
        仲間と協力し、一つの目標に向かって全力で取り組む。その経験は、きっとあなたの大きな財産になります。
      </p>
    </section>

    <!-- Events Section -->
    <section id="events">
      <h2>主なイベント</h2>
      <div class="grid-col-2">
        <div class="card">
          <h3>スポーツ大会</h3>
          <div class="event-tags">
            <span class="tag">春季スポーツ大会</span>
            <span class="tag">秋季スポーツ大会</span>
          </div>
          <h4>運営の流れ</h4>
          <div class="step-flow">
            <div class="step">行う競技を決定する(実現可能な競技を模索します)</div>
            <div class="step">競技長を決める</div>
            <div class="step">競技長が競技メンバーと相談しながら競技要項を作成する。</div>
            <div class="step">物品リストを作成する</div>
            <div class="step">競技メンバーのシフトを作成する</div>
            <div class="step">前日に大会準備を進める</div>
            <div class="step">大会当日の朝に最終的な準備をし開会式をする</div>
            <div class="step">大会結果を集計し、閉会式で結果発表を行う</div>
            <div class="step">後日の委員会で反省点を共有し、のちの大会運営に生かす</div>
          </div>
        </div>
        <div class="card">
          <h3>その他のイベント</h3>
          <div class="event-tags">
            <span class="tag">予餞会</span>
            <span class="tag">新入生歓迎会</span>
          </div>
          <h4>予餞会運営の流れ</h4>
          <p>予餞会では、各部活動や対象の教員からビデオレターを募りそれを流すという形式をとっている</p>
          <div class="step-flow">
            <div class="step">部活動と対象となる教員を整理する</div>
            <div class="step">委員でグループを作り、それぞれいくつかの部活動等を担当する</div>
            <div class="step">担当対象の部活動や教員にビデオレターをお願いする</div>
            <div class="step">ビデオレターを回収し、まとめる</div>
            <div class="step">予餞会当日にまとめてビデオレターを流す</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Roles Section -->
    <section id="roles">
      <h2>役職紹介</h2>
      <div class="grid-col-4">
        <div class="card">
          <h3>行事委員長</h3>
          <p>委員会全体のリーダーとして、活動をまとめます。</p>
        </div>
        <div class="card">
          <h3>副行事委員長</h3>
          <p>委員長を補佐し、円滑な運営をサポートします。</p>
        </div>
        <div class="card">
          <h3>会計</h3>
          <p>委員会の予算を管理し、活動資金を確保します。</p>
        </div>
        <div class="card">
          <h3>書記</h3>
          <p>会議の議事録を作成し、活動の記録を残します。</p>
        </div>
      </div>
    </section>

    <!-- Join Section -->
    <section id="join" class="card">
      <h2>君も仲間にならないか？</h2>
      <p>行事委員会は、各クラスから2名まで参加できます。</p>
      <p>学校を盛り上げたい、新しいことに挑戦したい、そんな熱い想いを持った君を待っています！</p>
      <p>興味のある方は、来年に行事委員会に立候補してみてください。</p>
    </section>
  </main>

  <!-- Footer -->
  <footer>
    <p>&copy; 2025 仙台高専広瀬キャンパス 行事委員会</p>
  </footer>
</div>

<style>
  :global(body) {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen,
      Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    background-color: #f0f2f5;
    color: #333;
  }

  .app-container {
    width: 100%;
  }

  header {
    background-color: #333;
    color: white;
    padding: 1rem 0;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    position: sticky;
    top: 0;
    z-index: 1000;
  }

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .header-title-group {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .header-icon {
    height:5rem;
  }

  header h1 {
    font-size: 1.8rem;
    margin: 0;
  }

  .desktop-nav {
    display: flex;
    gap: 1rem;
  }

  .desktop-nav button {
    background: none;
    border: none;
    color: white;
    cursor: pointer;
    font-size: 1rem;
    padding: 0.5rem;
    transition: color 0.3s;
  }

  .desktop-nav button:hover {
    color: #ddd;
  }

  .hamburger-menu {
    display: none;
    flex-direction: column;
    justify-content: space-around;
    width: 2rem;
    height: 2rem;
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 0;
    z-index: 10;
  }

  .hamburger-menu .line {
    width: 2rem;
    height: 0.25rem;
    background: white;
    border-radius: 10px;
    transition: all 0.3s linear;
  }
  
  .mobile-nav-container {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.85);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 999;
  }

  .mobile-nav-container nav {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2rem;
  }

  .mobile-nav-container button {
    background: none;
    border: none;
    color: white;
    cursor: pointer;
    font-size: 2rem;
  }

  .container {
    max-width: 1100px;
    margin: 0 auto;
    padding: 0 2rem;
  }

  .hero {
    background: linear-gradient(to right, #4a69bd, #6a89cc);
    color: white;
    text-align: center;
    padding: 6rem 2rem;
  }

  .hero h1 {
    font-size: 3rem;
    font-weight: bold;
    margin: 0;
  }

  .hero p {
    font-size: 1.25rem;
    margin-top: 1rem;
  }

  main {
    padding: 3rem 0;
  }

  section {
    margin-bottom: 4rem;
  }

  h2 {
    font-size: 2.5rem;
    font-weight: bold;
    text-align: center;
    margin-bottom: 2rem;
  }

  h3 {
    font-size: 1.75rem;
    font-weight: bold;
    margin-bottom: 1.5rem;
    color: #4a69bd;
  }
  
  h4 {
    font-size: 1.25rem;
    font-weight: bold;
    margin-top: 1.5rem;
    margin-bottom: 1rem;
  }

  .about-text {
    text-align: center;
    max-width: 800px;
    margin: 0 auto;
    font-size: 1.1rem;
    line-height: 1.6;
  }

  .card {
    background-color: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0,0,0,0.1);
    transition: transform 0.3s, box-shadow 0.3s;
  }

  .card:hover {
    transform: translateY(-5px);
    box-shadow: 0 6px 12px rgba(0,0,0,0.15);
  }

  .grid-col-2 {
    display: grid;
    grid-template-columns: 1fr;
    gap: 2rem;
  }

  .grid-col-4 {
    display: grid;
    grid-template-columns: 1fr;
    gap: 2rem;
    text-align: center;
  }
  
  .event-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 1.5rem;
  }

  .tag {
    background-color: #e0e7ff;
    color: #4a69bd;
    padding: 0.25rem 0.75rem;
    border-radius: 1rem;
    font-size: 0.9rem;
    font-weight: 500;
  }

  .step-flow {
    position: relative;
    border-left: 3px solid #dbe1f8;
    padding-left: 2rem;
    margin-left: 1rem;
  }

  .step {
    position: relative;
    padding-bottom: 2rem;
    line-height: 1.6;
  }

  .step:last-child {
    padding-bottom: 0;
  }

  .step::before {
    content: '';
    position: absolute;
    left: -2.2rem;
    top: 0.15rem;
    width: 1.25rem;
    height: 1.25rem;
    border-radius: 50%;
    background-color: #6a89cc;
    border: 3px solid #f0f2f5;
  }

  #join {
    text-align: center;
  }

  footer {
    background-color: #333;
    color: white;
    text-align: center;
    padding: 2rem 0;
    margin-top: 4rem;
  }

  @media (max-width: 768px) {
    .desktop-nav {
      display: none;
    }

    .hamburger-menu {
      display: flex;
    }

    header h1 {
      font-size: 1.5rem;
    }

    .container {
      padding: 0 1rem;
    }

    .hero {
      padding: 4rem 1rem;
    }

    .hero h1 {
      font-size: 2rem;
    }

    .hero p {
      font-size: 1rem;
    }

    h2 {
      font-size: 2rem;
    }

    h3 {
      font-size: 1.5rem;
    }

    .card {
      padding: 1.5rem;
    }
  }

  @media (min-width: 768px) {
    .grid-col-2 {
      grid-template-columns: repeat(2, 1fr);
    }
    .grid-col-4 {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  @media (min-width: 1024px) {
    .grid-col-4 {
      grid-template-columns: repeat(4, 1fr);
    }
  }
</style>