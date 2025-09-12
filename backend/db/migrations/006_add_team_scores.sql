-- teams テーブルに初期点(init_score)と出席点(attendance_score)を追加
ALTER TABLE teams
  ADD COLUMN init_score INT DEFAULT 0 AFTER entry_status,
  ADD COLUMN attendance_score INT DEFAULT 0 AFTER init_score;


