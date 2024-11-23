package database

import (
	"context"
	"database/sql"
	"fmt"
	"neuro-most/template-service/config"
	"neuro-most/template-service/internal/adapters/repo"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	db *gorm.DB
}

func NewGormDB(conf config.Config) *GormDB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Moscow", conf.DatabaseHost, conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseDB, conf.DatabasePort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &GormDB{db: db}
}

func (g *GormDB) AutoMigrate(models ...interface{}) {
	g.db.AutoMigrate(models...)
}

// Create создает новую запись
func (g *GormDB) Create(ctx context.Context, data interface{}) error {
	return g.db.WithContext(ctx).Create(data).Error
}

// Update обновляет запись
func (g *GormDB) Update(ctx context.Context, data interface{}) error {
	return g.db.WithContext(ctx).Updates(data).Error
}

// RawQuery выполняет сырой SQL запрос
func (g *GormDB) RawQuery(ctx context.Context, scanner interface{}, query string, args ...interface{}) error {
	return g.db.WithContext(ctx).Raw(query, args...).Scan(scanner).Error
}

// UpdateMany обновляет множество записей по условию
func (g *GormDB) UpdateMany(ctx context.Context, data, query interface{}, args ...interface{}) error {
	return g.db.WithContext(ctx).Where(query, args...).Updates(data).Error
}

