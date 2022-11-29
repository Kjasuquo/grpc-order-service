package services

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"order_svc/internal/helper"
	"order_svc/internal/models"
	"order_svc/internal/rabbitMQ"
	"order_svc/proto"
	"strings"
	"time"
	//"gitlab.com/grpc-buffer/proto/go/pkg/proto"
)

func (s *OrderServiceServer) CreateNewOrder(ctx context.Context, req *proto.NewOrderRequest) (*proto.NewOrderResponse, error) {
	orderId := req.GetOrderId()
	reference := req.GetReference()
	subtotal := req.GetSubtotal()
	deliveryCost := req.GetDeliveryCost()
	totalCost := req.GetTotalCost()
	deliveryType := req.GetDelivery_Type()
	userId := req.GetUserId()
	deliveryAddress := req.GetDeliveryAddress().GetAddress()
	latitude := req.GetDeliveryAddress().GetCoordinates().GetLatitude()
	longitude := req.GetDeliveryAddress().GetCoordinates().GetLongitude()
	orderItems := req.GetOrderItems()
	paymentMethod := req.GetPaymentMethod()
	currency := req.GetCurrency()

	if orderId == "" ||
		reference == "" ||
		subtotal == 0 ||
		deliveryCost == 0 ||
		totalCost == 0 ||
		deliveryType == "" ||
		userId == "" ||
		deliveryAddress == "" ||
		latitude == 0 ||
		longitude == 0 ||
		orderItems == nil {
		return nil, status.Error(codes.InvalidArgument,
			"Cannot create order with empty fields.")
	}

	var cartItems []models.CartItem
	var cartPackageItems []models.CartPackageItem

	for _, orderItem := range orderItems {
		for _, item := range orderItem.Items {
			cartPackageItems = append(cartPackageItems, models.CartPackageItem{
				ItemId:       item.ItemID,
				Name:         item.Name,
				Description:  item.Description,
				ItemCategory: item.ItemCategoryID,
				Image:        item.Image,
				Unit:         item.Unit,
				UserId:       orderItem.UserID,
			})
		}
		cartI := models.CartItem{
			PackageId:       orderItem.PackageID,
			PackageName:     orderItem.PackageName,
			Description:     orderItem.Description,
			BasePrice:       orderItem.BasePrice,
			ServiceId:       orderItem.ServiceAreaID,
			Image:           orderItem.Image,
			UserId:          orderItem.UserID,
			Quantity:        orderItem.Quantity,
			CartPackageItem: cartPackageItems,
			CartId:          orderItem.CartID,
		}

		in := &proto.CartRequestItem{
			RequestId: orderItem.CartID,
			PackageID: orderItem.PackageID,
			UserID:    orderItem.UserID,
		}

		_, err := s.DeleteFromCart(in)
		if err != nil {
			//return nil, status.Errorf(codes.Internal, fmt.Sprintf("Error Deleting from cart: %v\n", err))
			fmt.Println("Error Deleting from cart: ", err)
		}

		cartItems = append(cartItems, cartI)
	}

	deliverCode := helper.DeliveryCode()
	orderAcceptedTime := time.Now().String()

	order := &models.NewOrder{
		OrderID:      orderId,
		Reference:    reference,
		SubTotal:     subtotal,
		DeliveryCost: deliveryCost,
		TotalCost:    totalCost,
		DeliveryType: deliveryType,
		UserID:       userId,
		OrderDeliveryAddress: models.OrderDeliveryAddress{
			Address: deliveryAddress,
			OrderCoordinates: models.OrderCoordinates{
				Latitude:  latitude,
				Longitude: longitude,
			},
		},
		CartItems:         cartItems,
		Status:            "Order Accepted",
		DeliveryCode:      deliverCode,
		PaymentMethod:     paymentMethod,
		Currency:          currency,
		DeliveryTime:      "not delivered",
		OrderAcceptedTime: orderAcceptedTime,
	}

	err := s.Order.CreateNewOrder(order)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Error creating order: %v\n", err))
	}

	err = rabbitMQ.PublishToOrderNotificationQueue(order)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.NewOrderResponse{
		Message: "successful",
		Status:  codes.OK.String(),
	}, nil
}

