package pollution

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type PollutionRepo interface {
	GetPollutionValueByPosition(longitude, latitude float64) (float64, error)
	GetAnomaliesWithinTimeRange(from, to time.Time) ([]Pollution, error)
	GetPollutionDensityByRegion(ctx context.Context, latitude, longitude float64) (float64, error)

	InsertPollution(ctx context.Context, pollution Pollution) error
}

type PollutionRepoImpl struct {
	DB *pgx.Conn
}

func NewPollutionRepo(db *pgx.Conn) *PollutionRepoImpl {
	return &PollutionRepoImpl{
		DB: db,
	}
}

func (repo *PollutionRepoImpl) GetPollutionValueByPosition(ctx context.Context, latitude, longitude float64) (float64, error) {
	query := `
    SELECT value FROM air_pollution 
    WHERE latitude=$1 AND longitude=$2 
    `
	row := repo.DB.QueryRow(ctx, query, latitude, longitude)

	var val float64
	err := row.Scan(&val)
	if err != nil {
		return -1, fmt.Errorf("Unable to scan - %s", err.Error())
	}

	return val, nil
}

func (repo *PollutionRepoImpl) GetAnomaliesWithinTimeRange(ctx context.Context, from, to time.Time) ([]Pollution, error) {
	query := `
    SELECT latitude, longitude, region, value, is_anomaly, pollutant from air_pollution
    WHERE time >= $1 AND time <= $2 AND is_anomaly=true;
    `
	rows, err := repo.DB.Query(ctx, query, from, to)
	if err != nil {
		return nil, fmt.Errorf("Unable to query - %s", err.Error())
	}
	defer rows.Close()

	var pollutions []Pollution
	for rows.Next() {
		var pollution Pollution
		err = rows.Scan(&pollution.Latitude, &pollution.Longitude,
			&pollution.Region, &pollution.Value, &pollution.IsAnomaly, &pollution.Pollutant)
		if err != nil {
			return nil, fmt.Errorf("Unable to scan - %s", err.Error())
		}
		pollutions = append(pollutions, pollution)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows Error: %v\n", rows.Err())
	}

	return pollutions, nil
}

func (repo *PollutionRepoImpl) GetPollutionDensityByRegion(ctx context.Context, region string, from, to time.Time) (float64, error) {
	query := `
    SELECT AVG(value) 
    FROM air_pollution 
    WHERE region=$1 AND time >= $2 AND time <= $3;
    `
	row := repo.DB.QueryRow(ctx, query, region, from, to)

	var density float64
	err := row.Scan(&density)
	if err != nil {
		return -1, fmt.Errorf("Unable to scan %s", err.Error())
	}

	return density, nil
}

func (repo *PollutionRepoImpl) InsertPollution(ctx context.Context, pollution Pollution) error {
	query := `
    INSERT INTO air_pollution 
    (time, latitude, longitude, region, value, is_anomaly, pollutant) 
    VALUES ($1,$2,$3,$4,$5,$6,$7)
    `
	_, err := repo.DB.Exec(ctx, query,
		pollution.Time, pollution.Latitude, pollution.Longitude,
		pollution.Region, pollution.Value, pollution.IsAnomaly, pollution.Pollutant)
	if err != nil {
		return fmt.Errorf("Failed to insert into database - %s", err.Error())
	}

	return nil
}
