Mô Tả Chi Tiết Chức Năng và Luồng Hệ Thống Hỗ Trợ Tuyển Dụng
Hệ thống trực tuyến kết nối người tìm việc và nhà tuyển dụng, xây dựng bằng Go và React

Mục lục
Tổng Quan Hệ Thống
Chức Năng Chi Tiết và Luồng Xử Lý
Quản lý thông tin của người dùng, nhà tuyển dụng
Chức năng Chọn mẫu CV cho người dùng
Chức năng Ứng tuyển vào bài đăng của nhà tuyển dụng
Chức năng Đăng bài tuyển dụng của nhà tuyển dụng
Chức năng Quản lý người dùng hệ thống
Chức năng Cá nhân hóa cho người dùng
Bảng Tổng Hợp Phân Quyền
Sơ Đồ Luồng Quy Trình Chính
Kết Luận
1. Tổng Quan Hệ Thống
Hệ thống hỗ trợ tuyển dụng là một nền tảng trực tuyến kết nối người tìm việc và nhà tuyển dụng, giúp quản lý quy trình tuyển dụng một cách hiệu quả. Hệ thống được xây dựng với backend bằng Go và frontend bằng React, cung cấp giao diện người dùng trực quan và trải nghiệm mượt mà.

1.1. Mục tiêu hệ thống
Cung cấp nền tảng để người dùng dễ dàng tìm kiếm công việc phù hợp
Hỗ trợ người dùng tạo CV chuyên nghiệp từ mẫu có sẵn
Giúp nhà tuyển dụng tiếp cận ứng viên tiềm năng
Tự động hóa và tối ưu hóa quy trình tuyển dụng
Cá nhân hóa trải nghiệm người dùng sử dụng công nghệ học máy
1.2. Các đối tượng người dùng
Đối tượng	Mô tả
Khách	Người truy cập website chưa đăng nhập, có quyền truy cập hạn chế.
Người dùng thường	Người tìm việc đã đăng ký tài khoản và đăng nhập vào hệ thống.
Nhà tuyển dụng	Đơn vị/tổ chức có nhu cầu tuyển dụng, đăng bài và quản lý tuyển dụng.
Quản trị viên	Người có quyền quản lý toàn bộ hệ thống, người dùng và nội dung.
2. Chức Năng Chi Tiết và Luồng Xử Lý
2.1. Quản lý thông tin của người dùng, nhà tuyển dụng
Mô tả chức năng:
Cho phép người dùng và nhà tuyển dụng xem, quản lý và cập nhật thông tin cá nhân/tổ chức sau khi đăng ký tài khoản. Việc cập nhật thông tin giúp duy trì CV luôn được cập nhật chính xác cho người dùng và giúp nhà tuyển dụng xây dựng uy tín.

