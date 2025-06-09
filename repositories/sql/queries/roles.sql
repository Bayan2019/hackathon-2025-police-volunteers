-- name: GetRoles :many
SELECT * FROM roles;
--

-- name: GetRoleById :one
SELECT * FROM roles WHERE id = ?;
--

-- name: CreateRole :one
INSERT INTO roles(title)
VALUES (?)
RETURNING id;
--

-- name: UpdateRole :exec
UPDATE roles 
SET title = ?
WHERE id = ?;
--

-- name: GetRolesOfUser :many
SELECT r.*
FROM roles AS r
JOIN users_roles AS ur
ON r.id = ur.role_id
WHERE ur.user_id = ?;
--

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = ?;
--