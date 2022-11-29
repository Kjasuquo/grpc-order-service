package models

// NewOrder payload for creating order
type NewOrder struct {
	Models
	OrderID      string  `protobuf:"bytes,2,opt,name=orderID,proto3" json:"orderID,omitempty"`
	Reference    string  `protobuf:"bytes,2,opt,name=reference,proto3" json:"reference,omitempty"`
	SubTotal     float64 `protobuf:"bytes,2,opt,name=subTotal,proto3" json:"subTotal,omitempty"`
	DeliveryCost float64 `protobuf:"bytes,2,opt,name=deliveryCost,proto3" json:"deliveryCost,omitempty"`
	TotalCost    float64 `protobuf:"bytes,2,opt,name=totalCost,proto3" json:"totalCost,omitempty"`
	DeliveryType string  `protobuf:"bytes,2,opt,name=deliveryType,proto3" json:"deliveryType,omitempty"`
	UserID       string  `protobuf:"bytes,2,opt,name=userID,proto3" json:"userID,omitempty"`
	OrderDeliveryAddress
	CartItems             []CartItem `gorm:"foreignKey:CartId"`
	Status                string     `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	DeliveryCode          int32      `protobuf:"bytes,2,opt,name=deliveryCode,proto3" json:"deliveryCode,omitempty"`
	PaymentMethod         string     `protobuf:"bytes,2,opt,name=paymentMethod,proto3" json:"paymentMethod,omitempty"`
	Currency              string     `protobuf:"bytes,2,opt,name=currency,proto3" json:"currency,omitempty"`
	OrderAcceptedTime     string     `protobuf:"bytes,2,opt,name=orderAcceptedTime,proto3" json:"orderAcceptedTime,omitempty"`
	ShopperAssignedTime   string     `protobuf:"bytes,2,opt,name=shopperAssignedTime,proto3" json:"shopperAssignedTime,omitempty"`
	ShoppingCompletedTime string     `protobuf:"bytes,2,opt,name=shoppingCompletedTime,proto3" json:"shoppingCompletedTime,omitempty"`
	InProgressTime        string     `protobuf:"bytes,2,opt,name=inProgressTime,proto3" json:"inProgressTime,omitempty"`
	DeliveryTime          string     `protobuf:"bytes,2,opt,name=deliverTime,proto3" json:"deliverTime,omitempty"`
}

type CartItem struct {
	Models
	PackageId       string            `protobuf:"bytes,2,opt,name=packageId,proto3" json:"packageId,omitempty"`
	PackageName     string            `protobuf:"bytes,2,opt,name=packageName,proto3" json:"packageName,omitempty"`
	Description     string            `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	BasePrice       float64           `protobuf:"bytes,2,opt,name=basePrice,proto3" json:"basePrice,omitempty"`
	ServiceId       string            `protobuf:"bytes,2,opt,name=serviceId,proto3" json:"serviceId,omitempty"`
	Image           string            `protobuf:"bytes,2,opt,name=image,proto3" json:"image,omitempty"`
	UserId          string            `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
	Quantity        int32             `protobuf:"bytes,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	CartPackageItem []CartPackageItem `gorm:"foreignKey:ItemId" `
	CartId          string            `protobuf:"bytes,2,opt,name=cartId,proto3" json:"cartId,omitempty"`
}

type CartPackageItem struct {
	Models
	ItemId       string `protobuf:"bytes,2,opt,name=itemId,proto3" json:"itemId,omitempty"`
	Name         string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description  string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	ItemCategory string `protobuf:"bytes,2,opt,name=itemCategory,proto3" json:"itemCategory,omitempty"`
	Image        string `protobuf:"bytes,2,opt,name=image,proto3" json:"image,omitempty"`
	Unit         string `protobuf:"bytes,2,opt,name=unit,proto3" json:"unit,omitempty"`
	UserId       string `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
}

type OrderDeliveryAddress struct {
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	OrderCoordinates
}

type OrderCoordinates struct {
	Latitude  float64 `protobuf:"bytes,2,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude float64 `protobuf:"bytes,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
}
