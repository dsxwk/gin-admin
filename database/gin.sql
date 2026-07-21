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

 Date: 21/07/2026 17:17:30
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
INSERT INTO `article` VALUES (1, 1, '标题1', '<p>测试1</p>', 1, 2, 1, '[\"测试标签1\", \"测试标签2\", \"cs3\"]', '2023-09-19 11:43:58', '2026-07-21 16:43:33', NULL);
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
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '配置分类表' ROW_FORMAT = DYNAMIC;

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
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '字典表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of dict
-- ----------------------------
INSERT INTO `dict` VALUES (1, 0, 'gender', '性别', '', 1, 0, '{\"test\": 111, \"test2\": \"test222\"}', '', '2025-06-06 21:48:17', '2025-06-06 21:48:17', NULL);
INSERT INTO `dict` VALUES (2, 1, 'gender', '男', '1', 1, 0, '{\"a\": \"111\", \"b\": \"222\"}', '测试', '2025-06-06 21:49:00', '2026-07-15 15:37:07', NULL);
INSERT INTO `dict` VALUES (3, 1, 'gender', '女', '2', 1, 0, '{\"test\": 111, \"test2\": \"test222\"}', '性别女', '2025-06-06 21:49:10', '2026-07-15 15:36:48', NULL);
INSERT INTO `dict` VALUES (4, 1, 'gender', '保密', '0', 1, 0, '{}', '保密', '2026-07-03 13:43:30', '2026-07-21 16:43:12', NULL);

-- ----------------------------
-- Table structure for import_records
-- ----------------------------
DROP TABLE IF EXISTS `import_records`;
CREATE TABLE `import_records`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `type` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '导入类型',
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '类型名称',
  `data` json NULL COMMENT '导入数据',
  `created_user` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建人',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '导入记录表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of import_records
-- ----------------------------
INSERT INTO `import_records` VALUES (1, 1, '用户导入', '[{\"age\": 28, \"email\": \"zhangsan@example.com\", \"gender\": 1, \"status\": 1, \"fullName\": \"张三\", \"nickname\": \"小张\", \"password\": \"123456\", \"username\": \"zhangsan\"}, {\"age\": 28, \"email\": \"zhang1@example.com\", \"gender\": 1, \"status\": 1, \"fullName\": \"张1\", \"nickname\": \"小张1\", \"password\": \"123456\", \"username\": \"zhan\"}]', 1, '2026-07-16 16:37:40', '2026-07-16 16:37:40', NULL);

-- ----------------------------
-- Table structure for menu
-- ----------------------------
DROP TABLE IF EXISTS `menu`;
CREATE TABLE `menu`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `pid` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父级id',
  `type` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '类型 1=菜单 2=功能',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '路由名称|功能标识',
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态 1=启用 2=停用',
  `sort` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_pid`(`pid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 63 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of menu
