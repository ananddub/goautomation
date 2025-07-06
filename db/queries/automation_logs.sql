-- ==========================================
-- CREATE OPERATIONS
-- ==========================================


-- name: CreateAutomationLog :one
INSERT INTO automation_logs (
    name, title, action_type, x_position, y_position, pixel_color, success, error_message
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: CreateAutomationLogWithTimestamp :one
INSERT INTO automation_logs (
    name, title, action_type, x_position, y_position, pixel_color, create_at, success, error_message
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- ==========================================
-- READ OPERATIONS - Single Record
-- ==========================================

-- name: GetAutomationLog :one
SELECT * FROM automation_logs
WHERE id = ? LIMIT 1;

-- name: GetAutomationLogByName :one
SELECT * FROM automation_logs
WHERE name = ? 
ORDER BY create_at DESC 
LIMIT 1;

-- name: GetLatestAutomationLog :one
SELECT * FROM automation_logs
ORDER BY create_at DESC
LIMIT 1;

-- name: GetLatestSuccessfulLog :one
SELECT * FROM automation_logs
WHERE success = TRUE
ORDER BY create_at DESC
LIMIT 1;

-- name: GetLatestFailedLog :one
SELECT * FROM automation_logs
WHERE success = FALSE
ORDER BY create_at DESC
LIMIT 1;

-- ==========================================
-- READ OPERATIONS - Multiple Records
-- ==========================================

-- name: ListAutomationLogs :many
SELECT * FROM automation_logs
ORDER BY create_at DESC
LIMIT ?;

-- name: ListAutomationLogsByAction :many
SELECT * FROM automation_logs
WHERE action_type = ?
ORDER BY create_at DESC
LIMIT ?;

-- name: ListAutomationLogsByName :many
SELECT * FROM automation_logs
WHERE name = ?
ORDER BY create_at DESC
LIMIT ?;

-- name: ListAutomationLogsByTitle :many
SELECT * FROM automation_logs
WHERE title LIKE '%' || ? || '%'
ORDER BY create_at DESC
LIMIT ?;

-- name: ListSuccessfulLogs :many
SELECT * FROM automation_logs
WHERE success = TRUE
ORDER BY create_at DESC
LIMIT ?;

-- name: ListFailedLogs :many
SELECT * FROM automation_logs
WHERE success = FALSE
ORDER BY create_at DESC
LIMIT ?;

-- name: ListLogsByDateRange :many
SELECT * FROM automation_logs
WHERE create_at BETWEEN ? AND ?
ORDER BY create_at DESC;

-- name: ListLogsByPosition :many
SELECT * FROM automation_logs
WHERE x_position = ? AND y_position = ?
ORDER BY create_at DESC
LIMIT ?;

-- name: ListLogsByPixelColor :many
SELECT * FROM automation_logs
WHERE pixel_color = ?
ORDER BY create_at DESC
LIMIT ?;

-- ==========================================
-- READ OPERATIONS - Time-based Queries
-- ==========================================

-- name: GetRecentAutomationLogs :many
SELECT * FROM automation_logs
WHERE create_at >= datetime('now', '-1 hour')
ORDER BY create_at DESC;

-- name: GetLogsFromLastMinutes :many
SELECT * FROM automation_logs
WHERE create_at >= datetime('now', '-' || ? || ' minutes')
ORDER BY create_at DESC;

-- name: GetLogsFromLastHours :many
SELECT * FROM automation_logs
WHERE create_at >= datetime('now', '-' || ? || ' hours')
ORDER BY create_at DESC;

-- name: GetLogsFromLastDays :many
SELECT * FROM automation_logs
WHERE create_at >= datetime('now', '-' || ? || ' days')
ORDER BY create_at DESC;

-- name: GetTodaysLogs :many
SELECT * FROM automation_logs
WHERE date(create_at) = date('now')
ORDER BY create_at DESC;

-- name: GetYesterdaysLogs :many
SELECT * FROM automation_logs
WHERE date(create_at) = date('now', '-1 day')
ORDER BY create_at DESC;

-- ==========================================
-- COUNT OPERATIONS
-- ==========================================

-- name: CountAutomationLogs :one
SELECT COUNT(*) FROM automation_logs;

-- name: CountAutomationLogsByAction :one
SELECT COUNT(*) FROM automation_logs
WHERE action_type = ?;

-- name: CountAutomationLogsByName :one
SELECT COUNT(*) FROM automation_logs
WHERE name = ?;

-- name: CountSuccessfulLogs :one
SELECT COUNT(*) FROM automation_logs
WHERE success = TRUE;

-- name: CountFailedLogs :one
SELECT COUNT(*) FROM automation_logs
WHERE success = FALSE;

-- name: CountLogsByDateRange :one
SELECT COUNT(*) FROM automation_logs
WHERE create_at BETWEEN ? AND ?;

-- name: CountTodaysLogs :one
SELECT COUNT(*) FROM automation_logs
WHERE date(create_at) = date('now');

-- name: CountRecentLogs :one
SELECT COUNT(*) FROM automation_logs
WHERE create_at >= datetime('now', '-1 hour');

-- ==========================================
-- AGGREGATE OPERATIONS
-- ==========================================

-- name: GetActionTypeSummary :many
SELECT 
    action_type,
    COUNT(*) as total_count,
    SUM(CASE WHEN success = TRUE THEN 1 ELSE 0 END) as success_count,
    SUM(CASE WHEN success = FALSE THEN 1 ELSE 0 END) as failed_count,
    MIN(create_at) as first_occurrence,
    MAX(create_at) as last_occurrence
FROM automation_logs
GROUP BY action_type
ORDER BY total_count DESC;

-- name: GetDailySummary :many
SELECT 
    date(create_at) as log_date,
    COUNT(*) as total_logs,
    SUM(CASE WHEN success = TRUE THEN 1 ELSE 0 END) as successful_logs,
    SUM(CASE WHEN success = FALSE THEN 1 ELSE 0 END) as failed_logs
FROM automation_logs
WHERE create_at >= date('now', '-30 days')
GROUP BY date(create_at)
ORDER BY log_date DESC;

-- name: GetHourlySummary :many
SELECT 
    strftime('%Y-%m-%d %H:00:00', create_at) as hour_slot,
    COUNT(*) as total_logs,
    SUM(CASE WHEN success = TRUE THEN 1 ELSE 0 END) as successful_logs,
    SUM(CASE WHEN success = FALSE THEN 1 ELSE 0 END) as failed_logs
FROM automation_logs
WHERE create_at >= datetime('now', '-24 hours')
GROUP BY strftime('%Y-%m-%d %H', create_at)
ORDER BY hour_slot DESC;

-- name: GetMostUsedPositions :many
SELECT 
    x_position,
    y_position,
    COUNT(*) as usage_count,
    MAX(create_at) as last_used
FROM automation_logs
WHERE x_position IS NOT NULL AND y_position IS NOT NULL
GROUP BY x_position, y_position
ORDER BY usage_count DESC
LIMIT ?;

-- name: GetMostCommonPixelColors :many
SELECT 
    pixel_color,
    COUNT(*) as occurrence_count,
    MAX(create_at) as last_seen
FROM automation_logs
WHERE pixel_color IS NOT NULL
GROUP BY pixel_color
ORDER BY occurrence_count DESC
LIMIT ?;

-- ==========================================
-- UPDATE OPERATIONS
-- ==========================================

-- name: UpdateAutomationLogSuccess :exec
UPDATE automation_logs
SET success = ?, error_message = ?
WHERE id = ?;

-- name: UpdateAutomationLogName :exec
UPDATE automation_logs
SET name = ?
WHERE id = ?;

-- name: UpdateAutomationLogTitle :exec
UPDATE automation_logs
SET title = ?
WHERE id = ?;

-- name: UpdateAutomationLogPosition :exec
UPDATE automation_logs
SET x_position = ?, y_position = ?
WHERE id = ?;

-- name: UpdateAutomationLogPixelColor :exec
UPDATE automation_logs
SET pixel_color = ?
WHERE id = ?;

-- ==========================================
-- DELETE OPERATIONS
-- ==========================================

-- name: DeleteAutomationLog :exec
DELETE FROM automation_logs
WHERE id = ?;

-- name: DeleteAutomationLogsByAction :exec
DELETE FROM automation_logs
WHERE action_type = ?;

-- name: DeleteAutomationLogsByName :exec
DELETE FROM automation_logs
WHERE name = ?;

-- name: DeleteFailedLogs :exec
DELETE FROM automation_logs
WHERE success = FALSE;

-- name: DeleteOldAutomationLogs :exec
DELETE FROM automation_logs
WHERE create_at < datetime('now', '-7 days');

-- name: DeleteOldAutomationLogsByDays :exec
DELETE FROM automation_logs
WHERE create_at < datetime('now', '-' || ? || ' days');

-- name: DeleteLogsByDateRange :exec
DELETE FROM automation_logs
WHERE create_at BETWEEN ? AND ?;

-- name: ClearAllAutomationLogs :exec
DELETE FROM automation_logs;

-- ==========================================
-- SEARCH OPERATIONS
-- ==========================================

-- name: SearchLogsByErrorMessage :many
SELECT * FROM automation_logs
WHERE error_message LIKE '%' || ? || '%'
ORDER BY create_at DESC
LIMIT ?;

-- name: SearchLogsByNameOrTitle :many
SELECT * FROM automation_logs
WHERE name LIKE '%' || ? || '%' OR title LIKE '%' || ? || '%'
ORDER BY create_at DESC
LIMIT ?;

-- ==========================================
-- UTILITY OPERATIONS
-- ==========================================

-- name: GetDatabaseStats :one
SELECT 
    COUNT(*) as total_logs,
    COUNT(DISTINCT action_type) as unique_actions,
    COUNT(DISTINCT name) as unique_names,
    MIN(create_at) as oldest_log,
    MAX(create_at) as newest_log,
    SUM(CASE WHEN success = TRUE THEN 1 ELSE 0 END) as total_successful,
    SUM(CASE WHEN success = FALSE THEN 1 ELSE 0 END) as total_failed
FROM automation_logs;

-- name: GetLogSizeByAction :many
SELECT 
    action_type,
    COUNT(*) as log_count,
    ROUND(AVG(LENGTH(COALESCE(error_message, ''))), 2) as avg_error_msg_length
FROM automation_logs
GROUP BY action_type
ORDER BY log_count DESC;
