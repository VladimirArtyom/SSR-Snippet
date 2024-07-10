package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
  ID      int
  Title   string
  Content string
  Created time.Time
  Expires time.Time
}

type SnippetModel struct {
  DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
  var query string = "INSERT INTO snippets (title, content, created,  expires) VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ?  DAY))"
  result, err := m.DB.Exec(query, title, content, expires)
  if err != nil {
    return 0, err
  }

  newUserId, err := result.LastInsertId()
  if err != nil {
    return 0, err
  }

  return int(newUserId), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
  var query string = "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?"

  row :=  m.DB.QueryRow(query, id)
  var s *Snippet = &Snippet{}
  err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
  if err !=  nil {
    if errors.Is(err, sql.ErrNoRows) {
      return nil, ErrNoRecord	
    } else {
      return nil, err
    }
  }

  return s,  nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
  var query string = "SELECT  id, title, content,  created, expires FROM  snippets WHERE expires >  UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10"

  rows,  err := m.DB.Query(query)
  if  err != nil {
    return nil, err
  }
  
  defer rows.Close()

  snippets := []*Snippet{}
  for rows.Next() {
    snippet := &Snippet{}
    err := rows.Scan(&snippet.ID, &snippet.Title,
	&snippet.Content, &snippet.Created,
	&snippet.Expires)

    if err != nil {
      return nil, err
    }

    snippets = append(snippets, snippet)
  }
  
  if err = rows.Err(); err != nil {
    return nil, err
  }

  return snippets, nil
}
