# Orange Blog Backend API

基于 Gin 框架的博客系统后端 API 接口文档

## 基础信息

- **Base URL**: `http://localhost:8081`
- **API 版本**: v1
- **Content-Type**: `application/json`

## 认证方式

### 1. Cookie 认证（暗号系统）
- 用于注册和检查接口
- Cookie 名称: `暗号`
- Cookie 值: `暗号`
- 获取方式: 通过搜索接口输入暗号获取
- 关于暗号的值以及后续 cookie 设置请自行根据代码替换（懒得找让 trae 之类的 AI ide 帮你找），我的系统的暗号肯定不是这个，暂时没有试过中文能否成功，这只是个占位符

### 2. JWT Token 认证
- 用于需要登录的接口
- Token 在请求头中传递: `Authorization: Bearer {token}`
- Token 通过登录接口获取

## 接口列表

### 公开接口

#### 1. 获取文章列表

**接口地址**: `GET /api/v1/articles`

**接口描述**: 获取文章列表，支持按分类筛选和分页

**查询参数**:
| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| categoryId | int | 否 | 0 | 分类 ID，0 表示全部 |
| page | int | 否 | 1 | 页码 |
| pageSize | int | 否 | 10 | 每页数量 |

**请求示例**:
```http
GET /api/v1/articles?categoryId=0&page=1&pageSize=10 HTTP/1.1
Host: localhost:8081
```

**响应示例**:
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "title": "文章标题",
      "summary": "文章摘要",
      "content": "文章内容",
      "author": "作者",
      "date": "2026-03-29",
      "category": "分类名称"
    }
  ],
  "msg": "获取成功"
}
```

**错误响应**:
```json
{
  "code": 500,
  "data": null,
  "msg": "获取失败,数据库查询问题"
}
```

---

#### 2. 获取文章详情

**接口地址**: `GET /api/v1/article/:id`

**接口描述**: 根据文章 ID 获取文章详细信息

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 文章 ID |

**请求示例**:
```http
GET /api/v1/article/1 HTTP/1.1
Host: localhost:8081
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "title": "文章标题",
    "content": "文章内容",
    "author": "作者",
    "date": "2026-03-29",
    "category": "分类名称"
  },
  "msg": "查询成功"
}
```

**错误响应**:
```json
{
  "code": 400,
  "msg": "ID 格式错误，请输入数字"
}
```

```json
{
  "code": 500,
  "msg": "查询失败或文章不存在"
}
```

---

#### 3. 获取分类列表

**接口地址**: `GET /api/v1/categories`

**接口描述**: 获取所有分类列表

**请求示例**:
```http
GET /api/v1/categories HTTP/1.1
Host: localhost:8081
```

**响应示例**:
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "name": "技术"
    }
  ],
  "msg": "获取成功"
}
```

**错误响应**:
```json
{
  "code": 500,
  "data": null,
  "msg": "获取失败"
}
```

---

#### 4. 搜索接口（暗号接口）

**接口地址**: `POST /api/v1/search`

**接口描述**: 搜索功能接口，当搜索关键词为 "暗号" 时会设置 Cookie

**⚠️ 重要说明**: 这是一个暗号接口，输入正确的暗号后会设置 Cookie，用于注册接口的权限验证

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| query | string | 是 | 搜索关键词（暗号） |

**请求示例**:
```http
POST /api/v1/search HTTP/1.1
Host: localhost:8081
Content-Type: application/json

{
  "query": "暗号"
}
```

**响应示例（暗号正确）**:
```json
{
  "message": "成功啦"
}
```

**响应示例（暗号错误）**:
```json
{
  "message": "功能尚未开发"
}
```

**Cookie 设置信息**:
- 名称: `暗号`
- 值: `暗号`
- 有效期: 36000000 秒（约 416 天）
- 路径: `/`
- 域名: 当前域名
- Secure: `false`
- HttpOnly: `true`

**错误响应**:
```json
{
  "error": "错误信息"
}
```

---

#### 5. 检查认证状态

**接口地址**: `GET /api/v1/check`

**接口描述**: 检查用户认证状态，验证 Cookie 是否有效

**⚠️ 重要说明**: 此接口用于验证 Cookie `暗号` 的值是否为 `暗号`

