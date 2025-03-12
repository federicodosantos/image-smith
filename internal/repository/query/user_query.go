package query

const (
	InsertUserQuery = `INSERT INTO users(id, name, email, password, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6)`

	GetUserByIdQuery = `SELECT * FROM users WHERE id = $1`

	GetUserByEmailQuery = `SELECT * FROM users WHERE email = $1`

	CheckEmailExistQuery = `SELECT COUNT(*) FROM users WHERE email = $1`
)
