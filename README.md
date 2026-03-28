# Orange Blog Backend

一个基于 Go (Gin) 的博客后端项目，提供完整的博客功能 API，包括文章管理、评论系统、用户认证等功能。

# 后端仅提供示例，建议自行根据前端需要的接口实现后端

## 功能特性

### 核心功能
- **文章管理** - 文章的增删改查、分类管理
- **评论系统** - 支持根评论和子评论、评论删除
- **用户认证** - JWT Token 认证、用户注册登录
- **权限管理** - 基于角色的访问控制（管理员/普通用户）
- **暗号系统** - 通过搜索接口获取 Cookie 进行注册
- **邮件通知** - 新文章发布时发送邮件通知
- **限流保护** - 全局限流和严格限流中间件
- **日志记录** - 结构化日志记录

### 技术特性
- **RESTful API** - 标准的 REST 接口设计
- **Swagger 文档** - 自动生成的 API 文档
- **数据库连接池** - 优化的数据库连接管理
- **CORS 支持** - 跨域资源共享配置
- **中间件** - 认证、限流、日志等中间件

## 技术栈

- **语言**: Go 1.x
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **认证**: JWT (golang-jwt/jwt)
- **密码加密**: bcrypt
- **配置管理**: Viper
- **日志**: logrus
- **API 文档**: Swagger
- **限流**: go-redis/redis_rate

## 快速开始

### 环境要求
- Go >= 1.16
- MySQL >= 5.7
- Redis >= 5.0 (用于限流功能)

### 安装依赖

```bash
go mod download
```

### 配置文件

1. 复制配置文件模板：

```bash
cp config/config.yml.example config/config.yml
```

2. 编辑 `config/config.yml`，填入你的配置信息：

```yaml
app:
  name: your_app_name
  port: :8081

database:
  dsn: root:your_password@tcp(127.0.0.1:3306)/my_blog_db?charset=utf8mb4&parseTime=True&loc=Local
  MaxIdleCons: 11
  MaxOpenCons: 114

jwt:
  secret: your_jwt_secret_key_here

mail:
  host: smtp.qq.com
  port: 587
  user: your_email@qq.com
  password: your_email_authorization_code
```

#### 配置说明

**数据库配置 (database)**
- `dsn`: 数据库连接字符串
  - 格式: `username:password@tcp(host:port)/database_name?charset=utf8mb4&parseTime=True&loc=Local`
  - 请将 `your_password` 替换为你的数据库密码
  - 请将 `my_blog_db` 替换为你的数据库名称

**JWT 配置 (jwt)**
- `secret`: JWT 签名密钥
  - 请使用强随机字符串
  - 建议长度至少 32 字符
  - 示例: `your_jwt_secret_key_here`

**邮件配置 (mail)**
- `host`: SMTP 服务器地址
  - QQ 邮箱: `smtp.qq.com`
  - 163 邮箱: `smtp.163.com`
  - Gmail: `smtp.gmail.com`
- `port`: SMTP 端口
  - 通常为 587 (TLS) 或 465 (SSL)
- `user`: 你的邮箱地址
- `password`: 邮箱授权码（不是邮箱密码）
  - QQ 邮箱需要在设置中生成授权码
  - 163 邮箱需要在设置中开启 SMTP 并获取授权码

### 数据库初始化

1. 创建数据库：

```sql
CREATE DATABASE my_blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 运行程序，GORM 会自动创建表结构：

```bash
go run main.go
```

### 运行项目

```bash
# 开发模式
go run main.go

# 编译后运行
go build -o gblog
./gblog
```

服务将在 `http://localhost:8081` 启动

### 访问 API 文档

启动服务后，访问以下地址查看 Swagger API 文档：

```
http://localhost:8081/swagger/index.html
```

## API 接口文档

详细的 API 接口文档请查看 [BACKEND_API.md](./BACKEND_API.md)

### 主要接口

#### 公开接口
- `GET /api/v1/articles` - 获取文章列表
- `GET /api/v1/article/:id` - 获取文章详情
- `GET /api/v1/categories` - 获取分类列表
- `POST /api/v1/search` - 搜索接口（暗号接口）
- `GET /api/v1/check` - 检查认证状态

#### 认证接口
- `POST /api/v1/login` - 用户登录
- `POST /api/v1/register` - 用户注册（需要暗号 Cookie）

#### 需要认证的接口
- `GET /api/v1/comments` - 获取评论列表
- `POST /api/v1/comments` - 发表评论
- `DELETE /api/v1/comments/:id` - 删除评论
- `POST /api/v1/user/username` - 修改用户名
- `POST /api/v1/user/password` - 修改密码

#### 管理员接口
- `POST /api/v1/articles` - 创建文章

## 暗号系统说明

本系统使用暗号机制来控制用户注册权限：

1. **获取暗号 Cookie**：通过搜索接口输入正确的暗号
2. **注册用户**：使用获取到的 Cookie 进行注册
3. **检查权限**：检查接口验证 Cookie 有效性

