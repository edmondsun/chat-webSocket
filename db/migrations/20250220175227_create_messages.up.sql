CREATE TABLE IF NOT EXISTS messages (
    id BIGINT AUTO_INCREMENT NOT NULL COMMENT 'Message ID, primary key',
    sender_id VARCHAR(255) NOT NULL COMMENT 'ID of the sender',
    room_id VARCHAR(255) NOT NULL COMMENT 'ID of the room where the message was sent',
    content TEXT NOT NULL COMMENT 'Content of the message',
    action VARCHAR(50) NOT NULL DEFAULT 'message' COMMENT 'Type of action (e.g., message, join, leave)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Timestamp when the message was created',
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;