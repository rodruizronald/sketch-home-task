package dba_test

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sketch-home-task/src/pkg/dba"
	"github.com/sketch-home-task/src/pkg/illustrator"
	"github.com/stretchr/testify/suite"
)

type MockStorageSuite struct {
	suite.Suite
	storage dba.Storage
	mock    sqlmock.Sqlmock
	model   illustrator.CanvasModel
}

func (s *MockStorageSuite) SetupSuite() {
	db, mock, err := sqlmock.New()
	s.NoError(err)

	s.mock = mock
	s.storage = dba.Storage{DB: db}

	asteriskRune := '*'
	s.model.Name = "monalisa"
	s.model.Width = 20
	s.model.Height = 20
	s.model.Drawings = []illustrator.DrawingModel{
		{
			Coordinates: []int{4, 4},
			Width:       5,
			Height:      5,
			Fill:        &asteriskRune,
			Outline:     &asteriskRune,
		},
	}
	hashedName := sha1.Sum([]byte(s.model.Name))
	s.model.Name = hex.EncodeToString(hashedName[:])
}

func (s *MockStorageSuite) TearDownSuite() {
	s.mock.ExpectClose()
	s.Nil(s.storage.Close())
	s.Nil(s.mock.ExpectationsWereMet())
}

// ----------------- FIND TESTS -----------------

type StorageFindTestSuite struct {
	MockStorageSuite
}

func TestStorageFindTestSuite(t *testing.T) {
	suite.Run(t, new(StorageFindTestSuite))
}

func (s *StorageFindTestSuite) TestStorageFind() {
	rows := sqlmock.NewRows([]string{"width", "height", "drawings"}).AddRow(s.model.Width, s.model.Height, s.model.Drawings)
	query := regexp.QuoteMeta(`SELECT width, height, drawings FROM canvas WHERE name = $1`)
	s.mock.ExpectQuery(query).WithArgs(s.model.Name).WillReturnRows(rows)

	canvas, err := s.storage.FindByName(context.Background(), s.model.Name)
	s.NoError(err)

	s.model.Name = ""
	s.True(reflect.DeepEqual(*canvas, s.model))
}

// ----------------- CREATE TESTS -----------------

type StorageCreateTestSuite struct {
	MockStorageSuite
}

func TestStorageCreateTestSuite(t *testing.T) {
	suite.Run(t, new(StorageCreateTestSuite))
}

func (s *StorageCreateTestSuite) TestStorageCreate() {
	query := regexp.QuoteMeta(`INSERT INTO canvas (name, width, height, drawings) VALUES ($1, $2, $3, $4)`)
	prep := s.mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(s.model.Name, s.model.Width, s.model.Height, s.model.Drawings).WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.storage.Create(context.Background(), &s.model)
	s.NoError(err)
}

// ----------------- UPDATE TESTS -----------------

type StorageUpdateTestSuite struct {
	MockStorageSuite
}

func TestStorageUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(StorageUpdateTestSuite))
}

func (s *StorageUpdateTestSuite) TestStorageUpdate() {
	query := regexp.QuoteMeta(`UPDATE canvas SET width = $1, height = $2, drawings = $3 WHERE name = $4`)
	prep := s.mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(s.model.Width, s.model.Height, s.model.Drawings, s.model.Name).WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.storage.Update(context.Background(), &s.model)
	s.NoError(err)
}

// ----------------- DELETE TESTS -----------------

type StorageDeleteTestSuite struct {
	MockStorageSuite
}

func TestStorageDeleteTestSuite(t *testing.T) {
	suite.Run(t, new(StorageDeleteTestSuite))
}

func (s *StorageDeleteTestSuite) TestStorageDelete() {
	query := regexp.QuoteMeta(`DELETE FROM canvas WHERE name = $1`)
	prep := s.mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(s.model.Name).WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.storage.Delete(context.Background(), s.model.Name)
	s.NoError(err)
}