// UpdateOne обновляет одну запись по условию
func (g *GormDB) UpdateOne(ctx context.Context, data, query interface{}, args ...interface{}) error {
	result := g.db.WithContext(ctx).Where(query, args...).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

// FindAll получает все записи с пагинацией
func (g *GormDB) FindAll(ctx context.Context, total *int64, result interface{}, input repo.FindAllInput, query interface{}, args ...interface{}) error {
	db := g.db.WithContext(ctx)

	// Подсчет общего количества записей
	if err := db.Where(query, args...).Count(total).Error; err != nil {
		return err
	}

	// Применяем пагинацию
	if input.PageInput.Limit > 0 {
		offset := (input.PageInput.Current - 1) * input.PageInput.Limit
		db = db.Offset(offset).Limit(input.PageInput.Limit)
	}

	// Применяем сортировку
	if input.OrderBy != "" {
		db = db.Order(input.OrderBy)
	}

	return db.Where(query, args...).Find(result).Error
}

// FindAllWithJoins получает все записи с джойнами
func (g *GormDB) FindAllWithJoins(ctx context.Context, total *int64, result interface{}, input repo.FindAllInput, query interface{}, args ...interface{}) error {
	db := g.db.WithContext(ctx)

	// Применяем джойны
	for _, join := range input.JoinInput {
		switch join.JoinType {
		case "LEFT":
			db = db.Joins(fmt.Sprintf("LEFT JOIN %s ON %s", join.Table, join.Condition))
		case "RIGHT":
			db = db.Joins(fmt.Sprintf("RIGHT JOIN %s ON %s", join.Table, join.Condition))
		case "INNER":
			db = db.Joins(fmt.Sprintf("INNER JOIN %s ON %s", join.Table, join.Condition))
		default:
			db = db.Joins(fmt.Sprintf("JOIN %s ON %s", join.Table, join.Condition))
		}
	}

	// Подсчет общего количества записей
	if err := db.Where(query, args...).Count(total).Error; err != nil {
		return err
	}

	// Применяем пагинацию
	if input.PageInput.Limit > 0 {
		offset := (input.PageInput.Current - 1) * input.PageInput.Limit
		db = db.Offset(offset).Limit(input.PageInput.Limit)
	}

	// Применяем сортировку
	if input.OrderBy != "" {
		db = db.Order(input.OrderBy)
	}

	return db.Where(query, args...).Find(result).Error
}

// FindOne находит одну запись
func (g *GormDB) FindOne(ctx context.Context, result interface{}, query interface{}, args ...interface{}) error {
	return g.db.WithContext(ctx).Where(query, args...).First(result).Error
}

// BeginFind начинает построение запроса
func (g *GormDB) BeginFind(ctx context.Context, value interface{}) repo.Find {
	db := g.db.Model(value)
	return &GormFind{
		db: db.WithContext(ctx),
	}
}

// Delete удаляет записи по условию
func (g *GormDB) Delete(ctx context.Context, data interface{}, condition interface{}, args ...interface{}) error {
	return g.db.WithContext(ctx).Where(condition, args...).Delete(data).Error
}

// DeleteByQuery удаляет записи по запросу
func (g *GormDB) DeleteByQuery(ctx context.Context, data, query interface{}, args ...interface{}) error {
	return g.db.WithContext(ctx).Where(query, args...).Delete(data).Error
}

// GetInstance возвращает инстанс GORM DB
func (g *GormDB) GetInstance() interface{} {
	return g.db
}

// GormFind реализация интерфейса Find
type GormFind struct {
	db          *gorm.DB
	selectQuery interface{}
}

func (f *GormFind) Where(query interface{}, args ...interface{}) repo.Find {
	f.db = f.db.Where(query, args...)
	return f
}

func (f *GormFind) Having(query interface{}, args ...interface{}) repo.Find {
	f.db = f.db.Having(query, args...)
	return f
}

func (f *GormFind) Page(current int, limit int) repo.Find {
	offset := (current - 1) * limit
	f.db = f.db.Offset(offset).Limit(limit)
	return f
}

func (f *GormFind) Join(query string, args ...interface{}) repo.Find {
	f.db = f.db.Joins(query, args...)
	return f
}

func (f *GormFind) Or(query interface{}, args ...interface{}) repo.Find {
	f.db = f.db.Or(query, args...)
	return f
}

func (f *GormFind) Not(query interface{}, args ...interface{}) repo.Find {
	f.db = f.db.Not(query, args...)
	return f
}

func (f *GormFind) Count(total *int64) error {
	return f.db.Count(total).Error
}

func (f *GormFind) Find(result interface{}, args ...interface{}) error {
	query := f.db
	if f.selectQuery != nil {
		query = query.Select(f.selectQuery)
	}
	return query.Find(result, args...).Error
}

func (f *GormFind) First(result interface{}, args ...interface{}) error {
	query := f.db
	if f.selectQuery != nil {
		query = query.Select(f.selectQuery)
	}
	return query.First(result, args...).Error
}

func (f *GormFind) Select(query interface{}, args ...interface{}) repo.Find {
	f.selectQuery = query
	return f
}

func (f *GormFind) Scan(result interface{}) error {
	return f.db.Scan(result).Error
}

func (f *GormFind) OrderBy(query string) repo.Find {
	f.db = f.db.Order(query)
	return f
}

func (f *GormFind) Group(query string) repo.Find {
	f.db = f.db.Group(query)
	return f
}

func (f *GormFind) Limit(limit int) repo.Find {
	f.db = f.db.Limit(limit)
	return f
}

func (f *GormFind) Rows() (*sql.Rows, error) {
	return f.db.Rows()
}

// GormTransaction реализация транзакций
type GormTransaction struct {
	db *gorm.DB
}

func NewGormTransaction(db *gorm.DB) *GormTransaction {
	return &GormTransaction{db: db}
}

func (t *GormTransaction) Begin(ctx context.Context) context.Context {
	tx := t.db.WithContext(ctx).Begin()
	return context.WithValue(ctx, "tx", tx)
}

func (t *GormTransaction) Commit(ctx context.Context) error {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok {
		return tx.Commit().Error
	}
	return fmt.Errorf("no transaction found in context")
}

func (t *GormTransaction) Rollback(ctx context.Context) error {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok {
		return tx.Rollback().Error
	}
	return fmt.Errorf("no transaction found in context")
}

func (t *GormTransaction) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx := t.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	ctx = context.WithValue(ctx, "tx", tx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(ctx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
