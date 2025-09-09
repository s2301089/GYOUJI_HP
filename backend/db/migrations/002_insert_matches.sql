-- ----------------------------------------------------------------
-- トーナメント情報の登録
-- ----------------------------------------------------------------
INSERT INTO tournaments (name, sport, weather_condition) VALUES
('バレーボール', 'volleyball', 'any'),
('卓球（晴天時）', 'table_tennis', 'sunny'),
('卓球（雨天時・トーナメント）', 'table_tennis', 'rainy'),
('卓球（雨天時・敗者戦）', 'table_tennis', 'rainy'), -- 敗者戦用のトーナメント（試合の登録はしない）
('8人制サッカー', 'soccer', 'any');

-- ----------------------------------------------------------------
-- バレーボール: チームと試合の登録
-- ----------------------------------------------------------------
-- トーナメントIDを変数に格納
SET @volleyball_tour_id = (SELECT id FROM tournaments WHERE name = 'バレーボール');

-- チームの登録
INSERT INTO teams (name, tournament_id) VALUES
('専・教', @volleyball_tour_id), ('IE4', @volleyball_tour_id), ('IS5', @volleyball_tour_id), ('IT4', @volleyball_tour_id),
('IT3', @volleyball_tour_id), ('IT2', @volleyball_tour_id), ('1-1', @volleyball_tour_id), ('IE2', @volleyball_tour_id),
('IS3', @volleyball_tour_id), ('IS2', @volleyball_tour_id), ('IS4', @volleyball_tour_id), ('IE5', @volleyball_tour_id),
('1-2', @volleyball_tour_id), ('1-3', @volleyball_tour_id), ('IE3', @volleyball_tour_id), ('IT5', @volleyball_tour_id);

-- 1回戦の試合を登録
INSERT INTO matches (tournament_id, round, match_number_in_round, team1_id, team2_id) VALUES
(@volleyball_tour_id, 1, 1, (SELECT id FROM teams WHERE name = '専・教' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = 'IE4' AND tournament_id = @volleyball_tour_id)),
(@volleyball_tour_id, 1, 2, (SELECT id FROM teams WHERE name = 'IS5' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = 'IT4' AND tournament_id = @volleyball_tour_id)),
(@volleyball_tour_id, 1, 3, (SELECT id FROM teams WHERE name = 'IT3' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = 'IT2' AND tournament_id = @volleyball_tour_id)),
(@volleyball_tour_id, 1, 4, (SELECT id FROM teams WHERE name = '1-1' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = 'IE2' AND tournament_id = @volleyball_tour_id)),
(@volleyball_tour_id, 1, 5, (SELECT id FROM teams WHERE name = 'IS3' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = 'IS2' AND tournament_id = @volleyball_tour_id)),
(@volleyball_tour_id, 1, 6, (SELECT id FROM teams WHERE name = 'IS4' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = 'IE5' AND tournament_id = @volleyball_tour_id)),
(@volleyball_tour_id, 1, 7, (SELECT id FROM teams WHERE name = '1-2' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = '1-3' AND tournament_id = @volleyball_tour_id)),
(@volleyball_tour_id, 1, 8, (SELECT id FROM teams WHERE name = 'IE3' AND tournament_id = @volleyball_tour_id), (SELECT id FROM teams WHERE name = 'IT5' AND tournament_id = @volleyball_tour_id));

-- ----------------------------------------------------------------
-- 卓球（晴天時）: チームと試合の登録
-- ----------------------------------------------------------------
SET @tt_sunny_tour_id = (SELECT id FROM tournaments WHERE name = '卓球（晴天時）');

INSERT INTO teams (name, tournament_id) VALUES
('1-2', @tt_sunny_tour_id), ('IE5', @tt_sunny_tour_id), ('IE3', @tt_sunny_tour_id), ('IT2', @tt_sunny_tour_id),
('1-3', @tt_sunny_tour_id), ('IE4', @tt_sunny_tour_id), ('IT3', @tt_sunny_tour_id), ('IT4', @tt_sunny_tour_id),
('IS5', @tt_sunny_tour_id), ('IE2', @tt_sunny_tour_id), ('1-1', @tt_sunny_tour_id), ('IS3', @tt_sunny_tour_id),
('IS2', @tt_sunny_tour_id), ('IS4', @tt_sunny_tour_id), ('IT5', @tt_sunny_tour_id), ('専・教', @tt_sunny_tour_id);

INSERT INTO matches (tournament_id, round, match_number_in_round, team1_id, team2_id) VALUES
(@tt_sunny_tour_id, 1, 1, (SELECT id FROM teams WHERE name = '1-2' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = 'IE5' AND tournament_id = @tt_sunny_tour_id)),
(@tt_sunny_tour_id, 1, 2, (SELECT id FROM teams WHERE name = 'IE3' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = 'IT2' AND tournament_id = @tt_sunny_tour_id)),
(@tt_sunny_tour_id, 1, 3, (SELECT id FROM teams WHERE name = '1-3' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = 'IE4' AND tournament_id = @tt_sunny_tour_id)),
(@tt_sunny_tour_id, 1, 4, (SELECT id FROM teams WHERE name = 'IT3' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = 'IT4' AND tournament_id = @tt_sunny_tour_id)),
(@tt_sunny_tour_id, 1, 5, (SELECT id FROM teams WHERE name = 'IS5' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = 'IE2' AND tournament_id = @tt_sunny_tour_id)),
(@tt_sunny_tour_id, 1, 6, (SELECT id FROM teams WHERE name = '1-1' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = 'IS3' AND tournament_id = @tt_sunny_tour_id)),
(@tt_sunny_tour_id, 1, 7, (SELECT id FROM teams WHERE name = 'IS2' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = 'IS4' AND tournament_id = @tt_sunny_tour_id)),
(@tt_sunny_tour_id, 1, 8, (SELECT id FROM teams WHERE name = 'IT5' AND tournament_id = @tt_sunny_tour_id), (SELECT id FROM teams WHERE name = '専・教' AND tournament_id = @tt_sunny_tour_id));

