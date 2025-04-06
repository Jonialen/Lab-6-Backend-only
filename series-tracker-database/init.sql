-- Creación de la base de datos
CREATE DATABASE IF NOT EXISTS anime_db;

-- Usar la base de datos creada
USE anime_db;

-- Creación de la tabla series
CREATE TABLE IF NOT EXISTS series (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    status ENUM('Plan to Watch', 'Watching', 'Completed', 'Dropped') NOT NULL DEFAULT 'Plan to Watch',
    last_episode_watched INT DEFAULT 0,
    total_episodes INT DEFAULT 0,
    ranking INT DEFAULT 0
);

-- Insertando datos de prueba en la tabla 'series'
INSERT INTO series (title, status, last_episode_watched, total_episodes, ranking) VALUES
    ('fsd', 'Plan to Watch', 0, 0, 0),
    ('Attack on Titan', 'Completed', 75, 75, 1),
    ('One Piece', 'Watching', 1050, 1100, 2),
    ('Death Note', 'Completed', 37, 37, 3),
    ('Naruto', 'Completed', 220, 220, 4),
    ('Bleach', 'Dropped', 150, 366, 5),
    ('Dragon Ball Z', 'Completed', 291, 291, 6),
    ('Demon Slayer', 'Watching', 30, 50, 7),
    ('Steins;Gate', 'Completed', 24, 24, 8),
    ('Fullmetal Alchemist: Brotherhood', 'Completed', 64, 64, 9);

