# south-admin-gin

基于 Gin + GORM + JWT 的后台管理系统服务端。

## 技术栈

- Go + Gin
- GORM + MySQL
- JWT 认证（golang-jwt）
- bcrypt 密码加密

## 项目结构

```
├── api/v1/system/        # 接口处理
├── config/               # 配置加载
├── database/             # 数据库连接、迁移、初始化
├── middleware/            # JWT 认证中间件
├── models/system/        # 数据模型
├── router/system/        # 路由定义
├── services/system/      # 业务逻辑
├── main.go               # 入口
└── config.yaml           # 配置文件（不提交）
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 创建配置文件

在项目根目录创建 `config.yaml`：

```yaml
mysql:
  path: 127.0.0.1
  port: 3306
  config: charset=utf8mb4&parseTime=True&loc=Local
  db-name: south_admin
  username: root
  password: your_password
  max-idle-conns: 10
  max-open-conns: 100
  conn-max-lifetime: 3600

jwt:
  secret: your-jwt-secret
  expire-hour: 72
```

### 3. 初始化数据库

```bash
mysql -u root -p < database/init.sql
```

### 4. 启动项目

```bash
go run main.go
```

服务默认监听 `http://127.0.0.1:8081`。

## 默认账号

| 用户名 | 密码 |
|--------|------|
| admin  | admin123 |

## 权限模型

用户权限完全通过角色获取：`User → Role → Permission`。

## 前端项目

[react-admin](https://github.com/southliu/south-admin-react)

## 前后端联调

启动 Gin 项目后，将 react-admin 项目的 `.env.development` 改为：

```env
VITE_ENV = "development"
VITE_SERVER_PORT = 7000
VITE_PROXY = [["/api", "http://127.0.0.1:8081/"]]
```
