package repository

import (
	"context"
	"database/sql"
	"time"
)

type UserCustomEndpointRepoInterface interface {
	GetUserCustomData(ctx context.Context, userID int64) (*UserCustomData, error)
	GetLeaderboard(ctx context.Context, page, limit int, timeRange string) ([]LeaderboardEntry, int, error)
}

type userCustomEndpointRepo struct {
	db *sql.DB
}

func NewUserCustomEndpointRepo(db *sql.DB) UserCustomEndpointRepoInterface {
	return &userCustomEndpointRepo{db: db}
}

type UserCustomData struct {
	User                    UserDetail
	Vehicles                []VehicleWithCarbon
	Electronics             []ElectronicWithCarbon
	Missions                []MissionProgress
	Badges                  []UserBadge
	PointTransactions       []PointTransaction
	ActivityLogs            []ActivityLog
	Orders                  []Order
	UserPoints              UserPoints
	MonthlyVehicleCarbon    []MonthlyCarbon // Tambahan baru
	MonthlyElectronicCarbon []MonthlyCarbon // Tambahan baru
}

// Tambahkan struct MonthlyCarbon
type MonthlyCarbon struct {
	Month       time.Time
	TotalCarbon float64
}

type UserDetail struct {
	ID          int64      `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	FullName    string     `json:"full_name"`
	AvatarURL   string     `json:"avatar_url"`
	Birthdate   *time.Time `json:"birthdate"`
	Gender      string     `json:"gender"`
	TotalPoints int        `json:"total_points"`
	CreatedAt   time.Time  `json:"created_at"`
}

type VehicleWithCarbon struct {
	ID          int64     `json:"id"`
	VehicleType string    `json:"vehicle_type"`
	FuelType    string    `json:"fuel_type"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	TotalLogs   int       `json:"total_logs"`
	TotalCarbon float64   `json:"total_carbon_emission_g"`
}

type ElectronicWithCarbon struct {
	ID          int64     `json:"id"`
	DeviceName  string    `json:"device_name"`
	DeviceType  string    `json:"device_type"`
	PowerWatts  int       `json:"power_watts"`
	CreatedAt   time.Time `json:"created_at"`
	TotalLogs   int       `json:"total_logs"`
	TotalCarbon float64   `json:"total_carbon_emission_g"`
}

