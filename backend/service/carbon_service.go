// service/carbon_service.go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/Qodarrz/fiber-app/dto"
	models "github.com/Qodarrz/fiber-app/model"
	"github.com/Qodarrz/fiber-app/repository"
)

type CarbonServiceInterface interface {
	CreateVehicle(ctx context.Context, userID int64, req *dto.CreateVehicleDTO) (*models.CarbonVehicle, error)
	ListUserVehicles(ctx context.Context, userID int64) ([]*models.CarbonVehicleWithLog, error)
	AddVehicleLog(ctx context.Context, userID int64, req *dto.AddVehicleLogDTO) error
	GetVehicleLogs(ctx context.Context, userID, vehicleID int64) ([]*models.CarbonVehicleLog, error)
	GetVehicleLogByID(ctx context.Context, userID, logID int64) (*models.CarbonVehicleLog, error)

	CreateElectronic(ctx context.Context, userID int64, req *dto.CreateElectronicDTO) (*models.CarbonElectronic, error)
	ListUserElectronics(ctx context.Context, userID int64) ([]*models.CarbonElectronic, error)
	AddElectronicsLog(ctx context.Context, userID int64, req *dto.AddElectronicsLogDTO) error
	GetElectronicsLogs(ctx context.Context, userID, deviceID int64) ([]*models.CarbonElectronicLog, error)

	EditVehicle(ctx context.Context, userID, vehicleID int64, req *dto.EditVehicleDTO) (*models.CarbonVehicle, error)
	DeleteVehicle(ctx context.Context, userID, vehicleID int64) error
	GetAllVehicleLogs(ctx context.Context, userID int64) ([]*models.CarbonVehicleLog, error)

	EditElectronic(ctx context.Context, userID, deviceID int64, req *dto.EditElectronicDTO) (*models.CarbonElectronic, error)
	DeleteElectronic(ctx context.Context, userID, deviceID int64) error
	GetAllElectronicLogs(ctx context.Context, userID int64) ([]*models.CarbonElectronicLog, error)

	// Electronics methods would be similarly updated
}

type CarbonService struct {
	carbonRepo  repository.CarbonRepository
	missionRepo repository.CheckMissionRepositoryInterface
}

func NewCarbonService(carbonRepo repository.CarbonRepository, missionRepo repository.CheckMissionRepositoryInterface) *CarbonService {
	return &CarbonService{
		carbonRepo:  carbonRepo,
		missionRepo: missionRepo,
	}
}

// ======================== VEHICLE ========================

func (s *CarbonService) CreateVehicle(ctx context.Context, userID int64, req *dto.CreateVehicleDTO) (*models.CarbonVehicle, error) {
	// Check if vehicle with same name already exists for this user
	existing, err := s.carbonRepo.FindVehicleByUserAndName(ctx, userID, req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("vehicle with this name already exists for this user")
	}

	vehicle := &models.CarbonVehicle{
		UserID:      userID,
		VehicleType: models.VehicleType(req.VehicleType),
		FuelType:    models.FuelType(req.FuelType),
		Name:        req.Name,
	}

	return s.carbonRepo.CreateVehicle(ctx, vehicle)
}

func (s *CarbonService) ListUserVehicles(ctx context.Context, userID int64) ([]*models.CarbonVehicleWithLog, error) {
	return s.carbonRepo.ListUserVehicles(ctx, userID)
}

