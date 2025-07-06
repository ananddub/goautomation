

-- name: CheckRecentChanges5min :many
SELECT
    name,
    title,
    Count(pixel_color) as total_pixel,
    COUNT(DISTINCT pixel_color) as total_unique_pixel
FROM automation_logs
WHERE create_at >= datetime('now', '-5 minutes')
GROUP BY title,name;

-- name: DeleteAfter30min :exec
DELETE FROM automation_logs
WHERE create_at < datetime('now', '-30 minutes');