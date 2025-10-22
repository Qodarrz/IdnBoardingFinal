// repository/carbon_repository.go
package repository

import (
	"context"
	"database/sql"

	models "github.com/Qodarrz/fiber-app/model"
)

type CarbonRepository interface {
	// Vehicle
	FindVehicleByID(ctx context.Context, id int64) (*models.CarbonVehicle, error)
	FindVehicleByUserAndName(ctx context.Context, userID int64, name string) (*models.CarbonVehicle, error)
	CreateVehicle(ctx context.Context, v *models.CarbonVehicle) (*models.CarbonVehicle, error)
	ListUserVehicles(ctx context.Context, userID int64) ([]*models.CarbonVehicleWithLog, error)
	CreateVehicleLog(ctx context.Context, log *models.CarbonVehicleLog) error
	GetVehicleLogs(ctx context.Context, vehicleID int64) ([]*models.CarbonVehicleLog, error)
	GetVehicleLogByID(ctx context.Context, userID, logID int64) (*models.CarbonVehicleLog, error)

	FindElectronicsByID(ctx context.Context, id int64) (*models.CarbonElectronic, error)
	FindElectronicsByUserAndName(ctx context.Context, userID int64, deviceName string) (*models.CarbonElectronic, error)
	CreateElectronics(ctx context.Context, e *models.CarbonElectronic) (*models.CarbonElectronic, error)
	ListUserElectronics(ctx context.Context, userID int64) ([]*models.CarbonElectronic, error)
	CreateElectronicsLog(ctx context.Context, log *models.CarbonElectronicLog) error
	GetElectronicsLogs(ctx context.Context, deviceID int64) ([]*models.CarbonElectronicLog, error)

	UpdateVehicle(ctx context.Context, v *models.CarbonVehicle) error
	DeleteVehicle(ctx context.Context, id int64) error
	DeleteVehicleLogs(ctx context.Context, vehicleID int64) error
	GetAllVehicleLogsByUser(ctx context.Context, userID int64) ([]*models.CarbonVehicleLog, error)

	// Electronic methods
	UpdateElectronic(ctx context.Context, e *models.CarbonElectronic) error
	DeleteElectronic(ctx context.Context, id int64) error
	DeleteElectronicLogs(ctx context.Context, deviceID int64) error
	GetAllElectronicLogsByUser(ctx context.Context, userID int64) ([]*models.CarbonElectronicLog, error)
}

type carbonRepository struct {
	db *sql.DB
}

func NewCarbonRepository(db *sql.DB) CarbonRepository {
	return &carbonRepository{db: db}
}


func (r *carbonRepository) GetVehicleLogByID(ctx context.Context, userID, logID int64) (*models.CarbonVehicleLog, error) {
    var log models.CarbonVehicleLog
    query := `SELECT * FROM carbon_vehicle_logs WHERE id = $1 AND user_id = $2`
    err := r.db.QueryRowContext(ctx, query, logID, userID).Scan(&log.ID, &log.VehicleID, &log.StartLat, &log.StartLon, &log.EndLat, &log.EndLon, &log.DistanceKm, &log.DurationMinutes, &log.CarbonEmission, &log.LoggedAt)
    if err != nil {
        return nil, err
    }
    return &log, nil
}