func (s *CarbonService) AddVehicleLog(ctx context.Context, userID int64, req *dto.AddVehicleLogDTO) error {
	var vehicle *models.CarbonVehicle
	var err error

	if req.VehicleID != nil {
		vehicle, err = s.carbonRepo.FindVehicleByID(ctx, *req.VehicleID)
		if err != nil {
			return err
		}
		if vehicle == nil {
			return errors.New("vehicle not found")
		}
		if vehicle.UserID != userID {
			return errors.New("vehicle does not belong to user")
		}
	} else {
		existing, err := s.carbonRepo.FindVehicleByUserAndName(ctx, userID, req.VehicleName)
		if err != nil {
			return err
		}
		if existing != nil {
			vehicle = existing
		} else {
			vehicle, err = s.carbonRepo.CreateVehicle(ctx, &models.CarbonVehicle{
				UserID:      userID,
				VehicleType: models.VehicleType(req.VehicleType),
				FuelType:    models.FuelType(req.FuelType),
				Name:        req.VehicleName,
			})
			if err != nil {
				return err
			}
		}
	}

	var factor float64
	switch vehicle.FuelType {
	case models.FuelPetrol:
		factor = 0.161
	case models.FuelDiesel:
		factor = 0.162
	case models.FuelElectric:
		factor = 0.095
	case models.FuelNone:
		factor = 0
	default:
		factor = 0.16
	}

	carbon := req.DistanceKm * factor

	// Simpan log
	err = s.carbonRepo.CreateVehicleLog(ctx, &models.CarbonVehicleLog{
		VehicleID:       vehicle.ID,
		StartLat:        req.StartLat,
		StartLon:        req.StartLon,
		EndLat:          req.EndLat,
		EndLon:          req.EndLon,
		DistanceKm:      req.DistanceKm,
		DurationMinutes: req.DurationMinutes,
		CarbonEmission:  carbon,
	})

	if err != nil {
		return err
	}

	return s.missionRepo.CheckAllUserMissions(ctx, userID)
}

func (s *CarbonService) GetVehicleLogs(ctx context.Context, userID, vehicleID int64) ([]*models.CarbonVehicleLog, error) {
	vehicle, err := s.carbonRepo.FindVehicleByID(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, errors.New("vehicle not found")
	}
	if vehicle.UserID != userID {
		return nil, errors.New("vehicle does not belong to user")
	}

	return s.carbonRepo.GetVehicleLogs(ctx, vehicleID)
}


func (s *CarbonService) CreateElectronic(ctx context.Context, userID int64, req *dto.CreateElectronicDTO) (*models.CarbonElectronic, error) {
	existing, err := s.carbonRepo.FindElectronicsByUserAndName(ctx, userID, req.DeviceName)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("electronic device with this name already exists for this user")
	}

	electronic := &models.CarbonElectronic{
		UserID:     userID,
		DeviceName: req.DeviceName,
		DeviceType: req.DeviceType,
		PowerWatts: req.PowerWatts,
	}

	return s.carbonRepo.CreateElectronics(ctx, electronic)
}

func (s *CarbonService) ListUserElectronics(ctx context.Context, userID int64) ([]*models.CarbonElectronic, error) {
	return s.carbonRepo.ListUserElectronics(ctx, userID)
}

func (s *CarbonService) AddElectronicsLog(ctx context.Context, userID int64, req *dto.AddElectronicsLogDTO) error {
	var device *models.CarbonElectronic
	var err error

	// Cari atau buat device
	if req.DeviceID != nil {
		device, err = s.carbonRepo.FindElectronicsByID(ctx, *req.DeviceID)
		if err != nil {
			return err
		}
		if device == nil {
			return errors.New("electronic device not found")
		}
		if device.UserID != userID {
			return errors.New("electronic device does not belong to user")
		}
	} else {
		existing, err := s.carbonRepo.FindElectronicsByUserAndName(ctx, userID, req.DeviceName)
		if err != nil {
			return err
		}
		if existing != nil {
			device = existing
		} else {
			device, err = s.carbonRepo.CreateElectronics(ctx, &models.CarbonElectronic{
				UserID:     userID,
				DeviceName: req.DeviceName,
				DeviceType: req.DeviceType,
				PowerWatts: req.PowerWatts,
			})
			if err != nil {
				return err
			}
		}
	}

	carbon := float64(device.PowerWatts) / 1000.0 * req.DurationHours * 0.475

	err = s.carbonRepo.CreateElectronicsLog(ctx, &models.CarbonElectronicLog{
		DeviceID:       device.ID,
		DurationHours:  req.DurationHours,
		CarbonEmission: carbon,
		LoggedAt:       time.Now(),
	})
	if err != nil {
		return err
	}

	return s.missionRepo.CheckAllUserMissions(ctx, userID)
}

