# Cấu Trúc Dự Án - Hệ Thống Hỗ Trợ Tuyển DỤng

**Backend**: Go | **Frontend**: ReactJS

## Cấu trúc thư mục dự án

```
datn_backend/
├── .idea                # Thư mục cấu hình IDE (IntelliJ/GoLand)
├── config               # Cấu hình ứng dụng và môi trường
├── controller           # Các controllers xử lý HTTP requests
├── docs                 # Tài liệu API và hướng dẫn
├── domain               # Định nghĩa domain models
├── infrastructure       # Kết nối cơ sở dữ liệu, message queue
├── logs                 # File logs
├── message              # Định nghĩa message và thông báo hệ thống
├── migration            # Migrations cơ sở dữ liệu
├── payload              # Định nghĩa request/response payloads
├── public               # Tài nguyên tĩnh
├── response             # Định nghĩa response formats
├── router               # Router và routes configuration
├── service              # Business logic
├── utils                # Các tiện ích và helper
├── websocket            # Xử lý WebSocket cho thời gian thực
├── .dockerignore        # Cấu hình Docker ignore
├── .env                 # Biến môi trường
├── .env.example         # Mẫu biến môi trường
├── .gitignore           # Cấu hình Git ignore
├── .gitlab-ci.yml       # CI/CD pipeline configuration
├── Dockerfile           # Cấu hình Docker
├── go.mod               # Go modules
├── main.go              # Entry point
└── README.md            # Hướng dẫn sử dụng
```

## Chi tiết thư mục

- **config**: Chứa các file cấu hình:
  - `database.go`, `app.go`, `mail.go`, `auth.go`, `cors.go`
- **controller**: Xử lý HTTP requests:
  - `user_controller.go`, `auth_controller.go`, `employer_controller.go`, `job_post_controller.go`, `cv_controller.go`, `application_controller.go`, `admin_controller.go`
- **domain**: Định nghĩa core:
  - `entity/` (User, Employer, JobPost, CVTemplate, Application)
  - `repository/`, `service/`
- **infrastructure**: Triển khai repository, cache, messaging, storage, mail
- **message**: Định nghĩa lỗi và success messages (`error.go`, `success.go`, `en.go`, `vi.go`)
- **migration**: SQL migrations (`000001_create_users_table.up.sql`, ...)
- **payload**: Định nghĩa request/response structs
- **response**: Helper format response (`response.go`, `error_response.go`, `success_response.go`, `pagination.go`)
- **router**: Cấu hình routes và middleware
- **service**: Logic xử lý nghiệp vụ
- **utils**: JWT, hash, validation, file, date, string
- **websocket**: WebSocket client/hub/message

## Mô hình dữ liệu

