# Tài liệu mô tả chức năng Đăng ký và Đăng nhập

## 1. Tổng quan

Chức năng **Đăng ký** cho phép người dùng (người tìm việc hoặc nhà tuyển dụng) tạo tài khoản mới; sau khi đăng ký, họ sẽ nhận được email xác thực để kích hoạt.  
Chức năng **Đăng nhập** cho phép người dùng đã có tài khoản truy cập hệ thống bằng email và mật khẩu, sử dụng JWT để xác thực và quản lý phiên làm việc.

## 2. Luồng xử lý

### 2.1. Đăng ký tài khoản

1. Người dùng truy cập trang đăng ký (người tìm việc hoặc nhà tuyển dụng).  
2. Người dùng điền thông tin cơ bản:
   - Họ tên
   - Email
   - Mật khẩu  
3. Frontend gửi yêu cầu đến API `/api/Auth/register/user` (hoặc `/api/Auth/register/employer`).  
4. Backend xử lý:
   - Kiểm tra email chưa tồn tại
   - Validate mật khẩu theo quy tắc bảo mật
   - Mã hóa mật khẩu bằng bcrypt
   - Tạo bản ghi người dùng/chủ tuyển dụng với trạng thái `isEmailVerified = false`
   - Tạo token xác thực (GUID) lưu trong Redis (hoặc DB) với thời hạn 24 giờ
   - Gửi email chứa link xác thực tới người dùng  
5. Frontend hiển thị màn hình thông báo:
   - Tiêu đề “Registration Successful”
   - Hướng dẫn kiểm tra email để kích hoạt tài khoản
   - Thông tin về thời hạn link (24 giờ)
   - Nút “Gửi lại email xác thực” (nếu cần)
   - Nút “Quay lại trang đăng nhập”

### 2.2. Xác thực email

1. Người dùng nhấp vào link trong email (ví dụ: `/verify-email?token={token}&email={email}`).
2. Frontend gọi API `/api/Auth/verify-email` với token và email.
3. Backend xử lý:
   - Xác thực token còn hạn và khớp với email
   - Cập nhật `isEmailVerified = true`
   - Xóa token khỏi Redis
4. Frontend hiển thị thông báo:
   - “Email verified successfully”
   - Nút “Đăng nhập ngay”

### 2.3. Đăng nhập

1. Người dùng truy cập trang `/login`
2. Người dùng nhập:
   - Email
   - Mật khẩu  
3. Frontend gửi yêu cầu đến API `/api/Auth/login`.  
4. Backend xử lý:
   - Tìm người dùng theo email
   - Kiểm tra `isEmailVerified == true`; nếu chưa, trả về lỗi yêu cầu xác thực email
   - So sánh mật khẩu nhập với bcrypt-hash
   - Nếu đúng, tạo JWT access token (hết hạn 30 phút) và refresh token (hết hạn 7 ngày), lưu refresh token (DB hoặc Redis)
   - Trả về cả hai token cho frontend  
5. Frontend:
   - Lưu access token vào memory (hoặc HttpOnly cookie)
   - Lưu refresh token (nếu dùng cookie hoặc secure storage)
   - Chuyển hướng người dùng vào trang dashboard

## 3. Các API Endpoint

### 3.1. Đăng ký người tìm việc

- **Endpoint**: `POST /api/Auth/register/user`  
- **Request Body**:
  ```json
  {
    "fullName": "Nguyen Van A",
    "email": "ava@example.com",
    "password": "SecurePass123!"
  }
  ```
- **Response**:
  ```json
  {
    "message": "Registration successful. Please check your email to verify your account."
  }
  ```

### 3.2. Đăng ký nhà tuyển dụng

- **Endpoint**: `POST /api/Auth/register/employer`  
- **Request Body**:
  ```json
  {
    "companyName": "CÔNG TY ABC",
    "email": "hr@abc.com",
    "password": "EmployerPass123!"
  }
  ```
- **Response**:
  ```json
  {
    "message": "Employer registration successful. Please check your email to verify your account."
  }
  ```

### 3.3. Xác thực email

- **Endpoint**: `POST /api/Auth/verify-email`  
- **Request Body**:
  ```json
  {
    "token": "a1b2c3d4-e5f6-7890-ab12-cd34ef56gh78",
    "email": "ava@example.com"
  }
  ```
- **Response**:
  ```json
  {
    "message": "Email has been verified successfully."
  }
  ```

### 3.4. Đăng nhập

- **Endpoint**: `POST /api/Auth/login`  
- **Request Body**:
  ```json
  {
    "email": "ava@example.com",
    "password": "SecurePass123!"
  }
  ```
- **Response**:
  ```json
  {
    "accessToken": "<JWT_ACCESS_TOKEN>",
    "refreshToken": "<JWT_REFRESH_TOKEN>",
    "expiresIn": 1800
  }
  ```

## 4. Các trang Frontend

### 4.1. Trang Đăng ký

- **Route**:
  - `/register/user` → component `RegisterUser.tsx`
  - `/register/employer` → component `RegisterEmployer.tsx`
- **Chức năng**:
  - Form nhập thông tin (họ tên hoặc tên công ty, email, mật khẩu)
  - Validate client-side theo quy tắc mật khẩu
  - Gửi yêu cầu đăng ký
  - Hiển thị màn hình “Registration Successful”
  - Nút “Gửi lại email xác thực”
  - Nút “Quay lại đăng nhập”

### 4.2. Trang Xác thực email