**重要**：
- 暗号的值和 Cookie 设置需要根据实际代码修改
- 文档中的"暗号"只是占位符
- 实际暗号请查看后端代码中的 `controllers/search.go` 文件

## 项目结构

```
backend/
├── config/              # 配置文件
│   ├── config.go       # 配置结构定义
│   ├── config.yml      # 配置文件（不提交到 Git）
│   ├── config.yml.example  # 配置文件模板
│   ├── db.go           # 数据库连接
│   └── emailConfig.go  # 邮件配置
├── controllers/        # 控制器
│   ├── article.go      # 文章控制器
│   ├── auth.go         # 认证控制器
│   ├── check.go        # 检查控制器
│   ├── comment.go      # 评论控制器
│   ├── getArticleDetail.go  # 获取文章详情
│   ├── getArticles.go  # 获取文章列表
│   ├── getCategories.go  # 获取分类
│   └── search.go       # 搜索控制器（暗号接口）
├── docs/               # Swagger 文档
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── dto/                # 数据传输对象
│   └── comment_dto.go
├── global/             # 全局变量
│   └── global.go
├── logs/               # 日志文件（不提交到 Git）
├── middlewares/        # 中间件
│   ├── auth_middleware.go  # 认证中间件
│   └── rate_limiter.go     # 限流中间件
├── models/             # 数据模型
│   ├── article.go
│   ├── category.go
│   ├── comment.go
│   └── user.go
├── router/             # 路由配置
│   └── router.go
├── templates/          # 邮件模板
│   └── article_email.html
├── utils/              # 工具函数
│   ├── email.go        # 邮件发送
│   ├── jwt.go          # JWT 工具
│   └── logger.go       # 日志工具
├── .gitignore          # Git 忽略文件
├── BACKEND_API.md      # API 接口文档
├── go.mod              # Go 模块文件
├── go.sum              # Go 依赖锁定文件
└── main.go             # 程序入口
```

## 部署指南

### 1. 编译项目

```bash
go build -o gblog
```

### 2. 配置环境

确保服务器上已安装：
- MySQL
- Redis（用于限流功能）

### 3. 配置文件

将 `config/config.yml` 上传到服务器，并修改配置信息。

### 4. 运行服务

```bash
# 直接运行
./gblog

# 使用 nohup 后台运行
nohup ./gblog > app.log 2>&1 &

# 使用 systemd 管理服务（推荐）
sudo systemctl start gblog
sudo systemctl enable gblog
```

### 5. Systemd 配置示例

创建 `/etc/systemd/system/gblog.service`：

```ini
[Unit]
Description=Orange Blog Backend
After=network.target mysql.service redis.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/backend
ExecStart=/path/to/backend/gblog
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

### 6. Nginx 反向代理配置

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location /api/ {
        proxy_pass http://localhost:8081;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### 7. Docker 部署（可选）

创建 `Dockerfile`：

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o gblog

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /app/gblog .
COPY --from=builder /app/config/config.yml.example ./config/config.yml.example

EXPOSE 8081
CMD ["./gblog"]
```

构建和运行：

```bash
docker build -t blog-backend .
docker run -d -p 8081:8081 \
  -e DB_HOST=your_db_host \
  -e DB_PASSWORD=your_db_password \
  blog-backend
```

## 安全建议

1. **修改默认配置**
   - 修改 JWT secret 为强随机字符串
   - 修改数据库密码
   - 修改邮件授权码

2. **环境变量**
   - 敏感信息建议使用环境变量
   - 不要在代码中硬编码密钥

3. **HTTPS**
   - 生产环境必须使用 HTTPS
   - 配置 SSL/TLS 证书

4. **防火墙**
   - 限制数据库端口访问
   - 只允许必要的端口对外开放

5. **定期更新**
   - 定期更新依赖包
   - 及时修复安全漏洞

## 开发建议

### 代码规范
- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释

### 调试技巧
- 使用 Swagger 文档测试接口
- 查看日志文件排查问题
- 使用 Postman 或 curl 测试 API

### 常见问题

**Q: 数据库连接失败？**
A: 检查 `config/config.yml` 中的数据库配置是否正确，确保 MySQL 服务正在运行。

**Q: JWT Token 无效？**
A: 检查 JWT secret 配置是否正确，确保前后端使用相同的 secret。

**Q: 邮件发送失败？**
A: 检查邮件配置，确认 SMTP 服务器地址、端口、邮箱和授权码是否正确。

**Q: 限流过于严格？**
A: 调整 `middlewares/rate_limiter.go` 中的限流参数。

**Q: 注册时提示权限不足？**
A: 需要先通过搜索接口输入正确的暗号获取 Cookie，才能进行注册。

## 许可证

MIT License

## 联系方式

- Email: orange2006cn@foxmail.com
- GitHub: github.com/orangelgj
