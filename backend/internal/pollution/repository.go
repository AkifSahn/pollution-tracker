package pollution

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PollutionRepo interface {
	GetPollutionValueByPosition(ctx context.Context, latitude, longitude float64, from, to time.Time) ([]PollutionValueResponse, error)
	GetAnomaliesWithinTimeRange(ctx context.Context, from, to time.Time) ([]Pollution, error)
	GetAllPolutionWithinTimeRange(ctx context.Context, from, to time.Time) ([]Pollution, error)

	GetPollutionDensityOfRect(ctx context.Context, latFrom, latTo, longFrom, longTo float64, from, to time.Time, step time.Duration, pollutant string) ([]PollutionDensity, error)

	GetDistinctPollutants(ctx context.Context) ([]string, error)

	GetMeanAndStd(ctx context.Context, pollutant string, radius, latitude, longitude float64, from, to time.Time) (float64, float64, error)

	InsertPollution(ctx context.Context, pollution Pollution) error
}

type PollutionRepoImpl struct {
	DB *pgxpool.Pool
}

func NewPollutionRepo(db *pgxpool.Pool) *PollutionRepoImpl {
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

func (repo *PollutionRepoImpl) GetAllPolutionWithinTimeRange(ctx context.Context, from, to time.Time) ([]Pollution, error) {
	query := `
    SELECT latitude, longitude, value, is_anomaly, pollutant from air_pollution
    WHERE time BETWEEN $1 AND $2;
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

func (repo *PollutionRepoImpl) GetPollutionDensityOfRect(ctx context.Context, latFrom, latTo, longFrom, longTo float64, from, to time.Time, step time.Duration, pollutant string) ([]PollutionDensity, error) {
	query := `
    SELECT time_bucket($1, time) AS bucket, AVG(value) FROM air_pollution
    WHERE latitude BETWEEN $2 AND $3 
        AND longitude BETWEEN $4 AND $5 
        AND time BETWEEN $6 AND $7 
        AND pollutant = $8
    GROUP BY bucket 
    ORDER BY bucket;
    `

	rows, err := repo.DB.Query(ctx, query,
		step,
		latFrom,
		latTo,
		longFrom,
		longTo,
		from,
		to,
		pollutant,
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to query - %s", err.Error())
	}

	defer rows.Close()

	var result []PollutionDensity
	for rows.Next() {
		var density PollutionDensity
		density.Pollutant = pollutant
		err := rows.Scan(&density.Time, &density.Density)
		if err != nil {
			return nil, fmt.Errorf("Unable to scan - %s", err.Error())
		}
		result = append(result, density)
	}

	return result, nil

}

func (repo *PollutionRepoImpl) GetDistinctPollutants(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT pollutant FROM air_pollution;`

	rows, err := repo.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Unable to query - %s", err.Error())
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return nil, fmt.Errorf("Unable to scan - %s", err.Error())
		}
		result = append(result, s)
	}

	return result, nil
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
