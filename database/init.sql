-- 开始事务
START TRANSACTION;

-- 禁用外键约束检查
SET FOREIGN_KEY_CHECKS = 0;

-- 删除现有数据
DELETE FROM `user_role`;
DELETE FROM `role_permission`;
DELETE FROM `role_menu`;
DELETE FROM `sys_menu`;
DELETE FROM `sys_permission`;
DELETE FROM `sys_role`;
DELETE FROM `sys_user`;

-- 重新启用外键约束检查
SET FOREIGN_KEY_CHECKS = 1;

-- 插入用户表数据（密码已加密，明文为 admin123）
INSERT INTO `sys_user` (username, password, name, email, status, is_deleted, created_at, updated_at)
VALUES
    ('admin', '$2b$10$07h7npcIysHutrLYCY3yWOhEqtGTCR88pDp66ZztkAdG7RJT/4ZDO', '超级管理员', 'admin@example.com', 1, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
    ('user1', '$2b$10$07h7npcIysHutrLYCY3yWOhEqtGTCR88pDp66ZztkAdG7RJT/4ZDO', '普通用户', 'user1@example.com', 1, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 插入角色表数据
INSERT INTO `sys_role` (name, description, created_at, updated_at, is_deleted)
VALUES
    ('admin', '超级管理员', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('user', '普通用户', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0);

-- 插入权限表数据
INSERT INTO `sys_permission` (name, description, created_at, updated_at, is_deleted)
VALUES
    ('/dashboard', '查看仪表盘', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/demo', '查看示例菜单', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/demo/copy', '复制菜单', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/demo/editor', '编辑示例菜单', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/demo/wangEditor', 'WangEditor 示例', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/demo/virtualScroll', '虚拟滚动示例', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/demo/watermark', '水印示例', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/user', '用户管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/user/index', '用户列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/user/create', '创建用户', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/user/update', '修改用户', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/user/view', '查看用户', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/user/delete', '删除用户', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/user/authority', '用户权限配置', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/role', '角色管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/role/index', '角色列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/role/create', '创建角色', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/role/update', '修改角色', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/role/view', '查看角色', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/role/delete', '删除角色', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/menu', '菜单管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/menu/index', '菜单列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/menu/create', '创建菜单', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/menu/update', '修改菜单', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/menu/view', '查看菜单', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/authority/menu/delete', '删除菜单', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/content/article', '文章管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/content/article/index', '文章列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/content/article/create', '创建文章', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/content/article/update', '修改文章', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/content/article/view', '查看文章', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/content/article/delete', '删除文章', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
    ('/link', '外部链接', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0);

-- 关联用户与角色
INSERT INTO `user_role` (user_id, role_id)
VALUES
    ((SELECT id FROM `sys_user` WHERE username='admin'), (SELECT id FROM `sys_role` WHERE name='admin')),
    ((SELECT id FROM `sys_user` WHERE username='user1'), (SELECT id FROM `sys_role` WHERE name='user'));

-- 关联角色与权限
INSERT INTO `role_permission` (role_id, permission_id)
VALUES
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/dashboard')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/demo')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/demo/copy')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/demo/editor')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/demo/wangEditor')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/demo/virtualScroll')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/demo/watermark')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/user')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/user/index')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/user/create')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/user/update')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/user/view')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/user/delete')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/user/authority')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/role')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/role/index')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/role/create')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/role/update')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/role/view')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/role/delete')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/menu')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/menu/index')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/menu/create')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/menu/update')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/menu/view')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/authority/menu/delete')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/content/article')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/content/article/index')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/content/article/create')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/content/article/update')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/content/article/view')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_permission` WHERE name='/content/article/delete'));

-- 禁用外键约束检查
SET FOREIGN_KEY_CHECKS = 0;

-- 删除现有菜单数据
DELETE FROM `sys_menu` WHERE router IN (
    '/dashboard',
    '/demo',
    '/demo/copy',
    '/demo/watermark',
    '/demo/virtualScroll',
    '/demo/editor',
    '/demo/123/dynamic',
    '/demo/level1',
    '/demo/level1/level2',
    '/demo/level1/level2/level3',
    '/system',
    '/system/user',
    '/system/menu',
    '/system/role',
    '/content',
    '/content/article',
    'https://ant-design.antgroup.com'
);

