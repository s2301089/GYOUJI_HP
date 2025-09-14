-- ----------------------------------------------------------------
-- 6. 学年対抗リレーAブロック: チームと試合の登録
-- ----------------------------------------------------------------

-- Aブロックの処理
SET @relay_a_tour_id = (SELECT id FROM tournaments WHERE name = '学年対抗リレーAブロック');

-- チーム登録 (Aブロック)
INSERT INTO teams (class_id, name, tournament_id) VALUES
(1, '一年生', @relay_a_tour_id),
(2, '二年生', @relay_a_tour_id),
(3, '三年生', @relay_a_tour_id),
(4, '四年生', @relay_a_tour_id),
(5, '五年生', @relay_a_tour_id),
(6, '専・教', @relay_a_tour_id);

-- 試合枠の作成 (Aブロック決勝戦)
INSERT INTO matches (tournament_id, round, match_number_in_round, status) VALUES
(@relay_a_tour_id, 1, 1, 'scheduled');

-- Bブロックの処理
SET @relay_b_tour_id = (SELECT id FROM tournaments WHERE name = '学年対抗リレーBブロック');

-- チーム登録 (Bブロック)
INSERT INTO teams (class_id, name, tournament_id) VALUES
(1, '一年生', @relay_b_tour_id),
(2, '二年生', @relay_b_tour_id),
(3, '三年生', @relay_b_tour_id),
(4, '四年生', @relay_b_tour_id),
(5, '五年生', @relay_b_tour_id),
(6, '専・教', @relay_b_tour_id);

-- 試合枠の作成 (Bブロック決勝戦)
INSERT INTO matches (tournament_id, round, match_number_in_round, status) VALUES
(@relay_b_tour_id, 1, 1, 'scheduled');
