// service/store.go
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	dto "github.com/Qodarrz/fiber-app/dto"
	model "github.com/Qodarrz/fiber-app/model"
	"github.com/Qodarrz/fiber-app/repository"
)

type StoreServiceInterface interface {
	// Store Items
	GetAllStoreItems(ctx context.Context, status string) ([]dto.StoreItemDTO, error)
	GetStoreItemByID(ctx context.Context, id int64) (*dto.StoreItemDTO, error)
	CreateStoreItem(ctx context.Context, req *dto.CreateStoreItemDTO) (*dto.StoreItemDTO, error)
	UpdateStoreItem(ctx context.Context, id int64, req *dto.UpdateStoreItemDTO) (*dto.StoreItemDTO, error)
	DeleteStoreItem(ctx context.Context, id int64) error

	// Orders
	CreateOrder(ctx context.Context, userID int64, qty int64, req *dto.CreateOrderDTO) (*dto.OrderResponseDTO, error)
	GetOrderByID(ctx context.Context, orderID int64) (*dto.OrderDTO, error)
	GetUserOrders(ctx context.Context, userID int64) ([]dto.OrderDTO, error)
	CancelOrder(ctx context.Context, userID, orderID int64) error
	CreateOrderByItemID(ctx context.Context, userID, itemID int64, qty int) (*dto.OrderResponseDTO, error)
}

type storeService struct {
	storeRepo    repository.StoreRepositoryInterface
	pointsRepo   repository.PointsRepositoryInterface
	activityRepo repository.ActivityRepositoryInterface
	notificationRepo repository.NotificationRepository
}

func NewStoreService(
	storeRepo repository.StoreRepositoryInterface,
	pointsRepo repository.PointsRepositoryInterface,
	activityRepo repository.ActivityRepositoryInterface,
	notificationRepo repository.NotificationRepository,
) StoreServiceInterface {
	return &storeService{
		storeRepo:    storeRepo,
		pointsRepo:   pointsRepo,
		activityRepo: activityRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *storeService) GetAllStoreItems(ctx context.Context, status string) ([]dto.StoreItemDTO, error) {
	items, err := s.storeRepo.GetAllStoreItems(ctx, status)
	if err != nil {
		return nil, err
	}

	var result []dto.StoreItemDTO
	for _, item := range items {
		result = append(result, dto.StoreItemDTO{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			PricePoints: item.PricePoints,
			Stock:       item.Stock,
			Status:      item.Status,
			ImageURL:    item.ImageURL,
			CreatedAt:   item.CreatedAt,
		})
	}
	return result, nil
}

func (s *storeService) GetStoreItemByID(ctx context.Context, id int64) (*dto.StoreItemDTO, error) {
	item, err := s.storeRepo.GetStoreItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("store item not found")
	}

	return &dto.StoreItemDTO{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		PricePoints: item.PricePoints,
		Stock:       item.Stock,
		Status:      item.Status,
		ImageURL:    item.ImageURL,
		CreatedAt:   item.CreatedAt,
	}, nil
}

func (s *storeService) CreateStoreItem(ctx context.Context, req *dto.CreateStoreItemDTO) (*dto.StoreItemDTO, error) {
	item := &model.StoreItem{
		Name:        req.Name,
		Description: req.Description,
		PricePoints: req.PricePoints,
		Stock:       req.Stock,
		Status:      "active",
		ImageURL:    req.ImageURL,
		CreatedAt:   time.Now(),
	}

	if err := s.storeRepo.CreateStoreItem(ctx, item); err != nil {
		return nil, err
	}

	return &dto.StoreItemDTO{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		PricePoints: item.PricePoints,
		Stock:       item.Stock,
		Status:      item.Status,
		ImageURL:    item.ImageURL,
		CreatedAt:   item.CreatedAt,
	}, nil
}

func (s *storeService) UpdateStoreItem(ctx context.Context, id int64, req *dto.UpdateStoreItemDTO) (*dto.StoreItemDTO, error) {
	item, err := s.storeRepo.GetStoreItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("store item not found")
	}

	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.PricePoints > 0 {
		item.PricePoints = req.PricePoints
	}
	if req.Stock >= 0 {
		item.Stock = req.Stock
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	if req.ImageURL != "" {
		item.ImageURL = req.ImageURL
	}

	if err := s.storeRepo.UpdateStoreItem(ctx, item); err != nil {
		return nil, err
	}

	return &dto.StoreItemDTO{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		PricePoints: item.PricePoints,
		Stock:       item.Stock,
		Status:      item.Status,
		ImageURL:    item.ImageURL,
		CreatedAt:   item.CreatedAt,
	}, nil
}

func (s *storeService) DeleteStoreItem(ctx context.Context, id int64) error {
	item, err := s.storeRepo.GetStoreItemByID(ctx, id)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("store item not found")
	}

	item.Status = "inactive"
	return s.storeRepo.UpdateStoreItem(ctx, item)
}

