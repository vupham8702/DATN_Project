package message

type Message struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
}

var (
	Success             = Message{Code: 200, Message: "Thành công"}
	ValidationError     = Message{Code: 400, Message: "Yêu cầu thất bại"}
	AccountDoesNotExist = Message{Code: 403, Message: "Tài khoản không tồn tại"}
	UserHasBeenLocked   = Message{Code: 412, Message: "Tài khoản của bạn đã bị khóa. Vui lòng liên hệ quản trị viên để biết thêm chi tiết"}
	UserNotFound        = Message{Code: 400, Message: "Không tìm thấy thông tin người dùng!"}
	EmailNotExist       = Message{Code: 400, Message: "Tài khoản không tồn tại trên hệ thống"}
	PasswordNotCorrect  = Message{Code: 400, Message: "Mật khẩu không chính xác"}
	InternalServerError = Message{Code: 500, Message: "Lỗi hệ thống"}
	UnAuthorizedError   = Message{Code: 401, Message: "Phiên đăng nhập hết hạn hoặc tài khoản đã bị khóa"}
	TokenExpired        = Message{Code: 413, Message: "Phiên đăng nhập hết hạn"}
	NotFound            = Message{Code: 400, Message: "Không tìm thấy dữ liệu"}
	UserActionNotFound  = Message{Code: 400, Message: "Không tìm thấy thông tin người dùng hiện tại!"}
	FieldNotExist       = Message{Code: 500, Message: "Trường dữ liệu không tồn tại!"}
	ExcuteDatabaseError = Message{Code: 500, Message: "Thao tác với CSDL thất bại!"}
)