```go
// domain/entity/user.go
package entity

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Username       string     `gorm:"type:varchar(50);not null;unique"`
    Email          string     `gorm:"type:varchar(100);not null;unique"`
    PasswordHash   string     `gorm:"type:varchar(255);not null"`
    FullName       string     `gorm:"type:varchar(100);not null"`
    Phone          string     `gorm:"type:varchar(20)"`
    Gender         string     `gorm:"type:varchar(10)"`
    DateOfBirth    *time.Time `gorm:"type:date"`
    Address        string     `gorm:"type:varchar(255)"`
    Education      string     `gorm:"type:text"`
    Experience     string     `gorm:"type:text"`
    Skills         string     `gorm:"type:text"`
    Interests      string     `gorm:"type:text"`
    ProfilePicture string     `gorm:"type:varchar(255)"`
    ResumeURL      string     `gorm:"type:varchar(255)"`
    EmailVerified  bool       `gorm:"type:boolean;default:false"`
    IsActive       bool       `gorm:"type:boolean;default:true"`
    Role           string     `gorm:"type:varchar(20);default:'user'"`
    CreatedAt      time.Time  `gorm:"type:timestamp;default:current_timestamp"`
    UpdatedAt      time.Time  `gorm:"type:timestamp;default:current_timestamp"`
}

// domain/entity/employer.go
package entity

import (
    "time"
    "github.com/google/uuid"
)

type Employer struct {
    ID              uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    CompanyName     string    `gorm:"type:varchar(100);not null"`
    Email           string    `gorm:"type:varchar(100);not null;unique"`
    PasswordHash    string    `gorm:"type:varchar(255);not null"`
    Phone           string    `gorm:"type:varchar(20);not null"`
    Address         string    `gorm:"type:varchar(255);not null"`
    Website         string    `gorm:"type:varchar(100)"`
    Industry        string    `gorm:"type:varchar(50)"`
    Description     string    `gorm:"type:text"`
    Logo            string    `gorm:"type:varchar(255)"`
    ContactPerson   string    `gorm:"type:varchar(100);not null"`
    ContactPosition string    `gorm:"type:varchar(50)"`
    ContactEmail    string    `gorm:"type:varchar(100);not null"`
    ContactPhone    string    `gorm:"type:varchar(20);not null"`
    IsVerified      bool      `gorm:"type:boolean;default:false"`
    IsActive        bool      `gorm:"type:boolean;default:false"`
    CreatedAt       time.Time `gorm:"type:timestamp;default:current_timestamp"`
    UpdatedAt       time.Time `gorm:"type:timestamp;default:current_timestamp"`
}

// domain/entity/job_post.go
package entity

import (
    "time"
    "github.com/google/uuid"
)

type JobPost struct {
    ID                 uuid.UUID    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    EmployerID         uuid.UUID    `gorm:"type:uuid;not null"`
    Title              string       `gorm:"type:varchar(100);not null"`
    Description        string       `gorm:"type:text;not null"`
    Requirements       string       `gorm:"type:text;not null"`
    Benefits           string       `gorm:"type:text"`
    Location           string       `gorm:"type:varchar(100);not null"`
    JobType            string       `gorm:"type:varchar(50);not null"`
    Salary             string       `gorm:"type:varchar(50)"`
    ExperienceRequired string       `gorm:"type:varchar(50)"`
    Education          string       `gorm:"type:varchar(100)"`
    Skills             string       `gorm:"type:text"`
    Deadline           *time.Time   `gorm:"type:timestamp"`
    Positions          int          `gorm:"type:int;default:1"`
    IsActive           bool         `gorm:"type:boolean;default:true"`
    IsApproved         bool         `gorm:"type:boolean;default:false"`
    Views              int          `gorm:"type:int;default:0"`
    CreatedAt          time.Time    `gorm:"type:timestamp;default:current_timestamp"`
    UpdatedAt          time.Time    `gorm:"type:timestamp;default:current_timestamp"`
    Employer           Employer     `gorm:"foreignKey:EmployerID"`
    Applications       []Application `gorm:"foreignKey:JobPostID"`
}

// domain/entity/cv.go
package entity

import (
    "time"
    "github.com/google/uuid"
)

type CVTemplate struct {
    ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Name        string     `gorm:"type:varchar(100);not null"`
    Description string     `gorm:"type:text"`
    PreviewURL  string     `gorm:"type:varchar(255);not null"`
    TemplateURL string     `gorm:"type:varchar(255);not null"`
    IsActive    bool       `gorm:"type:boolean;default:true"`
    CreatedAt   time.Time  `gorm:"type:timestamp;default:current_timestamp"`
    UpdatedAt   time.Time  `gorm:"type:timestamp;default:current_timestamp"`
}

// domain/entity/application.go
package entity

import (
    "time"
    "github.com/google/uuid"
)

type Application struct {
    ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    UserID      uuid.UUID  `gorm:"type:uuid;not null"`
    JobPostID   uuid.UUID  `gorm:"type:uuid;not null"`
    ResumeURL   string     `gorm:"type:varchar(255);not null"`
    CoverLetter string     `gorm:"type:text"`
    Status      string     `gorm:"type:varchar(20);default:'pending'"`
    Notes       string     `gorm:"type:text"`
    CreatedAt   time.Time  `gorm:"type:timestamp;default:current_timestamp"`
    UpdatedAt   time.Time  `gorm:"type:timestamp;default:current_timestamp"`
    User        User       `gorm:"foreignKey:UserID"`
    JobPost     JobPost    `gorm:"foreignKey:JobPostID"`
}
```