-- ----------------------------
INSERT INTO `menu` VALUES (1, 0, 1, 'home', 1, 1, '2025-05-23 15:37:03', '2025-06-13 11:10:18', NULL);
INSERT INTO `menu` VALUES (2, 0, 1, 'system', 1, 2, '2025-05-23 15:39:37', '2025-05-27 16:49:52', NULL);
INSERT INTO `menu` VALUES (3, 2, 1, 'systemMenu', 1, 3, '2025-05-23 15:41:38', '2025-06-11 17:17:14', NULL);
INSERT INTO `menu` VALUES (4, 2, 1, 'systemUser', 1, 4, '2025-05-23 23:26:38', '2025-06-11 17:17:29', NULL);
INSERT INTO `menu` VALUES (5, 2, 1, 'systemRole', 1, 5, '2025-05-25 14:37:04', '2025-06-11 17:17:36', NULL);
INSERT INTO `menu` VALUES (6, 2, 1, 'systemDic', 1, 6, '2025-05-25 14:54:04', '2025-06-11 17:17:42', NULL);
INSERT INTO `menu` VALUES (10, 0, 1, 'article', 1, 7, '2025-06-16 15:34:11', '2025-06-16 15:34:11', NULL);
INSERT INTO `menu` VALUES (20, 2, 1, 'systemConfig', 1, 7, '2026-07-08 15:04:29', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (21, 20, 1, 'systemConfigList', 1, 7, '2026-07-08 15:51:11', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (22, 20, 1, 'systemConfigSetting', 1, 8, '2026-07-08 15:54:52', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (23, 20, 1, 'systemConfigCategory', 1, 1, '2026-07-09 10:56:40', '2026-07-09 10:56:40', NULL);
INSERT INTO `menu` VALUES (24, 3, 2, 'sys.menu.add', 1, 0, '2025-05-21 10:24:14', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (25, 3, 2, 'sys.menu.edit', 1, 1, '2025-05-21 10:30:24', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (27, 3, 2, 'sys.menu.del', 1, 2, '2025-05-21 10:30:49', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (32, 4, 2, 'sys.user.add', 1, 0, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (33, 4, 2, 'sys.user.batchDel', 1, 0, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (34, 4, 2, 'sys.user.edit', 1, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu` VALUES (35, 4, 2, 'sys.user.del', 1, 2, '2025-06-16 08:57:04', '2026-07-16 14:40:54', NULL);
INSERT INTO `menu` VALUES (37, 5, 2, 'sys.role.add', 1, 0, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (38, 5, 2, 'sys.role.edit', 1, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu` VALUES (39, 5, 2, 'sys.role.del', 1, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu` VALUES (40, 6, 2, 'sys.dic.add', 1, 0, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (41, 6, 2, 'sys.dic.edit', 1, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu` VALUES (42, 6, 2, 'sys.dic.del', 1, 2, '2025-06-16 08:57:04', '2026-07-15 15:59:14', NULL);
INSERT INTO `menu` VALUES (43, 10, 2, 'article.add', 1, 0, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (44, 10, 2, 'article.edit', 1, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu` VALUES (45, 10, 2, 'article.del', 1, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu` VALUES (48, 21, 2, 'sys.config.add', 1, 0, '2026-07-08 16:00:42', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (49, 21, 2, 'sys.config.edit', 1, 0, '2026-07-08 16:01:29', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (50, 21, 2, 'sys.config.del', 1, 0, '2026-07-08 16:03:15', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (51, 23, 2, 'sys.configCategory.add', 1, 0, '2026-07-09 11:11:01', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (52, 23, 2, 'sys.configCategory.edit', 1, 0, '2026-07-09 11:11:36', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (53, 23, 2, 'sys.configCategory.del', 1, 0, '2026-07-09 11:12:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu` VALUES (54, 4, 2, 'sys.user.import', 1, 0, '2026-07-10 13:34:50', '2026-07-10 13:34:50', NULL);
INSERT INTO `menu` VALUES (56, 54, 2, 'sys.user.importRecords', 1, 10, '2026-07-13 15:53:05', '2026-07-17 15:56:01', NULL);
INSERT INTO `menu` VALUES (58, 3, 2, 'sys.menu.addChildren', 1, 0, '2026-07-14 09:36:51', '2026-07-14 09:36:51', NULL);
INSERT INTO `menu` VALUES (59, 6, 2, 'sys.dic.addChildren', 1, 1, '2026-07-15 15:58:56', '2026-07-15 15:58:56', NULL);
INSERT INTO `menu` VALUES (60, 4, 2, 'sys.user.password', 1, 1, '2026-07-16 14:41:52', '2026-07-16 14:42:32', NULL);
INSERT INTO `menu` VALUES (61, 56, 2, 'sys.user.importRecords.detail', 1, 0, '2026-07-17 14:35:15', '2026-07-17 15:47:14', NULL);
INSERT INTO `menu` VALUES (62, 56, 2, 'sys.user.importRecords.delete', 1, 0, '2026-07-17 15:30:12', '2026-07-17 15:45:23', NULL);

-- ----------------------------
-- Table structure for menu_actions
-- ----------------------------
DROP TABLE IF EXISTS `menu_actions`;
CREATE TABLE `menu_actions`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `menu_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '菜单id',
  `type` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '类型 1=header 2=operation',
  `btn_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'btn' COMMENT '按钮类型 text|btn',
  `btn_style` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'primary' COMMENT '按钮样式',
  `btn_size` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'small' COMMENT '按钮尺寸',
  `is_confirm` tinyint(3) UNSIGNED NOT NULL DEFAULT 2 COMMENT '是否确认 1=是 2=否',
  `label` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '功能名称',
  `auth_value` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '权限标识',
  `is_link` tinyint(3) UNSIGNED NOT NULL DEFAULT 2 COMMENT '是否为链接 1=是 2=否',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_menu_id`(`menu_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 42 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单功能表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of menu_actions
-- ----------------------------
INSERT INTO `menu_actions` VALUES (1, 24, 1, 'btn', 'primary', 'default', 2, '新增菜单', 'sys.menu.add', 2, '2025-05-21 10:24:14', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (2, 25, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.menu.edit', 2, '2025-05-21 10:30:24', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (4, 27, 2, 'btn', 'danger', 'small', 1, '删除', 'sys.menu.del', 2, '2025-05-21 10:30:49', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (8, 32, 1, 'btn', 'primary', 'default', 2, '新增用户', 'sys.user.add', 2, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (9, 33, 1, 'btn', 'danger', 'default', 2, '批量删除', 'sys.user.batchDel', 2, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (10, 34, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.user.edit', 2, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (11, 35, 2, 'btn', 'danger', 'small', 2, '删除', 'sys.user.del', 2, '2025-06-16 08:57:04', '2026-07-16 14:40:54', NULL);
INSERT INTO `menu_actions` VALUES (13, 37, 1, 'btn', 'primary', 'default', 2, '新增角色', 'sys.role.add', 2, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (14, 38, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.role.edit', 2, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (15, 39, 2, 'btn', 'danger', 'small', 2, '删除', 'sys.role.del', 2, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (16, 40, 1, 'btn', 'primary', 'default', 2, '新增字典', 'sys.dic.add', 2, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (17, 41, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.dic.edit', 2, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (18, 42, 2, 'btn', 'danger', 'small', 2, '删除', 'sys.dic.del', 2, '2025-06-16 08:57:04', '2026-07-15 15:59:14', NULL);
INSERT INTO `menu_actions` VALUES (19, 43, 1, 'btn', 'primary', 'default', 2, '新增文章', 'article.add', 2, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (20, 44, 2, 'btn', 'primary', 'small', 2, '编辑', 'article.edit', 2, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (21, 45, 2, 'btn', 'danger', 'small', 2, '删除', 'article.del', 2, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (22, 48, 1, 'btn', 'primary', 'default', 2, '新增配置', 'sys.config.add', 2, '2026-07-08 16:00:42', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (23, 49, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.config.edit', 2, '2026-07-08 16:01:29', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (24, 50, 2, 'btn', 'danger', 'small', 2, '删除', 'sys.config.del', 2, '2026-07-08 16:03:15', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (25, 51, 1, 'btn', 'primary', 'default', 2, '新增分类', 'sys.configCategory.add', 2, '2026-07-09 11:11:01', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (26, 52, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.configCategory.edit', 2, '2026-07-09 11:11:36', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (27, 53, 2, 'btn', 'danger', 'small', 2, '删除', 'sys.configCategory.del', 2, '2026-07-09 11:12:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (28, 54, 1, 'btn', 'primary', 'default', 2, '用户导入', 'sys.user.import', 2, '2026-07-10 13:34:50', '2026-07-10 13:34:50', NULL);
INSERT INTO `menu_actions` VALUES (33, 58, 2, 'btn', 'primary', 'small', 2, '新增子集', 'sys.menu.addChildren', 2, '2026-07-14 09:36:51', '2026-07-14 09:36:51', NULL);
INSERT INTO `menu_actions` VALUES (34, 59, 2, 'btn', 'primary', 'small', 2, '新增子集', 'sys.dic.addChildren', 2, '2026-07-15 15:58:56', '2026-07-15 15:58:56', NULL);
INSERT INTO `menu_actions` VALUES (35, 60, 2, 'btn', 'primary', 'small', 2, '更新密码', 'sys.user.password', 2, '2026-07-16 14:41:52', '2026-07-16 14:42:32', NULL);
INSERT INTO `menu_actions` VALUES (39, 62, 2, 'btn', 'danger', 'small', 2, '删除', 'sys.user.importRecords.delete', 2, '2026-07-17 15:45:23', '2026-07-17 15:45:23', NULL);
INSERT INTO `menu_actions` VALUES (40, 61, 2, 'btn', 'primary', 'small', 2, '明细', 'sys.user.importRecords.detail', 2, '2026-07-17 15:47:14', '2026-07-17 15:47:14', NULL);
INSERT INTO `menu_actions` VALUES (41, 56, 1, 'btn', 'primary', 'default', 2, '导入记录', 'sys.user.importRecords', 2, '2026-07-17 15:56:01', '2026-07-17 15:56:01', NULL);

-- ----------------------------
-- Table structure for menu_meta
-- ----------------------------
DROP TABLE IF EXISTS `menu_meta`;
CREATE TABLE `menu_meta`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `menu_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '菜单id',
  `title` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '菜单名称',
  `icon` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '菜单图标',
  `path` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '路由路径',
  `redirect` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '重定向',
  `component` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '组件路径',
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
) ENGINE = InnoDB AUTO_INCREMENT = 18 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单元数据表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of menu_meta
-- ----------------------------
INSERT INTO `menu_meta` VALUES (1, 1, 'message.router.home', 'iconfont icon-shouye', '/home', '', 'home/index', 2, 1, 1, '', 2, '2025-05-23 15:37:03', '2025-06-13 11:10:18', NULL);
INSERT INTO `menu_meta` VALUES (2, 2, 'message.router.system', 'iconfont icon-xitongshezhi', '/system', '/system/menu', 'layouts/routerView/parent', 2, 1, 2, '', 2, '2025-05-23 15:39:37', '2025-05-27 16:49:52', NULL);
INSERT INTO `menu_meta` VALUES (3, 3, 'message.router.systemMenu', 'iconfont icon-caidan', '/system/menu', '', 'system/menu/index', 2, 1, 2, '', 2, '2025-05-23 15:41:38', '2025-06-11 17:17:14', NULL);
INSERT INTO `menu_meta` VALUES (4, 4, 'message.router.systemUser', 'iconfont icon-icon-', '/system/user', '', 'system/user/index', 2, 1, 2, '', 2, '2025-05-23 23:26:38', '2025-06-11 17:17:29', NULL);
INSERT INTO `menu_meta` VALUES (5, 5, 'message.router.systemRole', 'fa fa-user-circle-o', '/system/role', '', 'system/role/index', 2, 1, 2, '', 2, '2025-05-25 14:37:04', '2025-06-11 17:17:36', NULL);
INSERT INTO `menu_meta` VALUES (6, 6, 'message.router.systemDic', 'ele-Collection', '/system/dic', '', 'system/dic/index', 2, 1, 2, '', 2, '2025-05-25 14:54:04', '2025-06-11 17:17:42', NULL);
INSERT INTO `menu_meta` VALUES (7, 10, 'message.article.title', 'ele-Collection', '/article', '', 'article/index', 2, 1, 2, '', 2, '2025-06-16 15:34:11', '2025-06-16 15:34:11', NULL);
INSERT INTO `menu_meta` VALUES (8, 20, '配置管理', 'iconfont icon-ico', '/system/config', '', 'layouts/routerView/parent', 2, 1, 2, '', 2, '2026-07-08 15:04:29', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_meta` VALUES (9, 21, '配置列表', 'iconfont icon-quanjushezhi_o', '/system/config/index', '', 'system/config/index', 2, 1, 2, '', 2, '2026-07-08 15:51:11', '2026-07-08 15:51:11', NULL);
INSERT INTO `menu_meta` VALUES (10, 22, '系统配置', 'iconfont icon--chaifenhang', '/system/config/setting', '', 'system/config/setting', 2, 1, 2, '', 2, '2026-07-08 15:54:52', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_meta` VALUES (11, 23, '配置分类', 'iconfont icon--chaifenlie', '/system/config-category/index', '', 'system/config/category/index', 2, 1, 2, '', 2, '2026-07-09 10:56:40', '2026-07-09 10:56:40', NULL);

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
-- Table structure for permission
-- ----------------------------
DROP TABLE IF EXISTS `permission`;
CREATE TABLE `permission`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `key` varchar(130) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '权限标识',
  `method` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '请求方式',
  `uri` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '路由地址',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_key`(`key`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '权限表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of permission
-- ----------------------------
INSERT INTO `permission` VALUES (1, 'GET:/api/v1/article', 'GET', '/api/v1/article', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (2, 'GET:/api/v1/menu', 'GET', '/api/v1/menu', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (3, 'POST:/api/v1/article', 'POST', '/api/v1/article', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (4, 'PUT:/api/v1/system-config', 'PUT', '/api/v1/system-config', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (5, 'PUT:/api/v1/user/:id/password', 'PUT', '/api/v1/user/:id/password', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (6, 'GET:/api/v1/config-category', 'GET', '/api/v1/config-category', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (7, 'POST:/api/v1/system-config', 'POST', '/api/v1/system-config', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (8, 'PUT:/api/v1/article/:id', 'PUT', '/api/v1/article/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (9, 'DELETE:/api/v1/menu/:id', 'DELETE', '/api/v1/menu/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (10, 'PUT:/api/v1/dict/:id', 'PUT', '/api/v1/dict/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (11, 'DELETE:/api/v1/system-config/:id', 'DELETE', '/api/v1/system-config/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (12, 'GET:/api/v1/dict', 'GET', '/api/v1/dict', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (13, 'GET:/api/v1/user', 'GET', '/api/v1/user', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (14, 'GET:/api/v1/user/:id', 'GET', '/api/v1/user/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (15, 'POST:/api/v1/config-category', 'POST', '/api/v1/config-category', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (16, 'POST:/api/v1/dict', 'POST', '/api/v1/dict', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (17, 'PUT:/api/v1/config-category/:id', 'PUT', '/api/v1/config-category/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (18, 'DELETE:/api/v1/dict/:id', 'DELETE', '/api/v1/dict/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (19, 'GET:/api/v1/article/:id', 'GET', '/api/v1/article/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (20, 'GET:/api/v1/config-category/:id', 'GET', '/api/v1/config-category/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (21, 'GET:/api/v1/dict/:id', 'GET', '/api/v1/dict/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (22, 'POST:/api/v1/user/import', 'POST', '/api/v1/user/import', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (23, 'POST:/api/v1/menu', 'POST', '/api/v1/menu', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (24, 'POST:/api/v1/role', 'POST', '/api/v1/role', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (25, 'PUT:/api/v1/menu/:id', 'PUT', '/api/v1/menu/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (26, 'GET:/api/v1/system-config', 'GET', '/api/v1/system-config', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (27, 'GET:/api/v1/system-config/:id', 'GET', '/api/v1/system-config/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (28, 'GET:/api/v1/import-records', 'GET', '/api/v1/import-records', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (29, 'POST:/api/v1/user', 'POST', '/api/v1/user', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (30, 'PUT:/api/v1/user/:id', 'PUT', '/api/v1/user/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (31, 'DELETE:/api/v1/role/:id', 'DELETE', '/api/v1/role/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (32, 'GET:/api/v1/role/:id', 'GET', '/api/v1/role/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (33, 'GET:/api/v1/menu/:id', 'GET', '/api/v1/menu/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (34, 'PUT:/api/v1/role/:id', 'PUT', '/api/v1/role/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (35, 'GET:/api/v1/role/:id/menu', 'GET', '/api/v1/role/:id/menu', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (36, 'POST:/api/v1/user/batch-delete', 'POST', '/api/v1/user/batch-delete', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (37, 'PUT:/api/v1/system-config/:id', 'PUT', '/api/v1/system-config/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (38, 'DELETE:/api/v1/article/:id', 'DELETE', '/api/v1/article/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (39, 'DELETE:/api/v1/config-category/:id', 'DELETE', '/api/v1/config-category/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (40, 'DELETE:/api/v1/import-records/:id', 'DELETE', '/api/v1/import-records/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (41, 'DELETE:/api/v1/user/:id', 'DELETE', '/api/v1/user/:id', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (42, 'GET:/api/v1/role', 'GET', '/api/v1/role', '2026-07-21 12:39:57', '2026-07-21 12:39:57', NULL);
INSERT INTO `permission` VALUES (43, 'GET:/api/v1/permission', 'GET', '/api/v1/permission', '2026-07-21 13:33:14', '2026-07-21 13:33:14', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 310 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色菜单表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role_menus
-- ----------------------------
INSERT INTO `role_menus` VALUES (302, 2, 1, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `role_menus` VALUES (303, 2, 24, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `role_menus` VALUES (304, 2, 27, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `role_menus` VALUES (305, 2, 2, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `role_menus` VALUES (306, 2, 3, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `role_menus` VALUES (350, 1, 1, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (351, 1, 2, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (352, 1, 3, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (353, 1, 24, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (354, 1, 58, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (355, 1, 25, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (356, 1, 27, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (357, 1, 4, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (358, 1, 54, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (359, 1, 56, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (360, 1, 61, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (361, 1, 62, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (362, 1, 32, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (363, 1, 33, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (364, 1, 34, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (365, 1, 60, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (366, 1, 35, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (367, 1, 5, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (368, 1, 39, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (369, 1, 37, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (370, 1, 38, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (371, 1, 6, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (372, 1, 41, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (373, 1, 40, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (374, 1, 59, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (375, 1, 42, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (376, 1, 20, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (377, 1, 23, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (378, 1, 51, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (379, 1, 52, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (380, 1, 53, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (381, 1, 21, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (382, 1, 48, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (383, 1, 49, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (384, 1, 50, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (385, 1, 22, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (386, 1, 10, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (387, 1, 43, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (388, 1, 44, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `role_menus` VALUES (389, 1, 45, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);

-- ----------------------------
-- Table structure for role_permissions
-- ----------------------------
DROP TABLE IF EXISTS `role_permissions`;
CREATE TABLE `role_permissions`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '角色id',
  `permission_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '权限id',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_role_id`(`role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色权限表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role_permissions
-- ----------------------------
INSERT INTO `role_permissions` VALUES (1, 1, 1, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (2, 1, 3, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (3, 1, 38, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (4, 1, 19, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (5, 1, 8, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (6, 1, 6, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (7, 1, 15, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (8, 1, 39, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (9, 1, 20, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (10, 1, 17, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (11, 1, 12, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (12, 1, 16, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (13, 1, 18, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (14, 1, 21, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (15, 1, 10, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (16, 1, 28, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (17, 1, 40, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (18, 1, 2, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (19, 1, 23, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (20, 1, 9, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (21, 1, 33, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (22, 1, 25, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (23, 1, 43, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (24, 1, 42, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (25, 1, 24, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (26, 1, 31, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (27, 1, 32, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (28, 1, 34, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (29, 1, 35, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (30, 1, 26, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (31, 1, 7, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (32, 1, 4, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (33, 1, 11, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (34, 1, 27, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (35, 1, 37, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (36, 1, 13, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (37, 1, 29, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (38, 1, 41, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (39, 1, 14, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (40, 1, 30, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (41, 1, 5, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (42, 1, 36, '2026-07-21 14:34:18', NULL);
INSERT INTO `role_permissions` VALUES (43, 1, 22, '2026-07-21 14:34:18', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of roles
-- ----------------------------
INSERT INTO `roles` VALUES (1, 'admin', '超级管理员', 1, '2025-05-26 16:52:43', '2026-07-21 14:34:18', NULL);
INSERT INTO `roles` VALUES (2, 'test', '测试', 1, '2025-05-28 10:47:22', '2026-07-20 14:19:04', NULL);

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
INSERT INTO `system_config` VALUES (1, 'web_domain', '网站域名', 'www.a.com', '', 1, 1, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (2, 'is_open_site', '关闭站点', '开启', '关闭,开启', 2, 1, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (3, 'site_logo', '网站Logo', '', '', 6, 1, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (4, 'email_port', '邮件端口', '465', '', 1, 2, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (5, 'email_title', '邮件标题', '【xxx】验证码', '', 1, 2, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (6, 'send_user_info', '发件人信息', '【管理员】', '', 1, 2, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (7, 'email_content', '发送内容', '【xxx】你的验证码是：', '', 5, 2, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (8, 'web_keyword', '关键词', '关键词...', '', 5, 3, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (9, 'email', '邮箱账号', 'xxx@email.com', '', 1, 2, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (10, 'record_number', '备案编号', 'Copyright© 2014-2019 | Powered by ***1.1 | 粤ICP备****号', '', 1, 1, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (11, 'web_description', '网站描述', 'web', '', 1, 1, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (12, 'select', '下拉选项', '下拉3', '下拉1,下拉2,下拉3', 4, 1, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (13, 'checkbox', '复选框', 'HTML,CSS', 'AJAX,HTML,JS,CSS', 3, 1, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (14, 'textarea', '文本域', '文本域', '0', 5, 1, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (15, 'default_head_img', '默认头像', '', '', 6, 1, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);
INSERT INTO `system_config` VALUES (16, 'seo_description', '描述', '11', '', 5, 3, '2026-07-15 16:03:49', '2026-07-15 16:03:49', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 61 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, 'https://cdn.qitx.net/local/myblog/user_header_image/20230517/577a53d123bc4c4f19db0cb2c6c980a8.jpg', 'admin', '超级管理员', 'dsx.emil@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', '大师兄', 1, 31, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (2, '', 'test2', '李四1', 'ls@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2026-07-15 16:08:00', NULL);
INSERT INTO `user` VALUES (10, '', 'dsx', '大师兄111', 'dsx@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄', 1, 0, 1, '2024-07-22 17:34:36', '2026-07-15 16:08:00', NULL);
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
INSERT INTO `user` VALUES (25, '', 'dsx3', '大师兄3', 'dsx3@qq.com', '$2a$10$9O..2Bao08kzp4wWMS.7nOvxXwUfIS2infQPC8cI0WATZv4Dy2Gv2', '大师兄1', 1, 0, 1, '2024-07-22 17:34:36', '2024-07-22 17:34:36', NULL);
INSERT INTO `user` VALUES (26, '', 'admin5', '张三5', 'zs5@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (27, '', 'test11', '李四11', 'ls11@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', NULL);
INSERT INTO `user` VALUES (29, '', 'dsx4', '大师兄4', 'dsx4@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄1', 1, 0, 1, '2024-07-22 17:34:36', '2024-07-22 17:34:36', NULL);
INSERT INTO `user` VALUES (30, '', 'admin6', '张三6', 'zs6@qq.com', '$2a$10$aEnlH2qUqtMnRQ7edr4Z7eROMDfqUevbAfHDA.AMA66kwY3wRtxHO', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (31, '', 'test13', '李四13', 'ls13@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', NULL);
INSERT INTO `user` VALUES (33, '', 'dsx5', '大师兄5', 'dsx5@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄1', 1, 0, 1, '2024-07-22 17:34:36', '2024-07-22 17:34:36', NULL);
INSERT INTO `user` VALUES (34, '', 'admin7', '测试张三34', 'zs7@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2026-07-15 16:08:00', NULL);
INSERT INTO `user` VALUES (59, '', 'zhangsan', '张三', 'zhangsan@example.com', '$2a$10$2jdY9rqE0Vyeu9DTFgbF0uZ.JhPsnvePBj/QvcCv2ws8ipLUagWOW', '小张', 1, 28, 1, '2026-07-16 16:36:30', '2026-07-16 16:36:30', NULL);
INSERT INTO `user` VALUES (60, '', 'zhan', '张1', 'zhang1@example.com', '$2a$10$4Aol5C1x2vBTZYV8cUHvt.5JLmD1suHrYKg6rmJekxQFd4EE4Jnxi', '小张1', 1, 28, 1, '2026-07-16 16:36:30', '2026-07-16 16:36:30', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 39 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户角色表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user_roles
-- ----------------------------
INSERT INTO `user_roles` VALUES (34, 10, 2, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `user_roles` VALUES (35, 2, 2, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `user_roles` VALUES (41, 1, 1, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);
INSERT INTO `user_roles` VALUES (42, 10, 1, 'admin', '2026-07-21 14:34:18', '2026-07-21 14:34:18', NULL);

SET FOREIGN_KEY_CHECKS = 1;
