package dba

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/sketch-home-task/src/pkg/illustrator"
)

type Storage struct {
	*sql.DB
}

func NewStorage(dialect, dsn string, idleConn, maxConn int) (s illustrator.CanvasStorage, err error) {
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return
	}

	err = db.Ping()
	if err != nil {
		return
	}

	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)
	s = &Storage{db}
	return
}

func (s *Storage) Close() (err error) {
	return s.DB.Close()
}

func (s *Storage) FindByName(ctx context.Context, name string) (canvas *illustrator.CanvasModel, err error) {
	canvas = &illustrator.CanvasModel{}
	err = s.QueryRow("SELECT width, height, drawings FROM canvas WHERE name = $1", name).
		Scan(&canvas.Width, &canvas.Height, &canvas.Drawings)
	return
}

func (s *Storage) Create(ctx context.Context, canvas *illustrator.CanvasModel) (res sql.Result, err error) {
	stmt, err := s.PrepareContext(ctx, "INSERT INTO canvas (name, width, height, drawings) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return
	}
	defer stmt.Close()

	res, err = stmt.ExecContext(ctx, canvas.Name, canvas.Width, canvas.Height, canvas.Drawings)
	return
}

func (s *Storage) Update(ctx context.Context, canvas *illustrator.CanvasModel) (res sql.Result, err error) {
	stmt, err := s.PrepareContext(ctx, "UPDATE canvas SET width = $1, height = $2, drawings = $3 WHERE name = $4")
	if err != nil {
		return
	}
	defer stmt.Close()

	res, err = stmt.ExecContext(ctx, canvas.Width, canvas.Height, canvas.Drawings, canvas.Name)
	return
}

func (s *Storage) Delete(ctx context.Context, name string) (res sql.Result, err error) {
	stmt, err := s.PrepareContext(ctx, "DELETE FROM canvas WHERE name = $1")
	if err != nil {
		return
	}
	defer stmt.Close()

	res, err = stmt.ExecContext(ctx, name)
	return
}