func (s *OrderServiceServer) FetchOrdersByUser(req *proto.NewGetUserOrderRequest, stream proto.Order_FetchOrdersByUserServer) error {
	userId := strings.TrimSpace(req.GetUserId())
	if userId == "" {
		return status.Error(codes.InvalidArgument,
			"Cannot fetch orders with empty user id.")
	}

	var orderItems []*proto.OrderProperties
	var cartItems []*proto.InNewOrderItems
	var packageItem []*proto.InNewOrderPackageItem

	orders, err := s.Order.FetchOrdersByUser(userId)
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Error fetching orders: %v\n", err))
	}
	for _, order := range orders {
		for _, cartItem := range order.CartItems {
			for _, item := range cartItem.CartPackageItem {
				PackageI := &proto.InNewOrderPackageItem{
					ItemID:         item.ItemId,
					Name:           item.Name,
					Description:    item.Description,
					ItemCategoryID: item.ItemCategory,
					Image:          item.Image,
					Unit:           item.Unit,
				}
				packageItem = append(packageItem, PackageI)
			}
			cartI := &proto.InNewOrderItems{
				PackageID:     cartItem.PackageId,
				PackageName:   cartItem.PackageName,
				Description:   cartItem.Description,
				BasePrice:     cartItem.BasePrice,
				ServiceAreaID: cartItem.ServiceId,
				Image:         cartItem.Image,
				UserID:        cartItem.UserId,
				Quantity:      cartItem.Quantity,
				Items:         packageItem,
				CartID:        cartItem.CartId,
			}
			cartItems = append(cartItems, cartI)
		}
		response := &proto.OrderProperties{
			OrderId:       order.OrderID,
			Reference:     order.Reference,
			Subtotal:      order.SubTotal,
			DeliveryCost:  order.DeliveryCost,
			TotalCost:     order.TotalCost,
			Delivery_Type: order.DeliveryType,
			UserId:        order.UserID,
			DeliveryAddress: &proto.InNewOrderDeliveryAddress{
				Address: order.OrderDeliveryAddress.Address,
				Coordinates: &proto.InNewOrderCoordinates{
					Latitude:  order.OrderDeliveryAddress.OrderCoordinates.Latitude,
					Longitude: order.OrderDeliveryAddress.OrderCoordinates.Longitude,
				},
			},
			OrderItems:            cartItems,
			Status:                order.Status,
			DeliveryCode:          order.DeliveryCode,
			PaymentMethod:         order.PaymentMethod,
			Currency:              order.Currency,
			DeliveryTime:          order.DeliveryTime,
			CreatedAt:             order.CreatedAt.String(),
			UpdatedAt:             order.UpdatedAt.String(),
			AcceptanceTime:        order.OrderAcceptedTime,
			ShopperAssignedTime:   order.ShopperAssignedTime,
			ShoppingCompletedTime: order.ShoppingCompletedTime,
			InProgressTime:        order.InProgressTime,
		}

		orderItems = append(orderItems, response)
	}

	for len(orderItems) > 0 {
		// send order items in stream
		for _, orderItem := range orderItems {
			if err := stream.Send(&proto.NewGetUserOrderResponse{
				Data:   orderItem,
				Status: codes.OK.String(),
			}); err != nil {
				return status.Errorf(codes.Internal, fmt.Sprintf("Error sending order items: %v\n", err))
			}
		}
		break
	}

	return nil
}

func (s *OrderServiceServer) DeleteUserOrders(ctx context.Context, req *proto.DeleteUserOrdersRequest) (*proto.DeleteUserOrdersResponse, error) {
	userId := req.GetUserId()
	if userId == "" {
		return nil, status.Error(codes.InvalidArgument,
			"user ID is required")
	}

	err := s.Order.DeleteOrderByUserID(userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Error deleting orders: %v\n", err))
	}

	return &proto.DeleteUserOrdersResponse{
		Message: "successful",
		Status:  codes.OK.String(),
	}, nil
}

