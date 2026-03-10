package api

// ==================== User ====================

type UserRegisterRequest struct {
	Username string `json:"username" form:"username" query:"username" vd:"len($)>0"`
	Password string `json:"password" form:"password" query:"password" vd:"len($)>0"`
}

type UserRegisterResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	UserID     int64  `json:"user_id"`
	Token      string `json:"token"`
}

type UserLoginRequest struct {
	Username string `json:"username" form:"username" query:"username" vd:"len($)>0"`
	Password string `json:"password" form:"password" query:"password" vd:"len($)>0"`
}

type UserLoginResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	UserID     int64  `json:"user_id"`
	Token      string `json:"token"`
}

type UserRequest struct {
	UserID int64  `json:"user_id" form:"user_id" query:"user_id"`
	Token  string `json:"token" form:"token" query:"token"`
}

type UserResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	User       *User  `json:"user"`
}

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar,omitempty"`
	Signature string `json:"signature,omitempty"`
}

// ==================== Product ====================

type ProductCreateRequest struct {
	Token       string `json:"token" form:"token"`
	Name        string `json:"name" form:"name" vd:"len($)>0"`
	Description string `json:"description" form:"description"`
	Price       int64  `json:"price" form:"price"`
	Stock       int64  `json:"stock" form:"stock"`
	ImageUrl    string `json:"image_url" form:"image_url"`
	Category    string `json:"category" form:"category"`
}

type ProductCreateResponse struct {
	StatusCode int64    `json:"status_code"`
	StatusMsg  string   `json:"status_msg,omitempty"`
	Product    *Product `json:"product"`
}

type ProductDetailRequest struct {
	ProductID int64   `json:"product_id" form:"product_id" query:"product_id"`
	Token     *string `json:"token,omitempty" form:"token" query:"token"`
}

func (r *ProductDetailRequest) GetToken() string {
	if r.Token != nil {
		return *r.Token
	}
	return ""
}

type ProductDetailResponse struct {
	StatusCode int64    `json:"status_code"`
	StatusMsg  string   `json:"status_msg,omitempty"`
	Product    *Product `json:"product"`
}

type ProductListRequest struct {
	Token    *string `json:"token,omitempty" form:"token" query:"token"`
	Page     *int64  `json:"page,omitempty" form:"page" query:"page"`
	Size     *int64  `json:"size,omitempty" form:"size" query:"size"`
	Category *string `json:"category,omitempty" form:"category" query:"category"`
}

func (r *ProductListRequest) GetToken() string {
	if r.Token != nil {
		return *r.Token
	}
	return ""
}

func (r *ProductListRequest) GetPage() int64 {
	if r.Page != nil {
		return *r.Page
	}
	return 1
}

func (r *ProductListRequest) GetSize() int64 {
	if r.Size != nil {
		return *r.Size
	}
	return 10
}

func (r *ProductListRequest) GetCategory() string {
	if r.Category != nil {
		return *r.Category
	}
	return ""
}

type ProductListResponse struct {
	StatusCode  int64      `json:"status_code"`
	StatusMsg   string     `json:"status_msg,omitempty"`
	ProductList []*Product `json:"product_list"`
	Total       int64      `json:"total"`
}

type Product struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Stock       int64  `json:"stock"`
	ImageUrl    string `json:"image_url"`
	Category    string `json:"category"`
}

// ==================== Seckill ====================

type SeckillCreateRequest struct {
	Token        string `json:"token" form:"token"`
	ProductID    int64  `json:"product_id" form:"product_id"`
	SeckillPrice int64  `json:"seckill_price" form:"seckill_price"`
	TotalStock   int64  `json:"total_stock" form:"total_stock"`
	StartTime    string `json:"start_time" form:"start_time"`
	EndTime      string `json:"end_time" form:"end_time"`
}

type SeckillCreateResponse struct {
	StatusCode int64            `json:"status_code"`
	StatusMsg  string           `json:"status_msg,omitempty"`
	Activity   *SeckillActivity `json:"activity"`
}