func (s *storeService) CreateOrder(ctx context.Context, userID int64, qty int64, req *dto.CreateOrderDTO) (*dto.OrderResponseDTO, error) {
	var totalPoints int
	var orderItems []model.OrderItem
	var storeItems []*model.StoreItem

	for _, itemReq := range req.Items {
		// Get store item details
		storeItem, err := s.storeRepo.GetStoreItemByID(ctx, int64(itemReq.ItemID))
		if err != nil {
			return nil, err
		}
		if storeItem == nil {
			return nil, fmt.Errorf("item with ID %d not found", itemReq.ItemID)
		}
		if storeItem.Status != "active" {
			return nil, fmt.Errorf("item %s is not available", storeItem.Name)
		}
		if storeItem.Stock < itemReq.Qty {
			return nil, fmt.Errorf("insufficient stock for item %s", storeItem.Name)
		}

		itemTotal := storeItem.PricePoints * itemReq.Qty
		totalPoints += itemTotal

		orderItems = append(orderItems, model.OrderItem{
			ItemID:          int64(itemReq.ItemID),
			Qty:             itemReq.Qty,
			PriceEachPoints: storeItem.PricePoints,
			CreatedAt:       time.Now(),
		})
		storeItems = append(storeItems, storeItem)
	}

	// Check user points
	userPoints, err := s.pointsRepo.GetUserPoints(ctx, userID)
	if err != nil {
		return nil, err
	}
	if userPoints.TotalPoints < totalPoints {
		return nil, errors.New("insufficient points")
	}

	// Start transaction
	tx, err := s.storeRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Create order
	order := &model.Order{
		UserID:      userID,
		TotalPoints: totalPoints,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := s.storeRepo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	// Create order items
	if err := s.storeRepo.CreateOrderItems(ctx, order.ID, orderItems); err != nil {
		return nil, err
	}

	// Deduct points
	if err := s.pointsRepo.DeductPoints(ctx, userID, totalPoints, "store_purchase", order.ID); err != nil {
		return nil, err
	}

	// Update stock for each item
	for i, storeItem := range storeItems {
		if err := s.storeRepo.DecrementStoreItemStock(ctx, storeItem.ID, req.Items[i].Qty); err != nil {
			return nil, err
		}
	}

	// Update order status to completed
	order.Status = "completed"
	if err := s.storeRepo.UpdateOrderStatus(ctx, order.ID, "completed"); err != nil {
		return nil, err
	}

	// Log activity
	activityMsg := fmt.Sprintf("User %d purchased items for %d points", userID, totalPoints)
	if err := s.activityRepo.LogActivity(ctx, userID, activityMsg); err != nil {
		fmt.Printf("Failed to log activity: %v\n", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Get remaining points
	remainingPoints, err := s.pointsRepo.GetUserPoints(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get order details with items
	orderItemsResp, err := s.storeRepo.GetOrderItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	var itemsResponse []dto.OrderItemResponseDTO
	for _, item := range orderItemsResp {
		itemsResponse = append(itemsResponse, dto.OrderItemResponseDTO{
			ID:              item.ID,
			ItemID:          item.ItemID,
			Qty:             item.Qty,
			PriceEachPoints: item.PriceEachPoints,
			TotalPoints:     item.Qty * item.PriceEachPoints,
		})
	}

	orderDTO := dto.OrderDTO{
		ID:          order.ID,
		UserID:      order.UserID,
		TotalPoints: order.TotalPoints,
		Status:      order.Status,
		Items:       itemsResponse,
		CreatedAt:   order.CreatedAt,
	}

	return &dto.OrderResponseDTO{
		Order:           orderDTO,
		RemainingPoints: remainingPoints.TotalPoints,
	}, nil
}

func (s *storeService) GetOrderByID(ctx context.Context, orderID int64) (*dto.OrderDTO, error) {
	order, err := s.storeRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	items, err := s.storeRepo.GetOrderItems(ctx, orderID)
	if err != nil {
		return nil, err
	}

	var itemsResponse []dto.OrderItemResponseDTO
	for _, item := range items {
		itemsResponse = append(itemsResponse, dto.OrderItemResponseDTO{
			ID:              item.ID,
			ItemID:          item.ItemID,
			Qty:             item.Qty,
			PriceEachPoints: item.PriceEachPoints,
			TotalPoints:     item.Qty * item.PriceEachPoints,
		})
	}

	return &dto.OrderDTO{
		ID:          order.ID,
		UserID:      order.UserID,
		TotalPoints: order.TotalPoints,
		Status:      order.Status,
		Items:       itemsResponse,
		CreatedAt:   order.CreatedAt,
	}, nil
}

func (s *storeService) GetUserOrders(ctx context.Context, userID int64) ([]dto.OrderDTO, error) {
	orders, err := s.storeRepo.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []dto.OrderDTO
	for _, order := range orders {
		items, err := s.storeRepo.GetOrderItems(ctx, order.ID)
		if err != nil {
			return nil, err
		}

		var itemsResponse []dto.OrderItemResponseDTO
		for _, item := range items {
			itemsResponse = append(itemsResponse, dto.OrderItemResponseDTO{
				ID:              item.ID,
				ItemID:          item.ItemID,
				Qty:             item.Qty,
				PriceEachPoints: item.PriceEachPoints,
				TotalPoints:     item.Qty * item.PriceEachPoints,
			})
		}

		result = append(result, dto.OrderDTO{
			ID:          order.ID,
			UserID:      order.UserID,
			TotalPoints: order.TotalPoints,
			Status:      order.Status,
			Items:       itemsResponse,
			CreatedAt:   order.CreatedAt,
		})
	}

	return result, nil
}

func (s *storeService) CancelOrder(ctx context.Context, userID, orderID int64) error {
	order, err := s.storeRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}
	if order.UserID != userID {
		return errors.New("unauthorized to cancel this order")
	}
	if order.Status != "pending" {
		return errors.New("only pending orders can be cancelled")
	}

	// Start transaction
	tx, err := s.storeRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Refund points
	if err := s.pointsRepo.AddPoints(ctx, userID, order.TotalPoints, "order_refund", orderID); err != nil {
		return err
	}

	// Restore stock
	items, err := s.storeRepo.GetOrderItems(ctx, orderID)
	if err != nil {
		return err
	}

	for _, item := range items {
		if err := s.storeRepo.IncrementStoreItemStock(ctx, item.ItemID, item.Qty); err != nil {
			return err
		}
	}

	// Update order status
	if err := s.storeRepo.UpdateOrderStatus(ctx, orderID, "cancelled"); err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

func (s *storeService) CreateOrderByItemID(ctx context.Context, userID, itemID int64, qty int) (*dto.OrderResponseDTO, error) {
	if qty <= 0 {
		qty = 1
	}

	// Ambil data item
	storeItem, err := s.storeRepo.GetStoreItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if storeItem == nil {
		return nil, errors.New("item not found")
	}
	if storeItem.Status != "active" {
		return nil, errors.New("item is not available")
	}
	if storeItem.Stock < qty {
		return nil, errors.New("insufficient stock")
	}

	// Hitung total points
	totalPoints := storeItem.PricePoints * qty

	// Cek poin user
	userPoints, err := s.pointsRepo.GetUserPoints(ctx, userID)
	if err != nil {
		return nil, err
	}
	if userPoints.TotalPoints < totalPoints {
		return nil, errors.New("insufficient points")
	}

	// Start transaction
	tx, err := s.storeRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Create order
	order := &model.Order{
		UserID:      userID,
		TotalPoints: totalPoints,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}
	if err := s.storeRepo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	// Create order item
	orderItem := model.OrderItem{
		ItemID:          storeItem.ID,
		Qty:             qty,
		PriceEachPoints: storeItem.PricePoints,
		CreatedAt:       time.Now(),
	}
	if err := s.storeRepo.CreateOrderItems(ctx, order.ID, []model.OrderItem{orderItem}); err != nil {
		return nil, err
	}

	// Potong poin
	if err := s.pointsRepo.DeductPoints(ctx, userID, totalPoints, "store_purchase", order.ID); err != nil {
		return nil, err
	}

	// Kurangi stock
	if err := s.storeRepo.DecrementStoreItemStock(ctx, storeItem.ID, qty); err != nil {
		return nil, err
	}

	// Update order status jadi completed
	if err := s.storeRepo.UpdateOrderStatus(ctx, order.ID, "completed"); err != nil {
		return nil, err
	}
	order.Status = "completed"

	// Log aktivitas
	activityMsg := fmt.Sprintf("User %d purchased %s (x%d) for %d points", userID, storeItem.Name, qty, totalPoints)
	_ = s.activityRepo.LogActivity(ctx, userID, activityMsg)

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Setelah order berhasil dan commit tx
	notif := &model.Notification{
		UserID:    userID,
		Title:     "Order Successful",
		Message:   fmt.Sprintf("Your order for %s (x%d) has been completed. Total points spent: %d.", storeItem.Name, qty, totalPoints),
		CreatedAt: time.Now(),
	}
	_ = s.notificationRepo.Create(ctx, notif)

	// Sisa poin user
	remainingPoints, err := s.pointsRepo.GetUserPoints(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Response DTO
	itemResp := dto.OrderItemResponseDTO{
		ItemID:          storeItem.ID,
		Qty:             qty,
		PriceEachPoints: storeItem.PricePoints,
		TotalPoints:     totalPoints,
	}
	orderDTO := dto.OrderDTO{
		ID:          order.ID,
		UserID:      userID,
		TotalPoints: totalPoints,
		Status:      order.Status,
		Items:       []dto.OrderItemResponseDTO{itemResp},
		CreatedAt:   order.CreatedAt,
	}

	return &dto.OrderResponseDTO{
		Order:           orderDTO,
		RemainingPoints: remainingPoints.TotalPoints,
	}, nil
}