Các task chính:
Task	Mô tả	Vai trò thực hiện
Đăng ký tài khoản	Tạo tài khoản mới với thông tin cơ bản	Người dùng, Nhà tuyển dụng
Đăng nhập	Xác thực và truy cập vào hệ thống	Người dùng, Nhà tuyển dụng, Admin
Quản lý hồ sơ cá nhân	Xem và cập nhật thông tin cá nhân, học vấn, kinh nghiệm	Người dùng
Quản lý thông tin công ty	Cập nhật thông tin công ty, lĩnh vực hoạt động, quy mô	Nhà tuyển dụng
Đổi mật khẩu	Thay đổi mật khẩu đăng nhập	Người dùng, Nhà tuyển dụng, Admin
Luồng xử lý cho người dùng:
1. Đăng ký tài khoản: Người dùng điền thông tin cơ bản (họ tên, email, mật khẩu) và xác thực email
2. Đăng nhập: Người dùng đăng nhập bằng email và mật khẩu đã đăng ký
3. Vào trang cá nhân: Người dùng truy cập trang thông tin cá nhân để xem và quản lý thông tin
4. Cập nhật thông tin: Người dùng điền/cập nhật các thông tin chi tiết như học vấn, kinh nghiệm làm việc, kỹ năng, chứng chỉ, thông tin liên hệ
5. Lưu thông tin: Hệ thống lưu thông tin vào CSDL và cập nhật hiển thị cho người dùng
Luồng xử lý cho nhà tuyển dụng:
1. Đăng ký tài khoản: Nhà tuyển dụng điền thông tin cơ bản về công ty và người đại diện
2. Chờ phê duyệt: Admin kiểm tra và phê duyệt tài khoản nhà tuyển dụng
3. Nhận thông báo và đăng nhập: Nhà tuyển dụng nhận email xác nhận phê duyệt và đăng nhập vào hệ thống
4. Quản lý thông tin công ty: Cập nhật thông tin chi tiết về công ty: lĩnh vực hoạt động, quy mô, địa chỉ, website, logo, giới thiệu công ty
5. Lưu thông tin: Hệ thống lưu thông tin và hiển thị trong trang hồ sơ công ty
Các API Endpoint chính:
POST /api/Auth/register/user
Đăng ký tài khoản người dùng

POST /api/Auth/register/employer
Đăng ký tài khoản nhà tuyển dụng

POST /api/Auth/login
Đăng nhập vào hệ thống

PUT /api/User/profile
Cập nhật thông tin người dùng

Yêu cầu bảo mật:
Mật khẩu:
Độ dài tối thiểu: 8 ký tự
Phải chứa ít nhất 1 chữ hoa, 1 chữ thường, 1 số, 1 ký tự đặc biệt
Mật khẩu được lưu dưới dạng mã hóa (bcrypt)
Xác thực email:
Gửi email xác thực sau khi đăng ký
Link xác thực có thời hạn 24 giờ
Yêu cầu xác thực email trước khi sử dụng đầy đủ tính năng
Xác thực Token:
JWT Token với thời hạn ngắn (30 phút)
Refresh Token với thời hạn dài hơn (7 ngày)
Cơ chế làm mới token tự động
2.2. Chức năng Chọn mẫu CV cho người dùng
Mô tả chức năng:
Cho phép người dùng sau khi đăng ký và cập nhật thông tin cá nhân có thể chọn mẫu CV từ thư viện mẫu có sẵn, hệ thống sẽ tự động điền thông tin người dùng vào mẫu CV và cho phép tải về dưới dạng file Word hoặc PDF.

Các task chính:
Task	Mô tả	Vai trò thực hiện
Xem thư viện mẫu CV	Duyệt qua các mẫu CV có sẵn trong hệ thống	Người dùng
Chọn mẫu CV	Lựa chọn mẫu CV phù hợp với nhu cầu	Người dùng
Tùy chỉnh CV	Điều chỉnh một số thông tin hiển thị trên CV	Người dùng
Xem trước CV	Xem bản xem trước của CV đã được điền thông tin	Người dùng
Tải CV	Tải CV về máy dưới dạng file Word hoặc PDF	Người dùng
Quản lý mẫu CV	Thêm, sửa, xóa mẫu CV trong thư viện	Admin
Luồng xử lý chức năng CV:
1. Truy cập thư viện mẫu CV: Người dùng vào phần "Mẫu CV" từ menu chính hoặc trang cá nhân
2. Duyệt và chọn mẫu: Người dùng xem các mẫu CV có sẵn và chọn mẫu phù hợp
3. Xem trước và tùy chỉnh: Hệ thống tự động điền thông tin người dùng vào mẫu CV và hiển thị bản xem trước, người dùng có thể tùy chỉnh các thông tin hiển thị
4. Tải về: Người dùng nhấn nút "Tải về" để lưu CV dưới dạng file Word hoặc PDF về máy tính
5. Lưu lịch sử: Hệ thống lưu lại lịch sử CV đã tạo trong tài khoản người dùng để dễ dàng truy cập lại sau này
Các API Endpoint chính:
GET /api/CV/templates
Lấy danh sách mẫu CV

