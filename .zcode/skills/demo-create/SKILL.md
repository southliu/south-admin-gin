---
name: demo-create
description: 生成或修改增删改查(CRUD)模块时必须使用本skill。当用户要求新增业务模块、生成增删改查接口、创建xxx管理功能、或修改现有CRUD代码时，按本skill规定的文件结构、命名和接口契约生成，保证与项目现有格式一致。适用范围：新建 Model/Service/Handler/Router、批量删除、状态切换、分页过滤等。
---

# demo-create：gin 项目 CRUD 模块生成规范

为 `south-admin-gin`（Gin + GORM + MySQL）生成增删改查模块时，**必须**按本文件的结构和契约执行。下面以模块名 `{{name}}`（如 `article`）为例。

## 一、要创建/修改的文件（共 5 处）

| 文件 | 作用 |
|---|---|
| `models/system/{{name}}.go` | 数据模型 + 请求 DTO |
| `services/system/{{name}}_service.go` | 业务逻辑（分页/详情/增/改/删/批量删） |
| `api/v1/system/{{name}}.go` | HTTP Handler |
| `router/system/{{name}}.go` | 路由注册 |
| `router/router.go` + `database/init.sql` | 挂载路由 + 菜单权限种子数据 |

## 二、接口契约（与 react-admin 前端对齐，违反即为 bug）

1. **响应结构**：一律 `models.CommonResponse{Code, Message, Data}`，HTTP 状态码恒为 200，业务码放 `Code`（200 成功 / 400 参数 / 404 不存在 / 500 失败）。
2. **分页参数名是 `pageSize`**（不是 `page_size`）：
   ```go
   page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
   pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", c.DefaultQuery("page_size", "10")))
   ```
3. **分页返回**：`{items, page, pageSize, total, totalPages}`（用 `models.PageResult`）。
4. **JSON 字段驼峰**：`createdAt/updatedAt/parentId/permissionId/roleIds/labelEn`。禁止 `created_at/role_ids` 这类下划线键名（前端读不到）。
5. **int 零值合法**：`state=0`(隐藏)、`status=0`(禁用) 是有效值，DTO 上**禁止**给这类字段加 `binding:"required"`（Go 的 required 拒绝零值）。
6. **更新接口直接赋值**：前端编辑表单总是提交完整数据，字段直接 `entity.Field = dto.Field`，不要写 `if dto.Field != 0` 守卫（否则永远无法置零/清空）。
7. **接口面**（挂在 `withAuth` 下，登录/特殊接口除外）：
   `GET /page`、`GET /detail?id=`、`POST /create`、`PUT /update/:id`、`DELETE /:id`、`POST /batchDelete`、（有状态字段时）`PUT /changeState`、`GET /list`（下拉用，返回精简字段）。
8. **过滤参数**：列表接口支持模糊过滤（`LIKE %xx%`），参数名与前端搜索框 `name` 一致；total 计数必须跟随同一过滤条件。
9. **不泄露密码**：用户类响应不得包含 password（Model 上用 `json:"-"`）。
10. **软删除**：`is_deleted=1` + `deleted_at`（`models.CustomTime(time.Now().Unix())`），删除前检查关联（如子记录）。

## 三、文件模板

### 1. `models/system/{{name}}.go`

```go
package models

// {{Name}} 模型
type {{Name}} struct {
	ID        int64      `json:"id" gorm:"primaryKey"`
	// TODO: 业务字段，json 一律驼峰；字符串列写 gorm type，如 gorm:"type:varchar(100)"
	Status    int        `json:"status" gorm:"default:1;comment:状态 1=启用 0=禁用"`
	CreatedAt CustomTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt CustomTime `json:"updatedAt" gorm:"autoUpdateTime"`
	IsDeleted int        `json:"isDeleted" gorm:"default:0"`
	DeletedAt CustomTime `json:"-"`

	// 关联（多对多必须显式 through 表名时参照 user.go 的写法）
}

func ({{Name}}) TableName() string { return "{{table_name}}" }

// Create{{Name}}Dto 创建请求
type Create{{Name}}Dto struct {
	// TODO: 业务字段
}

// Update{{Name}}Dto 更新请求（字段不设 required，支持置零/清空）
type Update{{Name}}Dto struct {
	// TODO: 业务字段
}
```

**时间字段规则**：模型时间一律用 `CustomTime`（底层 int64 时间戳）。**绝不允许**修改 `models/system/time.go` 中 `Value()` 返回非 int64 的类型（表列是 BIGINT，返回 time.Time 会触发 Error 1265 拦截所有写操作）。

### 2. `services/system/{{name}}_service.go`

参照 `services/system/user_service.go` / `role_service.go`：

- `Get{{Name}}Page(page, pageSize int, filter ...)` → `*models.PageResult`，过滤条件同时作用于 Count 和 Find
- `Get{{Name}}ByID(id)` → 单条（含 Preload）
- `Create{{name}}(dto)` → 唯一性校验 → `database.DB.Create`
- `Update{{name}}(id, dto)` → 直接赋值 → `database.DB.Save`
- `Delete{{name}}(id)` → 软删（`Updates(map[string]interface{}{"is_deleted":1, "deleted_at": ...})`）
- `BatchDelete{{Name}}(ids)` → `Where("id IN ?", ids)` 批量软删；空 ids 报"请选择要删除的xx"

### 3. `api/v1/system/{{name}}.go`

参照 `api/v1/system/user.go`：每个 Handler 做参数解析 → 调 service → 按第二节契约返回 `models.CommonResponse`。路径参数 `strconv.ParseInt(c.Param("id"), 10, 64)`。

### 4. `router/system/{{name}}.go`

```go
package system

import (
	v1 "south-admin-gin/api/v1/system"
	"github.com/gin-gonic/gin"
)

func Init{{Name}}Routes(r *gin.RouterGroup) {
	auth := withAuth(r)
	auth.GET("/page", v1.Get{{Name}}Page)
	auth.GET("/detail", v1.Get{{Name}}Detail)
	auth.POST("/create", v1.Create{{Name}})
	auth.PUT("/update/:id", v1.Update{{Name}})
	auth.DELETE("/:id", v1.Delete{{Name}})
	auth.POST("/batchDelete", v1.BatchDelete{{Name}})
	auth.GET("/list", v1.Get{{Name}}List)
}
```

### 5. 挂载与种子数据

- `router/router.go`：`system.Init{{Name}}Routes(r.Group("/system"))`
- `database/init.sql`：追加权限（`/{{authority_prefix}}/{{name}}` 及 `/index|create|update|view|delete` 按钮）、菜单（type=2 页面菜单 + type=3 按钮菜单，参照"日志管理"三件套的写法）、`role_menu` 授权给 admin 角色。

## 四、完成前自查

- [ ] `go build ./...` 通过
- [ ] curl 验证：分页返回 `pageSize` 为数字、过滤生效、create/update/delete/batchDelete 均 200 且无 1265
- [ ] 新增列与 DTO 的 json 键全部驼峰
