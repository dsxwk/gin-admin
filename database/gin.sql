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

 Date: 10/08/2026 11:14:20
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
) ENGINE = InnoDB AUTO_INCREMENT = 14 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文章表' ROW_FORMAT = DYNAMIC;

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
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '分类表' ROW_FORMAT = DYNAMIC;

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
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '配置分类表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of config_category
-- ----------------------------
INSERT INTO `config_category` VALUES (1, '基本信息', '2026-07-08 13:42:00', '2006-01-02 15:04:05', NULL);
INSERT INTO `config_category` VALUES (2, '邮箱配置', '2026-07-08 13:42:00', '2006-01-02 15:04:05', NULL);
INSERT INTO `config_category` VALUES (3, 'seo设置 ', '2026-07-08 13:42:00', '2026-07-08 13:42:00', NULL);

-- ----------------------------
-- Table structure for department
-- ----------------------------
DROP TABLE IF EXISTS `department`;
CREATE TABLE `department`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `pid` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父级id',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '部门名称',
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态 1=启用 2=停用',
  `sort` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_pid`(`pid`) USING BTREE,
  INDEX `idx_name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '部门表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of department
-- ----------------------------
INSERT INTO `department` VALUES (1, 0, '测试公司', 1, 0, '2026-08-06 15:03:27', '2026-08-06 16:36:37', NULL);
INSERT INTO `department` VALUES (2, 1, 'IT技术部', 1, 0, '2026-08-06 15:03:27', '2026-08-06 16:48:17', NULL);
INSERT INTO `department` VALUES (3, 1, '综合管理部', 1, 0, '2026-08-06 15:03:27', '2026-08-06 15:03:27', NULL);
INSERT INTO `department` VALUES (4, 2, '前端开发组', 1, 0, '2026-08-06 15:03:27', '2026-08-06 15:03:27', NULL);
INSERT INTO `department` VALUES (5, 2, '后端开发组', 1, 0, '2026-08-06 15:03:27', '2026-08-06 15:03:27', NULL);
INSERT INTO `department` VALUES (6, 3, '行政部', 1, 0, '2026-08-06 15:03:27', '2026-08-06 15:03:27', NULL);
INSERT INTO `department` VALUES (7, 3, '人事部', 1, 0, '2026-08-06 15:03:27', '2026-08-06 15:03:27', NULL);

-- ----------------------------
-- Table structure for department_leaders
-- ----------------------------
DROP TABLE IF EXISTS `department_leaders`;
CREATE TABLE `department_leaders`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `department_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '部门id',
  `leader_user_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '领导id',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_dept_id`(`department_id`) USING BTREE,
  INDEX `idx_leader_user_id`(`leader_user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '部门领导表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of department_leaders