-- 重新启用外键约束检查
SET FOREIGN_KEY_CHECKS = 1;

-- 插入顶级菜单项
INSERT INTO `sys_menu` (label, label_en, type, icon, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id) VALUES
('仪表盘', 'Dashboard', 2, 'la:tachometer-alt', '/dashboard', 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), NULL, 0, (SELECT id FROM `sys_permission` WHERE name = '/dashboard')),
('组件', 'Components', 1, 'fluent:box-20-regular', '/demo', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), NULL, 0, NULL),
('系统管理', 'System Management', 1, 'ion:settings-outline', '/system', 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), NULL, 0, NULL),
('内容管理', 'Content Management', 1, 'majesticons:article-search-line', '/content', 3, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), NULL, 0, NULL),
('外部链接', 'External Link', 2, 'material-symbols:link', 'https://ant-design.antgroup.com', 4, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), NULL, 0, (SELECT id FROM `sys_permission` WHERE name = '/link'));

-- 插入组件子菜单
INSERT INTO `sys_menu` (label, label_en, type, icon, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '剪切板', 'Copy', 2, NULL, '/demo/copy', 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/demo/copy')
FROM (SELECT id FROM `sys_menu` WHERE router = '/demo') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, icon, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '水印', 'Watermark', 2, NULL, '/demo/watermark', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/demo/watermark')
FROM (SELECT id FROM `sys_menu` WHERE router = '/demo') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, icon, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '虚拟滚动', 'Virtual Scroll', 2, NULL, '/demo/virtualScroll', 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/demo/virtualScroll')
FROM (SELECT id FROM `sys_menu` WHERE router = '/demo') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, icon, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '富文本', 'Editor', 2, NULL, '/demo/editor', 3, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/demo/editor')
FROM (SELECT id FROM `sys_menu` WHERE router = '/demo') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, icon, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '动态路由参数', 'Dynamic', 2, NULL, '/demo/123/dynamic', 4, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, NULL
FROM (SELECT id FROM `sys_menu` WHERE router = '/demo') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, icon, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '层级1', 'Level1', 1, NULL, '/demo/level1', 5, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, NULL
FROM (SELECT id FROM `sys_menu` WHERE router = '/demo') AS parent_menu;

-- 插入层级子菜单
INSERT INTO `sys_menu` (label, label_en, type, icon, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '层级2', 'Level2', 1, NULL, '/demo/level1/level2', 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, NULL
FROM (SELECT id FROM `sys_menu` WHERE router = '/demo/level1') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, icon, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '层级3', 'Level3', 2, NULL, '/demo/level1/level2/level3', 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/demo/watermark')
FROM (SELECT id FROM `sys_menu` WHERE router = '/demo/level1/level2') AS parent_menu;

-- 插入系统管理子菜单
INSERT INTO `sys_menu` (label, label_en, type, icon, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '用户管理', 'User Management', 2, NULL, '/system/user', 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/user')
FROM (SELECT id FROM `sys_menu` WHERE router = '/system') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, icon, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '菜单管理', 'Menu Management', 2, NULL, '/system/menu', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/menu')
FROM (SELECT id FROM `sys_menu` WHERE router = '/system') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, icon, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '角色管理', 'Role Management', 2, NULL, '/system/role', 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/role')
FROM (SELECT id FROM `sys_menu` WHERE router = '/system') AS parent_menu;

-- 插入内容管理子菜单
INSERT INTO `sys_menu` (label, label_en, type, icon, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '文章管理', 'Article Management', 2, NULL, '/content/article', 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/content/article')
FROM (SELECT id FROM `sys_menu` WHERE router = '/content') AS parent_menu;

-- 插入用户管理按钮菜单
INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '用户列表', 'Index', 3, '/system/user', 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/user/index')
FROM (SELECT id FROM `sys_menu` WHERE label = '用户管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '查看用户', 'View', 3, '/system/user', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/user/view')
FROM (SELECT id FROM `sys_menu` WHERE label = '用户管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '创建用户', 'Create', 3, '/system/user', 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/user/create')
FROM (SELECT id FROM `sys_menu` WHERE label = '用户管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '修改用户', 'Update', 3, '/system/user', 3, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/user/update')
FROM (SELECT id FROM `sys_menu` WHERE label = '用户管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '删除用户', 'Delete', 3, '/system/user', 4, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/user/delete')
FROM (SELECT id FROM `sys_menu` WHERE label = '用户管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '用户权限按钮', 'Authority', 3, '/system/user', 5, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/user/authority')
FROM (SELECT id FROM `sys_menu` WHERE label = '用户管理') AS parent_menu;

-- 插入菜单管理按钮菜单
INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '菜单列表', 'Index', 3, '/system/menu', 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/menu/index')
FROM (SELECT id FROM `sys_menu` WHERE label = '菜单管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '查看菜单', 'View', 3, '/system/menu', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/menu/view')
FROM (SELECT id FROM `sys_menu` WHERE label = '菜单管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '创建菜单', 'Create', 3, '/system/menu', 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/menu/create')
FROM (SELECT id FROM `sys_menu` WHERE label = '菜单管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '修改菜单', 'Update', 3, '/system/menu', 3, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/menu/update')
FROM (SELECT id FROM `sys_menu` WHERE label = '菜单管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '删除菜单', 'Delete', 3, '/system/menu', 4, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/menu/delete')
FROM (SELECT id FROM `sys_menu` WHERE label = '菜单管理') AS parent_menu;

-- 插入角色管理按钮菜单
INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '角色列表', 'Index', 3, '/system/role', 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/role/index')
FROM (SELECT id FROM `sys_menu` WHERE label = '角色管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '查看角色', 'View', 3, '/system/role', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/role/view')
FROM (SELECT id FROM `sys_menu` WHERE label = '角色管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '创建角色', 'Create', 3, '/system/role', 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/role/create')
FROM (SELECT id FROM `sys_menu` WHERE label = '角色管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '修改角色', 'Update', 3, '/system/role', 3, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/role/update')
FROM (SELECT id FROM `sys_menu` WHERE label = '角色管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '删除角色', 'Delete', 3, '/system/role', 4, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/authority/role/delete')
FROM (SELECT id FROM `sys_menu` WHERE label = '角色管理') AS parent_menu;

-- 插入文章管理按钮菜单
INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '文章列表', 'Index', 3, '/content/article', 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/content/article/index')
FROM (SELECT id FROM `sys_menu` WHERE label = '文章管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '查看文章', 'View', 3, '/content/article', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/content/article/view')
FROM (SELECT id FROM `sys_menu` WHERE label = '文章管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '创建文章', 'Create', 3, '/content/article', 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/content/article/create')
FROM (SELECT id FROM `sys_menu` WHERE label = '文章管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '修改文章', 'Update', 3, '/content/article', 3, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/content/article/update')
FROM (SELECT id FROM `sys_menu` WHERE label = '文章管理') AS parent_menu;

INSERT INTO `sys_menu` (label, label_en, type, router, `order`, state, created_at, updated_at, parent_id, is_deleted, permission_id)
SELECT '删除文章', 'Delete', 3, '/content/article', 4, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), parent_menu.id, 0, (SELECT id FROM `sys_permission` WHERE name = '/content/article/delete')
FROM (SELECT id FROM `sys_menu` WHERE label = '文章管理') AS parent_menu;

-- 关联角色与菜单（顶级菜单）
INSERT INTO `role_menu` (role_id, menu_id)
VALUES
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='仪表盘')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='组件')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='系统管理')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='内容管理')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='外部链接')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='剪切板')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='水印')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='虚拟滚动')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='富文本')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='动态路由参数')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='层级1')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='层级2')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='层级3')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='用户管理')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='菜单管理')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='角色管理')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='文章管理'));

-- 关联角色与菜单（按钮菜单）
INSERT INTO `role_menu` (role_id, menu_id)
VALUES
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='用户列表')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='查看用户')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='创建用户')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='修改用户')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='删除用户')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='用户权限按钮')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='菜单列表')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='查看菜单')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='创建菜单')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='修改菜单')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='删除菜单')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='角色列表')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='查看角色')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='创建角色')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='修改角色')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='删除角色')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='文章列表')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='查看文章')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='创建文章')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='修改文章')),
    ((SELECT id FROM `sys_role` WHERE name='admin'), (SELECT id FROM `sys_menu` WHERE label='删除文章'));

-- 提交事务
COMMIT;
