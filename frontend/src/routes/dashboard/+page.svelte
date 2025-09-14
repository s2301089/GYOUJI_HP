<script>
	import { goto } from '$app/navigation';
	import { tick } from 'svelte';
	import { createBracket } from 'bracketry';
	import RelayTable from '../../lib/RelayTable.svelte';
	import AttendanceInput from '../../lib/AttendanceInput.svelte';
	import { dndzone } from 'svelte-dnd-action';

	// --- State Variables ---

	// General
	let userRole = '';
	let assignedSport = '';
	let activeTab = 'tournament'; // 'tournament' | 'input' | 'scores'
	let isLoading = false;
	let classNameMap = {};

	// 学年名のマッピング
	const gradeNames = {
		1: '1年生',
		2: '2年生', 
		3: '3年生',
		4: '4年生',
		5: '5年生',
		6: '専・教'
	};

	// Tournament
	let allTournaments = [];
	let selectedTournament = null;
	let bracketContainer;
	let selectedSport = '';
	let tableTennisWeather = (typeof window !== 'undefined' && localStorage.getItem('tableTennisWeather')) || 'sunny';
	let isRainyChecked = false;

	// Match Input
	let matchesBySport = {};
	let matchesLoading = false;
	let showConfirmModal = false;
	let editingMatch = null;
	let editingSport = '';
	let activeInputTab = 'tournament'; // 'tournament' | 'relay'

	// Score
	let scores = [];
	let scoresLoading = false;

	// Attendance
	let attendanceScores = [];
	let attendanceLoading = false;
	let attendanceUpdateStatus = '';

	// Relay
	let relayActive = false; // For tournament tab's relay view
	let relayType = 'A';
	let relayResults = [];
	let relayLoading = false;
	let relayError = '';
	let relayUpdateStatus = ''; // 'success' | 'error' | ''


	// --- Lifecycle & Initialisation ---

	import { onMount } from 'svelte';
	onMount(async () => {
		await fetchUser();
		await fetchTournament('volleyball');
		await fetchClassNameMap();
		if (userRole === 'superroot') {
			await fetchAttendanceScores();
		}
	});

	async function fetchUser() {
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
				}
			} catch (e) {
				console.error('Failed to fetch user', e);
			}
		}
	}

	function logout() {
		localStorage.removeItem('token');
		goto('/login');
	}


	// --- Data Fetching ---

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
			if (res.ok) {
				const data = await res.json();
				if (data && data.length > 0) {
					allTournaments = data;
					selectTournament(allTournaments[0]);
				} else {
					allTournaments = [];
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
				matchesBySport[sport] = await res.json();
			} else {
				matchesBySport[sport] = [];
			}
		} catch (e) {
			matchesBySport[sport] = [];
		} finally {
			matchesLoading = false;
		}
	}

	async function fetchScores() {
		scoresLoading = true;
		const token = localStorage.getItem('token');
		const classOrder = ['1-1', '1-2', '1-3', 'IS2', 'IT2', 'IE2', 'IS3', 'IT3', 'IE3', 'IS4', 'IT4', 'IE4', 'IS5', 'IT5', 'IE5', '専・教'];
		try {
			const res = await fetch('/api/score', {
				headers: token ? { Authorization: `Bearer ${token}` } : {}
			});
			if (res.ok) {
				let fetchedScores = await res.json();
				fetchedScores.sort((a, b) => {
                    const indexA = classOrder.indexOf(a.class_name);
                    const indexB = classOrder.indexOf(b.class_name);
                    if (indexA === -1) return 1;
                    if (indexB === -1) return -1;
                    return indexA - indexB;
                });
				scores = fetchedScores;
			}
		} catch (e) {
			scores = [];
		} finally {
			scoresLoading = false;
		}
	}

	async function fetchRelayResults(type) {
		relayLoading = true;
		relayError = '';
		relayType = type;
		const token = localStorage.getItem('token');
		try {
			const res = await fetch(`/api/relay?block=${type}`, {
				headers: token ? { Authorization: `Bearer ${token}` } : {}
			});
			if (res.ok) {
				let data = await res.json();
				
				// 新しいAPIレスポンス形式に対応
				if (data && data.rankings && Object.keys(data.rankings).length > 0) {
					// rankings: {1: 3, 2: 1, 3: 5, 4: 2, 5: 4, 6: 6} の形式
					// 順位 -> 学年のマッピングを順位順のリストに変換
					const formattedData = Object.entries(data.rankings)
						.sort(([rankA], [rankB]) => parseInt(rankA) - parseInt(rankB))
						.map(([rank, grade]) => ({
							id: `relay-${type}-${rank}-${grade}-${Date.now()}`, // より一意なID
							rank: parseInt(rank),
							grade: parseInt(grade),
							relay_type: type
						}));
					relayResults = formattedData;
				} else {
					// データがない場合の処理
					if (userRole === 'superroot' || userRole === 'admin_relay') {
						// 管理者の場合は編集用のデフォルトデータを作成
						relayResults = [1, 2, 3, 4, 5, 6].map(rank => ({
							id: `relay-${type}-${rank}-${Date.now()}`, // より一意なID
							rank: rank,
							grade: rank, // デフォルトで順番通り
							relay_type: type
						}));
					} else {
						// 一般ユーザーの場合は空のデータ（学年なし）を作成
						relayResults = [1, 2, 3, 4, 5, 6].map(rank => ({
							id: `relay-${type}-${rank}-${Date.now()}`,
							rank: rank,
							grade: null, // 学年は空
							relay_type: type
						}));
					}
				}
			} else {
				relayError = 'リレー結果の取得に失敗しました';
			}
		} catch (e) {
			console.error("Fetch error:", e);
			relayError = 'サーバー通信エラー、またはデータの形式が不正です。';
		} finally {
			relayLoading = false;
		}
	}

	async function fetchClassNameMap() {
		try {
			const res = await fetch('/api/score');
			if (res.ok) {
				const scoresData = await res.json();
				classNameMap = scoresData.reduce((acc, s) => {
					acc[s.class_id] = s.class_name;
					return acc;
				}, {});
			}
		} catch {}
	}


	// --- UI Logic & Event Handlers ---

	function handleInputTabClick() {
		activeTab = 'input';
		if (userRole === 'admin_relay') {
			activeInputTab = 'relay';
			fetchRelayResults(relayType);
		} else {
			activeInputTab = 'tournament';
			const sports = userRole === 'superroot' ? ['volleyball', 'table_tennis', 'soccer'] : (assignedSport ? [assignedSport] : []);
			sports.forEach(fetchMatches);
		}
	}

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

	function getFilteredTableTennisMatches() {
		const allTableTennisMatches = matchesBySport['table_tennis'] ?? [];
		if (tableTennisWeather === 'sunny') {
			return allTableTennisMatches.filter(m => m.tournament_name === '卓球（晴天時）');
		}
		return allTableTennisMatches.filter(m => m.tournament_name.startsWith('卓球（雨天時）'));
	}

	function getFilteredTableTennisTournaments() {
		if (tableTennisWeather === 'sunny') {
			return allTournaments.filter(t => t.name === '卓球（晴天時）');
		}
		return allTournaments.filter(t => t.name.startsWith('卓球（雨天時）'));
	}

	function selectTournament(tournamentData) {
		selectedTournament = tournamentData;
		tick().then(drawBracket);
	}

	async function drawBracket() {
		if (!selectedTournament || !bracketContainer) return;
		bracketContainer.innerHTML = '';
		await tick();
		createBracket(selectedTournament, bracketContainer);
	}


	// --- API Updates ---

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
				headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
				body: JSON.stringify(body)
			});
			match._updateStatus = res.ok ? 'success' : 'error';
		} catch (e) {
			match._updateStatus = 'error';
		}
		setTimeout(() => { match._updateStatus = ''; }, 2000);
		await fetchMatches(sport);
		if (sport === 'table_tennis') {
			matchesBySport['table_tennis'] = [...(matchesBySport['table_tennis'] ?? [])];
		}
	}

	function confirmUpdateMatchScore() {
		if (editingMatch && editingSport) {
			updateMatchScore(editingMatch, editingSport);
		}
		showConfirmModal = false;
		editingMatch = null;
		editingSport = '';
	}

	function handleDndConsider(event) {
		// ドラッグ中の一時的な状態更新
		if (event.detail.items && Array.isArray(event.detail.items)) {
			relayResults = event.detail.items.map((item, index) => ({
				...item,
				rank: index + 1 // 順位を更新
			}));
		}
	}

	function handleDndFinalize(event) {
		// ドロップ完了時の最終状態更新
		if (event.detail.items && Array.isArray(event.detail.items)) {
			relayResults = event.detail.items.map((item, index) => ({
				...item,
				rank: index + 1 // 順位を更新
			}));
		}
	}

	async function updateRelayRanks() {
		relayUpdateStatus = 'loading';
		const token = localStorage.getItem('token');

		if (relayResults.length !== 6) {
			alert('リレーのクラス数が6ではありません。');
			relayUpdateStatus = '';
			return;
		}

		// 新しいAPI形式に合わせて順位 -> 学年のマッピングを作成
		const rankings = {};
		relayResults.forEach((result, index) => {
			rankings[index + 1] = result.grade; // 順位(1-6) -> 学年(1-6)
		});

		const body = { rankings };

		try {
			const res = await fetch(`/api/relay?block=${relayType}`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
				body: JSON.stringify(body)
			});
			if (res.ok) {
				relayUpdateStatus = 'success';
				await fetchRelayResults(relayType);
			} else {
				const errorData = await res.json();
				alert(`更新に失敗しました: ${errorData.error}`);
				relayUpdateStatus = 'error';
			}
		} catch (e) {
			console.error('Relay update error:', e);
			alert('サーバーとの通信に失敗しました。');
			relayUpdateStatus = 'error';
		}
		setTimeout(() => { relayUpdateStatus = ''; }, 3000);
	}

	// 出席点関連の関数
	async function fetchAttendanceScores() {
		attendanceLoading = true;
		const token = localStorage.getItem('token');
		try {
			const res = await fetch('/api/attendance', {
				headers: token ? { Authorization: `Bearer ${token}` } : {}
			});
			if (res.ok) {
				attendanceScores = await res.json();
			} else {
				console.error('Failed to fetch attendance scores');
			}
		} catch (e) {
			console.error('Error fetching attendance scores:', e);
		} finally {
			attendanceLoading = false;
		}
	}

	function handleAttendanceScoreChange(event) {
		const { classId, score } = event.detail;
		attendanceScores = attendanceScores.map(item => 
			item.class_id === classId ? { ...item, score } : item
		);
	}

	async function handleAttendanceBatchUpdate() {
		attendanceUpdateStatus = 'loading';
		const token = localStorage.getItem('token');
		
		const scores = attendanceScores.map(item => ({
			class_id: item.class_id,
			score: item.score
		}));

		try {
			const res = await fetch('/api/attendance/batch', {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${token}`
				},
				body: JSON.stringify({ scores })
			});

			if (res.ok) {
				attendanceUpdateStatus = 'success';
				await fetchScores(); // スコア表を更新
			} else {
				attendanceUpdateStatus = 'error';
				const errorData = await res.json();
				alert(`更新に失敗しました: ${errorData.error}`);
			}
		} catch (e) {
			console.error('Error updating attendance scores:', e);
			attendanceUpdateStatus = 'error';
			alert('サーバーとの通信に失敗しました。');
		}

		setTimeout(() => { attendanceUpdateStatus = ''; }, 3000);
	}

	const scoreCategories = [
        { name: '春スポ体合計点', key: 'init_score' },
        { name: '出席点', key: 'attendance_score' },
        { name: 'バレーボール1勝点', key: 'volleyball1_score' },
        { name: 'バレーボール2勝点', key: 'volleyball2_score' },
        { name: 'バレーボール3勝点', key: 'volleyball3_score' },
        { name: 'バレーボール優勝点', key: 'volleyball_championship_score' },
        { name: '卓球1勝点', key: 'table_tennis1_score' },
        { name: '卓球2勝点', key: 'table_tennis2_score' },
        { name: '卓球3勝点', key: 'table_tennis3_score' },
        { name: '卓球優勝点', key: 'table_tennis_championship_score' },
        { name: '卓球雨天ボーナス', key: 'table_tennis_rainy_bonus_score' },
        { name: 'サッカー1勝点', key: 'soccer1_score' },
        { name: 'サッカー2勝点', key: 'soccer2_score' },
        { name: 'サッカー3勝点', key: 'soccer3_score' },
        { name: 'サッカー優勝点', key: 'soccer_championship_score' },
        { name: 'リレーAブロック得点', key: 'relay_A_score' },
        { name: 'リレーBブロック得点', key: 'relay_B_score' },
        { name: 'リレーボーナス得点', key: 'relay_bonus_score' },
        { name: '合計(春スポ体合計点除く)', key: 'total_excluding_init' },
        { name: '合計(春スポ体合計点含む)', key: 'total_including_init' },
        { name: '現在の順位', key: 'current_rank' },
    ];

	// 順位を計算する関数
	/**
 	* スコアの配列を受け取り、順位（current_rank）を追加して返す関数
 	* @param {Array} scoresData 各クラスのスコアオブジェクトの配列
 	* @returns {Array} `current_rank` プロパティが追加された新しい配列
 	*/
	function calculateRanks(scoresData) {
    	// データが空なら何もしない
    	if (!scoresData || scoresData.length === 0) {
        	return [];
    	}

    	// --- 手順1: 合計点（total_including_init）で全クラスを並び替える ---
    	// これが「各クラスで比較して」の部分です。
    	// 元の配列を壊さないようにコピー(`...scoresData`)してからソートします。
    	const sorted = [...scoresData].sort((a, b) => b.total_including_init - a.total_including_init);
    
    	// --- 手順2: 順位を決定し、Mapオブジェクトに保存する ---
    	// これが「順位付けすれば」の部分です。
    	const rankMap = new Map();
    	let rank = 1;

    	for (let i = 0; i < sorted.length; i++) {
        	// 前のクラスよりスコアが低い場合、順位をその位置（i + 1）に更新します。
        	// 同じスコア（同点）の場合は、このif文に入らないため、同じ順位が維持されます。
        	if (i > 0 && sorted[i].total_including_init < sorted[i - 1].total_including_init) {
            	rank = i + 1;
        	}
        	// 「クラス名」をキーにして、決定した順位を保存します。
        	// 例: rankMap.set('1-1', 5);
        	rankMap.set(sorted[i].class_name, rank);
    	}

    	// --- 手順3: 元の配列の順番を維持したまま、順位を合体させる ---
    	// 画面の表示順を変えずに、計算した順位だけを追加します。
    	return scoresData.map(s => ({
        	...s, // 元のスコアデータはそのまま
        	current_rank: rankMap.get(s.class_name) // Mapからクラス名に対応する順位を取得して追加
    	}));
	}

	// スコアデータに順位を追加
	$: scoresWithRanks = scores.length > 0 ? calculateRanks(scores) : [];

	let showTotalScores = (typeof window !== 'undefined' && localStorage.getItem('showTotalScores')) === 'false' ? false : true;

	// A reactive statement to update localStorage whenever the switch is toggled.
	$: if (typeof window !== 'undefined') {
    	localStorage.setItem('showTotalScores', String(showTotalScores));
	}
</script>

<div class="dashboard-container">
	<header>
		<h1>ダッシュボード</h1>
		<nav class="dashboard-tabs">
			<button class:active={activeTab === 'tournament'} on:click={() => { activeTab = 'tournament'; fetchTournament('volleyball'); }}>競技トーナメント</button>
			<button class:active={activeTab === 'scores'} on:click={() => { activeTab = 'scores'; fetchScores(); }}>現在の得点</button>
			{#if userRole === 'superroot' || userRole === 'admin' || userRole === 'admin_relay'}
				<button class:active={activeTab === 'input'} on:click={handleInputTabClick}>試合結果入力</button>
			{/if}
		</nav>
		<button class="logout-btn" on:click={logout}>ログアウト</button>
	</header>

	<main>
		{#if activeTab === 'tournament'}
			<div class="sports-buttons">
				<button on:click={() => { relayActive = false; fetchTournament('volleyball'); }} class:active={selectedSport === 'volleyball' && !relayActive}>バレーボール</button>
				<button on:click={() => { relayActive = false; fetchTournament('table_tennis'); }} class:active={selectedSport === 'table_tennis' && !relayActive}>卓球</button>
				<button on:click={() => { relayActive = false; fetchTournament('soccer'); }} class:active={selectedSport === 'soccer' && !relayActive}>サッカー</button>
				<button on:click={() => { relayActive = true; fetchRelayResults('A'); }} class:active={relayActive}>リレー</button>
				{#if (userRole === 'superroot' || (userRole === 'admin' && assignedSport === 'table_tennis')) && selectedSport === 'table_tennis' && !relayActive}
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
			{#if relayActive}
				<div class="relay-area">
					<h2>学年対抗リレー結果</h2>
					<div class="relay-info">
						<p><strong>得点システム:</strong> 1位30点、2位25点、3位20点、4位15点、5位10点、6位5点</p>
						<p><strong>ボーナス得点:</strong> 両ブロック合計で1位30点、2位20点、3位10点の追加得点</p>
					</div>
					<div class="relay-type-selector">
						<button on:click={() => fetchRelayResults('A')} class:active={relayType === 'A'}>リレーAブロック</button>
						<button on:click={() => fetchRelayResults('B')} class:active={relayType === 'B'}>リレーBブロック</button>
					</div>
					{#if relayLoading}
						<p>読み込み中...</p>
					{:else if relayError}
						<p style="color:red">{relayError}</p>
					{:else}
						<RelayTable {relayResults} {relayType} {classNameMap} />
					{/if}
				</div>
			{:else}
				{#if isLoading}<p>読み込み中...</p>{/if}
				{#if !isLoading && selectedSport === 'table_tennis'}
					{#if getFilteredTableTennisTournaments().length > 1}
						<div class="tournament-selector">
							{#each getFilteredTableTennisTournaments() as t}
								<button on:click={() => selectTournament(t)} class:active={selectedTournament?.name === t.name}>{t.name}</button>
							{/each}
						</div>
					{/if}
				{:else if !isLoading && allTournaments.length > 1}
					<div class="tournament-selector">
						{#each allTournaments as t}
							<button on:click={() => selectTournament(t)} class:active={selectedTournament?.name === t.name}>{t.name}</button>
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

		{:else if activeTab === 'scores'}
    	<div class="scores-area">
        	<h2>現在の得点</h2>
        
        	{#if userRole === 'superroot'}
            	<div class="visibility-switcher">
                	<span>合計・順位の表示:</span>
                	<label class="switch">
                    	<input type="checkbox" bind:checked={showTotalScores}>
                    	<span class="slider"></span>
                	</label>
                	<span>{showTotalScores ? '表示中' : '非表示'}</span>
            	</div>
        	{/if}
        	{#if scoresLoading}
            	<p>読み込み中...</p>
        	{:else if scores.length > 0}
            	<div class="scores-container">
                	<div class="score-category-column">
                    	<div class="score-header">得点項目</div>
                    	{#each scoreCategories as category, i}
                        	{#if showTotalScores || (category.key !== 'total_excluding_init' && category.key !== 'total_including_init' && category.key !== 'current_rank')}
                            	<div class="score-cell" class:odd-row={i % 2 === 0}><b>{category.name}</b></div>
                        	{/if}
                    	{/each}
                	</div>
                	<div class="scores-data-wrapper">
                    	{#each scoresWithRanks as s}
                        	<div class="score-column">
                            	<div class="score-header">{s.class_name}</div>
                            	{#each scoreCategories as category, i}
                                	{#if showTotalScores || (category.key !== 'total_excluding_init' && category.key !== 'total_including_init' && category.key !== 'current_rank')}
                                    	<div class="score-cell" class:odd-row={i % 2 === 0} class:rank-cell={category.key === 'current_rank'}>
                                        	{#if category.key === 'current_rank'}
                                            	<span class="rank-badge rank-{s[category.key]}">{s[category.key]}位</span>
                                        	{:else}
                                            	{s[category.key]}
                                        	{/if}
                                    	</div>
                                	{/if}
                            	{/each}
                	        </div>
                    	{/each}
                	</div>
            	</div>
        	{:else}
            	<p>スコアデータがありません。</p>
        	{/if}
    	</div>

		{:else if activeTab === 'input'}
			<div class="match-input-area">
				<h2>試合結果入力</h2>
				
				<div class="input-tabs">
					{#if userRole === 'superroot' || userRole === 'admin'}
						<button class:active={activeInputTab === 'tournament'} on:click={() => activeInputTab = 'tournament'}>トーナメント</button>
					{/if}
					{#if userRole === 'superroot' || assignedSport === 'relay'}
						<button class:active={activeInputTab === 'relay'} on:click={() => { activeInputTab = 'relay'; fetchRelayResults(relayType); }}>リレー</button>
					{/if}
					{#if userRole === 'superroot'}
						<button class:active={activeInputTab === 'attendance'} on:click={() => activeInputTab = 'attendance'}>出席点</button>
					{/if}
				</div>

				{#if activeInputTab === 'tournament'}
					{#if userRole === 'admin'}<p>あなたの担当競技: <b>{assignedSport}</b></p>{/if}
					{#if userRole === 'superroot'}<p>全競技の試合結果を編集できます</p>{/if}
					{#if matchesLoading}<p>試合データ取得中...</p>{/if}

					{#each (userRole === 'superroot' ? ['volleyball', 'table_tennis', 'soccer'] : (assignedSport ? [assignedSport] : [])) as sport}
						<div class="tournament-edit-card">
							<h3>{sport} の試合一覧</h3>
							<div class="matches-list">
								{#if sport === 'table_tennis'}
									{#each getFilteredTableTennisMatches() as m (m.id)}
										<div class="match-entry">
											{#if m.status !== 'finished'}
												<form class="match-edit-form" on:submit|preventDefault={() => { showConfirmModal = true; editingMatch = m; editingSport = sport; }}>
													<div class="match-info">
														<span>ID: {m.id}</span>
														<span>R: {m.round}</span>
														<span>{m.tournament_name}</span>
														<span>{m.team1_name || '-'} vs {m.team2_name || '-'}</span>
													</div>
													<div class="score-inputs">
														<span>{m.team1_name}</span>
														<input type="number" min="0" bind:value={m.team1_score} class="score-input" required/>
														<span> - </span>
														<input type="number" min="0" bind:value={m.team2_score} class="score-input" required />
														<span>{m.team2_name}</span>
													</div>
													<button type="submit" class="update-btn">更新</button>
													{#if m._updateStatus}
														<span class="update-status {m._updateStatus}">{m._updateStatus === 'success' ? '✓' : '✗'}</span>
													{/if}
												</form>
											{:else}
												<div class="match-info finished">
													<span>ID: {m.id}</span>
													<span>R: {m.round}</span>
													<span>{m.tournament_name}</span>
													<span>{m.team1_name} {m.team1_score} - {m.team2_score} {m.team2_name}</span>
													<span class="update-status success">試合終了</span>
												</div>
											{/if}
										</div>
									{:else}
										<p>試合データがありません。</p>
									{/each}
								{:else}
									{#each matchesBySport[sport] as m (m.id)}
										<div class="match-entry">
											{#if m.status !== 'finished'}
												<form class="match-edit-form" on:submit|preventDefault={() => { showConfirmModal = true; editingMatch = m; editingSport = sport; }}>
													<div class="match-info">
														<span>ID: {m.id}</span>
														<span>R: {m.round}</span>
														<span>{m.team1_name || '-'} vs {m.team2_name || '-'}</span>
													</div>
													<div class="score-inputs">
														<span>{m.team1_name}</span>
														<input type="number" min="0" bind:value={m.team1_score} class="score-input" required/>
														<span> - </span>
														<input type="number" min="0" bind:value={m.team2_score} class="score-input" required />
														<span>{m.team2_name}</span>
													</div>
													<button type="submit" class="update-btn">更新</button>
													{#if m._updateStatus}
														<span class="update-status {m._updateStatus}">{m._updateStatus === 'success' ? '✓' : '✗'}</span>
													{/if}
												</form>
											{:else}
												<div class="match-info finished">
													<span>ID: {m.id}</span>
													<span>R: {m.round}</span>
													<span>{m.tournament_name}</span>
													<span>{m.team1_name} {m.team1_score} - {m.team2_score} {m.team2_name}</span>
													<span class="update-status success">試合終了</span>
												</div>
											{/if}
										</div>
									{:else}
										<p>試合データがありません。</p>
									{/each}
								{/if}
							</div>
						</div>
					{/each}

				{:else if activeInputTab === 'relay'}
					<div class="relay-input-card">
						<h3>リレー結果入力 (ドラッグ&ドロップで順位変更)</h3>
						<div class="relay-type-selector">
							<button on:click={() => fetchRelayResults('A')} class:active={relayType === 'A'}>リレーA</button>
							<button on:click={() => fetchRelayResults('B')} class:active={relayType === 'B'}>リレーB</button>
						</div>
						{#if relayLoading}
							<p>読み込み中...</p>
						{:else if relayError}
							<p style="color:red;">{relayError}</p>
						{:else if relayResults.length > 0}
							<section 
								class="dnd-list" 
								use:dndzone={{
									items: relayResults, 
									flipDurationMs: 300, 
									dropTargetStyle: {},
									dragDisabled: relayLoading || relayUpdateStatus === 'loading'
								}} 
								on:consider={handleDndConsider} 
								on:finalize={handleDndFinalize}
							>
								{#each relayResults as result, i (result.id)}
									<div class="dnd-item" data-id={result.id}>
										<span class="rank-badge">{result.rank || i + 1}位</span>
										<span class="class-name">{gradeNames[result.grade] || `学年${result.grade}`}</span>
									</div>
								{/each}
							</section>
							<button 
								on:click={updateRelayRanks} 
								class="update-btn" 
								disabled={relayUpdateStatus === 'loading'}
							>
								{relayUpdateStatus === 'loading' ? '更新中...' : 'リレー結果を更新'}
							</button>
							{#if relayUpdateStatus === 'success'}
								<span class="update-status success">更新しました！</span>
							{:else if relayUpdateStatus === 'error'}
								<span class="update-status error">更新に失敗しました。</span>
							{:else if relayUpdateStatus === 'loading'}
								<span class="update-status loading">更新中...</span>
							{/if}
						{:else}
							<p>リレーデータがありません。</p>
						{/if}
					</div>

				{:else if activeInputTab === 'attendance'}
					<AttendanceInput 
						{attendanceScores}
						loading={attendanceLoading}
						updateStatus={attendanceUpdateStatus}
						on:scoreChange={handleAttendanceScoreChange}
						on:batchUpdate={handleAttendanceBatchUpdate}
					/>
				{/if}
			</div>
		{/if}
	</main>

	{#if showConfirmModal}
    	<div class="modal-overlay" on:click={() => showConfirmModal = false}>
        	<div class="modal-content" on:click|stopPropagation>
            	<h3>試合結果を更新しますか？</h3>
				<p>ID: {editingMatch.id}, R: {editingMatch.round}</p>
				<p>{editingMatch.team1_name || '-'} vs {editingMatch.team2_name || '-'}</p>
				<p>スコア: {editingMatch.team1_score} - {editingMatch.team2_score}</p>
            	<p>本当にこの内容で更新してよろしいですか？</p>
            	<div class="modal-actions">
                	<button on:click={() => { showConfirmModal = false; editingMatch = null; editingSport = ''; }}>キャンセル</button>
					<button on:click={confirmUpdateMatchScore} class="ok-btn">OK</button>
            	</div>
        	</div>
    	</div>
	{/if}
</div>

<style>
	/* General */
	.dashboard-container { padding: 2rem; font-family: sans-serif; }
	main { text-align: center; }
	h1 { font-size: 2.5rem; margin-right: auto; }
	h2 { margin-top: 2rem; }
	h3 { margin-bottom: 1rem; }
	button { padding: 0.5rem 1rem; color: white; background-color: #5a5a5a; border: none; border-radius: 4px; cursor: pointer; transition: background-color 0.3s; }
	button.active { background-color: #4285f4; }
	button:hover { opacity: 0.9; }
	.update-btn { 
		background: #4285f4; 
		margin-top: 1.5rem; 
		transition: all 0.2s ease;
	}
	.update-btn:hover:not(:disabled) { 
		background: #3367d6; 
		transform: translateY(-1px);
	}
	.update-btn:disabled { 
		background: #ccc; 
		cursor: not-allowed; 
		opacity: 0.6;
	}
	.update-status { 
		font-size: 0.95rem; 
		margin-left: 1rem; 
		font-weight: 500;
	}
	.update-status.success { color: #43a047; }
	.update-status.error { color: #e53935; }
	.update-status.loading { color: #ff9800; }

	/* Header */
	header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; padding-bottom: 1rem; border-bottom: 1px solid #ccc; flex-wrap: wrap; }
	.dashboard-tabs { margin: 0 auto; }
	.logout-btn { background-color: #d93025; margin-left: 1rem; }
	.logout-btn:hover { background-color: #c5221b; }

	/* Content Sections */
	.sports-buttons, .tournament-selector, .relay-type-selector, .input-tabs { margin-bottom: 1.5rem; display: flex; align-items: center; justify-content: center; gap: 1rem; flex-wrap: wrap; }
	.bracket-area { margin-top: 2rem; overflow-x: auto; }
	.bracket-wrapper { display: flex; align-items: center; justify-content: center; min-width: 800px; }
	.scores-area { width: 100%; }
	.match-input-area { margin-top: 2rem; display: flex; flex-direction: column; align-items: center; width: 100%; }
	.relay-area { margin-top: 2rem; }
	.relay-info { background: #f8f9fa; border-radius: 8px; padding: 1rem; margin: 1rem 0; max-width: 600px; margin-left: auto; margin-right: auto; }
	.relay-info p { margin: 0.5rem 0; font-size: 0.9rem; color: #666; }

	/* Weather Switcher */
	.weather-switcher { display: flex; align-items: center; gap: 8px; margin-left: 2rem; }
	.switch { position: relative; display: inline-block; width: 50px; height: 24px; }
	.switch input { opacity: 0; width: 0; height: 0; }
	.slider { position: absolute; cursor: pointer; top: 0; left: 0; right: 0; bottom: 0; background-color: #ccc; transition: .4s; border-radius: 24px; }
	.slider:before { position: absolute; content: ""; height: 18px; width: 18px; left: 3px; bottom: 3px; background-color: white; transition: .4s; border-radius: 50%; }
	input:checked + .slider { background-color: #2196F3; }
	input:checked + .slider:before { transform: translateX(26px); }

	/* Bracketry Global Styles */
	:global(.bracket-match-team) { background-color: #f0f0f0 !important; border: 1px solid #ccc !important; }
	:global(.bracket-match-winner .bracket-match-team) { background-color: #d4edda !important; font-weight: bold; }
	:global(.bracket-connector) { border-color: #999 !important; }

	/* Scores Table */
	.scores-container { display: flex; width: 100%; border: 1px solid #ddd; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
	.score-category-column { position: sticky; left: 0; z-index: 1; background-color: #fff; flex-shrink: 0; border-right: 1px solid #ddd; }
	.scores-data-wrapper { display: flex; overflow-x: auto; }
	.score-column { flex: 0 0 120px; display: flex; flex-direction: column; border-left: 1px solid #ddd; }
	.score-header { padding: 10px; text-align: center; font-weight: bold; background-color: #f2f2f2; border-bottom: 1px solid #ddd; }
	.score-cell { padding: 10px; text-align: center; border-bottom: 1px solid #ddd; min-width: 150px; }
	.score-column .score-cell { min-width: auto; }
	.score-cell.odd-row { background-color: #f9f9f9; }

	/* 順位表示のスタイル */
	.rank-cell {
		text-align: center;
		padding: 8px !important;
	}

	.rank-badge {
		display: inline-block;
		padding: 0.3rem 0.6rem;
		border-radius: 12px;
		font-weight: bold;
		color: white;
		font-size: 0.85rem;
		min-width: 40px;
	}

	.rank-badge.rank-1 {
		background: linear-gradient(135deg, #ffd700, #ffb300);
		box-shadow: 0 2px 4px rgba(255, 179, 0, 0.3);
	}

	.rank-badge.rank-2 {
		background: linear-gradient(135deg, #c0c0c0, #a0a0a0);
		box-shadow: 0 2px 4px rgba(160, 160, 160, 0.3);
	}

	.rank-badge.rank-3 {
		background: linear-gradient(135deg, #cd7f32, #b8860b);
		box-shadow: 0 2px 4px rgba(184, 134, 11, 0.3);
	}

	.rank-badge.rank-4,
	.rank-badge.rank-5,
	.rank-badge.rank-6,
	.rank-badge.rank-7,
	.rank-badge.rank-8,
	.rank-badge.rank-9,
	.rank-badge.rank-10,
	.rank-badge.rank-11,
	.rank-badge.rank-12,
	.rank-badge.rank-13,
	.rank-badge.rank-14,
	.rank-badge.rank-15,
	.rank-badge.rank-16 {
		background: linear-gradient(135deg, #6c757d, #5a6268);
		box-shadow: 0 2px 4px rgba(90, 98, 104, 0.3);
	}

	/* Match/Relay Input Cards */
	.tournament-edit-card, .relay-input-card { background: #f8f9fa; border-radius: 12px; box-shadow: 0 2px 8px rgba(66,133,244,0.08); padding: 2rem; margin-bottom: 2rem; width: 100%; max-width: 800px; }
	.match-edit-form { display: flex; flex-wrap: wrap; align-items: center; gap: 1rem; background: #fff; border-radius: 8px; padding: 1.5rem; margin-bottom: 1rem; }
	.match-info { font-size: 0.9rem; color: #333; display: flex; flex-wrap: wrap; gap: 1rem; flex-grow: 1; }
	.match-info.finished { background: #f0f0f0; padding: 1rem; border-radius: 6px; width: 100%; }
	.score-inputs { display: flex; align-items: center; gap: 0.5rem; }
	.score-input { width: 4rem; padding: 0.5rem; font-size: 1.1rem; border: 1px solid #4285f4; border-radius: 4px; text-align: center; }
	
	/* Relay D&D List */
	.dnd-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		max-width: 400px;
		margin: 1.5rem auto;
	}
	.dnd-item {
		display: flex;
		align-items: center;
		padding: 1rem;
		background-color: white;
		border: 1px solid #ddd;
		border-radius: 6px;
		box-shadow: 0 1px 3px rgba(0,0,0,0.05);
		cursor: grab;
		user-select: none;
		transition: all 0.2s ease;
		position: relative;
	}
	.dnd-item:hover { 
		box-shadow: 0 2px 8px rgba(0,0,0,0.1); 
		transform: translateY(-1px);
	}
	.dnd-item:active { 
		cursor: grabbing; 
		background-color: #f0f8ff; 
		transform: scale(1.02);
	}
	.rank-badge { font-weight: bold; font-size: 1rem; color: #fff; background-color: #6c757d; border-radius: 4px; padding: 0.3rem 0.6rem; margin-right: 1rem; }
	.class-name { font-size: 1.1rem; }

	/* Modal */
	.modal-overlay { position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
	.modal-content { background: #fff; border-radius: 8px; padding: 2rem; width: 90%; max-width: 500px; text-align: center; }
	.modal-actions { display: flex; gap: 1rem; justify-content: center; margin-top: 1.5rem; }
	.modal-actions button { background: #ccc; }
	.modal-actions .ok-btn { background: #4285f4; }

	/* Responsive */
	@media (max-width: 768px) {
		header { flex-direction: column; align-items: stretch; gap: 1rem; }
		h1 { text-align: center; margin-right: 0; }
		.dashboard-tabs { order: 2; width: 100%; display: flex; justify-content: center; }
		.logout-btn { order: 3; margin-left: 0; width: 100%; }
		.match-edit-form { flex-direction: column; align-items: stretch; }
	}
</style>