POST /api/CV/create
Tạo CV từ mẫu

GET /api/CV/download/{cvId}
Tải CV ở định dạng được chọn

2.3. Chức năng Ứng tuyển vào bài đăng của nhà tuyển dụng
Mô tả chức năng:
Cho phép người dùng ứng tuyển vào các bài đăng tuyển dụng trên hệ thống bằng cách sử dụng CV đã tạo. Sau khi ứng tuyển, người dùng sẽ nhận được thông báo về trạng thái đơn ứng tuyển từ nhà tuyển dụng.

Các task chính:
Task	Mô tả	Vai trò thực hiện
Tìm kiếm công việc	Tìm kiếm bài đăng tuyển dụng theo tiêu chí	Khách, Người dùng
Xem chi tiết bài đăng	Xem thông tin chi tiết về công việc và yêu cầu	Khách, Người dùng
Ứng tuyển công việc	Nộp CV để ứng tuyển vào công việc mong muốn	Người dùng
Xem trạng thái ứng tuyển	Kiểm tra trạng thái đơn ứng tuyển đã nộp	Người dùng
Nhận thông báo	Nhận thông báo khi có cập nhật về đơn ứng tuyển	Người dùng
Luồng xử lý chức năng ứng tuyển:
1. Tìm kiếm công việc: Người dùng tìm kiếm công việc phù hợp thông qua bộ lọc (vị trí, ngành nghề, địa điểm, mức lương, v.v.)
2. Xem chi tiết công việc: Người dùng xem thông tin chi tiết về công việc, yêu cầu ứng viên, phúc lợi, thông tin công ty
3. Nộp đơn ứng tuyển: Người dùng nhấn nút "Ứng tuyển", chọn CV đã tạo (hoặc tải lên CV mới), và có thể viết thư giới thiệu
4. Xác nhận ứng tuyển: Hệ thống xác nhận việc nộp đơn ứng tuyển và gửi thông báo đến người dùng và nhà tuyển dụng
5. Theo dõi trạng thái: Người dùng theo dõi trạng thái đơn ứng tuyển trong phần "Quản lý ứng tuyển" (Đã nộp, Đang xem xét, Được chọn phỏng vấn, Từ chối)
6. Nhận thông báo cập nhật: Khi nhà tuyển dụng cập nhật trạng thái đơn ứng tuyển, người dùng sẽ nhận được thông báo qua email và trong hệ thống
Các API Endpoint chính:
GET /api/Job/search
Tìm kiếm công việc

GET /api/Job/{jobId}
Xem chi tiết công việc

POST /api/Job/{jobId}/apply
Ứng tuyển công việc

GET /api/User/applications
Lấy danh sách đơn ứng tuyển của người dùng

2.4. Chức năng Đăng bài tuyển dụng của nhà tuyển dụng
Mô tả chức năng:
Cho phép nhà tuyển dụng đăng các bài tuyển dụng để tìm kiếm ứng viên phù hợp. Các bài đăng sẽ được Admin xét duyệt trước khi hiển thị công khai. Nhà tuyển dụng có thể quản lý bài đăng và xét duyệt ứng viên đã ứng tuyển.

