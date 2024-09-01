CREATE TABLE group_members {
    group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    member_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    added_by UUID NOT NULL REFERENCES users(id),
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, member_id)
}