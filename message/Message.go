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
)
