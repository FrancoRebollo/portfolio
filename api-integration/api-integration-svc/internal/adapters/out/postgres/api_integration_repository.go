package repository

import (
	"context"

	"github.com/FrancoRebollo/api-integration-svc/internal/domain"
	"github.com/FrancoRebollo/api-integration-svc/internal/platform/logger"
)

type ApiIntegrationRepository struct {
	dbPost *PostgresDB
}

func NewApiIntegrationRepository(dbPost *PostgresDB) *ApiIntegrationRepository {
	return &ApiIntegrationRepository{
		dbPost: dbPost,
	}
}

func (hr *ApiIntegrationRepository) GetDatabasesPing(ctx context.Context) ([]domain.Database, error) {
	databases := []domain.Database{}
	var fechaUltimaActividad string
	var mappedErr error
	var repoErr error

	query := `SELECT NOW()`

	rows, err := hr.dbPost.GetDB().QueryContext(ctx, query)
	if err != nil {
		mappedErr = hr.dbPost.MapPostgresError(err)
		repoErr = getRepoErr(mappedErr)
		logger.LoggerError().WithError(err).Error(repoErr)
		return databases, repoErr
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&fechaUltimaActividad)
		if err != nil {
			mappedErr = hr.dbPost.MapPostgresError(err)
			repoErr = getRepoErr(mappedErr)
			logger.LoggerError().WithError(err).Error(repoErr)
			return databases, repoErr
		}
	}

	if err = rows.Err(); err != nil {
		mappedErr = hr.dbPost.MapPostgresError(err)
		repoErr = getRepoErr(mappedErr)
		logger.LoggerError().WithError(err).Error(repoErr)
		return databases, repoErr
	}

	databases = append(databases, domain.Database{
		Base:                     "POSTGRES",
		FechaHoraUltimaActividad: fechaUltimaActividad,
	})

	return databases, nil
}

func (hr *ApiIntegrationRepository) CaptureEvent(ctx context.Context, reqCaptureEvent domain.Event) error {
	var idPersona int

	insert := `INSERT INTO api_int.event_example 
		(id_event_example,event_type,event_content,actualizado_por) VALUES ($1,$2,$3,$4,$5)`

	_, err := hr.dbPost.GetDB().ExecContext(ctx, insert, reqCaptureEvent.IdEvent, reqCaptureEvent.EventType, reqCaptureEvent.EventContent, idPersona)

	if err != nil {
		return err
	}

	return nil
}
