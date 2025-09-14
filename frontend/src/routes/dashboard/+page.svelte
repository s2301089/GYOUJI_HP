<script>
	import { goto } from '$app/navigation';
	import { getAllContexts, tick } from 'svelte';
	import { createBracket } from 'bracketry';
	import { onMount } from 'svelte';

	// 卓球天候切替用（スライド式）
	let tableTennisWeather = (typeof window !== 'undefined' && localStorage.getItem('tableTennisWeather')) || 'sunny';
	let isRainyChecked = false;
	$: isRainyChecked = tableTennisWeather === 'rainy';

	$: {
        if (typeof window !== 'undefined') {
            localStorage.setItem('tableTennisWeather', tableTennisWeather);
        }
    }

	function toggleTableTennisWeather() {
		tableTennisWeather = tableTennisWeather === 'sunny' ? 'rainy' : 'sunny';
		fetchTournament('table_tennis');
		fetchMatches('table_tennis');
	}

	// 卓球試合入力用フィルタ
	function getFilteredTableTennisMatches() {
		if (tableTennisWeather === 'sunny') {
			return matchesBySport['table_tennis']?.filter(m => m.tournament_name === '卓球（晴天時）') ?? [];
		} else {
			return matchesBySport['table_tennis']?.filter(m => m.tournament_name === '卓球（雨天時）' || m.tournament_name === '卓球（雨天時・敗者戦左側）' || m.tournament_name === '卓球（雨天時・敗者戦右側）') ?? [];
		}
	}


	let allTournaments = [];
	let selectedTournament = null;
	let bracketContainer;
	let isLoading = false;
	let selectedSport = '';

	// タブ管理
	let activeTab = 'tournament'; // 'tournament' | 'input' | 'scores'

	// ユーザー情報（localStorageのtokenからデコードするか、APIで取得する想定）
	let userRole = '';
	let assignedSport = '';


	onMount(async () => {
		// ユーザー情報をAPIから取得
		const token = localStorage.getItem('token');
		if (token) {
			try {
				const res = await fetch('/api/auth/me', {
					headers: { Authorization: `Bearer ${token}` }
				});
				if (res.ok) {
					const user = await res.json();
					userRole = user.role || '';
					assignedSport = user.assigned_sport || '';
					console.log('userRole:', userRole);
					console.log('assignedSport:', assignedSport);
				} else {
					userRole = '';
					assignedSport = '';
				}
			} catch (e) {
				userRole = '';
				assignedSport = '';
			}
		}
		fetchTournament('volleyball');
	});

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
			let url = `/api/tournaments/${sport}`;
			if (sport === 'table_tennis') {
				url += `?weather=${tableTennisWeather}`;
			}
			const res = await fetch(url, {
				headers: token ? { Authorization: `Bearer ${token}` } : {}
			});
			console.log(res);

			if (res.ok) {
				const data = await res.json();
				if (data && data.length > 0) {
					allTournaments = data; // APIからのレスポンスをそのまま格納
					console.log(allTournaments);
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
	

	// 卓球トーナメント表示用フィルタ
	function getFilteredTableTennisTournaments() {
		if (tableTennisWeather === 'sunny') {
			return allTournaments.filter(t => t.name === '卓球（晴天時）');
		} else {
			return allTournaments.filter(t => t.name === '卓球（雨天時）' || t.name === '卓球（雨天時・敗者戦左側）' || t.name === '卓球（雨天時・敗者戦右側）');
		}
	}

	// 表示するトーナメントを選択する関数
	function selectTournament(tournamentData) {
		selectedTournament = tournamentData;
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

	// 試合一覧（競技ごと）
	let matchesBySport = {};
	let matchesLoading = false;

	// スコア一覧
	let scores = [];
	let scoresLoading = false;

	async function fetchScores() {
		scoresLoading = true;
		const token = localStorage.getItem('token');
		const classOrder = [
			'1-1', '1-2', '1-3', 
			'IS2', 'IT2', 'IE2', 
			'IS3', 'IT3', 'IE3', 
			'IS4', 'IT4', 'IE4', 
			'IS5', 'IT5', 'IE5', 
			'専・教'
		];
		try {
			const res = await fetch('/api/score', {
				headers: token ? { Authorization: `Bearer ${token}` } : {}
			});
			if (res.ok) {
				let fetchedScores = await res.json();
				fetchedScores.sort((a, b) => {
                    const indexA = classOrder.indexOf(a.class_name);
                    const indexB = classOrder.indexOf(b.class_name);
                    
                    if (indexA === -1 && indexB === -1) return 0;
                    if (indexA === -1) return 1;
                    if (indexB === -1) return -1;

                    return indexA - indexB;
                });
				scores = fetchedScores;
			} else {
				scores = [];
			}
		} catch (e) {
			scores = [];
		} finally {
			scoresLoading = false;
		}
	}

	// モーダル制御用
	let showConfirmModal = false;
	let editingMatch = null;
	let editingSport = '';

	// 試合一覧取得API
	async function fetchMatches(sport) {
		matchesLoading = true;
		const token = localStorage.getItem('token');
		try {
			let url = `/api/matches/${sport}`;
			if (sport === 'table_tennis') {
				url += `?weather=${tableTennisWeather}`;
			}
			const res = await fetch(url, {
				headers: { Authorization: `Bearer ${token}` }
			});
			if (res.ok) {
				const data = await res.json();
				matchesBySport[sport] = data;
				console.log(matchesBySport);
			} else {
				matchesBySport[sport] = [];
			}
		} catch (e) {
			matchesBySport[sport] = [];
		} finally {
			matchesLoading = false;
		}
	}

	// 試合スコア更新API（スコア更新後に再取得）
	async function updateMatchScore(match, sport) {
		const token = localStorage.getItem('token');
		const body = {
			user: userRole,
			team1_score: match.team1_score,
			team2_score: match.team2_score,
			winner_team_id: match.team1_score > match.team2_score ? match.team1_id : (match.team2_score > match.team1_score ? match.team2_id : null),
			status: 'finished'
		};
		try {
			const res = await fetch(`/api/matches/${match.id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${token}`
				},
				body: JSON.stringify(body)
			});
			if (res.ok) {
				match._updateStatus = 'success';
			} else {
				match._updateStatus = 'error';
			}
		} catch (e) {
			match._updateStatus = 'error';
		}
		setTimeout(() => { match._updateStatus = ''; }, 2000);
		// 更新後に再取得
		await fetchMatches(sport);
		// 卓球のみ配列自体を再代入してUIを強制更新
		if (sport === 'table_tennis') {
			matchesBySport['table_tennis'] = [...(matchesBySport['table_tennis'] ?? [])];
		}
	}

	// モーダルから呼ばれる確認用関数
	function confirmUpdateMatchScore() {
		console.log(editingMatch)
		console.log(editingSport)
		if (editingMatch && editingSport) {
			updateMatchScore(editingMatch, editingSport);
		}
		showConfirmModal = false;
		editingMatch = null;
		editingSport = '';
	}
</script>

<div class="dashboard-container">
	<header>
		<h1>ダッシュボード</h1>
		<nav class="dashboard-tabs">
			<button class:active={activeTab === 'tournament'} on:click={() => { activeTab = 'tournament'; fetchTournament('volleyball'); }}>競技トーナメント</button>
			<button class:active={activeTab === 'scores'} on:click={async () => { activeTab = 'scores'; await fetchScores(); }}>現在の得点</button>
			{#if userRole === 'superroot' || userRole === 'admin'}
				<button class:active={activeTab === 'input'} on:click={async () => {
					activeTab = 'input';
					let sports = userRole === 'superroot' ? ['volleyball', 'table_tennis', 'soccer'] : [assignedSport];
					for (const sport of sports) {
						await fetchMatches(sport);
					}
				}}>試合結果入力</button>
			{/if}
		</nav>
		<button class="logout-btn" on:click={logout}>ログアウト</button>
	</header>
	<main>
		{#if activeTab === 'tournament'}
			<p>ようこそ！</p>
			<div class="sports-buttons">
				<button on:click={() => fetchTournament('volleyball')} class:active={selectedSport === 'volleyball'}>バレーボール</button>
				<button on:click={() => fetchTournament('table_tennis')} class:active={selectedSport === 'table_tennis'}>卓球</button>
				<button on:click={() => fetchTournament('soccer')} class:active={selectedSport === 'soccer'}>サッカー</button>
				{#if (userRole === 'superroot' || (userRole === 'admin' && assignedSport === 'table_tennis')) && selectedSport === 'table_tennis'}
				<div class="weather-switcher">
					<span>卓球トーナメント天候:</span>
					<label class="switch">
						<input type="checkbox" bind:checked={isRainyChecked} on:change={toggleTableTennisWeather}>
						<span class="slider"></span>
					</label>
					<span>{tableTennisWeather === 'sunny' ? '晴天' : '雨天'}</span>
				</div>
				{/if}
			</div>
			{#if isLoading}
				<p>読み込み中...</p>
			{/if}
			{#if !isLoading && selectedSport === 'table_tennis'}
				{#if getFilteredTableTennisTournaments().length > 1}
					<div class="tournament-selector">
						{#each getFilteredTableTennisTournaments() as t}
							<button on:click={() => selectTournament(t)} class:active={selectedTournament?.name === t.name}>
								{t.name}
							</button>
						{/each}
					</div>
				{/if}
			{:else if !isLoading && allTournaments.length > 1}
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
					<div bind:this={bracketContainer} class="bracket-wrapper"></div>
				{/if}
			</div>
		{/if}

		{#if activeTab === 'scores'}
			<div class="scores-area">
				<h2>現在の得点</h2>
				{#if scoresLoading}
					<p>読み込み中...</p>
				{:else}
					{#if scores.length > 0}
						 <table class="scores-table">
						 <thead>
						 <tr>
						 <th>クラス</th>
						 <th>春スポ体合計点</th>
						 <th>出席点</th>
						 <th>バレーボール1勝点</th>
						 <th>バレーボール2勝点</th>
						 <th>バレーボール3勝点</th>
						 <th>バレーボール優勝点</th>
						 <th>卓球1勝点</th>
						 <th>卓球2勝点</th>
						 <th>卓球3勝点</th>
						 <th>卓球優勝点</th>
						 <th>卓球雨天ボーナス</th>
						 <th>サッカー1勝点</th>
						 <th>サッカー2勝点</th>
						 <th>サッカー3勝点</th>
						 <th>サッカー優勝点</th>
						 <th>合計(春スポ体合計点除く)</th>
						 <th>合計(春スポ体合計点含む)</th>
						 </tr>
						 </thead>
						 <tbody>
						 {#each scores as s}
						 <tr>
						 <td>{s.class_name}</td>
						 <td>{s.init_score}</td>
						 <td>{s.attendance_score}</td>
						 <td>{s.volleyball1_score}</td>
						 <td>{s.volleyball2_score}</td>
						 <td>{s.volleyball3_score}</td>
						 <td>{s.volleyball_championship_score}</td>
						 <td>{s.table_tennis1_score}</td>
						 <td>{s.table_tennis2_score}</td>
						 <td>{s.table_tennis3_score}</td>
						 <td>{s.table_tennis_championship_score}</td>
						 <td>{s.table_tennis_rainy_bonus_score}</td>
						 <td>{s.soccer1_score}</td>
						 <td>{s.soccer2_score}</td>
						 <td>{s.soccer3_score}</td>
						 <td>{s.soccer_championship_score}</td>
						 <td><b>{s.total_excluding_init}</b></td>
						 <td><b>{s.total_including_init}</b></td>
						 </tr>
						 {/each}
						 </tbody>
						 </table>
					{:else}
						<p>スコアデータがありません。</p>
					{/if}
				{/if}
			</div>
		{/if}
		{#if activeTab === 'input' && (userRole === 'superroot' || userRole === 'admin')}
			<div class="match-input-area">
				<h2>試合結果入力</h2>
				{#if userRole === 'admin'}
					<p>あなたの担当競技: <b>{assignedSport}</b></p>
				{/if}
				{#if userRole === 'superroot'}
					<p>全競技の試合結果を編集できます</p>
				{/if}
				{#if matchesLoading}
					<p>試合データ取得中...</p>
				{/if}
				{#each (userRole === 'superroot' ? ['volleyball', 'table_tennis', 'soccer'] : [assignedSport]) as sport}
					<div class="tournament-edit-card">
						<h3>{sport} の試合一覧</h3>
						<div class="matches-list">
							{#if sport === 'table_tennis'}
								{#if getFilteredTableTennisMatches().length > 0}
									{#each getFilteredTableTennisMatches() as m}
										{#if m.status !== 'finished'}
											<form class="match-edit-form" on:submit|preventDefault={() => { showConfirmModal = true; editingMatch = m; editingSport = sport; }}>
												<div class="match-info">
													<span>試合ID: {m.id}</span>
													<span>ラウンド: {m.round}</span>
													<span>トーナメント: {m.tournament_name}</span>
													<span>チーム: {m.team1_name || '-'} vs {m.team2_name || '-'}</span>
												</div>
												<div class="score-inputs">
													<span>{m.team1_name}</span>
													<input type="number" min="0" bind:value={m.team1_score} class="score-input"  required/>
													<span> - </span>
													<input type="number" min="0" bind:value={m.team2_score} class="score-input" required />
													<span>{m.team2_name}</span>
												</div>
												<button type="submit" class="update-btn">更新</button>
												{#if m._updateStatus}
													<span class="update-status {m._updateStatus === 'success' ? 'success' : 'error'}">{m._updateStatus === 'success' ? '更新成功' : '更新失敗'}</span>
												{/if}
											</form>
										{:else}
											<div class="match-info">
												<span>試合ID: {m.id}</span>
												<span>ラウンド: {m.round}</span>
												<span>トーナメント: {m.tournament_name}</span>
												<span>チーム: {m.team1_name || '-'} vs {m.team2_name || '-'}</span>
												<span class="update-status success">試合終了</span>
											</div>
										{/if}
									{/each}
								{:else}
									<p>試合データがありません。</p>
								{/if}
							{:else}
								{#if matchesBySport[sport]?.length > 0}
									{#each matchesBySport[sport] as m}
										{#if m.status !== 'finished'}
											<form class="match-edit-form" on:submit|preventDefault={() => { showConfirmModal = true; editingMatch = m; editingSport = sport; }}>
												<div class="match-info">
													<span>試合ID: {m.id}</span>
													<span>ラウンド: {m.round}</span>
													<span>チーム: {m.team1_name || '-'} vs {m.team2_name || '-'}</span>
												</div>
												<div class="score-inputs">
													<span>{m.team1_name}</span>
													<input type="number" min="0" bind:value={m.team1_score} class="score-input si1"  required/>
													<span> - </span>
													<input type="number" min="0" bind:value={m.team2_score} class="score-input si2" required />
													<span>{m.team2_name}</span>
												</div>
												<button type="submit" class="update-btn">更新</button>
												{#if m._updateStatus}
													<span class="update-status {m._updateStatus === 'success' ? 'success' : 'error'}">{m._updateStatus === 'success' ? '更新成功' : '更新失敗'}</span>
												{/if}
											</form>
										{:else}
											<div class="match-info">
												<span>試合ID: {m.id}</span>
												<span>ラウンド: {m.round}</span>
												<span>チーム: {m.team1_name || '-'} vs {m.team2_name || '-'}</span>
												<span class="update-status success">試合終了</span>
											</div>
										{/if}
									{/each}
								{:else}
									<p>試合データがありません。</p>
								{/if}
							{/if}
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</main>
	{#if showConfirmModal}
    	<div class="modal-overlay">
        	<div class="modal-content">
            	<h3>試合結果を更新しますか？</h3>
				<p>試合ID: {editingMatch.id}</p>
				<p>ラウンド: {editingMatch.round}</p>
				<p>チーム: {editingMatch.team1_name || '-'} vs {editingMatch.team2_name || '-'}</p>
				<p>スコア: {editingMatch.team1_score} - {editingMatch.team2_score}</p>
            	<p>本当にこの内容で更新してよろしいですか？</p>
            	<div class="modal-actions">
                	<button on:click={() => { showConfirmModal = false; editingMatch = null; editingSport = ''; }} class="update-btn" style="background:#ccc;color:#333;">キャンセル</button>
					<button on:click={confirmUpdateMatchScore} class="update-btn">OK</button>
            	</div>
        	</div>
    	</div>
	{/if}
</div>

<style>
/* スライド式スイッチ */
.switch {
	position: relative;
	display: inline-block;
	width: 50px;
	height: 24px;
	margin: 0 8px;
}
.switch input {
	opacity: 0;
	width: 0;
	height: 0;
}
.slider {
	position: absolute;
	cursor: pointer;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background-color: #ccc;
	transition: .4s;
	border-radius: 24px;
}
.slider:before {
	position: absolute;
	content: "";
	height: 18px;
	width: 18px;
	left: 3px;
	bottom: 3px;
	background-color: white;
	transition: .4s;
	border-radius: 50%;
}
input:checked + .slider {
	background-color: #2196F3;
}
input:checked + .slider:before {
	transform: translateX(26px);
}
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

	header button.logout-btn {
		background-color: #d93025;
	}
	header button.logout-btn:hover {
		background-color: #c5221b;
	}

	.sports-buttons, .tournament-selector {
		margin-bottom: 2rem;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 1rem;
	}
		.modal-overlay {
			position: fixed;
			top: 0;
			left: 0;
			width: 100vw;
			height: 100vh;
			background: rgba(0,0,0,0.3); /* 背景のみ半透明 */
			display: flex;
			align-items: center;
			justify-content: center;
			z-index: 1000;
			pointer-events: auto;
		}
		.modal-content {
			background: #fff; /* モーダル本体は白色 */
			border-radius: 8px;
			box-shadow: 0 2px 8px rgba(66,133,244,0.18);
			padding: 2rem;
			min-width: 300px;
			text-align: center;
			margin: auto;
			pointer-events: auto;
			z-index: 1000;
		}
		.modal-overlay > * {
			pointer-events: auto;
		}
		.modal-overlay:before {
			content: '';
			position: absolute;
			top: 0;
			left: 0;
			width: 100vw;
			height: 100vh;
			background: rgba(0,0,0,0.3);
			z-index: 999;
			pointer-events: none;
		}
		.modal-actions {
			display: flex;
			gap: 1rem;
			justify-content: center;
			margin-top: 1.5rem;
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
	/* 試合結果入力デザイン */
	.match-input-area {
		margin-top: 2rem;
		display: flex;
		flex-direction: column;
		align-items: center;
	}
	.tournament-edit-card {
		background: #f8f9fa;
		border-radius: 12px;
		box-shadow: 0 2px 8px rgba(66,133,244,0.08);
		padding: 2rem;
		margin-bottom: 2rem;
		width: 100%;
		max-width: 600px;
	}
	.matches-list {
		margin-top: 1rem;
	}
	.match-edit-form {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		background: #fff;
		border-radius: 8px;
		box-shadow: 0 1px 4px rgba(66,133,244,0.06);
		padding: 1rem;
		margin-bottom: 1rem;
	}
	.match-info {
		font-size: 1rem;
		color: #333;
		display: flex;
		gap: 1.5rem;
	}
	.score-inputs {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}
	.score-input {
		width: 3rem;
		padding: 0.3rem;
		font-size: 1.1rem;
		border: 1px solid #4285f4;
		border-radius: 4px;
		text-align: center;
	}
	.update-btn {
		background: #4285f4;
		color: #fff;
		border: none;
		border-radius: 4px;
		padding: 0.5rem 1.2rem;
		font-size: 1rem;
		cursor: pointer;
		margin-top: 0.5rem;
		transition: background 0.2s;
	}
	.update-btn:hover {
		background: #3367d6;
	}
	.update-status {
		font-size: 0.95rem;
		margin-top: 0.2rem;
	}
	.update-status.success {
		color: #43a047;
	}
	.update-status.error {
		color: #e53935;
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
	.dashboard-tabs {
		display: flex;
		gap: 1rem;
		margin-top: 1rem;
		margin-bottom: 1rem;
	}
	.dashboard-tabs button {
		padding: 0.5rem 1.5rem;
		background: #4285f4;
		color: #fff;
		border: none;
		border-radius: 4px;
		font-size: 1.1rem;
		cursor: pointer;
		transition: background 0.2s;
	}
	.dashboard-tabs button {
		background-color: #3367d6 !important;
	}
	.dashboard-tabs button:hover {
		background-color: #3367d6 !important;
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