type SeckillListRequest struct {
	Token  *string `json:"token,omitempty" form:"token" query:"token"`
	Status *int64  `json:"status,omitempty" form:"status" query:"status"`
	Page   *int64  `json:"page,omitempty" form:"page" query:"page"`
	Size   *int64  `json:"size,omitempty" form:"size" query:"size"`
}

func (r *SeckillListRequest) GetToken() string {
	if r.Token != nil {
		return *r.Token
	}
	return ""
}

func (r *SeckillListRequest) GetStatus() int64 {
	if r.Status != nil {
		return *r.Status
	}
	return -1
}

func (r *SeckillListRequest) GetPage() int64 {
	if r.Page != nil {
		return *r.Page
	}
	return 1
}

func (r *SeckillListRequest) GetSize() int64 {
	if r.Size != nil {
		return *r.Size
	}
	return 10
}

type SeckillListResponse struct {
	StatusCode   int64              `json:"status_code"`
	StatusMsg    string             `json:"status_msg,omitempty"`
	ActivityList []*SeckillActivity `json:"activity_list"`
	Total        int64              `json:"total"`
}

type SeckillDetailRequest struct {
	ActivityID int64   `json:"activity_id" form:"activity_id" query:"activity_id"`
	Token      *string `json:"token,omitempty" form:"token" query:"token"`
}

func (r *SeckillDetailRequest) GetToken() string {
	if r.Token != nil {
		return *r.Token
	}
	return ""
}

type SeckillDetailResponse struct {
	StatusCode int64            `json:"status_code"`
	StatusMsg  string           `json:"status_msg,omitempty"`
	Activity   *SeckillActivity `json:"activity"`
}

type SeckillActivity struct {
	ID             int64  `json:"id"`
	ProductID      int64  `json:"product_id"`
	ProductName    string `json:"product_name"`
	SeckillPrice   int64  `json:"seckill_price"`
	TotalStock     int64  `json:"total_stock"`
	AvailableStock int64  `json:"available_stock"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Status         int64  `json:"status"`
}

type SeckillActionRequest struct {
	Token      string `json:"token" form:"token"`
	ActivityID int64  `json:"activity_id" form:"activity_id"`
}

type SeckillActionResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	OrderNo    string `json:"order_no"`
}

// ==================== Order ====================

type OrderListRequest struct {
	Token  string `json:"token" form:"token" query:"token"`
	Status *int64 `json:"status,omitempty" form:"status" query:"status"`
	Page   *int64 `json:"page,omitempty" form:"page" query:"page"`
	Size   *int64 `json:"size,omitempty" form:"size" query:"size"`
}

func (r *OrderListRequest) GetStatus() int64 {
	if r.Status != nil {
		return *r.Status
	}
	return -1
}

func (r *OrderListRequest) GetPage() int64 {
	if r.Page != nil {
		return *r.Page
	}
	return 1
}

func (r *OrderListRequest) GetSize() int64 {
	if r.Size != nil {
		return *r.Size
	}
	return 10
}

type OrderListResponse struct {
	StatusCode int64    `json:"status_code"`
	StatusMsg  string   `json:"status_msg,omitempty"`
	OrderList  []*Order `json:"order_list"`
	Total      int64    `json:"total"`
}

type OrderDetailRequest struct {
	Token   string `json:"token" form:"token" query:"token"`
	OrderNo string `json:"order_no" form:"order_no" query:"order_no"`
}

type OrderDetailResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	Order      *Order `json:"order"`
}

type OrderPayRequest struct {
	Token   string `json:"token" form:"token"`
	OrderNo string `json:"order_no" form:"order_no"`
}

type OrderPayResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type OrderCancelRequest struct {
	Token   string `json:"token" form:"token"`
	OrderNo string `json:"order_no" form:"order_no"`
}

type OrderCancelResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Order struct {
	ID          int64  `json:"id"`
	OrderNo     string `json:"order_no"`
	UserID      int64  `json:"user_id"`
	ProductID   int64  `json:"product_id"`
	ProductName string `json:"product_name"`
	ActivityID  int64  `json:"activity_id"`
	Amount      int64  `json:"amount"`
	Status      int64  `json:"status"`
	CreatedAt   string `json:"created_at"`
}