Các task chính:
Task	Mô tả	Vai trò thực hiện
Tạo bài đăng tuyển dụng	Tạo bài đăng mới với thông tin chi tiết về vị trí tuyển dụng	Nhà tuyển dụng
Chỉnh sửa bài đăng	Cập nhật thông tin bài đăng tuyển dụng	Nhà tuyển dụng
Duyệt bài đăng	Xét duyệt bài đăng tuyển dụng trước khi hiển thị	Admin
Quản lý ứng viên	Xem và đánh giá hồ sơ ứng viên đã ứng tuyển	Nhà tuyển dụng
Phản hồi ứng viên	Gửi phản hồi cho ứng viên (từ chối hoặc mời phỏng vấn)	Nhà tuyển dụng
Ẩn/Hiện bài đăng	Tạm dừng hoặc tiếp tục hiển thị bài đăng	Nhà tuyển dụng
Luồng xử lý đăng bài tuyển dụng:
1. Tạo bài đăng: Nhà tuyển dụng truy cập mục "Quản lý tuyển dụng" và chọn "Tạo bài đăng mới"
2. Nhập thông tin chi tiết: Điền các thông tin: tiêu đề công việc, mô tả công việc, yêu cầu ứng viên, phúc lợi, mức lương, địa điểm làm việc, hạn nộp hồ sơ, số lượng tuyển
3. Nộp bài đăng: Nhà tuyển dụng gửi bài đăng để chờ xét duyệt
4. Chờ Admin duyệt: Bài đăng ở trạng thái "Chờ duyệt" cho đến khi được Admin xét duyệt
5. Hiển thị bài đăng: Sau khi được duyệt, bài đăng xuất hiện trên bảng tin tuyển dụng và người dùng có thể xem và ứng tuyển
Luồng xử lý quản lý ứng viên:
1. Nhận thông báo ứng tuyển: Nhà tuyển dụng nhận thông báo khi có ứng viên ứng tuyển vào bài đăng
2. Xem danh sách ứng viên: Truy cập vào mục "Quản lý ứng viên" để xem danh sách ứng viên đã ứng tuyển cho từng bài đăng
3. Đánh giá hồ sơ: Xem chi tiết CV và thông tin ứng viên, đánh giá mức độ phù hợp
4. Phản hồi ứng viên: Gửi phản hồi cho ứng viên: từ chối, yêu cầu thông tin bổ sung, hoặc mời phỏng vấn
5. Cập nhật trạng thái: Cập nhật trạng thái ứng viên trong hệ thống (Đã xem, Từ chối, Mời phỏng vấn, Đã tuyển)
Các API Endpoint chính:
POST /api/Recruitment/job-postings
Tạo bài đăng tuyển dụng

PUT /api/Admin/job-postings/{id}/review
Xét duyệt bài đăng

PUT /api/Recruitment/applications/{id}/status
Cập nhật trạng thái ứng viên

2.5. Chức năng Quản lý người dùng hệ thống
Mô tả chức năng:
Cho phép Admin quản lý toàn bộ người dùng trong hệ thống, bao gồm người dùng thường và nhà tuyển dụng. Admin có quyền thêm, sửa, xóa, xem thông tin, chặn/mở khóa tài khoản và quản lý bài đăng tuyển dụng. Ngoài ra, Admin còn có thể theo dõi các hoạt động của hệ thống thông qua biểu đồ thống kê.