-- ----------------------------------------------------------------
-- 卓球（雨天時・トーナメント）: チームと試合の登録
-- ----------------------------------------------------------------
SET @tt_rainy_tour_id = (SELECT id FROM tournaments WHERE name = '卓球（雨天時・トーナメント）');

-- 元データではIS3が2回登場するため、晴天時トーナメントに合わせて IE3 vs IT2 の組み合わせと解釈してチームを登録します。
INSERT INTO teams (name, tournament_id) VALUES
('1-2', @tt_rainy_tour_id), ('IE5', @tt_rainy_tour_id), ('IE3', @tt_rainy_tour_id), ('IT2', @tt_rainy_tour_id),
('1-3', @tt_rainy_tour_id), ('IE4', @tt_rainy_tour_id), ('IT3', @tt_rainy_tour_id), ('IT4', @tt_rainy_tour_id),
('IS5', @tt_rainy_tour_id), ('IE2', @tt_rainy_tour_id), ('1-1', @tt_rainy_tour_id), ('IS3', @tt_rainy_tour_id),
('IS2', @tt_rainy_tour_id), ('IS4', @tt_rainy_tour_id), ('IT5', @tt_rainy_tour_id), ('専・教', @tt_rainy_tour_id);

INSERT INTO matches (tournament_id, round, match_number_in_round, team1_id, team2_id) VALUES
(@tt_rainy_tour_id, 1, 1, (SELECT id FROM teams WHERE name = '1-2' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = 'IE5' AND tournament_id = @tt_rainy_tour_id)),
(@tt_rainy_tour_id, 1, 2, (SELECT id FROM teams WHERE name = 'IE3' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = 'IT2' AND tournament_id = @tt_rainy_tour_id)),
(@tt_rainy_tour_id, 1, 3, (SELECT id FROM teams WHERE name = '1-3' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = 'IE4' AND tournament_id = @tt_rainy_tour_id)),
(@tt_rainy_tour_id, 1, 4, (SELECT id FROM teams WHERE name = 'IT3' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = 'IT4' AND tournament_id = @tt_rainy_tour_id)),
(@tt_rainy_tour_id, 1, 5, (SELECT id FROM teams WHERE name = 'IS5' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = 'IE2' AND tournament_id = @tt_rainy_tour_id)),
(@tt_rainy_tour_id, 1, 6, (SELECT id FROM teams WHERE name = '1-1' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = 'IS3' AND tournament_id = @tt_rainy_tour_id)),
(@tt_rainy_tour_id, 1, 7, (SELECT id FROM teams WHERE name = 'IS2' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = 'IS4' AND tournament_id = @tt_rainy_tour_id)),
(@tt_rainy_tour_id, 1, 8, (SELECT id FROM teams WHERE name = 'IT5' AND tournament_id = @tt_rainy_tour_id), (SELECT id FROM teams WHERE name = '専・教' AND tournament_id = @tt_rainy_tour_id));

-- ----------------------------------------------------------------
-- 8人制サッカー: チームと試合の登録
-- ----------------------------------------------------------------
SET @soccer_tour_id = (SELECT id FROM tournaments WHERE name = '8人制サッカー');

INSERT INTO teams (name, tournament_id) VALUES
('IS3', @soccer_tour_id), ('IE2', @soccer_tour_id), ('1-1', @soccer_tour_id), ('IS2', @soccer_tour_id),
('IS4', @soccer_tour_id), ('IT5', @soccer_tour_id), ('IS5', @soccer_tour_id), ('専・教', @soccer_tour_id),
('1-2', @soccer_tour_id), ('1-3', @soccer_tour_id), ('IE3', @soccer_tour_id), ('IT4', @soccer_tour_id),
('IT3', @soccer_tour_id), ('IE4', @soccer_tour_id), ('IT2', @soccer_tour_id), ('IE5', @soccer_tour_id);

INSERT INTO matches (tournament_id, round, match_number_in_round, team1_id, team2_id) VALUES
(@soccer_tour_id, 1, 1, (SELECT id FROM teams WHERE name = 'IS3' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = 'IE2' AND tournament_id = @soccer_tour_id)),
(@soccer_tour_id, 1, 2, (SELECT id FROM teams WHERE name = '1-1' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = 'IS2' AND tournament_id = @soccer_tour_id)),
(@soccer_tour_id, 1, 3, (SELECT id FROM teams WHERE name = 'IS4' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = 'IT5' AND tournament_id = @soccer_tour_id)),
(@soccer_tour_id, 1, 4, (SELECT id FROM teams WHERE name = 'IS5' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = '専・教' AND tournament_id = @soccer_tour_id)),
(@soccer_tour_id, 1, 5, (SELECT id FROM teams WHERE name = '1-2' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = '1-3' AND tournament_id = @soccer_tour_id)),
(@soccer_tour_id, 1, 6, (SELECT id FROM teams WHERE name = 'IE3' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = 'IT4' AND tournament_id = @soccer_tour_id)),
(@soccer_tour_id, 1, 7, (SELECT id FROM teams WHERE name = 'IT3' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = 'IE4' AND tournament_id = @soccer_tour_id)),
(@soccer_tour_id, 1, 8, (SELECT id FROM teams WHERE name = 'IT2' AND tournament_id = @soccer_tour_id), (SELECT id FROM teams WHERE name = 'IE5' AND tournament_id = @soccer_tour_id));
