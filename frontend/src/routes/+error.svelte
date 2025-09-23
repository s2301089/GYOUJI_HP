<script>
  import { page } from '$app/stores';
</script>

<div class="error-container">
  <h1 class="status-code">{$page.status}</h1>
  <p class="message">{$page.error?.message}</p>
  <a href="/" class="home-link">ホームに戻る</a>
</div>

<style>
  .error-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 80vh;
    font-family: 'Segoe UI', 'Hiragino Sans', 'Meiryo', sans-serif;
    color: #333;
  }

  .status-code {
    font-size: 10rem;
    font-weight: 700;
    color: #e0e0e0;
    margin: 0;
    letter-spacing: 0.05em;
  }

  .message {
    font-size: 1.5rem;
    margin: 0 0 2rem 0;
    color: #757575;
  }

  .home-link {
    font-size: 1.1rem;
    color: #fff;
    background-color: #1976d2;
    padding: 0.8rem 1.8rem;
    border-radius: 8px;
    text-decoration: none;
    transition: background-color 0.2s ease;
  }

  .home-link:hover {
    background-color: #1565c0;
  }
</style>