func (s *CarbonService) GetElectronicsLogs(ctx context.Context, userID, deviceID int64) ([]*models.CarbonElectronicLog, error) {
	device, err := s.carbonRepo.FindElectronicsByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, errors.New("electronic device not found")
	}
	if device.UserID != userID {
		return nil, errors.New("electronic device does not belong to user")
	}

	return s.carbonRepo.GetElectronicsLogs(ctx, deviceID)
}

func (s *CarbonService) EditVehicle(ctx context.Context, userID, vehicleID int64, req *dto.EditVehicleDTO) (*models.CarbonVehicle, error) {
	vehicle, err := s.carbonRepo.FindVehicleByID(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, errors.New("vehicle not found")
	}
	if vehicle.UserID != userID {
		return nil, errors.New("vehicle does not belong to user")
	}

	if req.Name != "" && req.Name != vehicle.Name {
		existing, err := s.carbonRepo.FindVehicleByUserAndName(ctx, userID, req.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("vehicle with this name already exists for this user")
		}
		vehicle.Name = req.Name
	}

	if req.VehicleType != "" {
		vehicle.VehicleType = models.VehicleType(req.VehicleType)
	}
	if req.FuelType != "" {
		vehicle.FuelType = models.FuelType(req.FuelType)
	}

	err = s.carbonRepo.UpdateVehicle(ctx, vehicle)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}

func (s *CarbonService) DeleteVehicle(ctx context.Context, userID, vehicleID int64) error {
	vehicle, err := s.carbonRepo.FindVehicleByID(ctx, vehicleID)
	if err != nil {
		return err
	}
	if vehicle == nil {
		return errors.New("vehicle not found")
	}
	if vehicle.UserID != userID {
		return errors.New("vehicle does not belong to user")
	}

	// Delete associated logs first
	err = s.carbonRepo.DeleteVehicleLogs(ctx, vehicleID)
	if err != nil {
		return err
	}

	// Delete the vehicle
	return s.carbonRepo.DeleteVehicle(ctx, vehicleID)
}

func (s *CarbonService) GetAllVehicleLogs(ctx context.Context, userID int64) ([]*models.CarbonVehicleLog, error) {
	return s.carbonRepo.GetAllVehicleLogsByUser(ctx, userID)
}
func (s *CarbonService) GetVehicleLogByID(ctx context.Context, userID, logID int64) (*models.CarbonVehicleLog, error) {
    return s.carbonRepo.GetVehicleLogByID(ctx, userID, logID)
}


func (s *CarbonService) EditElectronic(ctx context.Context, userID, deviceID int64, req *dto.EditElectronicDTO) (*models.CarbonElectronic, error) {
	device, err := s.carbonRepo.FindElectronicsByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, errors.New("electronic device not found")
	}
	if device.UserID != userID {
		return nil, errors.New("electronic device does not belong to user")
	}

	// Check if device name is being changed and if it conflicts with another device
	if req.DeviceName != "" && req.DeviceName != device.DeviceName {
		existing, err := s.carbonRepo.FindElectronicsByUserAndName(ctx, userID, req.DeviceName)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("electronic device with this name already exists for this user")
		}
		device.DeviceName = req.DeviceName
	}

	if req.DeviceType != "" {
		device.DeviceType = req.DeviceType
	}
	if req.PowerWatts != 0 {
		device.PowerWatts = req.PowerWatts
	}

	err = s.carbonRepo.UpdateElectronic(ctx, device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (s *CarbonService) DeleteElectronic(ctx context.Context, userID, deviceID int64) error {
	device, err := s.carbonRepo.FindElectronicsByID(ctx, deviceID)
	if err != nil {
		return err
	}
	if device == nil {
		return errors.New("electronic device not found")
	}
	if device.UserID != userID {
		return errors.New("electronic device does not belong to user")
	}

	// Delete associated logs first
	err = s.carbonRepo.DeleteElectronicLogs(ctx, deviceID)
	if err != nil {
		return err
	}

	// Delete the device
	return s.carbonRepo.DeleteElectronic(ctx, deviceID)
}

func (s *CarbonService) GetAllElectronicLogs(ctx context.Context, userID int64) ([]*models.CarbonElectronicLog, error) {
	return s.carbonRepo.GetAllElectronicLogsByUser(ctx, userID)
}