**请求示例**:
```http
GET /api/v1/check HTTP/1.1
Host: localhost:8081
Cookie: 暗号=暗号
```

**响应示例（认证成功）**:
```json
{
  "auth": true,
  "message": "成功啦"
}
```

**响应示例（认证失败）**:
```json
{
  "auth": false,
  "message": "嘻嘻嘻啦啦啦"
}
```

---

### 认证接口

#### 6. 用户登录

**接口地址**: `POST /api/v1/login`

**接口描述**: 用户登录接口，验证用户名和密码，返回 JWT Token

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**请求示例**:
```http
POST /api/v1/login HTTP/1.1
Host: localhost:8081
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "username": "testuser",
    "userId": 1,
    "role": 2
  }
}
```

**错误响应**:
```json
{
  "code": 400,
  "msg": "参数格式错误"
}
```

```json
{
  "code": 401,
  "msg": "用户名或密码错误"
}
```

```json
{
  "code": 500,
  "msg": "Token 生成失败"
}
```

**角色说明**:
- `role: 1` - 管理员
- `role: 2` - 普通用户

---

#### 7. 用户注册

**接口地址**: `POST /api/v1/register`

**接口描述**: 用户注册接口，创建新用户账号

**⚠️ 重要说明**: 需要先通过搜索接口获取 Cookie `暗号`，否则会返回 403 权限不足

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名（2-20 字符） |
| password | string | 是 | 密码（至少 6 字符） |
| email | string | 是 | 邮箱地址 |
| description | string | 否 | 个人简介（最多 200 字符） |

**请求示例**:
```http
POST /api/v1/register HTTP/1.1
Host: localhost:8081
Content-Type: application/json
Cookie: 暗号=暗号

{
  "username": "newuser",
  "password": "password123",
  "email": "user@example.com",
  "description": "这是我的个人简介"
}
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "欢迎你 newuser 呀！！！嘻嘻！"
}
```

**错误响应**:
```json
{
  "code": 400,
  "msg": "参数校验失败：..."
}
```

```json
{
  "code": 403,
  "msg": "权限不足，无法注册"
}
```

```json
{
  "code": 409,
  "msg": "用户名或邮箱已被注册"
}
```

```json
{
  "code": 500,
  "msg": "密码加密失败"
}
```

---

### 需要认证的接口

以下接口需要在请求头中携带 JWT Token:

```
Authorization: Bearer {token}
```

#### 8. 获取根评论列表

**接口地址**: `GET /api/v1/comments`

**接口描述**: 根据文章 ID 分页获取根评论列表（root_id = 0 的评论）

**查询参数**:
| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| articleId | uint | 是 | - | 文章 ID |
| pageNum | int | 否 | 1 | 页码 |
| pageSize | int | 否 | 10 | 每页数量 |

**请求示例**:
```http
GET /api/v1/comments?articleId=1&pageNum=1&pageSize=10 HTTP/1.1
Host: localhost:8081
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "获取成功",
  "data": [
    {
      "id": 1,
      "content": "这篇文章写得很好！",
      "userId": 1,
      "userName": "张三",
      "createTime": "2024-01-15T10:30:00Z",
      "likeCount": 10,
      "replyCount": 5,
      "ipLocation": "中国 北京",
      "previewReply": {
        "id": 10,
        "content": "感谢支持！",
        "userId": 2,
        "userName": "李四",
        "toUserId": 1,
        "toUserName": "张三",
        "createTime": "2024-01-15T11:00:00Z",
        "likeCount": 2,
        "ipLocation": "中国 上海"
      }
    }
  ],
  "total": 25
}
```

**字段说明**:
- `previewReply`: 预览的第一条回复，无则为 null
- `content`: 可能显示为 "已删除" 或 "审核中"（当 status 为 -1 或 0 时）

**错误响应**:
```json
{
  "code": 400,
  "msg": "缺少 articleId 参数"
}
```

```json
{
  "code": 400,
  "msg": "articleId 格式错误"
}
```

```json
{
  "code": 500,
  "msg": "查询评论失败"
}
```

---

#### 9. 获取子评论列表

**接口地址**: `GET /api/v1/comments/replies`

**接口描述**: 根据根评论 ID 分页获取子评论（回复）列表

