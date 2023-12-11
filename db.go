package torm

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Model struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" form:"-" json:"-"`
}

type DB struct {
	*gorm.DB

	pkg string
	err error
}

func (db *DB) Close() error {
	if db.DB == nil {
		return nil
	}
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		return err
	}
	db.DB = nil
	return nil
}

func (db *DB) Error() error {
	return db.err
}

func (db *DB) NoRecord() bool {
	return errors.Is(db.err, gorm.ErrRecordNotFound)
}

func (db *DB) CFirst(dest any, conds ...any) *DB {
	db.err = db.DB.First(dest, conds...).Error
	return db
}

func (db *DB) First(dest any, conds ...any) bool {
	return !db.CFirst(dest, conds...).NoRecord()
}

func (db *DB) CSelect(dest any, fields []string, conds ...any) *DB {
	db.err = db.DB.Select(fields).First(dest, conds...).Error
	return db
}

func (db *DB) Select(dest any, fields []string, conds ...any) bool {
	return !db.CSelect(dest, fields, conds...).NoRecord()
}

func (db *DB) CFind(dest any, conds ...any) *DB {
	db.err = db.DB.Find(dest, conds...).Error
	return db
}

func (db *DB) Find(dest any, conds ...any) bool {
	return !db.CFind(dest, conds...).NoRecord()
}

func (db *DB) FirstOrCreate(first, create func(), x any, conds ...any) {
	if db.First(x, conds...) {
		if first != nil {
			first()
		}
	} else {
		db.Create(x)
		if create != nil {
			create()
		}
	}
}

func (db *DB) PreloadDB(in any) *gorm.DB {
	r := db.Model(in)
	for _, s := range Ref.Get(in) {
		r.Preload(s)
	}
	return r.Preload(clause.Associations)
}

func (db *DB) CPreload(t any, conds ...any) *DB {
	db.err = db.PreloadDB(t).First(t, conds...).Error
	return db
}

func (db *DB) Preload(t any, conds ...any) bool {
	return !db.CPreload(t, conds...).NoRecord()
}

func (db *DB) CPreloads(t any, conds ...any) *DB {
	db.err = db.PreloadDB(t).Find(t, conds...).Error
	return db
}

func (db *DB) Preloads(t any, conds ...any) bool {
	return !db.CPreloads(t, conds...).NoRecord()
}

func (db *DB) AutoMigrate(dst ...any) *DB {
	Ref.Init(dst...)
	db.err = db.DB.AutoMigrate(dst...)
	return db
}
