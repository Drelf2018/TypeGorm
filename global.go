package torm

import (
	"path/filepath"
	"runtime"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	here     string
	GlobalDB = make(map[string]*DB)
)

func init() {
	_, here, _, _ = runtime.Caller(0)
}

func caller(skip int) string {
	file := here
	for ; file == here; skip++ {
		_, file, _, _ = runtime.Caller(skip)
	}
	return filepath.Dir(file)
}

func Where(relative string) *DB {
	return GlobalDB[filepath.Clean(filepath.Join(caller(2), relative))]
}

func Get() *DB {
	return GlobalDB[caller(3)]
}

func SetDB(gormDB *gorm.DB) (db *DB) {
	db = &DB{DB: gormDB, pkg: caller(2)}
	GlobalDB[db.pkg] = db
	return
}

func SetDialector(dialector gorm.Dialector) *DB {
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return SetDB(gormDB)
}

func SetSqlite(file string) *DB {
	return SetDialector(sqlite.Open(file))
}

func Close() error {
	return Get().Close()
}

func CFirst[T any](conds ...any) (*T, *DB) {
	dest := new(T)
	return dest, Get().CFirst(dest, conds...)
}

func First[T any](conds ...any) (*T, bool) {
	dest := new(T)
	return dest, Get().First(dest, conds...)
}

func MFirst[T any](conds ...any) *T {
	dest := new(T)
	Get().CFirst(dest, conds...)
	return dest
}

func CSelect[T any](fields []string, conds ...any) (*T, *DB) {
	dest := new(T)
	return dest, Get().CSelect(dest, fields, conds...)
}

func Select[T any](fields []string, conds ...any) (*T, bool) {
	dest := new(T)
	return dest, Get().Select(dest, fields, conds...)
}

func MSelect[T any](fields []string, conds ...any) *T {
	dest := new(T)
	Get().CSelect(dest, fields, conds...)
	return dest
}

func CFind[T any](conds ...any) ([]T, *DB) {
	dest := make([]T, 0)
	return dest, Get().CFind(&dest, conds...)
}

func Find[T any](conds ...any) ([]T, bool) {
	dest := make([]T, 0)
	return dest, Get().Find(&dest, conds...)
}

func MFind[T any](conds ...any) []T {
	dest := make([]T, 0)
	Get().CFind(&dest, conds...)
	return dest
}

func Exist[T any](conds ...any) bool {
	return Get().First(new(T), conds...)
}

func FirstOrCreate[T any](first, create func(), x *T, conds ...any) {
	Get().FirstOrCreate(first, create, x, conds...)
}

func CPreload[T any](conds ...any) (*T, *DB) {
	t := new(T)
	return t, Get().CPreload(t, conds...)
}

func Preload[T any](conds ...any) (*T, bool) {
	t := new(T)
	return t, Get().Preload(t, conds...)
}

func MPreload[T any](conds ...any) *T {
	t := new(T)
	Get().CPreload(t, conds...)
	return t
}

func CPreloads[T any](conds ...any) ([]T, *DB) {
	t := make([]T, 0)
	return t, Get().CPreloads(&t, conds...)
}

func Preloads[T any](conds ...any) ([]T, bool) {
	t := make([]T, 0)
	return t, Get().Preloads(&t, conds...)
}

func MPreloads[T any](conds ...any) []T {
	t := make([]T, 0)
	Get().CPreloads(&t, conds...)
	return t
}

func AutoMigrate(dst ...any) *DB {
	return Get().AutoMigrate(dst...)
}
