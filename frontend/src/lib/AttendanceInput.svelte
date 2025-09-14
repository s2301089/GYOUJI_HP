<script>
  export let attendanceScores = [];
  export let loading = false;
  export let updateStatus = '';

  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();

  // 学年別にグループ化
  $: groupedScores = groupByGrade(attendanceScores);

  function groupByGrade(scores) {
    const groups = {
      1: { name: '1年生', classes: [] },
      2: { name: '2年生', classes: [] },
      3: { name: '3年生', classes: [] },
      4: { name: '4年生', classes: [] },
      5: { name: '5年生', classes: [] },
      6: { name: '専・教', classes: [] }
    };

    scores.forEach(score => {
      const grade = Math.floor(score.class_id / 10);
      if (groups[grade]) {
        groups[grade].classes.push(score);
      }
    });

    return groups;
  }

  function handleScoreChange(classId, newScore) {
    const score = parseInt(newScore);
    if (score < 0 || score > 10) return;
    
    dispatch('scoreChange', { classId, score });
  }

  function handleBatchUpdate() {
    dispatch('batchUpdate');
  }

  function setGradeScore(grade, score) {
    const gradeClasses = groupedScores[grade].classes;
    gradeClasses.forEach(classData => {
      handleScoreChange(classData.class_id, score);
    });
  }
</script>

<div class="attendance-input-container">
  <div class="header">
    <h3>出席点入力</h3>
    <p class="description">各クラスの出席点を0〜10点の範囲で入力してください。</p>
  </div>

  {#if loading}
    <div class="loading">読み込み中...</div>
  {:else}
    <div class="grades-container">
      {#each Object.entries(groupedScores) as [grade, gradeData]}
        {#if gradeData.classes.length > 0}
          <div class="grade-section">
            <div class="grade-header">
              <h4>{gradeData.name}</h4>
              <div class="grade-controls">
                <span>一括設定:</span>
                {#each [0, 5, 8, 10] as score}
                  <button 
                    class="quick-set-btn" 
                    on:click={() => setGradeScore(grade, score)}
                    disabled={loading}
                  >
                    {score}点
                  </button>
                {/each}
              </div>
            </div>
            
            <div class="classes-grid">
              {#each gradeData.classes as classData}
                <div class="class-item">
                  <label for="score-{classData.class_id}">{classData.class_name}</label>
                  <input
                    id="score-{classData.class_id}"
                    type="number"
                    min="0"
                    max="10"
                    bind:value={classData.score}
                    on:input={(e) => handleScoreChange(classData.class_id, e.target.value)}
                    class="score-input"
                    disabled={loading}
                  />
                  <span class="unit">点</span>
                </div>
              {/each}
            </div>
          </div>
        {/if}
      {/each}
    </div>

    <div class="actions">
      <button 
        class="update-btn" 
        on:click={handleBatchUpdate}
        disabled={loading}
      >
        {loading ? '更新中...' : '一括更新'}
      </button>
      
      {#if updateStatus === 'success'}
        <span class="update-status success">✓ 更新完了</span>
      {:else if updateStatus === 'error'}
        <span class="update-status error">✗ 更新失敗</span>
      {/if}
    </div>
  {/if}
</div>

<style>
.attendance-input-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 1.5rem;
}

.header {
  text-align: center;
  margin-bottom: 2rem;
}

.header h3 {
  color: #333;
  margin-bottom: 0.5rem;
}

.description {
  color: #666;
  font-size: 0.9rem;
}

.loading {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.grades-container {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.grade-section {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 1.5rem;
  background: white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.grade-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #f0f0f0;
}

.grade-header h4 {
  margin: 0;
  color: #333;
  font-size: 1.2rem;
}

.grade-controls {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
}

.quick-set-btn {
  padding: 0.3rem 0.6rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: black;
  color: white;
  cursor: pointer;
  font-size: 0.8rem;
  transition: all 0.2s ease;
}

.quick-set-btn:hover:not(:disabled) {
  background: #e9ecef;
  border-color: #adb5bd;
}

.quick-set-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.classes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.class-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  background: #fafafa;
}

.class-item label {
  font-weight: 500;
  color: #495057;
  min-width: 60px;
}

.score-input {
  width: 60px;
  padding: 0.4rem;
  border: 1px solid #ced4da;
  border-radius: 4px;
  text-align: center;
  font-size: 1rem;
}

.score-input:focus {
  outline: none;
  border-color: #4285f4;
  box-shadow: 0 0 0 2px rgba(66, 133, 244, 0.2);
}

.score-input:disabled {
  background: #e9ecef;
  cursor: not-allowed;
}

.unit {
  color: #6c757d;
  font-size: 0.9rem;
}

.actions {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid #e0e0e0;
}

.update-btn {
  padding: 0.75rem 2rem;
  background: #4285f4;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.update-btn:hover:not(:disabled) {
  background: #3367d6;
  transform: translateY(-1px);
}

.update-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
}

.update-status {
  font-weight: 500;
  padding: 0.5rem 1rem;
  border-radius: 4px;
}

.update-status.success {
  color: #155724;
  background: #d4edda;
  border: 1px solid #c3e6cb;
}

.update-status.error {
  color: #721c24;
  background: #f8d7da;
  border: 1px solid #f5c6cb;
}

/* レスポンシブ対応 */
@media (max-width: 768px) {
  .grade-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
  
  .grade-controls {
    align-self: stretch;
    justify-content: space-between;
  }
  
  .classes-grid {
    grid-template-columns: 1fr;
  }
  
  .class-item {
    justify-content: space-between;
  }
}
</style>