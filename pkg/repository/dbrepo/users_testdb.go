package dbrepo

import (
	"database/sql"
	"personal-projects/webapp/pkg/data"
)

type TestDBRepo struct{}

func (m *TestDBRepo) Connection() *sql.DB {
	return nil
}

// AllUsers returns all users as a slice of *data.User
func (m *TestDBRepo) AllUsers() ([]*data.User, error) {
	var users []*data.User

	return users, nil
}

// GetUser returns one user by id
func (m *TestDBRepo) GetUser(id int) (*data.User, error) {
	var user = data.User{
		ID: 1,
	}
	return &user, nil
}

// GetUserByEmail returns one user by email address
func (m *TestDBRepo) GetUserByEmail(email string) (*data.User, error) {
	var user = data.User{
		ID: 1,
	}
	return &user, nil
}

// UpdateUser updates one user in the database
func (m *TestDBRepo) UpdateUser(u data.User) error {
	return nil
}

// DeleteUser deletes one user from the database, by id
func (m *TestDBRepo) DeleteUser(id int) error {
	return nil
}

// InsertUser inserts a new user into the database, and returns the ID of the newly inserted row
func (m *TestDBRepo) InsertUser(user data.User) (int, error) {

	return 2, nil
}

// ResetPassword is the method we will use to change a user's password.
func (m *TestDBRepo) ResetPassword(id int, password string) error {

	return nil
}

// InsertUserImage inserts a user profile image into the database.
func (m *TestDBRepo) InsertUserImage(i data.UserImage) (int, error) {
	return 1, nil
}