**查询参数**:
| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| rootId | uint | 是 | - | 根评论 ID |
| pageNum | int | 否 | 1 | 页码 |
| pageSize | int | 否 | 10 | 每页数量 |

**请求示例**:
```http
GET /api/v1/comments/replies?rootId=1&pageNum=1&pageSize=5 HTTP/1.1
Host: localhost:8081
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "获取成功",
  "data": [
    {
      "id": 10,
      "content": "感谢支持！",
      "userId": 2,
      "userName": "李四",
      "toUserId": 1,
      "toUserName": "张三",
      "createTime": "2024-01-15T11:00:00Z",
      "likeCount": 2,
      "ipLocation": "中国 上海"
    }
  ],
  "total": 5
}
```

**字段说明**:
- `toUserId`: 被回复用户 ID，0 表示无
- `toUserName`: 被回复用户名
- `content`: 可能显示为 "已删除" 或 "审核中"

**错误响应**:
```json
{
  "code": 400,
  "msg": "缺少 rootId 参数"
}
```

```json
{
  "code": 400,
  "msg": "rootId 格式错误"
}
```

```json
{
  "code": 500,
  "msg": "查询回复失败"
}
```

---

#### 10. 发表评论

**接口地址**: `POST /api/v1/comments`

**接口描述**: 发表评论或回复

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| articleId | uint | 是 | 文章 ID |
| content | string | 是 | 评论内容（最多 200 字，不能为空） |
| rootId | uint | 否 | 根评论 ID，发表根评论时为 0 或不传 |
| parentId | uint | 否 | 父评论 ID，回复时可选 |
| toUserId | uint | 否 | 被回复用户 ID，回复时可选 |

**请求示例（发表根评论）**:
```http
POST /api/v1/comments HTTP/1.1
Host: localhost:8081
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "articleId": 1,
  "content": "这篇文章写得很好！"
}
```

**请求示例（回复评论）**:
```http
POST /api/v1/comments HTTP/1.1
Host: localhost:8081
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "articleId": 1,
  "content": "感谢支持！",
  "rootId": 5,
  "parentId": 5,
  "toUserId": 3
}
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "评论发表成功",
  "data": {
    "id": 10,
    "content": "这篇文章写得很好！",
    "userId": 1,
    "userName": "张三",
    "toUserId": 0,
    "toUserName": "",
    "createTime": "2026-03-18T10:30:00Z",
    "likeCount": 0,
    "ipLocation": "中国 北京"
  }
}
```

**错误响应**:
```json
{
  "code": 400,
  "msg": "参数错误：..."
}
```

```json
{
  "code": 400,
  "msg": "评论内容不能为空"
}
```

```json
{
  "code": 400,
  "msg": "评论内容不能超过200字"
}
```

```json
{
  "code": 401,
  "msg": "用户未登录"
}
```

```json
{
  "code": 500,
  "msg": "评论发表失败"
}
```

---

#### 11. 删除评论

**接口地址**: `DELETE /api/v1/comments/:id`

**接口描述**: 删除评论（只能删除自己的评论，管理员可以删除所有评论）

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | uint | 是 | 评论 ID |

**请求示例**:
```http
DELETE /api/v1/comments/10 HTTP/1.1
Host: localhost:8081
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "删除成功"
}
```

**错误响应**:
```json
{
  "code": 400,
  "msg": "缺少评论ID"
}
```

```json
{
  "code": 400,
  "msg": "评论ID格式错误"
}
```

```json
{
  "code": 401,
  "msg": "用户未登录"
}
```

```json
{
  "code": 403,
  "msg": "权限不足"
}
```

```json
{
  "code": 404,
  "msg": "评论不存在"
}
```

```json
{
  "code": 500,
  "msg": "删除评论失败"
}
```

---

#### 12. 修改用户名

**接口地址**: `POST /api/v1/user/username`

**接口描述**: 修改用户名，需要验证旧密码

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| old_password | string | 是 | 旧密码 |
| new_username | string | 是 | 新用户名（2-20 字符） |

**请求示例**:
```http
POST /api/v1/user/username HTTP/1.1
Host: localhost:8081
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "old_password": "oldpassword123",
  "new_username": "newusername"
}
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "用户名修改成功，请重新登录"
}
```

**错误响应**:
```json
{
  "code": 400,
  "msg": "参数校验失败：..."
}
```

