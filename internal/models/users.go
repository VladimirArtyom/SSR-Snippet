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
  return 0, nil
}

func (u *UserModel) Exists(id int) (bool, error){
  return false , nil
}





