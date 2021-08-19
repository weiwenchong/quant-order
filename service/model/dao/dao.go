package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/manager"
	"github.com/didi/gendry/scanner"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var DB *sql.DB

func init() {
	db, err := manager.New(
		"quantorder", "root", "Wwcwwc123", "172.17.0.1",
	).Set(manager.SetCharset("utf8"),
		manager.SetAllowCleartextPasswords(true),
		manager.SetInterpolateParams(true),
		manager.SetTimeout(1*time.Second),
		manager.SetReadTimeout(1*time.Second),
	).Port(3306).Open(true)
	if err != nil {
		log.Panicf("sql.Open err:%v", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	DB = db
}

func SelectOne(ctx context.Context, db *sql.DB, tableName string, where map[string]interface{}, target interface{}) (err error) {
	if db == nil {
		return errors.New("DB nil")
	}
	where["_limit"] = []uint{0, 1}
	cond, vals, err := builder.BuildSelect(tableName, where, nil)
	if err != nil {
		return nil
	}
	row, err := db.QueryContext(ctx, cond, vals...)
	if err != nil || row == nil {
		return err
	}
	defer row.Close()
	err = scanner.Scan(row, &target)
	return err
}

func SelectList(ctx context.Context, db *sql.DB, tableName string, where map[string]interface{}, target interface{}) (err error) {
	if db == nil {
		return errors.New("DB nil")
	}
	cond, vals, err := builder.BuildSelect(tableName, where, nil)
	if err != nil {
		return nil
	}
	row, err := db.QueryContext(ctx, cond, vals...)
	if err != nil || row == nil {
		return err
	}
	defer row.Close()
	err = scanner.Scan(row, &target)
	return err
}

func Insert(ctx context.Context, db *sql.DB, tableName string, data []map[string]interface{}) (id int64, err error) {
	if db == nil {
		return 0, errors.New("DB nil")
	}
	cond, vals, err := builder.BuildInsert(tableName, data)
	if err != nil {
		return 0, err
	}
	res, err := db.ExecContext(ctx, cond, vals...)
	if err != nil || res == nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Update(ctx context.Context, db *sql.DB, tableName string, where, data map[string]interface{}) (int64, error) {
	if db == nil {
		return 0, errors.New("DB nil")
	}
	cond, vals, err := builder.BuildUpdate(tableName, where, data)
	if nil != err {
		return 0, err
	}
	result, err := db.ExecContext(ctx, cond, vals...)
	if nil != err || result == nil {
		return 0, err
	}
	return result.RowsAffected()
}

func Delete(ctx context.Context, db *sql.DB, tableName string, where map[string]interface{}) (int64, error) {
	if db == nil {
		return 0, errors.New("DB nil")
	}
	cond, vals, err := builder.BuildDelete(tableName, where)
	if nil != err {
		return 0, err
	}
	result, err := db.ExecContext(ctx, cond, vals...)
	if nil != err || result == nil {
		return 0, err
	}
	return result.RowsAffected()
}
