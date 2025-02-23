CREATE TABLE IF NOT EXISTS clients (
    client_id VARCHAR(255) NOT NULL COMMENT 'Unique client identifier',
    email VARCHAR(255) NOT NULL COMMENT 'Client email address',
    name VARCHAR(255) COMMENT 'Client name',
    status INT NOT NULL DEFAULT 1 COMMENT 'Client status (1 for active)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Timestamp when the client record was created',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Timestamp when the client record was last updated',
    PRIMARY KEY (client_id),
    UNIQUE KEY idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;