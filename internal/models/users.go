package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

//For interacting with the app
type User struct {
  ID int
  Name string
  Email string
  hashed_password []byte
  Created time.Time
}

// For interacting with the DB
type UserModel struct {
  DB *sql.DB
}

func (u *UserModel) Insert(name string, email string,
  password string) error {
 
  hashed_pass, err := bcrypt.GenerateFromPassword([]byte(password),12)
  if err != nil {
    return err
  }
  
  var stmt string = "INSERT INTO users (name, email, hashed_password, created) VALUES (?, ?, ?, UTC_TIMESTAMP())" 
  _, err = u.DB.Exec(stmt, name, email, hashed_pass)
  if err != nil {
    var mysqlError *mysql.MySQLError


    // Always make the errors specific
    if errors.As(err, &mysqlError) {
      
      if mysqlError.Number == 1062 && strings.Contains(mysqlError.Message, "users.email") {
	return ErrDuplicateEmail
      }

    }
    return err
  }

  return nil
}


func (u *UserModel) Authenticate(email string, password string) (int, error){

  var id int 
  var hashed_password []byte

  var stmt string = "SELECT id, hashed_password FROM users WHERE email=?"

  err := u.DB.QueryRow(stmt, email).Scan(&id, &hashed_password)
  if err != nil {
    
    if errors.Is(err, sql.ErrNoRows) {
      return 0, ErrInvalidCredentials
    }

    return 0, err
  }

  err = bcrypt.CompareHashAndPassword(hashed_password, []byte(password))

  if err != nil {
    if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
      return 0, ErrInvalidCredentials
    }

    return 0, err
  }

  return id, nil
}

func (u *UserModel) Exists(id int) (bool, error){
  var exists bool
  var stmt string = "SELECT EXISTS(SELECT true FROM users WHERE id=?)"

  err := u.DB.QueryRow(stmt, id).Scan(&exists)

  return exists, err
}