- **Route**: `/verify-email`
- **Component**: `VerifyEmail.tsx`
- **Chức năng**:
  - Đọc token & email từ query params
  - Gọi API verify-email
  - Hiển thị thông báo thành công hoặc lỗi

### 4.3. Trang Đăng nhập

- **Route**: `/login`
- **Component**: `Login.tsx`
- **Chức năng**:
  - Form nhập email, mật khẩu
  - Validate client-side (email hợp lệ, mật khẩu không rỗng)
  - Gửi yêu cầu login
  - Xử lý error (chưa xác thực email, sai mật khẩu)
  - Lưu token, chuyển hướng

## 5. Yêu cầu bảo mật

### 5.1. Mật khẩu

- Độ dài tối thiểu: 8 ký tự  
- Phải chứa:
  - Ít nhất 1 chữ hoa  
  - Ít nhất 1 chữ thường  
  - Ít nhất 1 chữ số  
  - Ít nhất 1 ký tự đặc biệt  
- Lưu trữ: mã hóa bằng **bcrypt**

### 5.2. Xác thực email

- Token xác thực:
  - Định dạng: GUID  
  - Thời hạn: 24 giờ  
  - Lưu trữ: Redis hoặc DB  
- Không cho phép sử dụng chức năng đầy đủ nếu chưa xác thực email

### 5.3. Xác thực Token

- **Access Token** (JWT):
  - Hết hạn sau 30 phút  
  - Truyền qua header `Authorization: Bearer <token>`  
- **Refresh Token**:
  - Hết hạn sau 7 ngày  
  - Lưu trữ server-side (Redis/DB)  
  - Frontend có thể tự động gọi API làm mới khi access token hết hạn

## 6. Validate dữ liệu nhập

### 6.1. Trang Đăng ký

| Field         | Quy tắc validate                                         | Thông báo lỗi                              |
|---------------|----------------------------------------------------------|---------------------------------------------|
| Họ tên/Tên CT | Không để trống                                           | “Vui lòng nhập họ tên” hoặc “Vui lòng nhập tên công ty” |
| Email         | Phải là địa chỉ email hợp lệ                             | “Vui lòng nhập địa chỉ email hợp lệ”        |
| Password      | Theo quy tắc bảo mật (5.1)                               | “Mật khẩu không hợp lệ”                     |

### 6.2. Trang Đăng nhập

| Field    | Quy tắc validate               | Thông báo lỗi                              |
|----------|--------------------------------|---------------------------------------------|
| Email    | Phải là địa chỉ email hợp lệ   | “Vui lòng nhập email hợp lệ”               |
| Password | Không để trống                 | “Vui lòng nhập mật khẩu”                   |

## 7. Thông báo lỗi và thành công

### 7.1. Thông báo lỗi chung

| Tình huống                                   | Thông báo lỗi                                                             |
|-----------------------------------------------|----------------------------------------------------------------------------|
| Email đã tồn tại                              | “Email này đã được đăng ký”                                               |
| Token xác thực không hợp lệ/đã hết hạn        | “Link xác thực không hợp lệ hoặc đã hết hạn. Vui lòng đăng ký lại.”       |
| Chưa xác thực email khi login                 | “Vui lòng xác thực email trước khi đăng nhập.”                            |
| Sai email hoặc mật khẩu                       | “Email hoặc mật khẩu không đúng.”                                         |
| Lỗi server/API                                | “Có lỗi xảy ra. Vui lòng thử lại sau.”                                     |

### 7.2. Thông báo thành công

| Tính năng            | Thông báo thành công                                                                                          |
|----------------------|---------------------------------------------------------------------------------------------------------------|
| Đăng ký thành công   | “Registration successful. Please check your email to verify your account.”                                     |
| Xác thực email       | “Email has been verified successfully. You can now log in.”                                                   |
| Đăng nhập thành công | Không hiển thị message, chuyển hướng vào dashboard với access token và refresh token được lưu tự động.         |

## 8. Lưu ý triển khai

1. **URL email xác thực**:
   - `/verify-email?token={token}&email={email}`  
2. **Xử lý token**:
   - Frontend ưu tiên lấy token từ query params  
   - Hết hạn hoặc sai token → hiển thị màn hình lỗi và hướng dẫn gửi lại  
3. **Cơ chế làm mới token**:
   - Khi access token hết hạn, frontend tự động gọi API `/api/Auth/refresh-token`  
   - Backend kiểm tra refresh token, phát access token mới  
4. **Bảo mật**:
   - Không trả về quá chi tiết thông tin lý do thất bại (tránh so sánh email có tồn tại hay không)  
   - Dùng HttpOnly cookie cho refresh token (nếu có thể)

## 9. Kiểm thử

### 9.1. Chức năng Đăng ký

- Đăng ký với email mới → nhận email xác thực  
- Nhập sai định dạng email → hiển thị lỗi  
- Nhập mật khẩu không đủ mạnh → hiển thị lỗi  
- Đăng ký lại với email đã tồn tại → hiển thị lỗi  
- Gửi lại email xác thực từ màn hình thông báo  

### 9.2. Chức năng Xác thực email

- Nhấp link đúng token → kích hoạt thành công  
- Nhấp link sai/hết hạn → hiển thị lỗi và cho phép gửi lại  

### 9.3. Chức năng Đăng nhập

- Đăng nhập sau khi xác thực email → thành công  
- Đăng nhập trước khi xác thực → hiển thị lỗi  
- Đăng nhập sai mật khẩu hoặc email → hiển thị lỗi  
- Access token hết hạn → frontend tự động làm mới qua refresh token  
- Refresh token hết hạn → người dùng phải đăng nhập lại  