type MissionProgress struct {
	ID            int64      `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	MissionType   string     `json:"mission_type"`
	PointsReward  int        `json:"points_reward"`
	TargetValue   float64    `json:"target_value"`
	ProgressValue float64    `json:"progress_value"`
	CompletedAt   *time.Time `json:"completed_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type UserBadge struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	RedeemedAt  time.Time `json:"redeemed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type PointTransaction struct {
	ID            int64     `json:"id"`
	Amount        int       `json:"amount"`
	Direction     string    `json:"direction"`
	Source        string    `json:"source"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   int64     `json:"reference_id"`
	Note          string    `json:"note"`
	CreatedAt     time.Time `json:"created_at"`
}

type ActivityLog struct {
	ID        int64     `json:"id"`
	Activity  string    `json:"activity"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID          int64       `json:"id"`
	TotalPoints int         `json:"total_points"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	Items       []OrderItem `json:"items"`
}

type OrderItem struct {
	ID              int64  `json:"id"`
	ItemName        string `json:"item_name"`
	Qty             int    `json:"qty"`
	PriceEachPoints int    `json:"price_each_points"`
}

type LeaderboardEntry struct {
	User              UserSimple
	TotalPoints       int     `json:"total_points"`
	CompletedMissions int     `json:"completed_missions"`
	Score             float64 `json:"score"`
	CarbonReduction   float64 `json:"carbon_reduction_g"`
}

type UserSimple struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url"`
}

type UserPoints struct {
	TotalPoints int `json:"total_points"`
}

func (r *userCustomEndpointRepo) GetUserCustomData(ctx context.Context, userID int64) (*UserCustomData, error) {
	data := &UserCustomData{}

	// Query user dan profile
	userQuery := `
	SELECT u.id, u.username, u.email, u.role, u.created_at, 
	       COALESCE(up.full_name, '')   AS full_name,
	       COALESCE(up.avatar_url, '')  AS avatar_url,
	       up.birthdate,
	       COALESCE(up.gender, '')      AS gender
	FROM users u
	LEFT JOIN user_profiles up ON u.id = up.user_id
	WHERE u.id = $1
	`
	var role sql.NullString
	err := r.db.QueryRowContext(ctx, userQuery, userID).Scan(
		&data.User.ID,
		&data.User.Username,
		&data.User.Email,
		&role,
		&data.User.CreatedAt,
		&data.User.FullName,
		&data.User.AvatarURL,
		&data.User.Birthdate,
		&data.User.Gender,
	)
	if err != nil {
		return nil, err
	}
	if role.Valid {
		data.User.Role = role.String
	}

	// Query total points terpisah
	pointsQuery := `
	SELECT COALESCE(total_points, 10) 
	FROM points 
	WHERE user_id = $1
	`
	err = r.db.QueryRowContext(ctx, pointsQuery, userID).Scan(&data.UserPoints.TotalPoints)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Query untuk vehicles
	vehicleQuery := `
		SELECT cv.id, cv.vehicle_type, cv.fuel_type, cv.name, cv.created_at,
		       COUNT(cvl.id) as total_logs,
		       COALESCE(SUM(cvl.carbon_emission_g), 0) as total_carbon
		FROM carbon_vehicles cv
		LEFT JOIN carbon_vehicle_logs cvl ON cv.id = cvl.vehicle_id
		WHERE cv.user_id = $1
		GROUP BY cv.id
	`
	rows, _ := r.db.QueryContext(ctx, vehicleQuery, userID)
	for rows.Next() {
		var vehicle VehicleWithCarbon
		rows.Scan(&vehicle.ID, &vehicle.VehicleType, &vehicle.FuelType, &vehicle.Name, &vehicle.CreatedAt, &vehicle.TotalLogs, &vehicle.TotalCarbon)
		data.Vehicles = append(data.Vehicles, vehicle)
	}
	rows.Close()

	electronicQuery := `
		SELECT ce.id, ce.device_name, ce.device_type, ce.power_watts, ce.created_at,
		       COUNT(cel.id) as total_logs,
		       COALESCE(SUM(cel.carbon_emission_g), 0) as total_carbon
		FROM carbon_electronics ce
		LEFT JOIN carbon_electronics_logs cel ON ce.id = cel.device_id
		WHERE ce.user_id = $1
		GROUP BY ce.id
	`
	rows, _ = r.db.QueryContext(ctx, electronicQuery, userID)
	for rows.Next() {
		var electronic ElectronicWithCarbon
		rows.Scan(&electronic.ID, &electronic.DeviceName, &electronic.DeviceType, &electronic.PowerWatts, &electronic.CreatedAt, &electronic.TotalLogs, &electronic.TotalCarbon)
		data.Electronics = append(data.Electronics, electronic)
	}
	rows.Close()

	// Query untuk missions
	missionQuery := `
		SELECT m.id, 
       m.title, 
       m.description, 
       m.mission_type, 
       m.points_reward, 
       m.target_value,
       ump.progress_value, 
       um.completed_at, 
       m.created_at
FROM missions m
LEFT JOIN user_mission_progress ump 
       ON m.id = ump.mission_id AND ump.user_id = $1
INNER JOIN user_missions um 
       ON m.id = um.mission_id AND um.user_id = $1
WHERE um.completed_at IS NOT NULL
  AND (m.expired_at IS NULL OR m.expired_at > NOW());
	`
	rows, _ = r.db.QueryContext(ctx, missionQuery, userID)
	for rows.Next() {
		var mission MissionProgress
		rows.Scan(&mission.ID, &mission.Title, &mission.Description, &mission.MissionType, &mission.PointsReward, &mission.TargetValue, &mission.ProgressValue, &mission.CompletedAt, &mission.CreatedAt)
		data.Missions = append(data.Missions, mission)
	}
	rows.Close()

	// Query untuk badges
	badgeQuery := `
		SELECT b.id, b.name, b.image_url, b.description, ub.redeemed_at, b.created_at
		FROM badges b
		JOIN user_badges ub ON b.id = ub.badge_id
		WHERE ub.user_id = $1
	`
	rows, _ = r.db.QueryContext(ctx, badgeQuery, userID)
	for rows.Next() {
		var badge UserBadge
		rows.Scan(&badge.ID, &badge.Name, &badge.ImageURL, &badge.Description, &badge.RedeemedAt, &badge.CreatedAt)
		data.Badges = append(data.Badges, badge)
	}
	rows.Close()

	// Query untuk point transactions
	pointQuery := `
		SELECT id, amount, direction, source, reference_type, reference_id, note, created_at
		FROM point_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 50
	`
	rows, _ = r.db.QueryContext(ctx, pointQuery, userID)
	for rows.Next() {
		var transaction PointTransaction
		rows.Scan(&transaction.ID, &transaction.Amount, &transaction.Direction, &transaction.Source, &transaction.ReferenceType, &transaction.ReferenceID, &transaction.Note, &transaction.CreatedAt)
		data.PointTransactions = append(data.PointTransactions, transaction)
	}
	rows.Close()

	// Query untuk activity logs
	activityQuery := `
		SELECT id, activity, created_at
		FROM activity_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 50
	`
	rows, _ = r.db.QueryContext(ctx, activityQuery, userID)
	for rows.Next() {
		var activity ActivityLog
		rows.Scan(&activity.ID, &activity.Activity, &activity.CreatedAt)
		data.ActivityLogs = append(data.ActivityLogs, activity)
	}
	rows.Close()

	// Query untuk orders
	orderQuery := `
		SELECT o.id, o.total_points, o.status, o.created_at
		FROM orders o
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
		LIMIT 20
	`
	rows, _ = r.db.QueryContext(ctx, orderQuery, userID)
	for rows.Next() {
		var order Order
		rows.Scan(&order.ID, &order.TotalPoints, &order.Status, &order.CreatedAt)

		// Query untuk order items
		itemQuery := `
			SELECT oi.id, si.name, oi.qty, oi.price_each_points
			FROM order_items oi
			JOIN store_items si ON oi.item_id = si.id
			WHERE oi.order_id = $1
		`
		itemRows, _ := r.db.QueryContext(ctx, itemQuery, order.ID)
		for itemRows.Next() {
			var item OrderItem
			itemRows.Scan(&item.ID, &item.ItemName, &item.Qty, &item.PriceEachPoints)
			order.Items = append(order.Items, item)
		}
		itemRows.Close()

		data.Orders = append(data.Orders, order)
	}
	rows.Close()

	carbonVehicleMonthlyQuery := `
		SELECT 
			DATE_TRUNC('month', cvl.logged_at) as month,
			COALESCE(SUM(cvl.carbon_emission_g), 0) as total_carbon
		FROM carbon_vehicle_logs cvl
		JOIN carbon_vehicles cv ON cvl.vehicle_id = cv.id
		WHERE cv.user_id = $1 
			AND cvl.logged_at >= DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '5 months'
		GROUP BY DATE_TRUNC('month', cvl.logged_at)
		ORDER BY month DESC
		LIMIT 6
	`
	rows, err = r.db.QueryContext(ctx, carbonVehicleMonthlyQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data.MonthlyVehicleCarbon = make([]MonthlyCarbon, 0)
	for rows.Next() {
		var monthlyCarbon MonthlyCarbon
		if err := rows.Scan(&monthlyCarbon.Month, &monthlyCarbon.TotalCarbon); err != nil {
			return nil, err
		}
		data.MonthlyVehicleCarbon = append(data.MonthlyVehicleCarbon, monthlyCarbon)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Query untuk total karbon per bulan (6 bulan terakhir) - ELECTRONIC
	carbonElectronicMonthlyQuery := `
		SELECT 
			DATE_TRUNC('month', cel.logged_at) as month,
			COALESCE(SUM(cel.carbon_emission_g), 0) as total_carbon
		FROM carbon_electronics_logs cel
		JOIN carbon_electronics ce ON cel.device_id = ce.id
		WHERE ce.user_id = $1 
			AND cel.logged_at >= DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '5 months'
		GROUP BY DATE_TRUNC('month', cel.logged_at)
		ORDER BY month DESC
		LIMIT 6
	`
	rows, err = r.db.QueryContext(ctx, carbonElectronicMonthlyQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data.MonthlyElectronicCarbon = make([]MonthlyCarbon, 0)
	for rows.Next() {
		var monthlyCarbon MonthlyCarbon
		if err := rows.Scan(&monthlyCarbon.Month, &monthlyCarbon.TotalCarbon); err != nil {
			return nil, err
		}
		data.MonthlyElectronicCarbon = append(data.MonthlyElectronicCarbon, monthlyCarbon)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func (r *userCustomEndpointRepo) GetLeaderboard(ctx context.Context, page, limit int, timeRange string) ([]LeaderboardEntry, int, error) {
	offset := (page - 1) * limit
	var entries []LeaderboardEntry

	// Build time range condition
	var timeCondition string
	switch timeRange {
	case "day":
		timeCondition = "AND um.completed_at >= CURRENT_DATE"
	case "week":
		timeCondition = "AND um.completed_at >= DATE_TRUNC('week', CURRENT_DATE)"
	case "month":
		timeCondition = "AND um.completed_at >= DATE_TRUNC('month', CURRENT_DATE)"
	default:
		timeCondition = ""
	}

	query := `
		SELECT u.id, u.username, 
       COALESCE(up.full_name, '') as full_name,
       COALESCE(up.avatar_url, '') as avatar_url,
       COALESCE(p.total_points, 0) as total_points,
       COUNT(DISTINCT um.mission_id) as completed_missions,
       (COALESCE(p.total_points, 0) * 0.7 + COUNT(DISTINCT um.mission_id) * 0.3) as score,
       COALESCE((
           SELECT SUM(carbon_emission_g) 
           FROM (
               SELECT carbon_emission_g FROM carbon_vehicle_logs cvl
               JOIN carbon_vehicles cv ON cvl.vehicle_id = cv.id
               WHERE cv.user_id = u.id
               UNION ALL
               SELECT carbon_emission_g FROM carbon_electronics_logs cel
               JOIN carbon_electronics ce ON cel.device_id = ce.id
               WHERE ce.user_id = u.id
           ) emissions
       ), 0) as carbon_reduction
FROM users u
LEFT JOIN user_profiles up ON u.id = up.user_id
LEFT JOIN points p ON u.id = p.user_id
LEFT JOIN user_missions um ON u.id = um.user_id AND um.completed_at IS NOT NULL ` + timeCondition + `
GROUP BY u.id, u.username, up.full_name, up.avatar_url, p.total_points
ORDER BY score DESC, total_points DESC, completed_missions DESC
LIMIT $1 OFFSET $2


	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry LeaderboardEntry
		err := rows.Scan(
			&entry.User.ID, &entry.User.Username, &entry.User.FullName, &entry.User.AvatarURL,
			&entry.TotalPoints, &entry.CompletedMissions, &entry.Score, &entry.CarbonReduction,
		)
		if err != nil {
			return nil, 0, err
		}
		entries = append(entries, entry)
	}

	// Get total count
	countQuery := `SELECT COUNT(DISTINCT u.id) FROM users u`
	var total int
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}
