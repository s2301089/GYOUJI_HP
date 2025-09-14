-- 学年対抗リレー順位登録用テーブル
CREATE TABLE relay_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    relay_type ENUM('A', 'B') NOT NULL,
    relay_rank INT NOT NULL, 
    class_id INT NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
