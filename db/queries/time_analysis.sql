

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
WHERE create_at < datetime('now', '-100 minutes');

-- name: GetRecentPixelStats :many
WITH recent_logs AS (
    SELECT
        name,
        title,
        pixel_color,
        create_at
    FROM automation_logs
    WHERE create_at >= datetime('now', '-5 minutes')
)
SELECT
    name,
    COUNT(pixel_color) as total_pixel,
    COUNT(DISTINCT pixel_color) as total_unique_pixel
FROM recent_logs
GROUP BY name
HAVING total_unique_pixel > 4;