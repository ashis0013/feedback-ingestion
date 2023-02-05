package repository

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ashis0013/feedback-ingestion/models"
	. "github.com/ashis0013/gollections"
	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

const (
	host  = "localhost"
	port  = 5432
	delim = ","
)

type PostgresRepository struct {
	db             *sqlx.DB
	dataSourceName string
	sourceCache    map[string][]string
}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{
		db: nil,
		dataSourceName: fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, os.Getenv("PG_USER"), os.Getenv("PASS"), os.Getenv("DB_NAME")),
		sourceCache: make(map[string][]string),
	}
}

func (r *PostgresRepository) Init() {
	go r.cacheClearRoutine()
	var err error
	r.db, err = sqlx.Connect("postgres", r.dataSourceName)

	if err != nil {
		log.Fatalln("Failed to connect")
		log.Fatalln(err)
		return
	}
	r.createTables()
}

func (r *PostgresRepository) Terminate() {
	if r.db != nil {
		r.db.Close()
	}
}

func (r *PostgresRepository) AddRecord(records []*models.Feedback) error {
	for _, record := range records {
		record.RecordId = uuid.New().String()
		record.CreatedOn = time.Now()
		err := transactExcec(r.db, func(tx *sqlx.Tx) error {
			_, err := tx.NamedExec(insertRecord, record)
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresRepository) GetRecords(filter *models.QueryFilter) (*models.GetFeedbacksResponse, error) {
	q := filter.BuildSelectquery()
	feedbacks := []*models.Feedback{}
	err := transactQuery(r.db, &feedbacks, func(tx *sqlx.Tx) (*sqlx.Rows, error) {
		args := filter.GetQueryArgs()
		if len(args) == 0 {
			return tx.Queryx(q)
		}
		return tx.Queryx(q, args...)
	})
	if err != nil {
		return nil, err
	}

	r.getSourceDataCached(filter.SourceId)
	response := new(models.GetFeedbacksResponse)
	response.Feedbacks = Map(feedbacks, func(it *models.Feedback) *models.FeedbackResponse {
		return it.ToFeedbackResponse(r.sourceCache)
	})
	return response, nil
}

func (r *PostgresRepository) AddTenant(tenantName string, tags []string) error {
	tenant := &models.TenantRecord{
		TenantId:   uuid.New().String(),
		TenantName: tenantName,
		Tags:       strings.Join(tags, delim),
	}

	return transactExcec(r.db, func(tx *sqlx.Tx) error {
		_, err := tx.NamedExec(insertTenant, tenant)
		return err
	})
}

func (r *PostgresRepository) FetchTags() (map[string]string, error) {
	tenants := []*models.TenantRecord{}
	err := transactQuery(r.db, &tenants, func(tx *sqlx.Tx) (*sqlx.Rows, error) {
		return tx.Queryx(selectTenants)
	})
	if err != nil {
		return nil, err
	}

	tags := make(map[string]string)
	ForEach(tenants, func(t *models.TenantRecord) {
		ForEach(strings.Split(t.Tags, delim), func(tag string) {
			tags[tag] = t.TenantId
		})
	})
	return tags, nil
}

func (r *PostgresRepository) AddSource(sourceName string, metadata string) (string, error) {
	id := uuid.New().String()
	source := &models.SourceRecord{
		SourceId:       id,
		SourceName:     sourceName,
		SourceMetadata: metadata,
	}

	err := transactExcec(r.db, func(tx *sqlx.Tx) error {
		_, err := tx.NamedExec(insertSource, source)
		return err
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *PostgresRepository) GetSources() (map[string]string, error) {
	sources := []*models.SourceRecord{}
	err := transactQuery(r.db, &sources, func(tx *sqlx.Tx) (*sqlx.Rows, error) {
		return tx.Queryx(selectSources)
	})
	if err != nil {
		return nil, err
	}
	sourceMap := make(map[string]string)
	ForEach(sources, func(source *models.SourceRecord) {
		sourceMap[source.SourceName] = source.SourceId
	})
	return sourceMap, nil
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

func (r *PostgresRepository) getSourceDataCached(sourceId string) {
	_, sourceExists := r.sourceCache[sourceId]
	if len(r.sourceCache) != 0 && (sourceId == "" || sourceExists) {
		return
	}
	sources := []*models.SourceRecord{}
	transactQuery(r.db, &sources, func(tx *sqlx.Tx) (*sqlx.Rows, error) {
		return tx.Queryx(selectSources)
	})

	for _, source := range sources {
		r.sourceCache[source.SourceId] = []string{source.SourceName, source.SourceMetadata}
	}
}

func (r *PostgresRepository) cacheClearRoutine() {
	if r == nil {
		return
	}
	for {
		if r == nil {
			break
		}
		r.sourceCache = make(map[string][]string)
		time.Sleep(time.Second * 10)
	}
}
