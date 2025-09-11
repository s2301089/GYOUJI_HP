<script>
	import { goto } from '$app/navigation';
	import { tick } from 'svelte';
	import { createBracket } from 'bracketry';
	import { onMount } from 'svelte';

	let allTournaments = []; // APIから取得したbracketry用データの配列
	let selectedTournament = null; // 選択中のbracketry用データ
	let bracketContainer;
	let isLoading = false;
	let selectedSport = '';

	function logout() {
		localStorage.removeItem('token');
		goto('/login');
	}

	async function fetchTournament(sport) {
		selectedSport = sport;
		isLoading = true;
		allTournaments = [];
		selectedTournament = null;

		const token = localStorage.getItem('token');
		try {
			const res = await fetch(`/api/tournaments/${sport}`, {
				headers: token ? { Authorization: `Bearer ${token}` } : {}
			});

			if (res.ok) {
				const data = await res.json();
				if (data && data.length > 0) {
					allTournaments = data; // APIからのレスポンスをそのまま格納
					selectTournament(allTournaments[0]); // 最初のトーナメントをデフォルトで表示
				} else {
					allTournaments = []; // データがない場合は空にする
					alert('この競技のトーナメント情報が見つかりませんでした。');
				}
			} else {
				alert('トーナメント情報の取得に失敗しました。');
			}
		} catch (error) {
			console.error('Fetch error:', error);
			alert('サーバーとの通信に失敗しました。');
		} finally {
			isLoading = false;
		}
	}
	
	// 表示するトーナメントを選択する関数
	function selectTournament(tournamentData) {
		selectedTournament = tournamentData;
		// データが更新された後にDOMの再描画を待ってからbracketを描画
		tick().then(drawBracket);
	}

	// bracketryライブラリを呼び出してトーナメント表を描画する関数
	async function drawBracket() {
		if (!selectedTournament || !bracketContainer) return;
		
		// 描画前に中身をクリア
		bracketContainer.innerHTML = '';
        
		// DOMが確実に存在することを保証
		await tick();

		// 整形済みのデータをそのままライブラリに渡す
		createBracket(selectedTournament, bracketContainer);
	}

	// ダッシュボード初期表示時にバレーボールのトーナメントを表示
	onMount(() => {
		fetchTournament('volleyball');
	});
</script>

<div class="dashboard-container">
	<header>
		<h1>ダッシュボード</h1>
		<button on:click={logout}>ログアウト</button>
	</header>
	<main>
		<p>ようこそ！</p>
		<div class="sports-buttons">
			<button on:click={() => fetchTournament('volleyball')} class:active={selectedSport === 'volleyball'}>バレーボール</button>
			<button on:click={() => fetchTournament('table_tennis')} class:active={selectedSport === 'table_tennis'}>卓球</button>
			<button on:click={() => fetchTournament('soccer')} class:active={selectedSport === 'soccer'}>サッカー</button>
		</div>

		{#if isLoading}
			<p>読み込み中...</p>
		{/if}

		<!-- 複数のトーナメント（例：卓球の晴天時・雨天時）を選択するボタン -->
		{#if !isLoading && allTournaments.length > 1}
			<div class="tournament-selector">
				{#each allTournaments as t}
					<button on:click={() => selectTournament(t)} class:active={selectedTournament?.name === t.name}>
						{t.name}
					</button>
				{/each}
			</div>
		{/if}

		<div class="bracket-area">
			{#if selectedTournament}
				<h2>{selectedTournament.name} トーナメント表</h2>
				<div bind:this={bracketContainer} class="bracket-wrapper" />
			{/if}
		</div>
	</main>
</div>

<style>
	:global(.bracket-match-team) {
		background-color: #f0f0f0 !important;
		border: 1px solid #ccc !important;
	}
	:global(.bracket-match-winner .bracket-match-team) {
		background-color: #d4edda !important;
		font-weight: bold;
	}
	:global(.bracket-connector) {
		border-color: #999 !important;
	}

	.dashboard-container {
		align-items: center;
		padding: 2rem;
		font-family: sans-serif;
	}

	header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
		padding-bottom: 1rem;
		border-bottom: 1px solid #ccc;
	}

	main {
		text-align: center;
		align-items: center;
		justify-content: center;
	}

	h1 {
		font-size: 2.5rem;
	}

	button {
		padding: 0.5rem 1rem;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		transition: background-color 0.3s;
	}

	header button {
		background-color: #d93025;
	}
	header button:hover {
		background-color: #c5221b;
	}

	.sports-buttons, .tournament-selector {
		margin-bottom: 2rem;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 1rem;
	}

	.sports-buttons button {
		padding: 0.5rem 1.5rem;
		background: #f1f3f4;
		color: #5f6368;
		border: 1px solid #dadce0;
		font-size: 1.1rem;
	}
	.sports-buttons button.active {
		background: #4285f4;
		color: #fff;
		border-color: #4285f4;
	}
	
	.tournament-selector button {
		background: #e8f0fe;
		color: #1967d2;
		border: 1px solid #d2e3fc;
	}
	.tournament-selector button.active {
		background: #1967d2;
		color: white;
		font-weight: bold;
	}

	.bracket-area {
		margin-top: 2rem;
		overflow-x: auto;
	}

	.bracket-wrapper {
		display: flex;
		align-items: center;
		justify-content: center;
		min-width: 800px;
	}

@media (max-width: 900px) {
	.dashboard-container {
		padding: 1rem;
	}
	h1 {
		font-size: 2rem;
	}
	.bracket-wrapper {
		min-width: 600px;
	}
}

@media (max-width: 600px) {
	.dashboard-container {
		padding: 0.5rem;
	}
	header {
		flex-direction: column;
		align-items: flex-start;
		gap: 0.5rem;
		padding-bottom: 0.5rem;
	}
	h1 {
		font-size: 1.3rem;
	}
	.sports-buttons, .tournament-selector {
		flex-direction: column;
		gap: 0.5rem;
		width: 100%;
	}
	.sports-buttons button, .tournament-selector button {
		width: 100%;
		font-size: 1rem;
		padding: 0.7rem 0.5rem;
	}
	.bracket-area {
		margin-top: 1rem;
		overflow-x: auto;
	}
	.bracket-wrapper {
		min-width: 350px;
	}
}
</style>