func (s *OrderServiceServer) UpdateOrderStatus(ctx context.Context, req *proto.UpdateOrderStatusRequest) (*proto.UpdateOrderStatusResponse, error) {
	orderId := req.GetOrderId()
	sta := req.GetStatus()

	if orderId == "" {
		return nil, status.Error(codes.InvalidArgument,
			"order ID is required")
	}

	if sta == "Shopper Assigned" {
		shopperAssignedTime := time.Now().String()
		_, err := s.Order.UpdateShopperTime(orderId, shopperAssignedTime)
		if err != nil {
			return nil, status.Error(codes.Internal, "Error updating shopper assigned time")
		}
	} else if sta == "Shopping Completed" {
		shoppingCompletedTime := time.Now().String()
		_, err := s.Order.UpdateShoppingCompletedTime(orderId, shoppingCompletedTime)
		if err != nil {
			return nil, status.Error(codes.Internal, "Error updating shopping completed time")
		}
	} else if sta == "Delivery In Progress" {
		deliveryInProgressTime := time.Now().String()
		_, err := s.Order.UpdateDeliveryInProgressTime(orderId, deliveryInProgressTime)
		if err != nil {
			return nil, status.Error(codes.Internal, "Error updating delivery in progress time")
		}
	} else if sta == "Delivered" {
		deliveryTime := time.Now().String()
		_, err := s.Order.UpdateDeliveryTime(orderId, deliveryTime)
		if err != nil {
			return nil, status.Error(codes.Internal, "Error updating delivery time")
		}
	} else {
		return nil, status.Error(codes.InvalidArgument,
			"status is invalid")
	}

	_, err := s.Order.UpdateOrder(orderId, sta)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error updating order status")
	}

	return &proto.UpdateOrderStatusResponse{
		Message: "successfully  updated order status",
		Status:  codes.OK.String(),
	}, nil
}

func (s *OrderServiceServer) FetchAllOrders(req *emptypb.Empty, stream proto.Order_FetchAllOrdersServer) error {
	var orderItems []*proto.OrderProperties
	var cartItems []*proto.InNewOrderItems
	var packageItem []*proto.InNewOrderPackageItem

	orders, err := s.Order.FetchAllOrders()
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Error fetching orders: %v\n", err))
	}
	for _, order := range orders {
		for _, cartItem := range order.CartItems {
			for _, item := range cartItem.CartPackageItem {
				PackageI := &proto.InNewOrderPackageItem{
					ItemID:         item.ItemId,
					Name:           item.Name,
					Description:    item.Description,
					ItemCategoryID: item.ItemCategory,
					Image:          item.Image,
					Unit:           item.Unit,
				}
				packageItem = append(packageItem, PackageI)
			}
			cartI := &proto.InNewOrderItems{
				PackageID:     cartItem.PackageId,
				PackageName:   cartItem.PackageName,
				Description:   cartItem.Description,
				BasePrice:     cartItem.BasePrice,
				ServiceAreaID: cartItem.ServiceId,
				Image:         cartItem.Image,
				UserID:        cartItem.UserId,
				Quantity:      cartItem.Quantity,
				Items:         packageItem,
				CartID:        cartItem.CartId,
			}
			cartItems = append(cartItems, cartI)
		}
		response := &proto.OrderProperties{
			OrderId:       order.OrderID,
			Reference:     order.Reference,
			Subtotal:      order.SubTotal,
			DeliveryCost:  order.DeliveryCost,
			TotalCost:     order.TotalCost,
			Delivery_Type: order.DeliveryType,
			UserId:        order.UserID,
			DeliveryAddress: &proto.InNewOrderDeliveryAddress{
				Address: order.OrderDeliveryAddress.Address,
				Coordinates: &proto.InNewOrderCoordinates{
					Latitude:  order.OrderDeliveryAddress.OrderCoordinates.Latitude,
					Longitude: order.OrderDeliveryAddress.OrderCoordinates.Longitude,
				},
			},
			OrderItems:            cartItems,
			Status:                order.Status,
			DeliveryCode:          order.DeliveryCode,
			PaymentMethod:         order.PaymentMethod,
			Currency:              order.Currency,
			DeliveryTime:          order.DeliveryTime,
			CreatedAt:             order.CreatedAt.String(),
			UpdatedAt:             order.UpdatedAt.String(),
			AcceptanceTime:        order.OrderAcceptedTime,
			ShopperAssignedTime:   order.ShopperAssignedTime,
			ShoppingCompletedTime: order.ShoppingCompletedTime,
			InProgressTime:        order.InProgressTime,
		}

		orderItems = append(orderItems, response)
	}

	for len(orderItems) > 0 {
		// send order items in stream
		for _, orderItem := range orderItems {
			if err := stream.Send(&proto.AllOrdersResponse{
				Data:   orderItem,
				Status: codes.OK.String(),
			}); err != nil {
				return status.Errorf(codes.Internal, fmt.Sprintf("Error sending order items: %v\n", err))
			}
		}
		break
	}

	return nil
}

