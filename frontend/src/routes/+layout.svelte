<script>
  import { onMount } from 'svelte';
  import '../app.css';
  import { page } from '$app/stores';
  import { derived } from 'svelte/store';

  onMount(async () => {
    const { registerSW } = await import('virtual:pwa-register');
    let notified = false;
    registerSW({
      onNeedRefresh() {
        if (confirm('新しいバージョンがあります。更新しますか？')) {
          window.location.reload();
        }
      },
      onOfflineReady() {
        if (!notified) {
          alert('オフラインで利用可能です。');
          notified = true;
        }
      },
    });
  });

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

  // ダッシュボード配下か判定
  const isDashboard = derived(page, $page => $page.url.pathname.startsWith('/dashboard'));
</script>

<div class="app-container">
  {#if !$isDashboard}
    <!-- Header -->
    <header>
      <div class="container header-content">
        <div class="header-title-group">
          <a href="/">
            <h1>行事委員会</h1>
          </a>
        </div>
        <nav class="desktop-nav">
          <a href="/#home">ホーム</a>
          <a href="/#about">概要</a>
          <a href="/#events">イベント</a>
          <a href="/#roles">役職</a>
          <a href="/#join">参加方法</a>
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
          <a href="/#home" on:click={() => closeMenuAndScroll('home')}>ホーム</a>
          <a href="/#about" on:click={() => closeMenuAndScroll('about')}>概要</a>
          <a href="/#events" on:click={() => closeMenuAndScroll('events')}>イベント</a>
          <a href="/#roles" on:click={() => closeMenuAndScroll('roles')}>役職</a>
          <a href="/#join" on:click={() => closeMenuAndScroll('join')}>参加方法</a>
        </nav>
      </div>
    {/if}
  {/if}

  <main>
    <slot />
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

  .header-title-group a {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    color: white;
    text-decoration: none;
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

  .desktop-nav a {
    background: none;
    border: none;
    color: white;
    cursor: pointer;
    font-size: 1rem;
    padding: 0.5rem;
    transition: color 0.3s;
    text-decoration: none;
  }

  .desktop-nav a:hover {
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

  .mobile-nav-container a {
    background: none;
    border: none;
    color: white;
    cursor: pointer;
    font-size: 2rem;
    text-decoration: none;
  }

  .container {
    max-width: 1100px;
    margin: 0 auto;
    padding: 0 2rem;
  }

  main {
    /* The main content area for the pages */
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
  }
</style>
