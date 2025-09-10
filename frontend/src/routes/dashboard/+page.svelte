<script>
  import { goto } from "$app/navigation";

  function logout() {
    localStorage.removeItem("token");
    goto("/login");
  }
</script>

<div class="dashboard-container">
  <header>
    <h1>ダッシュボード</h1>
    <button on:click={logout}>ログアウト</button>
  </header>
  <main>
    <p>ようこそ！</p>
    <!-- ここにダッシュボードのコンテンツを追加します -->
  </main>
</div>

<style>
  .dashboard-container {
    padding: 2rem;
  }

  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid #ccc;
  }

  h1 {
    font-size: 2.5rem;
  }

  button {
    padding: 0.5rem 1rem;
    background-color: #d93025;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.3s;
  }

  button:hover {
    background-color: #c5221b;
  }
</style>
