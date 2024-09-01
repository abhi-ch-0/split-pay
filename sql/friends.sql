CREATE TABLE friends (
    friendship_id INT PRIMARY KEY AUTO_INCREMENT,
    user1_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    user2_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user1_id, user2_id),
    CHECK (user1_id < user2_id)
);