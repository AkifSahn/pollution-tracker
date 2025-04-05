package pollution

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type PollutionRepo interface {
	GetPollutionValueByPosition(ctx context.Context, latitude, longitude float64, from, to time.Time) ([]PollutionValueResponse, error)
	GetAnomaliesWithinTimeRange(ctx context.Context, from, to time.Time) ([]Pollution, error)
	GetPollutionDensityByRegion(ctx context.Context, radius, latitude, longitude float64, from, to time.Time) (float64, error)

	GetMeanAndStd(ctx context.Context, pollutant string, radius, latitude, longitude float64, from, to time.Time) (float64, float64, error)

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

func (repo *PollutionRepoImpl) GetPollutionValueByPosition(ctx context.Context, latitude, longitude float64, from, to time.Time) ([]PollutionValueResponse, error) {
	query := `
    SELECT time, value, pollutant FROM air_pollution 
    WHERE latitude=$1 AND longitude=$2 
    AND time BETWEEN $3 AND $4 
    ORDER BY pollutant, time DESC;
    `
	rows, err := repo.DB.Query(ctx, query, latitude, longitude, from, to)
	if err != nil {
		return nil, fmt.Errorf("Unable to query - %s", err.Error())
	}
	defer rows.Close()

	var pollutions []PollutionValueResponse
	for rows.Next() {
		var pollution PollutionValueResponse
		var t time.Time
		err = rows.Scan(&pollution.Time, &pollution.Value, &pollution.Pollutant)
		if err != nil {
			return nil, fmt.Errorf("Unable to scan - %s", err.Error())
		}

		fmt.Println("time: ", t)

		pollutions = append(pollutions, pollution)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows Error: %v\n", rows.Err())
	}

	return pollutions, nil
}

func (repo *PollutionRepoImpl) GetAnomaliesWithinTimeRange(ctx context.Context, from, to time.Time) ([]Pollution, error) {
	query := `
    SELECT latitude, longitude, value, is_anomaly, pollutant from air_pollution
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
			&pollution.Value, &pollution.IsAnomaly, &pollution.Pollutant)
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

func (repo *PollutionRepoImpl) GetPollutionDensityByRegion(ctx context.Context, radius, latitude, longitude float64, from, to time.Time) (float64, error) {
	query := `
    SELECT AVG(value) 
    FROM air_pollution 
    WHERE time >= $1 AND time <= $2 AND 
    ST_DWithin(
        geog,
        ST_MakePoint($3,$4)::geography,
        $5*1000
    );
    `
	row := repo.DB.QueryRow(ctx, query, from, to, longitude, latitude, radius)

	var density float64
	err := row.Scan(&density)
	if err != nil {
		return -1, fmt.Errorf("Unable to scan %s", err.Error())
	}

	return density, nil
}

func (repo *PollutionRepoImpl) GetMeanAndStd(ctx context.Context, pollutant string, radius, latitude, longitude float64, from, to time.Time) (float64, float64, error) {
	query := `
        SELECT COALESCE(AVG(value), 0), COALESCE(STDDEV_POP(value), 0)
        FROM air_pollution
        WHERE pollutant = $1
          AND time BETWEEN $2 AND $3 
          AND ST_DWithin(
              geog,
              ST_MakePoint($4,$5)::geography,
              $6*1000
        );
    `
	row := repo.DB.QueryRow(ctx, query, pollutant, from, to, longitude, latitude, radius)
	var mean, stddev float64
	if err := row.Scan(&mean, &stddev); err != nil {
		return 0, 0, fmt.Errorf("Unable to scan %s", err.Error())
	}
	return mean, stddev, nil
}

func (repo *PollutionRepoImpl) InsertPollution(ctx context.Context, pollution Pollution) error {
	query := `
    INSERT INTO air_pollution 
    (time, pollutant, value, is_anomaly, latitude, longitude, geog) 
    VALUES ($1,$2,$3,$4,$5,$6,
    ST_SetSRID(ST_MakePoint($6, $5), 4326)
    );
    `
	_, err := repo.DB.Exec(ctx, query,
		pollution.Time, pollution.Pollutant, pollution.Value,
		pollution.IsAnomaly, pollution.Latitude, pollution.Longitude)
	if err != nil {
		return fmt.Errorf("Failed to insert into database - %s", err.Error())
	}

	return nil
}
