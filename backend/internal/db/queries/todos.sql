-- name: CreateTodo :one
INSERT INTO todos (user_id, description, position)
VALUES ($1, $2, 
    (SELECT COALESCE(MAX(position) + 1, 0) FROM todos WHERE user_id = $1)
)
RETURNING *;

-- name: ListTodos :many
SELECT * FROM todos WHERE user_id = $1 ORDER BY completed, position;

-- name: SearchTodos :many
SELECT *
FROM todos
WHERE user_id = $1
  AND description @@ to_tsquery('english', $2)
ORDER BY completed, position;

-- name: UpdateTodo :one
UPDATE todos
SET description = $2, 
    completed = $3, 
    position = $4,
    updated_at = NOW()
WHERE id = $1 AND user_id = $5
RETURNING *;

-- name: UpdateTodoPosition :one
UPDATE todos
SET position = $2,
    updated_at = NOW()
WHERE id = $1 AND user_id = $3
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1 AND user_id = $2;