func (s *OrderServiceServer) FetchOderByStatus(req *proto.GetStatusRequest, stream proto.Order_FetchOderByStatusServer) error {
	userId := req.GetUserId()
	sta := req.GetStatus()

	var orderItems []*proto.OrderProperties
	var cartItems []*proto.InNewOrderItems
	var packageItem []*proto.InNewOrderPackageItem

	orders, err := s.Order.FetchUserOrderByStatus(userId, sta)
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Error fetching orders: %v\n", err))
	}
	for _, order := range orders {
		for _, cartItem := range order.CartItems {
			for _, item := range cartItem.CartPackageItem {
				PackageI := &proto.InNewOrderPackageItem{
					ItemID:         item.ItemId,
					Name:           item.Name,
					Description:    item.Description,
					ItemCategoryID: item.ItemCategory,
					Image:          item.Image,
					Unit:           item.Unit,
				}
				packageItem = append(packageItem, PackageI)
			}
			cartI := &proto.InNewOrderItems{
				PackageID:     cartItem.PackageId,
				PackageName:   cartItem.PackageName,
				Description:   cartItem.Description,
				BasePrice:     cartItem.BasePrice,
				ServiceAreaID: cartItem.ServiceId,
				Image:         cartItem.Image,
				UserID:        cartItem.UserId,
				Quantity:      cartItem.Quantity,
				Items:         packageItem,
				CartID:        cartItem.CartId,
			}
			cartItems = append(cartItems, cartI)
		}
		response := &proto.OrderProperties{
			OrderId:       order.OrderID,
			Reference:     order.Reference,
			Subtotal:      order.SubTotal,
			DeliveryCost:  order.DeliveryCost,
			TotalCost:     order.TotalCost,
			Delivery_Type: order.DeliveryType,
			UserId:        order.UserID,
			DeliveryAddress: &proto.InNewOrderDeliveryAddress{
				Address: order.OrderDeliveryAddress.Address,
				Coordinates: &proto.InNewOrderCoordinates{
					Latitude:  order.OrderDeliveryAddress.OrderCoordinates.Latitude,
					Longitude: order.OrderDeliveryAddress.OrderCoordinates.Longitude,
				},
			},
			OrderItems:            cartItems,
			Status:                order.Status,
			DeliveryCode:          order.DeliveryCode,
			PaymentMethod:         order.PaymentMethod,
			Currency:              order.Currency,
			DeliveryTime:          order.DeliveryTime,
			CreatedAt:             order.CreatedAt.String(),
			UpdatedAt:             order.UpdatedAt.String(),
			AcceptanceTime:        order.OrderAcceptedTime,
			ShopperAssignedTime:   order.ShopperAssignedTime,
			ShoppingCompletedTime: order.ShoppingCompletedTime,
			InProgressTime:        order.InProgressTime,
		}

		orderItems = append(orderItems, response)
	}

	for len(orderItems) > 0 {
		// send order items in stream
		for _, orderItem := range orderItems {
			if err := stream.Send(&proto.GetStatusResponse{
				Data:   orderItem,
				Status: codes.OK.String(),
			}); err != nil {
				return status.Errorf(codes.Internal, fmt.Sprintf("Error sending order items: %v\n", err))
			}
		}
		break
	}

	return nil
}

func (s *OrderServiceServer) GetDeliveryCode(ctx context.Context, req *proto.GetDeliveryCodeRequest) (*proto.GetDeliveryCodeResponse, error) {
	orderId := req.GetOrderId()
	deliveryCode, err := s.Order.GetDeliveryCode(orderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Error fetching delivery code: %v\n", err))
	}

	return &proto.GetDeliveryCodeResponse{
		DeliveryCode: deliveryCode,
		Status:       codes.OK.String(),
	}, nil
}
