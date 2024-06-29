USE testdb;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL
);

INSERT INTO users (username, password) VALUES 
    ('alice', SUBSTRING(UUID(), 1, 8)),
    ('bob', SUBSTRING(UUID(), 1, 8)),
    ('charlie', SUBSTRING(UUID(), 1, 8)),
    ('dave', SUBSTRING(UUID(), 1, 8)),
    ('eve', SUBSTRING(UUID(), 1, 8)),
    ('frank', SUBSTRING(UUID(), 1, 8)),
    ('grace', SUBSTRING(UUID(), 1, 8)),
    ('hannah', SUBSTRING(UUID(), 1, 8)),
    ('isaac', SUBSTRING(UUID(), 1, 8)),
    ('jane', SUBSTRING(UUID(), 1, 8)),
    ('kyle', SUBSTRING(UUID(), 1, 8)),
    ('linda', SUBSTRING(UUID(), 1, 8)),
    ('mike', SUBSTRING(UUID(), 1, 8)),
    ('nancy', SUBSTRING(UUID(), 1, 8)),
    ('oliver', SUBSTRING(UUID(), 1, 8)),
    ('pamela', SUBSTRING(UUID(), 1, 8)),
    ('quentin', SUBSTRING(UUID(), 1, 8)),
    ('rachel', SUBSTRING(UUID(), 1, 8)),
    ('steve', SUBSTRING(UUID(), 1, 8)),
    ('tina', SUBSTRING(UUID(), 1, 8));
