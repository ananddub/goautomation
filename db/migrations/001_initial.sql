

-- Initial schema for automation logs
CREATE TABLE automation_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    title TEXT NOT NULL,
    action_type TEXT NOT NULL,
    x_position REAL,
    y_position REAL,
    pixel_color TEXT,
    create_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    success BOOLEAN DEFAULT TRUE,
    error_message TEXT
);

-- Index for faster queries by created_at
CREATE INDEX idx_automation_logs_created_at ON automation_logs(create_at);

-- Index for faster queries by action type
CREATE INDEX idx_automation_logs_action_type ON automation_logs(action_type);
