CREATE TABLE friend_requests (
    request_id INT PRIMARY KEY AUTO_INCREMENT,
    sender_id INT NOT NULL REFERENCES Users(user_id) ON DELETE CASCADE,
    receiver_id INT NOT NULL REFERENCES Users(user_id) ON DELETE CASCADE,
    status ENUM('PENDING', 'ACCEPTED', 'DECLINED') DEFAULT 'PENDING',
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    responded_at TIMESTAMP NULL,
    UNIQUE (sender_id, receiver_id),
    CHECK (sender_id <> receiver_id)
);