-- ----------------------------
INSERT INTO `department_leaders` VALUES (7, 2, 1, '2026-08-06 16:48:17', '2026-08-06 16:48:17', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '导入记录表' ROW_FORMAT = DYNAMIC;

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
) ENGINE = InnoDB AUTO_INCREMENT = 71 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of menu
-- ----------------------------
INSERT INTO `menu` VALUES (1, 0, 1, 'home', 1, 1, '2025-05-23 15:37:03', '2026-08-04 17:15:38', NULL);
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
INSERT INTO `menu` VALUES (25, 3, 2, 'sys.menu.edit', 1, 1, '2025-05-21 10:30:24', '2026-08-04 17:21:54', NULL);
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
INSERT INTO `menu` VALUES (63, 0, 1, 'operatorLog', 1, 8, '2026-07-27 10:58:51', '2026-07-27 14:22:13', NULL);
INSERT INTO `menu` VALUES (64, 63, 2, 'sys.operatorLog.detail', 1, 0, '2026-07-27 14:17:11', '2026-07-27 14:17:11', NULL);
INSERT INTO `menu` VALUES (65, 63, 2, 'sys.operatorLog.del', 1, 0, '2026-07-27 14:29:56', '2026-07-27 14:29:56', NULL);
INSERT INTO `menu` VALUES (66, 63, 2, 'sys.operatorLog.batchDel', 1, 0, '2026-07-27 14:30:40', '2026-07-27 14:31:31', NULL);
INSERT INTO `menu` VALUES (67, 2, 1, 'systemDepartment', 1, 0, '2026-08-06 15:45:45', '2026-08-06 15:45:45', NULL);
INSERT INTO `menu` VALUES (68, 67, 2, 'sys.dept.add', 1, 0, '2026-08-06 15:55:23', '2026-08-06 15:55:47', NULL);
INSERT INTO `menu` VALUES (69, 67, 2, 'sys.dept.edit', 1, 0, '2026-08-06 15:58:32', '2026-08-06 16:01:58', NULL);
INSERT INTO `menu` VALUES (70, 67, 2, 'sys.dept.addChildren', 1, 0, '2026-08-06 16:03:34', '2026-08-06 16:03:34', NULL);
INSERT INTO `menu` VALUES (71, 67, 2, 'sys.dept.del', 1, 0, '2026-08-06 16:04:07', '2026-08-06 16:04:07', NULL);

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
  `trans_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '翻译键',
  `auth_value` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '权限标识',
  `is_link` tinyint(3) UNSIGNED NOT NULL DEFAULT 2 COMMENT '是否为链接 1=是 2=否',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_menu_id`(`menu_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 52 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单功能表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of menu_actions
-- ----------------------------
INSERT INTO `menu_actions` VALUES (1, 24, 1, 'btn', 'primary', 'default', 2, '新增菜单', '', 'sys.menu.add', 2, '2025-05-21 10:24:14', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (4, 27, 2, 'btn', 'danger', 'small', 1, '删除', '', 'sys.menu.del', 2, '2025-05-21 10:30:49', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (8, 32, 1, 'btn', 'primary', 'default', 2, '新增用户', '', 'sys.user.add', 2, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (9, 33, 1, 'btn', 'danger', 'default', 2, '批量删除', '', 'sys.user.batchDel', 2, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (10, 34, 2, 'btn', 'primary', 'small', 2, '编辑', '', 'sys.user.edit', 2, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (11, 35, 2, 'btn', 'danger', 'small', 2, '删除', '', 'sys.user.del', 2, '2025-06-16 08:57:04', '2026-07-16 14:40:54', NULL);
INSERT INTO `menu_actions` VALUES (13, 37, 1, 'btn', 'primary', 'default', 2, '新增角色', '', 'sys.role.add', 2, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (14, 38, 2, 'btn', 'primary', 'small', 2, '编辑', '', 'sys.role.edit', 2, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (15, 39, 2, 'btn', 'danger', 'small', 2, '删除', '', 'sys.role.del', 2, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (16, 40, 1, 'btn', 'primary', 'default', 2, '新增字典', '', 'sys.dic.add', 2, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (17, 41, 2, 'btn', 'primary', 'small', 2, '编辑', '', 'sys.dic.edit', 2, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (18, 42, 2, 'btn', 'danger', 'small', 2, '删除', '', 'sys.dic.del', 2, '2025-06-16 08:57:04', '2026-07-15 15:59:14', NULL);
INSERT INTO `menu_actions` VALUES (19, 43, 1, 'btn', 'primary', 'default', 2, '新增文章', '', 'article.add', 2, '2025-06-16 08:57:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (20, 44, 2, 'btn', 'primary', 'small', 2, '编辑', '', 'article.edit', 2, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (21, 45, 2, 'btn', 'danger', 'small', 2, '删除', '', 'article.del', 2, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_actions` VALUES (22, 48, 1, 'btn', 'primary', 'default', 2, '新增配置', '', 'sys.config.add', 2, '2026-07-08 16:00:42', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (23, 49, 2, 'btn', 'primary', 'small', 2, '编辑', '', 'sys.config.edit', 2, '2026-07-08 16:01:29', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (24, 50, 2, 'btn', 'danger', 'small', 2, '删除', '', 'sys.config.del', 2, '2026-07-08 16:03:15', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (25, 51, 1, 'btn', 'primary', 'default', 2, '新增分类', '', 'sys.configCategory.add', 2, '2026-07-09 11:11:01', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (26, 52, 2, 'btn', 'primary', 'small', 2, '编辑', '', 'sys.configCategory.edit', 2, '2026-07-09 11:11:36', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (27, 53, 2, 'btn', 'danger', 'small', 2, '删除', '', 'sys.configCategory.del', 2, '2026-07-09 11:12:04', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_actions` VALUES (28, 54, 1, 'btn', 'primary', 'default', 2, '用户导入', '', 'sys.user.import', 2, '2026-07-10 13:34:50', '2026-07-10 13:34:50', NULL);
INSERT INTO `menu_actions` VALUES (33, 58, 2, 'btn', 'primary', 'small', 2, '新增子集', '', 'sys.menu.addChildren', 2, '2026-07-14 09:36:51', '2026-07-14 09:36:51', NULL);
INSERT INTO `menu_actions` VALUES (34, 59, 2, 'btn', 'primary', 'small', 2, '新增子集', '', 'sys.dic.addChildren', 2, '2026-07-15 15:58:56', '2026-07-15 15:58:56', NULL);
INSERT INTO `menu_actions` VALUES (35, 60, 2, 'btn', 'primary', 'small', 2, '更新密码', '', 'sys.user.password', 2, '2026-07-16 14:41:52', '2026-07-16 14:42:32', NULL);
INSERT INTO `menu_actions` VALUES (39, 62, 2, 'btn', 'danger', 'small', 2, '删除', '', 'sys.user.importRecords.delete', 2, '2026-07-17 15:45:23', '2026-07-17 15:45:23', NULL);
INSERT INTO `menu_actions` VALUES (40, 61, 2, 'btn', 'primary', 'small', 2, '明细', '', 'sys.user.importRecords.detail', 2, '2026-07-17 15:47:14', '2026-07-17 15:47:14', NULL);
INSERT INTO `menu_actions` VALUES (41, 56, 1, 'btn', 'primary', 'default', 2, '导入记录', '', 'sys.user.importRecords', 2, '2026-07-17 15:56:01', '2026-07-17 15:56:01', NULL);
INSERT INTO `menu_actions` VALUES (42, 64, 2, 'btn', 'primary', 'small', 2, '详情', '', 'sys.operatorLog.detail', 2, '2026-07-27 14:17:11', '2026-07-27 14:17:11', NULL);
INSERT INTO `menu_actions` VALUES (43, 65, 2, 'btn', 'danger', 'small', 2, '删除', '', 'sys.operatorLog.del', 2, '2026-07-27 14:29:56', '2026-07-27 14:29:56', NULL);
INSERT INTO `menu_actions` VALUES (45, 66, 1, 'btn', 'danger', 'default', 2, '批量删除', '', 'sys.operatorLog.batchDel', 2, '2026-07-27 14:31:31', '2026-07-27 14:31:31', NULL);
INSERT INTO `menu_actions` VALUES (46, 25, 2, 'btn', 'primary', 'small', 2, '编辑', '', 'sys.menu.edit', 2, '2026-08-04 17:21:54', '2026-08-04 17:21:54', NULL);
INSERT INTO `menu_actions` VALUES (48, 68, 1, 'btn', 'primary', 'default', 2, '新增部门', '', 'sys.dept.add', 2, '2026-08-06 15:55:47', '2026-08-06 15:55:47', NULL);
INSERT INTO `menu_actions` VALUES (50, 69, 2, 'btn', 'primary', 'small', 2, '编辑', '', 'sys.dept.edit', 2, '2026-08-06 16:01:58', '2026-08-06 16:01:58', NULL);
INSERT INTO `menu_actions` VALUES (51, 70, 2, 'btn', 'primary', 'small', 2, '新增子集', '', 'sys.dept.addChildren', 2, '2026-08-06 16:03:34', '2026-08-06 16:03:34', NULL);
INSERT INTO `menu_actions` VALUES (52, 71, 2, 'btn', 'danger', 'small', 2, '删除', '', 'sys.dept.del', 2, '2026-08-06 16:04:07', '2026-08-06 16:04:07', NULL);

-- ----------------------------
-- Table structure for menu_meta
-- ----------------------------
DROP TABLE IF EXISTS `menu_meta`;
CREATE TABLE `menu_meta`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `menu_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '菜单id',
  `title` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '菜单名称',
  `trans_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '翻译键',
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
) ENGINE = InnoDB AUTO_INCREMENT = 28 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单元数据表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of menu_meta
-- ----------------------------
INSERT INTO `menu_meta` VALUES (2, 2, '系统设置', 'message.router.system', 'iconfont icon-xitongshezhi', '/system', '/system/menu', 'layouts/routerView/parent', 2, 1, 2, '', 2, '2025-05-23 15:39:37', '2025-05-27 16:49:52', NULL);
INSERT INTO `menu_meta` VALUES (3, 3, '菜单管理', 'message.router.systemMenu', 'iconfont icon-caidan', '/system/menu', '', 'system/menu/index', 2, 1, 2, '', 2, '2025-05-23 15:41:38', '2025-06-11 17:17:14', NULL);
INSERT INTO `menu_meta` VALUES (4, 4, '用户管理', 'message.router.systemUser', 'iconfont icon-icon-', '/system/user', '', 'system/user/index', 2, 1, 2, '', 2, '2025-05-23 23:26:38', '2025-06-11 17:17:29', NULL);
INSERT INTO `menu_meta` VALUES (5, 5, '角色管理', 'message.router.systemRole', 'fa fa-user-circle-o', '/system/role', '', 'system/role/index', 2, 1, 2, '', 2, '2025-05-25 14:37:04', '2025-06-11 17:17:36', NULL);
INSERT INTO `menu_meta` VALUES (6, 6, '字典管理', 'message.router.systemDic', 'ele-Collection', '/system/dic', '', 'system/dic/index', 2, 1, 2, '', 2, '2025-05-25 14:54:04', '2025-06-11 17:17:42', NULL);
INSERT INTO `menu_meta` VALUES (7, 10, '文章管理', 'message.article.title', 'ele-Collection', '/article', '', 'article/index', 2, 1, 2, '', 2, '2025-06-16 15:34:11', '2025-06-16 15:34:11', NULL);
INSERT INTO `menu_meta` VALUES (8, 20, '配置管理', '', 'iconfont icon-ico', '/system/config', '', 'layouts/routerView/parent', 2, 1, 2, '', 2, '2026-07-08 15:04:29', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_meta` VALUES (9, 21, '配置列表', '', 'iconfont icon-quanjushezhi_o', '/system/config/index', '', 'system/config/index', 2, 1, 2, '', 2, '2026-07-08 15:51:11', '2026-07-08 15:51:11', NULL);
INSERT INTO `menu_meta` VALUES (10, 22, '系统配置', '', 'iconfont icon--chaifenhang', '/system/config/setting', '', 'system/config/setting', 2, 1, 2, '', 2, '2026-07-08 15:54:52', '2026-07-15 16:08:00', NULL);
INSERT INTO `menu_meta` VALUES (11, 23, '配置分类', '', 'iconfont icon--chaifenlie', '/system/config-category/index', '', 'system/config/category/index', 2, 1, 2, '', 2, '2026-07-09 10:56:40', '2026-07-09 10:56:40', NULL);
INSERT INTO `menu_meta` VALUES (24, 63, '操作日志', '', 'ele-AlarmClock', '/operator-log', '', 'operator_log/index', 2, 1, 2, '', 2, '2026-07-27 14:22:13', '2026-07-27 14:22:13', NULL);
INSERT INTO `menu_meta` VALUES (27, 1, '首页', 'message.router.home', 'iconfont icon-shouye', '/home', '', 'home/index', 2, 1, 1, '', 2, '2026-08-04 17:15:38', '2026-08-04 17:15:38', NULL);
INSERT INTO `menu_meta` VALUES (28, 67, '部门管理', '', 'ele-Avatar', '/system/department', '', 'system/department/index', 2, 1, 2, '', 2, '2026-08-06 15:45:45', '2026-08-06 15:45:45', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of migrations
-- ----------------------------
INSERT INTO `migrations` VALUES (1, '20251212_create_user_table', '2025-12-12 17:04:27.313');

-- ----------------------------
-- Table structure for operator_log
-- ----------------------------
DROP TABLE IF EXISTS `operator_log`;
CREATE TABLE `operator_log`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'ip地址',
  `method` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '请求方式',
  `uri` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '路由地址',
  `lang` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '语言',
  `params` json NULL COMMENT '请求参数',
  `user_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作用户ID',
  `trace_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '追踪ID',
  `status_code` int(10) NOT NULL DEFAULT 0 COMMENT '响应状态码',
  `user_agent` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户代理',
  `cost_ms` float(10, 4) NOT NULL DEFAULT 0.0000 COMMENT '耗时(毫秒)',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id`) USING BTREE,
  INDEX `idx_trace_id`(`trace_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 25 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '操作日志表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of operator_log
-- ----------------------------
INSERT INTO `operator_log` VALUES (1, '127.0.0.1', 'GET', '/api/v1/role/1/menu', 'zh', '{}', 1, 'bb11add5-03e3-4786-9739-aad51c242eab', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36', 59.0000, '2026-07-27 15:06:31', '2026-07-27 15:06:31', NULL);
INSERT INTO `operator_log` VALUES (2, '127.0.0.1', 'GET', '/api/v1/operator-log', 'zh', '{\"page\": \"1\", \"pageSize\": \"10\"}', 1, '5cd0523e-4590-47ed-8b47-fa7ce266ce42', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36', 38.0000, '2026-07-27 15:06:36', '2026-07-27 15:06:36', NULL);
INSERT INTO `operator_log` VALUES (3, '127.0.0.1', 'GET', '/api/v1/operator-log', 'zh', '{\"page\": \"1\", \"pageSize\": \"10\"}', 1, '1526d79f-858b-46e2-86c7-7f7f41d4bc8b', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36', 3.0000, '2026-07-27 15:07:34', '2026-07-27 15:07:34', NULL);
INSERT INTO `operator_log` VALUES (4, '127.0.0.1', 'GET', '/api/v1/user/60', 'zh', '{}', 1, '1347b5cc-7a5a-4a3b-9af3-944ed1120362', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 5.0000, '2026-08-03 09:26:44', '2026-08-03 09:26:44', NULL);
INSERT INTO `operator_log` VALUES (5, '127.0.0.1', 'GET', '/api/v1/role', 'zh', '{\"page\": \"1\", \"notPage\": \"true\", \"pageSize\": \"10\"}', 1, '52f26573-92d1-4098-80e7-aab2862ec360', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 10.0000, '2026-08-03 09:26:44', '2026-08-03 09:26:44', NULL);
INSERT INTO `operator_log` VALUES (6, '127.0.0.1', 'PUT', '/api/v1/user/60', 'zh', '{\"id\": 60, \"age\": 28, \"email\": \"zhang1@example.com\", \"avatar\": \"\", \"gender\": 1, \"status\": 1, \"fullName\": \"张1\", \"nickname\": \"小张1\", \"username\": \"zhan\", \"userRoles\": []}', 1, '1157d8a1-ec59-4e4c-93df-2a62c3945861', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 75.0000, '2026-08-03 09:26:46', '2026-08-03 09:26:46', NULL);
INSERT INTO `operator_log` VALUES (7, '127.0.0.1', 'GET', '/api/v1/user', 'zh', '{\"page\": \"1\", \"pageSize\": \"10\"}', 1, '1b879aef-3a25-4383-b6d2-cef4e996d568', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 163.0000, '2026-08-03 09:26:46', '2026-08-03 09:26:46', NULL);
INSERT INTO `operator_log` VALUES (8, '127.0.0.1', 'DELETE', '/api/v1/user/60', 'zh', '{}', 1, 'e7dd3f16-3241-42ef-9c95-1e986184f625', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 64.0000, '2026-08-03 09:26:51', '2026-08-03 09:26:51', NULL);
INSERT INTO `operator_log` VALUES (9, '127.0.0.1', 'GET', '/api/v1/menu', 'zh', '{\"page\": \"1\", \"notPage\": \"true\", \"pageSize\": \"10\"}', 1, '70b7fb7d-c18f-4f7d-86c8-853d763d30cf', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 21.0000, '2026-08-04 15:31:22', '2026-08-04 15:31:22', NULL);
INSERT INTO `operator_log` VALUES (10, '127.0.0.1', 'GET', '/api/v1/permission', 'zh', '{\"page\": \"1\", \"notPage\": \"true\", \"pageSize\": \"100\"}', 1, '2764d9b8-fe5e-4f9d-9a28-7d83c7562601', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 8.0000, '2026-08-04 15:31:22', '2026-08-04 15:31:22', NULL);
INSERT INTO `operator_log` VALUES (11, '127.0.0.1', 'GET', '/api/v1/role/1', 'zh', '{}', 1, '680edda3-d9bd-4b4f-8de2-9cd909b2a78c', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 13.0000, '2026-08-04 15:31:22', '2026-08-04 15:31:22', NULL);
INSERT INTO `operator_log` VALUES (12, '127.0.0.1', 'PUT', '/api/v1/role/1', 'zh', '{\"id\": 1, \"desc\": \"超级管理员\", \"name\": \"admin\", \"status\": 1, \"roleMenus\": [{\"name\": \"admin\", \"menuId\": 1, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 2, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 3, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 58, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 24, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 25, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 27, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 4, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 54, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 56, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 62, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 61, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 32, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 33, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 34, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 60, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 35, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 5, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 37, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 38, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 39, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 6, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 41, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 40, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 59, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 42, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 20, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 23, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 51, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 52, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 53, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 21, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 49, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 50, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 48, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 22, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 10, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 45, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 43, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 44, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 63, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 64, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 65, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 66, \"roleId\": 1}], \"userRoles\": [{\"name\": \"admin\", \"roleId\": 1, \"userId\": 1}, {\"name\": \"admin\", \"roleId\": 1, \"userId\": 10}], \"rolePermissions\": [{\"roleId\": 1, \"permissionId\": 1}, {\"roleId\": 1, \"permissionId\": 3}, {\"roleId\": 1, \"permissionId\": 38}, {\"roleId\": 1, \"permissionId\": 19}, {\"roleId\": 1, \"permissionId\": 8}, {\"roleId\": 1, \"permissionId\": 6}, {\"roleId\": 1, \"permissionId\": 15}, {\"roleId\": 1, \"permissionId\": 39}, {\"roleId\": 1, \"permissionId\": 20}, {\"roleId\": 1, \"permissionId\": 17}, {\"roleId\": 1, \"permissionId\": 12}, {\"roleId\": 1, \"permissionId\": 16}, {\"roleId\": 1, \"permissionId\": 18}, {\"roleId\": 1, \"permissionId\": 21}, {\"roleId\": 1, \"permissionId\": 10}, {\"roleId\": 1, \"permissionId\": 28}, {\"roleId\": 1, \"permissionId\": 40}, {\"roleId\": 1, \"permissionId\": 2}, {\"roleId\": 1, \"permissionId\": 23}, {\"roleId\": 1, \"permissionId\": 9}, {\"roleId\": 1, \"permissionId\": 33}, {\"roleId\": 1, \"permissionId\": 25}, {\"roleId\": 1, \"permissionId\": 46}, {\"roleId\": 1, \"permissionId\": 47}, {\"roleId\": 1, \"permissionId\": 44}, {\"roleId\": 1, \"permissionId\": 45}, {\"roleId\": 1, \"permissionId\": 51}, {\"roleId\": 1, \"permissionId\": 50}, {\"roleId\": 1, \"permissionId\": 43}, {\"roleId\": 1, \"permissionId\": 42}, {\"roleId\": 1, \"permissionId\": 24}, {\"roleId\": 1, \"permissionId\": 31}, {\"roleId\": 1, \"permissionId\": 32}, {\"roleId\": 1, \"permissionId\": 34}, {\"roleId\": 1, \"permissionId\": 35}, {\"roleId\": 1, \"permissionId\": 26}, {\"roleId\": 1, \"permissionId\": 7}, {\"roleId\": 1, \"permissionId\": 4}, {\"roleId\": 1, \"permissionId\": 11}, {\"roleId\": 1, \"permissionId\": 27}, {\"roleId\": 1, \"permissionId\": 37}, {\"roleId\": 1, \"permissionId\": 13}, {\"roleId\": 1, \"permissionId\": 29}, {\"roleId\": 1, \"permissionId\": 41}, {\"roleId\": 1, \"permissionId\": 14}, {\"roleId\": 1, \"permissionId\": 30}, {\"roleId\": 1, \"permissionId\": 5}, {\"roleId\": 1, \"permissionId\": 36}, {\"roleId\": 1, \"permissionId\": 22}]}', 1, 'b79f74d4-2b54-4b31-aa8c-f44457e18834', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 96.0000, '2026-08-04 15:31:29', '2026-08-04 15:31:29', NULL);
INSERT INTO `operator_log` VALUES (13, '127.0.0.1', 'GET', '/api/v1/role', 'zh', '{\"page\": \"1\", \"pageSize\": \"10\"}', 1, '4c6157e3-e585-4f23-901c-c59fc75c9c02', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 17.0000, '2026-08-04 15:31:29', '2026-08-04 15:31:29', NULL);
INSERT INTO `operator_log` VALUES (14, '127.0.0.1', 'GET', '/api/v1/config-category', 'zh', '{\"page\": \"1\", \"pageSize\": \"10\"}', 1, 'bb58db02-87be-4ff9-b326-5011ce6805e1', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 8.0000, '2026-08-04 15:31:44', '2026-08-04 15:31:44', NULL);
INSERT INTO `operator_log` VALUES (15, '127.0.0.1', 'GET', '/api/v1/system-config', 'zh', '{\"page\": \"1\", \"pageSize\": \"10\"}', 1, '33df760a-e862-4e38-b8fb-3e79def760f4', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 4.0000, '2026-08-04 15:31:47', '2026-08-04 15:31:47', NULL);
INSERT INTO `operator_log` VALUES (16, '127.0.0.1', 'GET', '/api/v1/config-category', 'zh', '{\"page\": \"1\", \"notPage\": \"true\", \"pageSize\": \"100\"}', 1, '5690b62c-6121-45a9-97a9-fa1caabcd986', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 7.0000, '2026-08-04 15:31:47', '2026-08-04 15:31:47', NULL);
INSERT INTO `operator_log` VALUES (17, '127.0.0.1', 'GET', '/api/v1/dashboard/cards', 'zh', '{}', 1, '391907d3-093b-4c22-b6a4-d0e73b28157e', 403, 'PostmanRuntime/7.55.1', 2.0000, '2026-08-04 16:20:53', '2026-08-04 16:20:53', NULL);
INSERT INTO `operator_log` VALUES (18, '127.0.0.1', 'GET', '/api/v1/menu', 'zh', '{\"page\": \"1\", \"notPage\": \"true\", \"pageSize\": \"10\"}', 1, '6fa7c047-bab0-4e64-a95c-b2ddbdf32d11', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 6.0000, '2026-08-04 16:21:01', '2026-08-04 16:21:01', NULL);
INSERT INTO `operator_log` VALUES (19, '127.0.0.1', 'GET', '/api/v1/permission', 'zh', '{\"page\": \"1\", \"notPage\": \"true\", \"pageSize\": \"100\"}', 1, '639c0fab-caba-4f68-b85d-41a89a7b77cc', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 3.0000, '2026-08-04 16:21:01', '2026-08-04 16:21:01', NULL);
INSERT INTO `operator_log` VALUES (20, '127.0.0.1', 'GET', '/api/v1/role/1', 'zh', '{}', 1, '768b04c0-cd8b-458f-9462-7b2bdc99950e', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 8.0000, '2026-08-04 16:21:01', '2026-08-04 16:21:01', NULL);
INSERT INTO `operator_log` VALUES (21, '127.0.0.1', 'PUT', '/api/v1/role/1', 'zh', '{\"id\": 1, \"desc\": \"超级管理员\", \"name\": \"admin\", \"status\": 1, \"roleMenus\": [{\"name\": \"admin\", \"menuId\": 1, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 2, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 3, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 58, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 24, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 25, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 27, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 4, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 54, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 56, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 62, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 61, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 32, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 33, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 34, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 60, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 35, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 5, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 37, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 38, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 39, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 6, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 41, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 40, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 59, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 42, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 20, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 23, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 51, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 52, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 53, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 21, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 49, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 50, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 48, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 22, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 10, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 45, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 43, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 44, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 63, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 64, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 65, \"roleId\": 1}, {\"name\": \"admin\", \"menuId\": 66, \"roleId\": 1}], \"userRoles\": [{\"name\": \"admin\", \"roleId\": 1, \"userId\": 1}, {\"name\": \"admin\", \"roleId\": 1, \"userId\": 10}], \"rolePermissions\": [{\"roleId\": 1, \"permissionId\": 1}, {\"roleId\": 1, \"permissionId\": 3}, {\"roleId\": 1, \"permissionId\": 38}, {\"roleId\": 1, \"permissionId\": 19}, {\"roleId\": 1, \"permissionId\": 8}, {\"roleId\": 1, \"permissionId\": 6}, {\"roleId\": 1, \"permissionId\": 15}, {\"roleId\": 1, \"permissionId\": 39}, {\"roleId\": 1, \"permissionId\": 20}, {\"roleId\": 1, \"permissionId\": 17}, {\"roleId\": 1, \"permissionId\": 12}, {\"roleId\": 1, \"permissionId\": 16}, {\"roleId\": 1, \"permissionId\": 18}, {\"roleId\": 1, \"permissionId\": 21}, {\"roleId\": 1, \"permissionId\": 10}, {\"roleId\": 1, \"permissionId\": 28}, {\"roleId\": 1, \"permissionId\": 40}, {\"roleId\": 1, \"permissionId\": 2}, {\"roleId\": 1, \"permissionId\": 23}, {\"roleId\": 1, \"permissionId\": 9}, {\"roleId\": 1, \"permissionId\": 33}, {\"roleId\": 1, \"permissionId\": 25}, {\"roleId\": 1, \"permissionId\": 46}, {\"roleId\": 1, \"permissionId\": 47}, {\"roleId\": 1, \"permissionId\": 44}, {\"roleId\": 1, \"permissionId\": 45}, {\"roleId\": 1, \"permissionId\": 43}, {\"roleId\": 1, \"permissionId\": 42}, {\"roleId\": 1, \"permissionId\": 24}, {\"roleId\": 1, \"permissionId\": 31}, {\"roleId\": 1, \"permissionId\": 32}, {\"roleId\": 1, \"permissionId\": 34}, {\"roleId\": 1, \"permissionId\": 35}, {\"roleId\": 1, \"permissionId\": 26}, {\"roleId\": 1, \"permissionId\": 7}, {\"roleId\": 1, \"permissionId\": 4}, {\"roleId\": 1, \"permissionId\": 11}, {\"roleId\": 1, \"permissionId\": 27}, {\"roleId\": 1, \"permissionId\": 37}, {\"roleId\": 1, \"permissionId\": 13}, {\"roleId\": 1, \"permissionId\": 29}, {\"roleId\": 1, \"permissionId\": 41}, {\"roleId\": 1, \"permissionId\": 14}, {\"roleId\": 1, \"permissionId\": 30}, {\"roleId\": 1, \"permissionId\": 5}, {\"roleId\": 1, \"permissionId\": 36}, {\"roleId\": 1, \"permissionId\": 22}, {\"roleId\": 1, \"permissionId\": 54}, {\"roleId\": 1, \"permissionId\": 52}, {\"roleId\": 1, \"permissionId\": 53}]}', 1, 'ba17bdc1-f99e-4a2a-9fdb-d9bde0b59244', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 88.0000, '2026-08-04 16:21:09', '2026-08-04 16:21:09', NULL);
INSERT INTO `operator_log` VALUES (22, '127.0.0.1', 'GET', '/api/v1/role', 'zh', '{\"page\": \"1\", \"pageSize\": \"10\"}', 1, 'a2a65754-88f4-4008-adb8-897a86845778', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 5.0000, '2026-08-04 16:21:10', '2026-08-04 16:21:10', NULL);
INSERT INTO `operator_log` VALUES (23, '127.0.0.1', 'GET', '/api/v1/role', 'zh', '{\"page\": \"1\", \"notPage\": \"true\", \"pageSize\": \"10\"}', 1, '1f78d107-31be-46fd-9880-b7515b22e545', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 34.0000, '2026-08-04 16:28:09', '2026-08-04 16:28:09', NULL);
INSERT INTO `operator_log` VALUES (24, '127.0.0.1', 'POST', '/api/v1/user', 'zh', '{\"age\": 30, \"email\": \"test321@qq.com\", \"avatar\": \"\", \"gender\": 1, \"status\": 1, \"fullName\": \"测试321\", \"nickname\": \"测试321\", \"password\": \"123456\", \"username\": \"测试321\", \"userRoles\": [{\"name\": \"test\", \"roleId\": 2, \"userId\": 0}]}', 1, 'bd1a589e-a21a-4168-a953-f86622c1a036', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 168.0000, '2026-08-04 16:28:39', '2026-08-04 16:28:39', NULL);
INSERT INTO `operator_log` VALUES (25, '127.0.0.1', 'GET', '/api/v1/user', 'zh', '{\"page\": \"1\", \"pageSize\": \"10\"}', 1, 'fa7263d3-0b85-4bbf-88b1-06972036b36c', 200, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36', 27.0000, '2026-08-04 16:28:39', '2026-08-04 16:28:39', NULL);

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
  UNIQUE INDEX `uk_key`(`key`) USING BTREE,
  INDEX `idx_key`(`key`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 59 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '权限表' ROW_FORMAT = DYNAMIC;

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
INSERT INTO `permission` VALUES (44, 'GET:/api/v1/operator-log/:id', 'GET', '/api/v1/operator-log/:id', '2026-07-27 11:01:52', '2026-07-27 11:01:52', NULL);
INSERT INTO `permission` VALUES (45, 'POST:/api/v1/operator-log/batch-delete', 'POST', '/api/v1/operator-log/batch-delete', '2026-07-27 11:01:52', '2026-07-27 11:01:52', NULL);
INSERT INTO `permission` VALUES (46, 'GET:/api/v1/operator-log', 'GET', '/api/v1/operator-log', '2026-07-27 11:01:52', '2026-07-27 11:01:52', NULL);
INSERT INTO `permission` VALUES (47, 'DELETE:/api/v1/operator-log/:id', 'DELETE', '/api/v1/operator-log/:id', '2026-07-27 11:01:52', '2026-07-27 11:01:52', NULL);
INSERT INTO `permission` VALUES (52, 'GET:/api/v1/dashboard/statistics', 'GET', '/api/v1/dashboard/statistics', '2026-08-04 16:04:29', '2026-08-04 16:04:29', NULL);
INSERT INTO `permission` VALUES (53, 'GET:/api/v1/dashboard/system-resource', 'GET', '/api/v1/dashboard/system-resource', '2026-08-04 16:04:29', '2026-08-04 16:04:29', NULL);
INSERT INTO `permission` VALUES (54, 'GET:/api/v1/dashboard/cards', 'GET', '/api/v1/dashboard/cards', '2026-08-04 16:04:29', '2026-08-04 16:04:29', NULL);
INSERT INTO `permission` VALUES (55, 'POST:/api/v1/department', 'POST', '/api/v1/department', '2026-08-06 15:07:45', '2026-08-06 15:07:45', NULL);
INSERT INTO `permission` VALUES (56, 'PUT:/api/v1/department/:id', 'PUT', '/api/v1/department/:id', '2026-08-06 15:07:45', '2026-08-06 15:07:45', NULL);
INSERT INTO `permission` VALUES (57, 'GET:/api/v1/department/:id', 'GET', '/api/v1/department/:id', '2026-08-06 15:07:45', '2026-08-06 15:07:45', NULL);
INSERT INTO `permission` VALUES (58, 'DELETE:/api/v1/department/:id', 'DELETE', '/api/v1/department/:id', '2026-08-06 15:07:45', '2026-08-06 15:07:45', NULL);
INSERT INTO `permission` VALUES (59, 'GET:/api/v1/department', 'GET', '/api/v1/department', '2026-08-06 15:07:45', '2026-08-06 15:07:45', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 719 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色菜单表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role_menus
-- ----------------------------
INSERT INTO `role_menus` VALUES (303, 2, 24, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `role_menus` VALUES (304, 2, 27, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `role_menus` VALUES (305, 2, 2, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `role_menus` VALUES (306, 2, 3, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `role_menus` VALUES (666, 2, 1, 'test', '2026-08-04 17:15:38', '2026-08-04 17:15:38', NULL);
INSERT INTO `role_menus` VALUES (669, 1, 1, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (670, 1, 2, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (671, 1, 3, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (672, 1, 58, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (673, 1, 24, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (674, 1, 25, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (675, 1, 27, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (676, 1, 4, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (677, 1, 54, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (678, 1, 56, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (679, 1, 62, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (680, 1, 61, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (681, 1, 32, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (682, 1, 33, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (683, 1, 34, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (684, 1, 60, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (685, 1, 35, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (686, 1, 5, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (687, 1, 37, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (688, 1, 38, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (689, 1, 39, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (690, 1, 6, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (691, 1, 41, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (692, 1, 40, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (693, 1, 59, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (694, 1, 42, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (695, 1, 20, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (696, 1, 23, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (697, 1, 51, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (698, 1, 52, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (699, 1, 53, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (700, 1, 21, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (701, 1, 49, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (702, 1, 50, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (703, 1, 48, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (704, 1, 22, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (705, 1, 10, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (706, 1, 45, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (707, 1, 43, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (708, 1, 44, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (709, 1, 63, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (710, 1, 64, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (711, 1, 65, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (712, 1, 66, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `role_menus` VALUES (713, 1, 67, 'admin', '2026-08-06 15:45:45', '2026-08-06 15:45:45', NULL);
INSERT INTO `role_menus` VALUES (715, 1, 68, 'admin', '2026-08-06 15:55:47', '2026-08-06 15:55:47', NULL);
INSERT INTO `role_menus` VALUES (717, 1, 69, 'admin', '2026-08-06 16:01:58', '2026-08-06 16:01:58', NULL);
INSERT INTO `role_menus` VALUES (718, 1, 70, 'admin', '2026-08-06 16:03:34', '2026-08-06 16:03:34', NULL);
INSERT INTO `role_menus` VALUES (719, 1, 71, 'admin', '2026-08-06 16:04:07', '2026-08-06 16:04:07', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 390 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色权限表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role_permissions
-- ----------------------------
INSERT INTO `role_permissions` VALUES (336, 1, 1, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (337, 1, 3, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (338, 1, 38, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (339, 1, 19, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (340, 1, 8, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (341, 1, 6, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (342, 1, 15, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (343, 1, 39, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (344, 1, 20, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (345, 1, 17, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (346, 1, 54, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (347, 1, 52, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (348, 1, 53, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (349, 1, 12, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (350, 1, 16, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (351, 1, 18, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (352, 1, 21, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (353, 1, 10, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (354, 1, 28, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (355, 1, 40, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (356, 1, 2, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (357, 1, 23, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (358, 1, 9, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (359, 1, 33, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (360, 1, 25, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (361, 1, 46, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (362, 1, 47, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (363, 1, 44, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (364, 1, 45, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (365, 1, 43, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (366, 1, 42, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (367, 1, 24, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (368, 1, 31, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (369, 1, 32, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (370, 1, 34, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (371, 1, 35, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (372, 1, 26, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (373, 1, 7, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (374, 1, 4, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (375, 1, 11, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (376, 1, 27, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (377, 1, 37, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (378, 1, 13, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (379, 1, 29, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (380, 1, 41, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (381, 1, 14, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (382, 1, 30, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (383, 1, 5, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (384, 1, 36, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (385, 1, 22, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (386, 1, 59, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (387, 1, 55, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (388, 1, 58, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (389, 1, 57, '2026-08-06 15:12:16', NULL);
INSERT INTO `role_permissions` VALUES (390, 1, 56, '2026-08-06 15:12:16', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of roles
-- ----------------------------
INSERT INTO `roles` VALUES (1, 'admin', '超级管理员', 1, '2025-05-26 16:52:43', '2026-08-06 15:12:16', NULL);
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
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统配置表' ROW_FORMAT = DYNAMIC;

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
) ENGINE = InnoDB AUTO_INCREMENT = 63 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, 'https://cdn.qitx.net/local/myblog/user_header_image/20230517/577a53d123bc4c4f19db0cb2c6c980a8.jpg', 'admin', '超级管理员', 'dsx.emil@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', '大师兄', 1, 31, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (2, '', 'test2', '李四1', 'ls@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2026-07-15 16:08:00', NULL);
INSERT INTO `user` VALUES (10, '', 'dsx', '大师兄111', 'dsx@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄', 1, 0, 1, '2024-07-22 17:34:36', '2026-08-07 16:06:05', NULL);
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
INSERT INTO `user` VALUES (61, '', '测试321', '测试321', '', '$2a$10$7sBfSDNQULz1zOFvuxYUHOJa/eSfkzMx1XnafpkpNjjJ8TC1kH1jC', '测试321', 1, 30, 1, '2026-08-04 16:28:39', '2026-08-04 16:28:39', NULL);

-- ----------------------------
-- Table structure for user_departments
-- ----------------------------
DROP TABLE IF EXISTS `user_departments`;
CREATE TABLE `user_departments`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
  `department_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '部门id',
  `is_main` tinyint(3) UNSIGNED NOT NULL DEFAULT 2 COMMENT '是否是主部门 1=是 2=否',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id`) USING BTREE,
  INDEX `idx_dept_id`(`department_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户部门表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user_departments
-- ----------------------------
INSERT INTO `user_departments` VALUES (1, 1, 2, 1, '2026-08-06 15:03:27', '2026-08-06 15:03:27', NULL);
INSERT INTO `user_departments` VALUES (2, 1, 5, 2, '2026-08-06 15:03:27', '2026-08-06 15:03:27', NULL);
INSERT INTO `user_departments` VALUES (5, 10, 2, 1, '2026-08-07 16:06:06', '2026-08-07 16:06:06', NULL);
INSERT INTO `user_departments` VALUES (6, 10, 5, 2, '2026-08-07 16:06:06', '2026-08-07 16:06:06', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 59 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户角色表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user_roles
-- ----------------------------
INSERT INTO `user_roles` VALUES (35, 2, 2, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `user_roles` VALUES (55, 61, 2, 'test', '2026-08-04 16:28:39', '2026-08-04 16:28:39', NULL);
INSERT INTO `user_roles` VALUES (56, 1, 1, 'admin', '2026-08-06 15:12:16', '2026-08-06 15:12:16', NULL);
INSERT INTO `user_roles` VALUES (58, 10, 2, 'test', '2026-08-07 16:06:05', '2026-08-07 16:06:05', NULL);
INSERT INTO `user_roles` VALUES (59, 10, 1, 'admin', '2026-08-07 16:06:05', '2026-08-07 16:06:05', NULL);

SET FOREIGN_KEY_CHECKS = 1;
