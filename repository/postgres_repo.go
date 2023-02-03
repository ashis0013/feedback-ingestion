package repository

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ashis0013/feedback-ingestion/models"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

const (
	host = "localhost"
	port = 5432
)

type PostgresRepository struct {
	db             *sqlx.DB
	dataSourceName string
}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{
		db: nil,
		dataSourceName: fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, os.Getenv("PG_USER"), os.Getenv("PASS"), os.Getenv("DB_NAME")),
	}
}

func (r *PostgresRepository) Init() {
	var err error
	r.db, err = sqlx.Connect("postgres", r.dataSourceName)

	if err != nil {
		log.Fatalln("Failed to connect")
		log.Fatalln(err)
		return
	}
	r.createTables()
}

func (r *PostgresRepository) createTables() {
	if r.db == nil {
		log.Fatalln("Cannot connect to db")
	}

	transactExcec(r.db, func(tx *sqlx.Tx) error {
		for _, query := range createTables {
			_, err := tx.Exec(query)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PostgresRepository) Terminate() {
	if r.db != nil {
		r.db.Close()
	}
}

func (r *PostgresRepository) AddRecord(record *models.Feedback) {
	if r.db == nil {
		log.Fatalln("Cannot connect to db")
		return
	}

	record.CreatedOn = time.Now()
	transactExcec(r.db, func(tx *sqlx.Tx) error {
		_, err := tx.NamedExec(insertRecord, record)
		return err
	})
}

func (r *PostgresRepository) GetRecords(filter *models.QueryFilter) []*models.Feedback {
	if r.db == nil {
		log.Fatalln("Cannot connect to db")
		return []*models.Feedback{}
	}

	q := filter.BuildSelectquery()
	feedbacks := []*models.Feedback{}
	transactQuery(r.db, &feedbacks, func(tx *sqlx.Tx) (*sqlx.Rows, error) {
		args := getQueryArgs(filter)
		if len(args) == 0 {
			return tx.Queryx(q)
		}
		return tx.Queryx(q, args...)
	})
	for _, f := range feedbacks {
		println(f.RecordId)
	}
	return feedbacks
}

func getQueryArgs(filter *models.QueryFilter) []any {
	args := []any{}
	if filter.TenantId != "" {
		args = append(args, filter.TenantId)
	}
	if filter.SourceId != "" {
		args = append(args, filter.SourceId)
	}
	if filter.StartTime != nil {
		args = append(args, filter.StartTime.Format("2006-01-02 15:04:05"))
	}
	if filter.EndTime != nil {
		args = append(args, filter.EndTime.Format("2006-01-02 15:04:05"))
	}
	return args
}

func transformRecord(record *models.Feedback) []interface{} {
	return []interface{}{
		record.RecordId,
		record.SourceId,
		record.TenantId,
		record.PersonName,
		record.PersonEmail,
		record.FeedbackType,
		record.FeedbackContent,
		record.RecordLang,
		record.AdditionalData,
	}
}
