<script>
  export let relayResults = [];
  export let relayType = '';
  export let classNameMap = {};
  
  // 学年名のマッピング
  const gradeNames = {
    1: '1年生',
    2: '2年生', 
    3: '3年生',
    4: '4年生',
    5: '5年生',
    6: '専・教'
  };

  // 得点のマッピング
  const scoreMapping = {
    1: 30, // 1位: 30点
    2: 25, // 2位: 25点
    3: 20, // 3位: 20点
    4: 15, // 4位: 15点
    5: 10, // 5位: 10点
    6: 5   // 6位: 5点
  };

  // リレー結果が入力されているかチェック（gradeがnullでない、かつ有効な値）
  $: hasRelayData = relayResults.length > 0 && relayResults.some(result => result.grade !== null && result.grade !== undefined);
</script>

<div class="relay-results-container">
  <h3>リレー{relayType}ブロック結果</h3>
  <table class="relay-table">
    <thead>
      <tr>
        <th>順位</th>
        <th>学年</th>
        <th>獲得得点</th>
      </tr>
    </thead>
    <tbody>
      {#if hasRelayData}
        {#each relayResults as result}
          <tr class="rank-{result.rank}">
            <td class="rank-cell">
              <span class="rank-badge rank-{result.rank}">{result.rank}位</span>
            </td>
            <td class="grade-cell">{gradeNames[result.grade] || `学年${result.grade}`}</td>
            <td class="score-cell">{scoreMapping[result.rank] || 0}点</td>
          </tr>
        {/each}
      {:else}
        {#each Array(6) as _, i}
          <tr class="rank-{i + 1} no-data">
            <td class="rank-cell">
              <span class="rank-badge rank-{i + 1}">{i + 1}位</span>
            </td>
            <td class="grade-cell empty">-</td>
            <td class="score-cell">{scoreMapping[i + 1] || 0}点</td>
          </tr>
        {/each}
      {/if}
    </tbody>
  </table>
  {#if !hasRelayData}
    <p class="no-data-message">リレー結果はまだ入力されていません。</p>
  {/if}
</div>

<style>
.relay-results-container {
  max-width: 600px;
  margin: 0 auto;
  padding: 1rem;
}

.relay-table {
  border-collapse: collapse;
  width: 100%;
  margin-top: 1rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  border-radius: 8px;
  overflow: hidden;
}

.relay-table th {
  background: linear-gradient(135deg, #4285f4, #34a853);
  color: white;
  padding: 1rem;
  text-align: center;
  font-weight: bold;
  font-size: 1.1rem;
}

.relay-table td {
  padding: 1rem;
  text-align: center;
  border-bottom: 1px solid #e0e0e0;
}

.relay-table tbody tr:hover {
  background-color: #f8f9fa;
}

.relay-table tbody tr:last-child td {
  border-bottom: none;
}

.rank-cell {
  width: 100px;
}

.rank-badge {
  display: inline-block;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-weight: bold;
  color: white;
  font-size: 0.9rem;
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
.rank-badge.rank-6 {
  background: linear-gradient(135deg, #6c757d, #5a6268);
  box-shadow: 0 2px 4px rgba(90, 98, 104, 0.3);
}

.grade-cell {
  font-size: 1.1rem;
  font-weight: 500;
}

.score-cell {
  font-size: 1.2rem;
  font-weight: bold;
  color: #4285f4;
}

.rank-1 .score-cell {
  color: #ffd700;
}

.rank-2 .score-cell {
  color: #c0c0c0;
}

.rank-3 .score-cell {
  color: #cd7f32;
}

/* データが未入力の場合のスタイル */
.no-data {
  opacity: 0.6;
}

.grade-cell.empty {
  color: #999;
  font-style: italic;
  font-weight: normal;
}

.no-data-message {
  text-align: center;
  color: #666;
  font-style: italic;
  margin-top: 1rem;
  padding: 1rem;
  background-color: #f8f9fa;
  border-radius: 4px;
}

/* レスポンシブ対応 */
@media (max-width: 768px) {
  .relay-table th,
  .relay-table td {
    padding: 0.75rem 0.5rem;
    font-size: 0.9rem;
  }
  
  .rank-badge {
    padding: 0.4rem 0.8rem;
    font-size: 0.8rem;
  }
  
  .grade-cell {
    font-size: 1rem;
  }
  
  .score-cell {
    font-size: 1.1rem;
  }
  
  .no-data-message {
    font-size: 0.9rem;
    padding: 0.75rem;
  }
}
</style>