func (r *carbonRepository) FindVehicleByID(ctx context.Context, id int64) (*models.CarbonVehicle, error) {
	var v models.CarbonVehicle
	err := r.db.QueryRowContext(ctx, `SELECT id, user_id, vehicle_type, fuel_type, name FROM carbon_vehicles WHERE id = $1`, id).
		Scan(&v.ID, &v.UserID, &v.VehicleType, &v.FuelType, &v.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *carbonRepository) FindVehicleByUserAndName(ctx context.Context, userID int64, name string) (*models.CarbonVehicle, error) {
	var v models.CarbonVehicle
	err := r.db.QueryRowContext(ctx, `SELECT id, user_id, vehicle_type, fuel_type, name FROM carbon_vehicles WHERE user_id = $1 AND name = $2`, userID, name).
		Scan(&v.ID, &v.UserID, &v.VehicleType, &v.FuelType, &v.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *carbonRepository) CreateVehicle(ctx context.Context, v *models.CarbonVehicle) (*models.CarbonVehicle, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `INSERT INTO carbon_vehicles (user_id, vehicle_type, fuel_type, name) VALUES ($1, $2, $3, $4) RETURNING id`,
		v.UserID, v.VehicleType, v.FuelType, v.Name).Scan(&id)
	if err != nil {
		return nil, err
	}

	v.ID = id
	return v, nil
}

func (r *carbonRepository) ListUserVehicles(ctx context.Context, userID int64) ([]*models.CarbonVehicleWithLog, error) {
    query := `
    SELECT 
        v.id, v.user_id, v.vehicle_type, v.fuel_type, v.name,
        l.id AS log_id, l.start_lat, l.start_lon, l.end_lat, l.end_lon,
        l.distance_km, l.duration_minutes, l.carbon_emission_g, l.logged_at
    FROM carbon_vehicles v
    LEFT JOIN LATERAL (
        SELECT * 
        FROM carbon_vehicle_logs 
        WHERE vehicle_id = v.id 
        ORDER BY logged_at DESC 
        LIMIT 1
    ) l ON true
    WHERE v.user_id = $1
    `

    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var vehicles []*models.CarbonVehicleWithLog
    for rows.Next() {
        var v models.CarbonVehicleWithLog
        var logID sql.NullInt64
        var startLat, startLon, endLat, endLon, distanceKm, carbonEmission sql.NullFloat64
        var durationMinutes sql.NullInt64
        var loggedAt sql.NullTime

        if err := rows.Scan(
            &v.ID, &v.UserID, &v.VehicleType, &v.FuelType, &v.Name,
            &logID, &startLat, &startLon, &endLat, &endLon,
            &distanceKm, &durationMinutes, &carbonEmission, &loggedAt,
        ); err != nil {
            return nil, err
        }

        if logID.Valid {
            v.LatestLog = &models.CarbonVehicleLog{
                ID:              logID.Int64,
                VehicleID:       v.ID,
                StartLat:        startLat.Float64,
                StartLon:        startLon.Float64,
                EndLat:          endLat.Float64,
                EndLon:          endLon.Float64,
                DistanceKm:      distanceKm.Float64,
                DurationMinutes: int(durationMinutes.Int64),
                CarbonEmission:  carbonEmission.Float64,
                LoggedAt:        loggedAt.Time,
            }
        }

        vehicles = append(vehicles, &v)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return vehicles, nil
}


func (r *carbonRepository) CreateVehicleLog(ctx context.Context, log *models.CarbonVehicleLog) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO carbon_vehicle_logs 
			(vehicle_id, start_lat, start_lon, end_lat, end_lon, distance_km, duration_minutes, carbon_emission_g) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`,
		log.VehicleID, log.StartLat, log.StartLon, log.EndLat, log.EndLon,
		log.DistanceKm, log.DurationMinutes, log.CarbonEmission,
	)
	return err
}

func (r *carbonRepository) GetVehicleLogs(ctx context.Context, vehicleID int64) ([]*models.CarbonVehicleLog, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, vehicle_id, start_lat, start_lon, end_lat, end_lon, 
		       distance_km, duration_minutes, carbon_emission_g, logged_at
		FROM carbon_vehicle_logs 
		WHERE vehicle_id = $1
	`, vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.CarbonVehicleLog
	for rows.Next() {
		var log models.CarbonVehicleLog
		if err := rows.Scan(
			&log.ID, &log.VehicleID, &log.StartLat, &log.StartLon,
			&log.EndLat, &log.EndLon, &log.DistanceKm, &log.DurationMinutes,
			&log.CarbonEmission, &log.LoggedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	return logs, rows.Err()
}

func (r *carbonRepository) FindElectronicsByID(ctx context.Context, id int64) (*models.CarbonElectronic, error) {
	var e models.CarbonElectronic
	err := r.db.QueryRowContext(ctx, `SELECT id, user_id, device_name, device_type, power_watts FROM carbon_electronics WHERE id = $1`, id).
		Scan(&e.ID, &e.UserID, &e.DeviceName, &e.DeviceType, &e.PowerWatts)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *carbonRepository) FindElectronicsByUserAndName(ctx context.Context, userID int64, deviceName string) (*models.CarbonElectronic, error) {
	var e models.CarbonElectronic
	err := r.db.QueryRowContext(ctx, `SELECT id, user_id, device_name, device_type, power_watts FROM carbon_electronics WHERE user_id = $1 AND device_name = $2`, userID, deviceName).
		Scan(&e.ID, &e.UserID, &e.DeviceName, &e.DeviceType, &e.PowerWatts)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *carbonRepository) CreateElectronics(ctx context.Context, e *models.CarbonElectronic) (*models.CarbonElectronic, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `INSERT INTO carbon_electronics (user_id, device_name, device_type, power_watts) VALUES ($1, $2, $3, $4) RETURNING id`,
		e.UserID, e.DeviceName, e.DeviceType, e.PowerWatts).Scan(&id)
	if err != nil {
		return nil, err
	}

	e.ID = id
	return e, nil
}

func (r *carbonRepository) ListUserElectronics(ctx context.Context, userID int64) ([]*models.CarbonElectronic, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, device_name, device_type, power_watts FROM carbon_electronics WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var electronics []*models.CarbonElectronic
	for rows.Next() {
		var e models.CarbonElectronic
		if err := rows.Scan(&e.ID, &e.UserID, &e.DeviceName, &e.DeviceType, &e.PowerWatts); err != nil {
			return nil, err
		}
		electronics = append(electronics, &e)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return electronics, nil
}

func (r *carbonRepository) CreateElectronicsLog(ctx context.Context, log *models.CarbonElectronicLog) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO carbon_electronics_logs (device_id, duration_hours, carbon_emission_g) VALUES ($1, $2, $3)`,
		log.DeviceID, log.DurationHours, log.CarbonEmission)
	return err
}

func (r *carbonRepository) GetElectronicsLogs(ctx context.Context, deviceID int64) ([]*models.CarbonElectronicLog, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, device_id, duration_hours, carbon_emission_g, logged_at FROM carbon_electronics_logs WHERE device_id = $1`, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.CarbonElectronicLog
	for rows.Next() {
		var log models.CarbonElectronicLog
		if err := rows.Scan(&log.ID, &log.DeviceID, &log.DurationHours, &log.CarbonEmission, &log.LoggedAt); err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *carbonRepository) UpdateVehicle(ctx context.Context, v *models.CarbonVehicle) error {
	_, err := r.db.ExecContext(ctx, `UPDATE carbon_vehicles SET vehicle_type = $1, fuel_type = $2, name = $3 WHERE id = $4`,
		v.VehicleType, v.FuelType, v.Name, v.ID)
	return err
}

func (r *carbonRepository) DeleteVehicle(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM carbon_vehicles WHERE id = $1`, id)
	return err
}

func (r *carbonRepository) DeleteVehicleLogs(ctx context.Context, vehicleID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM carbon_vehicle_logs WHERE vehicle_id = $1`, vehicleID)
	return err
}

func (r *carbonRepository) GetAllVehicleLogsByUser(ctx context.Context, userID int64) ([]*models.CarbonVehicleLog, error) {
	query := `
		SELECT cvl.id, cvl.vehicle_id, cvl.start_lat, cvl.start_lon, cvl.end_lat, cvl.end_lon,
		       cvl.distance_km, cvl.duration_minutes, cvl.carbon_emission_g, cvl.logged_at
		FROM carbon_vehicle_logs cvl
		JOIN carbon_vehicles cv ON cvl.vehicle_id = cv.id
		WHERE cv.user_id = $1
		ORDER BY cvl.id DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.CarbonVehicleLog
	for rows.Next() {
		var log models.CarbonVehicleLog
		if err := rows.Scan(
			&log.ID, &log.VehicleID, &log.StartLat, &log.StartLon, 
			&log.EndLat, &log.EndLon, &log.DistanceKm, &log.DurationMinutes,
			&log.CarbonEmission, &log.LoggedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	return logs, rows.Err()
}


func (r *carbonRepository) UpdateElectronic(ctx context.Context, e *models.CarbonElectronic) error {
	_, err := r.db.ExecContext(ctx, `UPDATE carbon_electronics SET device_name = $1, device_type = $2, power_watts = $3 WHERE id = $4`,
		e.DeviceName, e.DeviceType, e.PowerWatts, e.ID)
	return err
}

func (r *carbonRepository) DeleteElectronic(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM carbon_electronics WHERE id = $1`, id)
	return err
}

func (r *carbonRepository) DeleteElectronicLogs(ctx context.Context, deviceID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM carbon_electronics_logs WHERE device_id = $1`, deviceID)
	return err
}

func (r *carbonRepository) GetAllElectronicLogsByUser(ctx context.Context, userID int64) ([]*models.CarbonElectronicLog, error) {
	query := `
		SELECT cel.id, cel.device_id, cel.duration_hours, cel.carbon_emission_g, cel.logged_at
		FROM carbon_electronics_logs cel
		JOIN carbon_electronics ce ON cel.device_id = ce.id
		WHERE ce.user_id = $1
		ORDER BY cel.id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.CarbonElectronicLog
	for rows.Next() {
		var log models.CarbonElectronicLog
		if err := rows.Scan(&log.ID, &log.DeviceID, &log.DurationHours, &log.CarbonEmission, &log.LoggedAt); err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
