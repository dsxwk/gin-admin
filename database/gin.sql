/*
 Navicat Premium Data Transfer

 Source Server         : 127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 80012
 Source Host           : 127.0.0.1:3306
 Source Schema         : gin

 Target Server Type    : MySQL
 Target Server Version : 80012
 File Encoding         : 65001

 Date: 09/07/2026 11:23:25
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for article
-- ----------------------------
DROP TABLE IF EXISTS `article`;
CREATE TABLE `article`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` int(11) NOT NULL DEFAULT 0 COMMENT '用户id',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标题',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '内容',
  `category_id` int(11) NOT NULL DEFAULT 0 COMMENT '分类id',
  `data_source` tinyint(3) UNSIGNED NOT NULL DEFAULT 2 COMMENT '数据来源 1=文章库 2=自建',
  `is_publish` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '是否发布 0=待发布 1=已发布 2=已下架',
  `tag` json NULL COMMENT '标签',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文章表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of article
-- ----------------------------
INSERT INTO `article` VALUES (1, 1, '标题1', '<p>测试1</p>', 1, 2, 1, '[\"测试标签1\", \"测试标签2\", \"cs3\"]', '2023-09-19 11:43:58', '2006-01-02 15:04:05', NULL);
INSERT INTO `article` VALUES (13, 1, '标题1', '<p>内容13</p>', 0, 2, 1, '[\"测试标签11\", \"测试标签22\"]', '2024-07-22 11:21:18', '2025-06-17 10:27:27', NULL);
INSERT INTO `article` VALUES (14, 1, 'Go语言', '内容', 0, 0, 2, 'null', '2025-07-03 17:32:08', '2025-07-03 17:32:08', NULL);

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `pid` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父级id',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '分类表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of category
-- ----------------------------
INSERT INTO `category` VALUES (1, 0, '分类名称', '2023-09-19 11:43:43', '2023-09-19 11:43:43', NULL);

-- ----------------------------
-- Table structure for config_category
-- ----------------------------
DROP TABLE IF EXISTS `config_category`;
CREATE TABLE `config_category`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '配置分类表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of config_category
-- ----------------------------
INSERT INTO `config_category` VALUES (1, '基本信息', '2026-07-08 13:42:00', '2006-01-02 15:04:05', NULL);
INSERT INTO `config_category` VALUES (2, '邮箱配置', '2026-07-08 13:42:00', '2006-01-02 15:04:05', NULL);
INSERT INTO `config_category` VALUES (3, 'seo设置 ', '2026-07-08 13:42:00', '2026-07-08 13:42:00', NULL);

-- ----------------------------
-- Table structure for dict
-- ----------------------------
DROP TABLE IF EXISTS `dict`;
CREATE TABLE `dict`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `pid` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父级id',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标识',
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '名称',
  `value` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '映射值',
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态 1=启用 2=停用',
  `sort` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序',
  `extend` json NULL COMMENT '扩展字段',
  `desc` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '字段描述',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_pid`(`pid`) USING BTREE,
  INDEX `idx_name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '字典表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of dict
-- ----------------------------
INSERT INTO `dict` VALUES (1, 0, 'gender', '性别', '', 1, 0, '{\"test\": 111, \"test2\": \"test222\"}', '', '2025-06-06 21:48:17', '2025-06-06 21:48:17', NULL);
INSERT INTO `dict` VALUES (2, 1, 'gender', '男', '1', 1, 0, '{\"a\": \"111\", \"b\": \"222\"}', '测试', '2025-06-06 21:49:00', '2006-01-02 15:04:05', NULL);
INSERT INTO `dict` VALUES (3, 1, 'gender', '女', '2', 1, 0, '{\"test\": 111, \"test2\": \"test222\"}', '性别女', '2025-06-06 21:49:10', '2025-06-08 20:39:03', NULL);
INSERT INTO `dict` VALUES (4, 1, 'gender', '保密', '0', 1, 0, '{}', '保密', '2026-07-03 13:43:30', '2006-01-02 15:04:05', NULL);

-- ----------------------------
-- Table structure for menu
-- ----------------------------
DROP TABLE IF EXISTS `menu`;
CREATE TABLE `menu`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `pid` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父级id',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '路由名称',
  `path` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '路由路径',
  `redirect` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '重定向',
  `component` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '组件路径',
  `is_link` tinyint(3) UNSIGNED NOT NULL DEFAULT 2 COMMENT '是否外链 1=是 2=否 默认=2',
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态 1=启用 2=停用',
  `sort` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_pid`(`pid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of menu
-- ----------------------------
INSERT INTO `menu` VALUES (1, 0, 'home', '/home', '', 'home/index', 2, 1, 1, '2025-05-23 15:37:03', '2025-06-13 11:10:18', NULL);
INSERT INTO `menu` VALUES (2, 0, 'system', '/system', '/system/menu', 'layouts/routerView/parent', 2, 1, 2, '2025-05-23 15:39:37', '2025-05-27 16:49:52', NULL);
INSERT INTO `menu` VALUES (3, 2, 'systemMenu', '/system/menu', '', 'system/menu/index', 2, 1, 3, '2025-05-23 15:41:38', '2025-06-11 17:17:14', NULL);
INSERT INTO `menu` VALUES (4, 2, 'systemUser', '/system/user', '', 'system/user/index', 2, 1, 4, '2025-05-23 23:26:38', '2025-06-11 17:17:29', NULL);
INSERT INTO `menu` VALUES (5, 2, 'systemRole', '/system/role', '', 'system/role/index', 2, 1, 5, '2025-05-25 14:37:04', '2025-06-11 17:17:36', NULL);
INSERT INTO `menu` VALUES (6, 2, 'systemDic', '/system/dic', '', 'system/dic/index', 2, 1, 6, '2025-05-25 14:54:04', '2025-06-11 17:17:42', NULL);
INSERT INTO `menu` VALUES (10, 0, 'article', '/article', '', 'article/index', 2, 1, 7, '2025-06-16 15:34:11', '2025-06-16 15:34:11', NULL);
INSERT INTO `menu` VALUES (20, 2, 'systemConfig', '/system/config', '', 'layouts/routerView/parent', 2, 1, 7, '2026-07-08 15:04:29', '2006-01-02 15:04:05', NULL);
INSERT INTO `menu` VALUES (21, 20, 'systemConfigList', '/system/config/index', '', 'system/config/index', 2, 1, 7, '2026-07-08 15:51:11', '2026-07-08 15:51:11', NULL);
INSERT INTO `menu` VALUES (22, 20, 'systemConfigSetting', '/system/config/setting', '', 'system/config/setting', 2, 1, 8, '2026-07-08 15:54:52', '2006-01-02 15:04:05', NULL);
INSERT INTO `menu` VALUES (23, 20, 'systemConfigCategory', '/system/config-category/index', '', 'system/config/category/index', 2, 1, 1, '2026-07-09 10:56:40', '2026-07-09 10:56:40', NULL);

-- ----------------------------
-- Table structure for menu_actions
-- ----------------------------
DROP TABLE IF EXISTS `menu_actions`;
CREATE TABLE `menu_actions`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `pid` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父级id',
  `menu_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '菜单id',
  `type` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '类型 1=header 2=operation',
  `btn_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'btn' COMMENT '按钮类型 text|btn',
  `btn_style` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'primary' COMMENT '按钮样式',
  `btn_size` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'small' COMMENT '按钮尺寸',
  `is_confirm` tinyint(3) UNSIGNED NOT NULL DEFAULT 2 COMMENT '是否确认 1=是 2=否',
  `label` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '功能名称',
  `auth_value` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '权限标识',
  `is_link` tinyint(3) UNSIGNED NOT NULL DEFAULT 2 COMMENT '是否为链接 1=是 2=否',
  `sort` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_menu_id`(`menu_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 25 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单功能表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of menu_actions
-- ----------------------------
INSERT INTO `menu_actions` VALUES (1, 0, 3, 1, 'btn', 'primary', 'small', 2, '新增菜单', 'sys.menu.add', 2, 0, '2025-05-21 10:24:14', '2025-06-13 16:26:48', NULL);
INSERT INTO `menu_actions` VALUES (2, 0, 3, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.menu.edit', 2, 0, '2025-05-21 10:30:24', '2025-06-13 17:17:32', NULL);
INSERT INTO `menu_actions` VALUES (3, 0, 3, 2, 'btn', 'primary', 'small', 2, '功能', 'sys.menu.action', 2, 0, '2025-05-21 10:30:37', '2025-06-13 16:27:43', NULL);
INSERT INTO `menu_actions` VALUES (4, 0, 3, 2, 'btn', 'danger', 'small', 1, '删除', 'sys.menu.del', 2, 0, '2025-05-21 10:30:49', '2025-06-13 16:27:51', NULL);
INSERT INTO `menu_actions` VALUES (6, 3, 3, 1, 'btn', 'primary', 'small', 2, '新增功能', 'sys.menu.action.add', 2, 0, '2025-06-11 11:46:16', '2025-06-13 16:28:07', NULL);
INSERT INTO `menu_actions` VALUES (7, 3, 3, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.menu.action.edit', 2, 0, '2025-06-13 11:36:56', '2025-06-13 16:28:14', NULL);
INSERT INTO `menu_actions` VALUES (8, 3, 3, 2, 'btn', 'danger', 'small', 2, '删除', 'sys.menu.action.del', 2, 0, '2025-06-13 11:37:07', '2025-06-13 16:28:21', NULL);
INSERT INTO `menu_actions` VALUES (9, 0, 4, 1, 'btn', 'primary', 'small', 2, '新增用户', 'sys.user.add', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (10, 0, 4, 1, 'btn', 'danger', 'small', 2, '批量删除', 'sys.user.batchDel', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (11, 0, 4, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.user.edit', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (12, 0, 4, 2, 'btn', 'danger', 'small', 2, '删除', 'sys.user.del', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (13, 0, 5, 1, 'btn', 'danger', 'small', 2, '批量删除', 'sys.role.batchDel', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (14, 0, 5, 1, 'btn', 'primary', 'small', 2, '新增角色', 'sys.role.add', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (15, 0, 5, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.role.edit', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (16, 0, 5, 2, 'btn', 'danger', 'small', 2, '删除', 'sys.role.del', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (17, 0, 6, 1, 'btn', 'primary', 'small', 2, '新增字典', 'sys.dic.add', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (18, 0, 6, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.dic.edit', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (19, 0, 6, 2, 'btn', 'danger', 'small', 2, '删除', 'sys.dic.del', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (20, 0, 10, 1, 'btn', 'primary', 'small', 2, '新增文章', 'article.add', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (21, 0, 10, 2, 'btn', 'primary', 'small', 2, '编辑', 'article.edit', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (22, 0, 10, 2, 'btn', 'danger', 'small', 2, '删除', 'article.del', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (25, 0, 21, 1, 'btn', 'primary', 'small', 2, '新增配置', 'sys.config.add', 2, 1, '2026-07-08 16:00:42', '2026-07-08 16:02:30', NULL);
INSERT INTO `menu_actions` VALUES (26, 0, 21, 2, 'btn', 'primary', 'small', 2, '编辑配置', 'sys.config.edit', 2, 2, '2026-07-08 16:01:29', '2026-07-08 16:02:37', NULL);
INSERT INTO `menu_actions` VALUES (27, 0, 21, 2, 'btn', 'danger', 'small', 2, '删除配置', 'sys.config.del', 2, 3, '2026-07-08 16:03:15', '2026-07-08 16:03:15', NULL);
INSERT INTO `menu_actions` VALUES (28, 0, 23, 1, 'btn', 'primary', 'default', 2, '新增配置分类', 'sys.configCategory.add', 2, 1, '2026-07-09 11:11:01', '2026-07-09 11:13:11', NULL);
INSERT INTO `menu_actions` VALUES (29, 0, 23, 2, 'btn', 'primary', 'small', 2, '编辑配置分类', 'sys.configCategory.edit', 2, 2, '2026-07-09 11:11:36', '2026-07-09 11:13:17', NULL);
INSERT INTO `menu_actions` VALUES (30, 0, 23, 2, 'btn', 'primary', 'small', 2, '删除配置分类', 'sys.configCategory.del', 2, 3, '2026-07-09 11:12:04', '2026-07-09 11:13:25', NULL);

-- ----------------------------
-- Table structure for menu_meta
-- ----------------------------
DROP TABLE IF EXISTS `menu_meta`;
CREATE TABLE `menu_meta`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `menu_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '菜单id',
  `title` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '菜单名称',
  `icon` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '菜单图标',
  `is_hide` tinyint(3) UNSIGNED NOT NULL DEFAULT 2 COMMENT '是否隐藏 1=是 2=否',
  `is_keep_alive` tinyint(3) UNSIGNED NOT NULL DEFAULT 2 COMMENT '是否缓存 1=是 2=否',
  `is_affix` tinyint(3) UNSIGNED NOT NULL DEFAULT 2 COMMENT '是否固定 1=是 2=否',
  `is_link` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '外链/内嵌时链接地址(http:xxx.com),开启外链条件`1 isLink:链接地址不为空`',
  `is_iframe` tinyint(3) UNSIGNED NOT NULL DEFAULT 2 COMMENT '是否内嵌 1=是 2=否 开启条件`1 isIframe:true 2 isLink:链接地址不为空`',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_menu_id`(`menu_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单元数据表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of menu_meta
-- ----------------------------
INSERT INTO `menu_meta` VALUES (1, 1, 'message.router.home', 'iconfont icon-shouye', 2, 1, 1, '', 2, '2025-05-23 15:37:03', '2025-06-13 11:10:18', NULL);
INSERT INTO `menu_meta` VALUES (2, 2, 'message.router.system', 'iconfont icon-xitongshezhi', 2, 1, 2, '', 2, '2025-05-23 15:39:37', '2025-05-27 16:49:52', NULL);
INSERT INTO `menu_meta` VALUES (3, 3, 'message.router.systemMenu', 'iconfont icon-caidan', 2, 1, 2, '', 2, '2025-05-23 15:41:38', '2025-06-11 17:17:14', NULL);
INSERT INTO `menu_meta` VALUES (4, 4, 'message.router.systemUser', 'iconfont icon-icon-', 2, 1, 2, '', 2, '2025-05-23 23:26:38', '2025-06-11 17:17:29', NULL);
INSERT INTO `menu_meta` VALUES (5, 5, 'message.router.systemRole', 'fa fa-user-circle-o', 2, 1, 2, '', 2, '2025-05-25 14:37:04', '2025-06-11 17:17:36', NULL);
INSERT INTO `menu_meta` VALUES (6, 6, 'message.router.systemDic', 'ele-Collection', 2, 1, 2, '', 2, '2025-05-25 14:54:04', '2025-06-11 17:17:42', NULL);
INSERT INTO `menu_meta` VALUES (7, 10, 'message.article.title', 'ele-Collection', 2, 1, 2, '', 2, '2025-06-16 15:34:11', '2025-06-16 15:34:11', NULL);
INSERT INTO `menu_meta` VALUES (17, 20, '配置管理', 'iconfont icon-ico', 2, 1, 2, '', 2, '2026-07-08 15:04:29', '2006-01-02 15:04:05', NULL);
INSERT INTO `menu_meta` VALUES (18, 21, '配置列表', 'iconfont icon-quanjushezhi_o', 2, 1, 2, '', 2, '2026-07-08 15:51:11', '2026-07-08 15:51:11', NULL);
INSERT INTO `menu_meta` VALUES (19, 22, '系统配置', 'iconfont icon--chaifenhang', 2, 1, 2, '', 2, '2026-07-08 15:54:52', '2006-01-02 15:04:05', NULL);
INSERT INTO `menu_meta` VALUES (20, 23, '配置分类', 'iconfont icon--chaifenlie', 2, 1, 2, '', 2, '2026-07-09 10:56:40', '2026-07-09 10:56:40', NULL);

-- ----------------------------
-- Table structure for migrations
-- ----------------------------
DROP TABLE IF EXISTS `migrations`;
CREATE TABLE `migrations`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `migration` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `migration`(`migration`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of migrations
-- ----------------------------
INSERT INTO `migrations` VALUES (1, '20251212_create_user_table', '2025-12-12 17:04:27.313');

-- ----------------------------
-- Table structure for role_actions
-- ----------------------------
DROP TABLE IF EXISTS `role_actions`;
CREATE TABLE `role_actions`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '角色id',
  `action_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '功能id',
  `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '角色名称',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 58 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色功能表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role_actions
-- ----------------------------
INSERT INTO `role_actions` VALUES (14, 1, 1, 'admin', '2025-06-13 16:26:48', '2025-06-13 16:26:48', NULL);
INSERT INTO `role_actions` VALUES (18, 1, 3, 'admin', '2025-06-13 16:27:43', '2025-06-13 16:27:43', NULL);
INSERT INTO `role_actions` VALUES (20, 1, 4, 'admin', '2025-06-13 16:27:51', '2025-06-13 16:27:51', NULL);
INSERT INTO `role_actions` VALUES (22, 1, 6, 'admin', '2025-06-13 16:28:07', '2025-06-13 16:28:07', NULL);
INSERT INTO `role_actions` VALUES (34, 1, 2, 'admin', '2025-06-13 17:17:32', '2025-06-13 17:17:32', NULL);
INSERT INTO `role_actions` VALUES (35, 1, 7, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (36, 1, 8, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (37, 1, 9, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (38, 1, 10, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (39, 1, 11, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (40, 1, 12, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (41, 1, 13, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (42, 1, 14, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (43, 1, 15, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (44, 1, 16, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (45, 1, 17, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (46, 1, 18, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (47, 1, 19, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (48, 1, 20, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (49, 1, 21, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (50, 1, 22, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `role_actions` VALUES (60, 1, 25, 'admin', '2026-07-08 16:02:30', '2026-07-08 16:02:30', NULL);
INSERT INTO `role_actions` VALUES (61, 1, 26, 'admin', '2026-07-08 16:02:37', '2026-07-08 16:02:37', NULL);
INSERT INTO `role_actions` VALUES (62, 1, 27, 'admin', '2026-07-08 16:03:15', '2026-07-08 16:03:15', NULL);
INSERT INTO `role_actions` VALUES (66, 1, 28, 'admin', '2026-07-09 11:13:11', '2026-07-09 11:13:11', NULL);
INSERT INTO `role_actions` VALUES (67, 1, 29, 'admin', '2026-07-09 11:13:17', '2026-07-09 11:13:17', NULL);
INSERT INTO `role_actions` VALUES (68, 1, 30, 'admin', '2026-07-09 11:13:25', '2026-07-09 11:13:25', NULL);

-- ----------------------------
-- Table structure for role_menus
-- ----------------------------
DROP TABLE IF EXISTS `role_menus`;
CREATE TABLE `role_menus`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '角色id',
  `menu_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '菜单id',
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '角色名称',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_role_id`(`role_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 105 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色菜单表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role_menus
-- ----------------------------
INSERT INTO `role_menus` VALUES (69, 1, 1, '', '2026-07-06 15:42:29', '2026-07-06 15:42:29', NULL);
INSERT INTO `role_menus` VALUES (70, 1, 2, '', '2026-07-06 15:42:29', '2026-07-06 15:42:29', NULL);
INSERT INTO `role_menus` VALUES (71, 1, 3, '', '2026-07-06 15:42:29', '2026-07-06 15:42:29', NULL);
INSERT INTO `role_menus` VALUES (72, 1, 4, '', '2026-07-06 15:42:29', '2026-07-06 15:42:29', NULL);
INSERT INTO `role_menus` VALUES (73, 1, 5, '', '2026-07-06 15:42:29', '2026-07-06 15:42:29', NULL);
INSERT INTO `role_menus` VALUES (74, 1, 6, '', '2026-07-06 15:42:29', '2026-07-06 15:42:29', NULL);
INSERT INTO `role_menus` VALUES (75, 1, 10, '', '2026-07-06 15:42:29', '2026-07-06 15:42:29', NULL);
INSERT INTO `role_menus` VALUES (108, 1, 21, '', '2026-07-08 15:51:11', '2026-07-08 15:51:11', NULL);
INSERT INTO `role_menus` VALUES (109, 1, 20, '', '2026-07-08 15:51:29', '2026-07-08 15:51:29', NULL);
INSERT INTO `role_menus` VALUES (112, 1, 22, '', '2026-07-08 15:57:35', '2026-07-08 15:57:35', NULL);
INSERT INTO `role_menus` VALUES (113, 1, 23, '', '2026-07-09 10:56:40', '2026-07-09 10:56:40', NULL);

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '角色名称',
  `desc` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '角色描述',
  `status` tinyint(1) UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态 1=启用 2=停用',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 14 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of roles
-- ----------------------------
INSERT INTO `roles` VALUES (1, 'admin', '超级管理员', 1, '2025-05-26 16:52:43', '2006-01-02 15:04:05', NULL);
INSERT INTO `roles` VALUES (2, 'test', '测试', 1, '2025-05-28 10:47:22', '2006-01-02 15:04:05', NULL);

-- ----------------------------
-- Table structure for system_config
-- ----------------------------
DROP TABLE IF EXISTS `system_config`;
CREATE TABLE `system_config`  (
  `id` smallint(5) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
  `key` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标识',
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '名称',
  `default_value` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '默认值',
  `option_value` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '可选值',
  `type` tinyint(1) UNSIGNED NOT NULL DEFAULT 1 COMMENT '配置类型 1=输入框 2=单选 3=复选 4=下拉菜单 5=文本域 6=附件',
  `config_category_id` tinyint(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '配置分类Id',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unq_key`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统配置表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of system_config
-- ----------------------------
INSERT INTO `system_config` VALUES (1, 'web_domain', '网站域名', 'www.a.com', '', 1, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (2, 'is_open_site', '关闭站点', '开启', '关闭,开启', 2, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (3, 'site_logo', '网站Logo', '', '', 6, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (4, 'email_port', '邮件端口', '465', '', 1, 2, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (5, 'email_title', '邮件标题', '【xxx】验证码', '', 1, 2, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (6, 'send_user_info', '发件人信息', '【管理员】', '', 1, 2, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (7, 'email_content', '发送内容', '【xxx】你的验证码是：', '', 5, 2, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (8, 'web_keyword', '关键词', '关键词...', '', 5, 3, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (9, 'email', '邮箱账号', 'xxx@email.com', '', 1, 2, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (10, 'record_number', '备案编号', 'Copyright© 2014-2019 | Powered by ***1.1 | 粤ICP备****号', '', 1, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (11, 'web_description', '网站描述', 'web', '', 1, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (12, 'select', '下拉选项', '下拉3', '下拉1,下拉2,下拉3', 4, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (13, 'checkbox', '复选框', 'HTML,CSS', 'AJAX,HTML,JS,CSS', 3, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (14, 'textarea', '文本域', '文本域', '0', 5, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (15, 'default_head_img', '默认头像', '', '', 6, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (16, 'seo_description', '描述', '11', '', 5, 3, NULL, NULL, NULL);

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '头像',
  `username` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `full_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '姓名',
  `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '邮箱',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '密码',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '昵称',
  `gender` tinyint(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '性别 1=男 2=女',
  `age` tinyint(3) UNSIGNED NOT NULL DEFAULT 0 COMMENT '年龄',
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态 1=启用 2=停用',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 56 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, 'https://cdn.qitx.net/local/myblog/user_header_image/20230517/577a53d123bc4c4f19db0cb2c6c980a8.jpg', 'admin', '超级管理员', 'dsx.emil@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', '大师兄', 1, 31, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (2, '', 'test2', '李四1', 'ls@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2006-01-02 15:04:05', NULL);
INSERT INTO `user` VALUES (10, '', 'dsx', '大师兄111', 'dsx@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄', 1, 0, 1, '2024-07-22 17:34:36', '2006-01-02 15:04:05', NULL);
INSERT INTO `user` VALUES (11, '', 'admin1', '张三1', 'zs1@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (12, '', 'test3', '李四1', 'ls3@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', NULL);
INSERT INTO `user` VALUES (14, '', 'dsx1', '大师兄1', 'dsx1@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄1', 1, 0, 1, '2024-07-22 17:34:36', '2024-07-22 17:34:36', NULL);
INSERT INTO `user` VALUES (15, '', 'admin2', '张三2', 'zs2@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (16, '', 'test5', '李四5', 'ls5@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', NULL);
INSERT INTO `user` VALUES (18, '', 'dsx2', '大师兄2', 'dsx2@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄1', 1, 0, 1, '2024-07-22 17:34:36', '2024-07-22 17:34:36', NULL);
INSERT INTO `user` VALUES (19, '', 'admin3', '张三3', 'zs3@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (20, '', 'test7', '李四7', 'ls7@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', NULL);
INSERT INTO `user` VALUES (22, '', 'admin4', '张三4', 'zs4@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (23, '', 'test9', '李四9', 'ls9@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', NULL);
INSERT INTO `user` VALUES (25, '', 'dsx3', '大师兄3', 'dsx3@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄1', 1, 0, 1, '2024-07-22 17:34:36', '2024-07-22 17:34:36', NULL);
INSERT INTO `user` VALUES (26, '', 'admin5', '张三5', 'zs5@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (27, '', 'test11', '李四11', 'ls11@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', NULL);
INSERT INTO `user` VALUES (29, '', 'dsx4', '大师兄4', 'dsx4@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄1', 1, 0, 1, '2024-07-22 17:34:36', '2024-07-22 17:34:36', NULL);
INSERT INTO `user` VALUES (30, '', 'admin6', '张三6', 'zs6@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (31, '', 'test13', '李四13', 'ls13@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', NULL);
INSERT INTO `user` VALUES (33, '', 'dsx5', '大师兄5', 'dsx5@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄1', 1, 0, 1, '2024-07-22 17:34:36', '2024-07-22 17:34:36', NULL);
INSERT INTO `user` VALUES (34, '', 'admin7', '测试张三34', 'zs7@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2006-01-02 15:04:05', NULL);
INSERT INTO `user` VALUES (54, '', '测试123', '测试123', '', '$2a$10$kckcEhxk/7pppnM6XCT.cu.xQyoxlnhdY.nhdpk5CmzzCFY32l2xm', '测试123', 1, 30, 1, '2026-06-24 16:52:38', '2006-01-02 15:04:05', NULL);

-- ----------------------------
-- Table structure for user_roles
-- ----------------------------
DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
  `role_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '角色id',
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '角色名称',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 32 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户角色表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user_roles
-- ----------------------------
INSERT INTO `user_roles` VALUES (1, 1, 1, 'admin', '2025-05-26 17:53:10', '2025-05-26 17:53:10', NULL);
INSERT INTO `user_roles` VALUES (11, 10, 2, 'test', '2025-05-29 14:37:18', '2025-05-29 14:37:18', NULL);
INSERT INTO `user_roles` VALUES (28, 54, 1, 'admin', '2026-06-24 16:55:38', '2026-06-24 16:55:38', NULL);
INSERT INTO `user_roles` VALUES (29, 2, 9, 'cs', '2026-07-06 16:17:28', '2026-07-06 16:17:28', NULL);
INSERT INTO `user_roles` VALUES (31, 55, 2, 'test', '2026-07-07 14:03:56', '2026-07-07 14:03:56', NULL);

SET FOREIGN_KEY_CHECKS = 1;