```json
{
  "code": 401,
  "msg": "用户未登录"
}
```

```json
{
  "code": 401,
  "msg": "旧密码错误"
}
```

```json
{
  "code": 409,
  "msg": "用户名已被占用"
}
```

```json
{
  "code": 500,
  "msg": "用户名更新失败"
}
```

---

#### 13. 修改密码

**接口地址**: `POST /api/v1/user/password`

**接口描述**: 修改密码，需要验证旧密码

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| old_password | string | 是 | 旧密码 |
| new_password | string | 是 | 新密码（至少 6 字符） |

**请求示例**:
```http
POST /api/v1/user/password HTTP/1.1
Host: localhost:8081
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "old_password": "oldpassword123",
  "new_password": "newpassword123"
}
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "密码修改成功，请重新登录"
}
```

**错误响应**:
```json
{
  "code": 400,
  "msg": "参数校验失败：..."
}
```

```json
{
  "code": 401,
  "msg": "用户未登录"
}
```

```json
{
  "code": 401,
  "msg": "旧密码错误"
}
```

```json
{
  "code": 500,
  "msg": "密码加密失败"
}
```

```json
{
  "code": 500,
  "msg": "密码更新失败"
}
```

---

### 管理员接口

以下接口需要 JWT Token 认证 + 管理员权限（role = 1）

#### 14. 创建文章

**接口地址**: `POST /api/v1/articles`

**接口描述**: 创建新文章（仅管理员可用）

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| title | string | 是 | 文章标题 |
| author | string | 是 | 作者 |
| category_id | uint | 是 | 分类 ID |
| summary | string | 否 | 文章摘要 |
| content | string | 是 | 文章内容 |

**请求示例**:
```http
POST /api/v1/articles HTTP/1.1
Host: localhost:8081
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "title": "我的第一篇博客",
  "author": "Orange",
  "category_id": 1,
  "summary": "这是文章摘要",
  "content": "这是文章内容"
}
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "文章创建成功"
}
```

**错误响应**:
```json
{
  "code": 400,
  "msg": "参数格式错误"
}
```

```json
{
  "code": 401,
  "msg": "未认证"
}
```

```json
{
  "code": 403,
  "msg": "权限不足"
}
```

```json
{
  "code": 500,
  "msg": "文章创建失败"
}
```

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 / Token 无效 / 密码错误 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突（如用户名已存在） |
| 500 | 服务器内部错误 |

---

## CORS 配置

- 允许的源: `http://localhost:5173`
- 允许的方法: `GET`, `POST`, `OPTIONS`, `DELETE`
- 允许的请求头: `Origin`, `Content-Type`, `Authorization`
- 允许携带凭证: `true`
- 预检请求缓存时间: 12 小时

---

## 限流配置

- **全局限流**: 所有接口都有基础限流
- **严格限流**: 登录和注册接口使用更严格的限流策略

---

## IP 归属地查询

评论接口会自动查询并记录评论者的 IP 归属地信息：
- 本地 IP（127.0.0.1、内网 IP）显示为 "本地"
- 公网 IP 通过 ip-api.com 查询归属地
- 查询超时时间: 2 秒

---

## 邮件通知

创建文章后，系统会自动向所有非管理员用户发送邮件通知：
- 邮件内容包含文章标题和链接
- 发送失败的邮件会被记录到日志

---

## 日志记录

以下操作会被记录到日志：
- 用户登录
- 用户注册
- 发表评论
- 删除评论

日志包含以下信息：
- 用户 ID 和用户名
- 操作类型
- IP 地址
- 时间戳

---

## Swagger 文档

访问 `http://localhost:8081/swagger/index.html` 查看在线 API 文档

---

## 注意事项

1. **注册流程**: 必须先调用搜索接口输入暗号 "暗号" 获取 Cookie，才能调用注册接口
2. **Token 有效期**: Token 有效期与用户密码更新时间相关，修改密码后旧 Token 会失效
3. **评论审核**: 评论状态字段 `status`：1=正常，0=审核中，-1=已删除
4. **密码加密**: 所有密码使用 bcrypt 加密存储
5. **时间格式**: 所有时间使用 RFC3339 格式（如：2024-01-15T10:30:00Z）
