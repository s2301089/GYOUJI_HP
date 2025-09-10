<script>
  import { goto } from "$app/navigation";

  let username = "";
  let password = "";
  let errorMessage = "";

  async function login() {
    try {
      const response = await fetch("/api/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });

      if (response.ok) {
        const data = await response.json();
        localStorage.setItem("token", data.token);
        goto("/dashboard");
      } else {
        errorMessage = "ログインに失敗しました。";
      }
    } catch (error) {
      errorMessage = "エラーが発生しました。";
    }
  }
</script>

<div class="login-container">
  <div class="login-card">
    <h2>ログイン</h2>
    <form on:submit|preventDefault={login}>
      <div class="input-group">
        <label for="username">ユーザー名</label>
        <input type="text" id="username" bind:value={username} required />
      </div>
      <div class="input-group">
        <label for="password">パスワード</label>
        <input type="password" id="password" bind:value={password} required />
      </div>
      {#if errorMessage}
        <p class="error-message">{errorMessage}</p>
      {/if}
      <button type="submit">ログイン</button>
    </form>
  </div>
</div>

<style>
  .login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background-color: #f0f2f5;
  }

  .login-card {
    background-color: white;
    padding: 2.5rem;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    width: 100%;
    max-width: 400px;
  }

  h2 {
    text-align: center;
    font-size: 2rem;
    margin-bottom: 2rem;
    color: #333;
  }

  .input-group {
    margin-bottom: 1.5rem;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
  }

  input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 1rem;
  }

  button {
    width: 100%;
    padding: 0.75rem;
    background-color: #4a69bd;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 1.1rem;
    cursor: pointer;
    transition: background-color 0.3s;
  }

  button:hover {
    background-color: #3b5998;
  }

  .error-message {
    color: #d93025;
    margin-bottom: 1.5rem;
    text-align: center;
  }
</style>
