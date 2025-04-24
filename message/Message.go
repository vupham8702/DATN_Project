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
	TooManyRequests     = Message{Code: 429, Message: "Quá nhiều yêu cầu. Vui lòng thử lại sau."}

	// Registration messages
	EmailAlreadyExists     = Message{Code: 400, Message: "Email này đã được đăng ký"}
	RegistrationSuccess    = Message{Code: 200, Message: "Đăng ký thành công. Vui lòng kiểm tra email để xác thực tài khoản."}
	InvalidVerifyToken     = Message{Code: 400, Message: "Link xác thực không hợp lệ hoặc đã hết hạn. Vui lòng đăng ký lại."}
	EmailVerifySuccess     = Message{Code: 200, Message: "Email đã được xác thực thành công. Bạn có thể đăng nhập ngay bây giờ."}
	EmailNotVerified       = Message{Code: 400, Message: "Vui lòng xác thực email trước khi đăng nhập."}
	PasswordRequirements   = Message{Code: 400, Message: "Mật khẩu phải có ít nhất 8 ký tự, bao gồm chữ hoa, chữ thường, số và ký tự đặc biệt."}
	EmailExistsNotVerified = Message{Code: 409, Message: "Email này đã được đăng ký nhưng chưa xác thực. Vui lòng kiểm tra email hoặc yêu cầu gửi lại email xác thực."}
	VerificationEmailSent  = Message{Code: 200, Message: "Email xác thực đã được gửi lại. Vui lòng kiểm tra hộp thư của bạn."}
	TokenExpiredOrInvalid  = Message{Code: 400, Message: "Link xác thực đã hết hạn hoặc không hợp lệ. Vui lòng yêu cầu gửi lại email xác thực."}
	ApprovalAccountPenning = Message{Code: 400, Message: "Tài khoản của bạn đang chờ phê duyệt. Vui long liên hệ quả trị viên để hỗ trợ."}
)
