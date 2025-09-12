-- チームごとのポイント内訳を記録するテーブル
CREATE TABLE IF NOT EXISTS team_points (
  id INT AUTO_INCREMENT PRIMARY KEY,
  team_id INT NOT NULL,
  sport VARCHAR(255) NOT NULL,
  tournament_name VARCHAR(255) NOT NULL,
  point_type ENUM('win','final_bonus_winner','final_bonus_runnerup','bronze_bonus_winner','bronze_bonus_runnerup','rainy_loser_champion') NOT NULL,
  points INT NOT NULL,
  source_match_id INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uniq_team_match_type (team_id, source_match_id, point_type),
  FOREIGN KEY (team_id) REFERENCES teams(id),
  FOREIGN KEY (source_match_id) REFERENCES matches(id)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


