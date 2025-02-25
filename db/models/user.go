package models

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64
	Username     string
	Name         string
	Email        string
	PasswordHash string
	Admin        bool
	CreatedAt    int
	UpdatedAt    int
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func UserCheckCredentials(db *sql.DB, username, password string) (int64, error) {
	var user User
	row := db.QueryRow("SELECT id, username, password_hash FROM Users WHERE username = ?", username)
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash); err != nil {
		return -1, fmt.Errorf("User %s not found", username)
	}

	hashedPassword := user.PasswordHash
	if !checkPassword(hashedPassword, password) {
		return -1, fmt.Errorf("Wrong password")
	}
	return user.ID, nil
}

func UserGetByUsername(db *sql.DB, username string) (*User, error) {
	var user User
	var admin int
	row := db.QueryRow("SELECT id, username, name, email, password_hash, admin, created_at, updated_at FROM Users WHERE username = ?", username)
	if err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.PasswordHash, &admin, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, fmt.Errorf("User %s not found", username)
	}
	user.Admin = admin == 1
	return &user, nil
}

func UserGetByID(db *sql.DB, userID int64) (*User, error) {
	var user User
	var admin int
	row := db.QueryRow("SELECT id, username, name, email, password_hash, admin, created_at, updated_at FROM users WHERE id = ?", userID)
	if err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.PasswordHash, &admin, &user.CreatedAt, &user.UpdatedAt); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("User#%d not found", userID)
	}
	user.Admin = admin == 1
	return &user, nil
}

func UserCreate(db *sql.DB, username, name, email, plain_password string) (int64, error) {
	password, err := hashPassword(plain_password)
	if err != nil {
		return -1, err
	}

	row, err := db.Exec(
		"INSERT INTO users(username, name, email, password_hash) VALUES(?, ?, ?, ?)",
		username, name, email, password,
	)

	if err != nil {
		return -1, err
	}

	return row.LastInsertId()
}