## Ví dụ triển khai API

### 1. Payload
```go
// payload/job_post_payload.go
package payload

import "time"

type CreateJobPostRequest struct {
    Title              string    `json:"title" binding:"required"`
    Description        string    `json:"description" binding:"required"`
    Requirements       string    `json:"requirements" binding:"required"`
    Benefits           string    `json:"benefits"`
    Location           string    `json:"location" binding:"required"`
    JobType            string    `json:"job_type" binding:"required"`
    Salary             string    `json:"salary"`
    ExperienceRequired string    `json:"experience_required"`
    Education          string    `json:"education"`
    Skills             string    `json:"skills"`
    Deadline           time.Time `json:"deadline"`
    Positions          int       `json:"positions" binding:"required,min=1"`
}

type JobPostResponse struct {
    ID                 string    `json:"id"`
    Title              string    `json:"title"`
    Description        string    `json:"description"`
    Requirements       string    `json:"requirements"`
    Benefits           string    `json:"benefits"`
    Location           string    `json:"location"`
    JobType            string    `json:"job_type"`
    Salary             string    `json:"salary"`
    ExperienceRequired string    `json:"experience_required"`
    Education          string    `json:"education"`
    Skills             string    `json:"skills"`
    Deadline           time.Time `json:"deadline"`
    Positions          int       `json:"positions"`
    IsActive           bool      `json:"is_active"`
    IsApproved         bool      `json:"is_approved"`
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
    Employer           struct {
        ID          string `json:"id"`
        CompanyName string `json:"company_name"`
        Logo        string `json:"logo"`
    } `json:"employer"`
}
```

### 2. Service
```go
// service/job_post_service.go
package service

import (
    "errors"
    "time"
    "github.com/google/uuid"
    "recruitment-system/domain/entity"
    "recruitment-system/infrastructure/persistence"
    "recruitment-system/payload"
)

type JobPostService struct {
    jobPostRepo  persistence.JobPostRepository
    employerRepo persistence.EmployerRepository
}

func NewJobPostService(j job_repo persistence.JobPostRepository, e persistence.EmployerRepository) *JobPostService {
    return &JobPostService{jobPostRepo: j, employerRepo: e}
}

func (s *JobPostService) CreateJobPost(employerID uuid.UUID, req payload.CreateJobPostRequest) (*payload.JobPostResponse, error) {
    emp, err := s.employerRepo.GetByID(employerID)
    if err != nil {
        return nil, err
    }
    if !emp.IsVerified {
        return nil, errors.New("employer not verified")
    }
    jp := &entity.JobPost{
        ID:                 uuid.New(),
        EmployerID:         employerID,
        Title:              req.Title,
        Description:        req.Description,
        Requirements:       req.Requirements,
        Benefits:           req.Benefits,
        Location:           req.Location,
        JobType:            req.JobType,
        Salary:             req.Salary,
        ExperienceRequired: req.ExperienceRequired,
        Education:          req.Education,
        Skills:             req.Skills,
        Deadline:           &req.Deadline,
        Positions:          req.Positions,
        IsActive:           true,
        IsApproved:         false,
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
    }
    if err := s.jobPostRepo.Create(jp); err != nil {
        return nil, err
    }
    resp := &payload.JobPostResponse{ /* map fields */ }
    // map emp to resp.Employer
    return resp, nil
}
```