Các task chính:
Task	Mô tả	Vai trò thực hiện
Quản lý tài khoản người dùng	Thêm, sửa, xóa, chặn/mở khóa tài khoản người dùng	Admin
Quản lý tài khoản nhà tuyển dụng	Phê duyệt, sửa, xóa, chặn/mở khóa tài khoản nhà tuyển dụng	Admin
Quản lý bài đăng tuyển dụng	Phê duyệt, sửa, xóa, ẩn/hiện bài đăng tuyển dụng	Admin
Theo dõi hoạt động hệ thống	Xem thống kê, báo cáo về hoạt động của hệ thống	Admin
Gửi thông báo hệ thống	Gửi thông báo đến người dùng hoặc nhà tuyển dụng	Admin
Quản lý mẫu CV	Thêm, sửa, xóa mẫu CV trong thư viện	Admin
Luồng xử lý quản lý người dùng:
1. Truy cập trang quản trị: Admin đăng nhập và truy cập vào trang quản trị hệ thống
2. Xem danh sách người dùng: Mở mục "Quản lý người dùng" để xem danh sách tất cả người dùng với các thông tin cơ bản
3. Thực hiện các thao tác quản lý: Admin có thể xem chi tiết, thêm mới, chỉnh sửa thông tin, hoặc xóa/chặn tài khoản người dùng
4. Xác nhận thao tác: Hệ thống yêu cầu xác nhận trước khi thực hiện các thao tác quan trọng như xóa/chặn tài khoản
5. Hoàn tất thao tác: Hệ thống thực hiện thao tác và cập nhật trạng thái mới
Luồng xử lý quản lý bài đăng tuyển dụng:
1. Truy cập mục quản lý bài đăng: Admin mở mục "Quản lý bài đăng tuyển dụng" từ trang quản trị
2. Xem danh sách bài đăng: Xem tất cả bài đăng trong hệ thống với thông tin cơ bản và trạng thái
3. Duyệt bài đăng mới: Kiểm tra và phê duyệt hoặc từ chối các bài đăng mới từ nhà tuyển dụng
4. Quản lý bài đăng hiện có: Xem chi tiết, chỉnh sửa nếu cần, hoặc ẩn/xóa các bài đăng vi phạm quy định
5. Thông báo cho nhà tuyển dụng: Gửi thông báo về kết quả xét duyệt hoặc thông báo vi phạm đến nhà tuyển dụng
Các API Endpoint chính:
GET /api/Admin/users
Lấy danh sách người dùng

PUT /api/Admin/users/{id}/status
Cập nhật trạng thái người dùng

PUT /api/Admin/employers/{id}/approve
Phê duyệt tài khoản nhà tuyển dụng

GET /api/Admin/stats
Lấy thống kê hệ thống

2.6. Chức năng Cá nhân hóa cho người dùng
Mô tả chức năng:
Sử dụng công nghệ học máy để phân tích thông tin của người dùng và đề xuất các công việc phù hợp, tăng trải nghiệm người dùng và hiệu quả tuyển dụng. Hệ thống sẽ thu thập dữ liệu từ hồ sơ người dùng, lịch sử ứng tuyển, phân tích và đề xuất công việc phù hợp dựa trên thuật toán so sánh.

Các task chính:
Task	Mô tả	Vai trò thực hiện
Thu thập dữ liệu người dùng	Thu thập thông tin từ hồ sơ người dùng và lịch sử ứng tuyển	Hệ thống
Phân tích dữ liệu	Phân tích dữ liệu để tìm mẫu và đặc điểm	Hệ thống
So sánh dữ liệu	So sánh thông tin người dùng với yêu cầu công việc	Hệ thống
Đề xuất công việc	Đề xuất công việc phù hợp cho người dùng	Hệ thống
Xem đề xuất	Xem danh sách công việc được đề xuất	Người dùng
Tương tác với đề xuất	Phản hồi về đề xuất (quan tâm, bỏ qua)	Người dùng
Luồng xử lý cá nhân hóa:
1. Thu thập và phân tích dữ liệu: Hệ thống thu thập dữ liệu từ hồ sơ người dùng, lịch sử ứng tuyển, và tương tác với bài đăng
2. Xây dựng hồ sơ người dùng: Hệ thống tạo hồ sơ cá nhân hóa dựa trên kỹ năng, kinh nghiệm, sở thích của người dùng
3. So sánh với bài đăng: Thuật toán so sánh hồ sơ của người dùng với các yêu cầu của bài đăng tuyển dụng
4. Tạo danh sách đề xuất: Hệ thống sắp xếp và lọc các công việc phù hợp nhất với hồ sơ người dùng
5. Hiển thị đề xuất: Người dùng xem danh sách công việc được đề xuất khi đăng nhập hoặc truy cập trang chủ
6. Cải thiện thuật toán: Hệ thống học từ phản hồi của người dùng để cải thiện chất lượng đề xuất trong tương lai
Các API Endpoint chính:
GET /api/User/recommendations
Lấy danh sách công việc được đề xuất cho người dùng

