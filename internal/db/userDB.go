package db

// UserExists checks if a user with given username and password exists in the database.
func UserExists(username, password string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT id FROM users WHERE username = ? AND password = ?);`

	err := Client.QueryRow(query, username, password).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