### 3. Controller
```go
// controller/job_post_controller.go
package controller

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "recruitment-system/payload"
    "recruitment-system/response"
    "recruitment-system/service"
    "recruitment-system/utils"
)

type JobPostController struct { svc *service.JobPostService }

func NewJobPostController(s *service.JobPostService) *JobPostController { return &JobPostController{svc: s} }

func (ctr *JobPostController) CreateJobPost(ctx *gin.Context) {
    claims, ok := ctx.Get("claims")
    if !ok { response.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized"); return }
    user := claims.(*utils.JWTClaims)
    if user.Role != "employer" { response.ErrorResponse(ctx, http.StatusForbidden, "only employers"); return }
    id, _ := uuid.Parse(user.UserID)
    var req payload.CreateJobPostRequest
    if err := ctx.ShouldBindJSON(&req); err != nil { response.ValidationErrorResponse(ctx, err); return }
    res, err := ctr.svc.CreateJobPost(id, req)
    if err != nil { response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error()); return }
    response.SuccessResponse(ctx, http.StatusCreated, "created", res)
}
```

### 4. Router
```go
// router/api_routes.go
// Trong group employer:
employerRoutes.POST("/job-posts", controllers.JobPost.CreateJobPost)
```

## Ví dụ triển khai chức năng cá nhân hóa

### Recommendation Service
```go
// service/recommendation_service.go
package service

import (
    "encoding/json"
    "net/http"
    "bytes"
    "github.com/google/uuid"
    "recruitment-system/domain/entity"
    "recruitment-system/infrastructure/persistence"
    "recruitment-system/payload"
)

type RecommendationService struct { /* ... */ }

func (s *RecommendationService) GetPersonalizedJobRecommendations(userID uuid.UUID) ([]payload.JobPostResponse, error) {
    // tương tự code gốc từ HTML
}
```

### Recommendation Controller
```go
// controller/recommendation_controller.go
package controller

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "recruitment-system/response"
    "recruitment-system/service"
    "recruitment-system/utils"
)

func (ctr *RecommendationController) GetPersonalizedJobRecommendations(ctx *gin.Context) {
    // tương tự code gốc từ HTML
}
```

## Deployment

### Dockerfile
```dockerfile
# Dockerfile
FROM golang:1.17-alpine AS builder

WORKDIR /app

# Cài đặt các dependencies
RUN apk add --no-cache git

# Copy go mod và sum
COPY go.mod go.sum ./

# Tải dependencies
RUN go mod download

# Copy source code
COPY . .

# Build ứng dụng
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o recruitment-system .

# Sử dụng alpine để tạo image nhẹ hơn
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary từ builder
COPY --from=builder /app/recruitment-system .
COPY --from=builder /app/.env.example ./.env

# Tạo thư mục logs
RUN mkdir -p ./logs

# Expose cổng
EXPOSE 8080

# Chạy ứng dụng
CMD ["./recruitment-system"]
```

### Docker Compose
```yaml
# docker-compose.yml
version: '3'

services:
  postgres:
    image: postgres:13
    container_name: recruitment_postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: recruitment_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - recruitment-network

  redis:
    image: redis:alpine
    container_name: recruitment_redis
    ports:
      - "6379:6379"
    networks:
      - recruitment-network

  ml-service:
    build:
      context: ./ml-service
      dockerfile: Dockerfile
    container_name: recruitment_ml_service
    ports:
      - "5000:5000"
    networks:
      - recruitment-network

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: recruitment_backend
    depends_on:
      - postgres
      - redis
      - ml-service
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=recruitment_db
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - ML_SERVICE_URL=http://ml-service:5000/recommend
    networks:
      - recruitment-network

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: recruitment_frontend
    depends_on:
      - backend
    ports:
      - "80:80"
    networks:
      - recruitment-network

networks:
  recruitment-network:
    driver: bridge

volumes:
  postgres_data:
```

