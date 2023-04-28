package model

import (
	"LibManSys/utils"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type ConnectConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

type DatabaseConnector struct {
	Config ConnectConfig
	DB     *sql.DB
}

type SQLExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func (c *DatabaseConnector) Connect() error {
	var err error
	c.DB, err = sql.Open("mysql", c.Config.databaseLoginInfo())
	if err != nil {
		logrus.Panic(err)
	}
	c.DB.SetMaxOpenConns(2000)
	c.DB.SetMaxIdleConns(1000)
	c.DB.SetConnMaxLifetime(time.Minute * 60)
	return nil
}

func (c *DatabaseConnector) Close() error {
	return c.DB.Close()
}

func (c *DatabaseConnector) ResetDatabase() error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	dropSQL := "DROP TABLE IF EXISTS %s"

	dbNames := []string{"borrow", "book", "card"}
	for _, dbName := range dbNames {
		_, err := tx.Exec(fmt.Sprintf(dropSQL, dbName))
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = AutoMigrateInTx(
		tx,
		Book{},
		Card{},
		Borrow{},
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (config *ConnectConfig) databaseLoginInfo() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)
}

func autoMigrate(executor SQLExecutor, args ...any) error {
	for _, table := range args {
		createTableSQL := "CREATE TABLE IF NOT EXISTS "
		uniqueFields := make(map[string][]string)
		primaryFields := make([]string, 0)
		checkConstraints := make([]string, 0)
		foreignKeys := make([]string, 0)
		createTableSQL += utils.ToSnake(reflect.ValueOf(table).Type().Name()) + " ("
		for i := 0; i < reflect.ValueOf(table).NumField(); i++ {
			field := reflect.ValueOf(table).Type().Field(i)
			createTableSQL += utils.ToSnake(field.Name) + " " + getType(field)
			tag := field.Tag.Get("sql")
			tags := strings.Split(tag, ";")
			for _, t := range tags {
				if t == "not null" {
					createTableSQL += " NOT NULL"
				}
				if t == "autoIncrement" {
					createTableSQL += " AUTO_INCREMENT"
				}
				if t == "unique" {
					createTableSQL += " UNIQUE"
				}
				if strings.HasPrefix(t, "default:") {
					createTableSQL += " DEFAULT " + strings.TrimPrefix(t, "default:")
				}
				if strings.HasPrefix(t, "check:") {
					checkConstraint := "CHECK(" + strings.TrimPrefix(t, "check:") + ")"
					checkConstraints = append(checkConstraints, checkConstraint)
				}
				if strings.HasPrefix(t, "constraint:") {
					foreignKey := ""
					constraints := strings.Split(strings.TrimPrefix(t, "constraint:"), ",")
					if len(constraints) == 0 {
						logrus.WithField("field", reflect.ValueOf(table).Type().Name()+"."+field.Name).
							Fatal("empty constraint")
						return errors.New("empty constraint")
					}
					references := strings.Split(constraints[0], ".")
					if len(references) != 2 {
						logrus.WithField("reference", references).Fatal("invalid reference")
						return errors.New("invalid reference")
					}
					foreignKey = "FOREIGN KEY (" + utils.ToSnake(field.Name) + ") REFERENCES " + utils.ToSnake(references[0]) + "(" + utils.ToSnake(references[1]) + ")"
					for _, val := range constraints {
						if strings.HasPrefix(val, "OnUpdate:") {
							foreignKey += " ON UPDATE " + strings.TrimPrefix(val, "OnUpdate:")
						}
						if strings.HasPrefix(val, "OnDelete:") {
							foreignKey += " ON DELETE " + strings.TrimPrefix(val, "OnDelete:")
						}
					}
					foreignKeys = append(foreignKeys, foreignKey)
				}
				if strings.HasPrefix(t, "unique:") {
					name := strings.TrimPrefix(t, "unique:")
					if _, ok := uniqueFields[name]; !ok {
						uniqueFields[name] = make([]string, 0)
					}
					uniqueFields[name] = append(uniqueFields[name], utils.ToSnake(field.Name))
				}
				if t == "primaryKey" {
					primaryFields = append(primaryFields, utils.ToSnake(field.Name))
				}
			}
			if i != reflect.ValueOf(table).NumField()-1 {
				createTableSQL += ", "
			}
		}
		for name, fields := range uniqueFields {
			createTableSQL += ", CONSTRAINT " + name + " UNIQUE("
			for i, field := range fields {
				createTableSQL += field
				if i != len(fields)-1 {
					createTableSQL += ", "
				}
			}
			createTableSQL += ")"
		}
		if len(primaryFields) != 0 {
			createTableSQL += ", PRIMARY KEY("
			for i, field := range primaryFields {
				createTableSQL += field
				if i != len(primaryFields)-1 {
					createTableSQL += ", "
				}
			}
			createTableSQL += ")"
		}
		if len(checkConstraints) != 0 {
			for _, constraint := range checkConstraints {
				createTableSQL += ", " + constraint
			}
		}
		if len(foreignKeys) != 0 {
			for _, foreignKey := range foreignKeys {
				createTableSQL += ", " + foreignKey
			}
		}
		createTableSQL += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"
		_, err := executor.Exec(createTableSQL)
		if err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}

func (c *DatabaseConnector) AutoMigrate(args ...any) error {
	return autoMigrate(c.DB, args...)
}

func AutoMigrateInTx(tx *sql.Tx, args ...any) error {
	return autoMigrate(tx, args...)
}

func getType(field reflect.StructField) string {
	switch field.Type.Name() {
	case "string":
		tag := field.Tag.Get("sql")
		tags := strings.Split(tag, ";")
		for _, t := range tags {
			if strings.HasPrefix(t, "size:") {
				size, err := strconv.Atoi(strings.TrimPrefix(t, "size:"))
				if err != nil {
					logrus.Fatal(err)
					return ""
				}
				if size <= 0 {
					logrus.Fatal("size must be greater than 0")
					return ""
				} else if size > 0 && size <= 255 {
					return "VARCHAR(" + strconv.Itoa(size) + ")"
				} else if size > 255 && size <= 65535 {
					return "TEXT"
				} else if size > 65535 && size <= 16777215 {
					return "MEDIUMTEXT"
				} else if size > 16777215 && size <= 4294967295 {
					return "LONGTEXT"
				} else {
					logrus.Fatal("size must be less than 4294967295")
					return ""
				}
			}
			if strings.HasPrefix(t, "char:") {
				size, err := strconv.Atoi(strings.TrimPrefix(t, "char:"))
				if err != nil {
					logrus.Fatal(err)
					return ""
				}
				if size <= 0 {
					logrus.Fatal("size must be greater than 0")
					return ""
				} else if size > 0 && size <= 255 {
					return "CHAR(" + strconv.Itoa(size) + ")"
				} else {
					logrus.Fatal("size must be less than 255")
					return ""
				}
			}
		}
	case "int8":
		return "TINYINT"
	case "uint8":
		return "CHAR(1)"
	case "int16":
		return "SMALLINT"
	case "uint16":
		return "SMALLINT UNSIGNED"
	case "int":
		return "INT"
	case "int32":
		return "INT"
	case "uint":
		return "INT UNSIGNED"
	case "uint32":
		return "INT UNSIGNED"
	case "int64":
		return "BIGINT"
	case "uint64":
		return "BIGINT UNSIGNED"
	case "float32", "myFloat":
		tag := field.Tag.Get("sql")
		tags := strings.Split(tag, ";")
		for _, t := range tags {
			if strings.HasPrefix(t, "decimal:") {
				nums := strings.Split(strings.TrimPrefix(t, "decimal:"), ",")
				if len(nums) != 2 {
					logrus.WithField("field", field.Name).Fatal("invalid decimal")
					return ""
				}
				precision, err := strconv.Atoi(nums[0])
				if err != nil {
					logrus.WithField("field", field.Name).Fatal(err)
					return ""
				}
				scale, err := strconv.Atoi(nums[1])
				if err != nil {
					logrus.WithField("field", field.Name).Fatal(err)
					return ""
				}
				if precision <= 0 {
					logrus.WithField("field", field.Name).Fatal("precision must be greater than 0")
					return ""
				}
				if scale <= 0 {
					logrus.WithField("field", field.Name).Fatal("scale must be greater than 0")
					return ""
				}
				return "DECIMAL(" + strconv.Itoa(precision) + "," + strconv.Itoa(scale) + ")"
			}
		}
		return "FLOAT"
	case "float64":
		tag := field.Tag.Get("sql")
		tags := strings.Split(tag, ";")
		for _, t := range tags {
			if strings.HasPrefix(t, "decimal:") {
				nums := strings.Split(strings.TrimPrefix(t, "decimal:"), ",")
				if len(nums) != 2 {
					logrus.WithField("field", field.Name).Fatal("invalid decimal")
					return ""
				}
				precision, err := strconv.Atoi(nums[0])
				if err != nil {
					logrus.WithField("field", field.Name).Fatal(err)
					return ""
				}
				scale, err := strconv.Atoi(nums[1])
				if err != nil {
					logrus.WithField("field", field.Name).Fatal(err)
					return ""
				}
				if precision <= 0 {
					logrus.WithField("field", field.Name).Fatal("precision must be greater than 0")
					return ""
				}
				if scale <= 0 {
					logrus.WithField("field", field.Name).Fatal("scale must be greater than 0")
					return ""
				}
				return "DECIMAL(" + strconv.Itoa(precision) + "," + strconv.Itoa(scale) + ")"
			}
		}
		return "DOUBLE"
	case "time.Time":
		return "DATETIME"
	default:
		return "VARCHAR(255)"
	}
	return "VARCHAR(255)"
}