POST /api/User/recommendations/{jobId}/feedback
Gửi phản hồi về đề xuất công việc

3. Bảng Tổng Hợp Phân Quyền
Chức năng			Khách	Người dùng  Nhà tuyển dụng    Admin
Xem bài đăng tuyển dụng		✓	     ✓		✓		✓
Tìm kiếm công việc		✓	     ✓		✓		✓
Đăng ký tài khoản		✓	     -		-		-
Quản lý thông tin cá nhân	-	     ✓		✓		✓
Tạo và quản lý CV		-	     ✓		-		-
Ứng tuyển công việc		-	     ✓		-		-
Xem đề xuất công việc		-	     ✓		-		-
Đăng bài tuyển dụng		-	     -		✓		✓
Quản lý ứng viên		-	     -		✓		✓
Quản lý người dùng		-	     -		-		✓
Duyệt tài khoản nhà tuyển dụng	-	     -		-		✓
Duyệt bài đăng tuyển dụng	-	     -		-		✓
Quản lý mẫu CV			-	     -		-		✓
Xem thống kê hệ thống		-	     -		-		✓
4. Sơ Đồ Luồng Quy Trình Chính
4.1. Luồng đăng ký và đăng nhập
1. Người dùng/Nhà tuyển dụng đăng ký tài khoản với email và mật khẩu
2. Hệ thống gửi email xác thực tài khoản
3. Người dùng/Nhà tuyển dụng xác thực email
4. Đối với nhà tuyển dụng, Admin xét duyệt tài khoản
5. Người dùng/Nhà tuyển dụng đăng nhập vào hệ thống
6. Cập nhật thông tin cá nhân/công ty
4.2. Luồng tìm kiếm và ứng tuyển công việc
1. Người dùng tìm kiếm công việc theo tiêu chí
2. Xem chi tiết công việc quan tâm
3. Tạo/Chọn CV từ thư viện mẫu
4. Ứng tuyển vào công việc với CV đã chọn
5. Theo dõi trạng thái ứng tuyển
6. Nhận phản hồi từ nhà tuyển dụng
4.3. Luồng đăng và quản lý bài tuyển dụng
1. Nhà tuyển dụng tạo bài đăng tuyển dụng mới
2. Gửi bài đăng để chờ xét duyệt
3. Admin xét duyệt bài đăng
4. Bài đăng được hiển thị công khai sau khi được duyệt
5. Nhà tuyển dụng nhận thông báo khi có người ứng tuyển
6. Nhà tuyển dụng xem xét hồ sơ và phản hồi ứng viên
5. Kết Luận
Hệ thống hỗ trợ tuyển dụng được thiết kế để tối ưu hóa quy trình tuyển dụng từ cả phía người tìm việc và nhà tuyển dụng. Với các chức năng đa dạng như quản lý thông tin cá nhân, tạo CV chuyên nghiệp, tìm kiếm và ứng tuyển công việc, đăng và quản lý bài tuyển dụng, cùng với công nghệ cá nhân hóa dựa trên học máy, hệ thống cung cấp giải pháp toàn diện cho nhu cầu tuyển dụng hiện đại.

Luồng xử lý của các chức năng được thiết kế đơn giản, trực quan và hiệu quả, đảm bảo trải nghiệm người dùng tốt nhất trên nền tảng. Các API và giao diện Frontend được xây dựng đồng bộ, tạo nên một hệ thống liền mạch và dễ sử dụng.

Việc phân quyền rõ ràng giữa các đối tượng người dùng đảm bảo tính bảo mật và quyền riêng tư của dữ liệu, trong khi vẫn cho phép Admin quản lý và giám sát toàn diện hoạt động của hệ thống.

Với kiến trúc backend xây dựng trên Go và frontend trên React, hệ thống đảm bảo hiệu suất cao, khả năng mở rộng tốt và trải nghiệm người dùng mượt mà.