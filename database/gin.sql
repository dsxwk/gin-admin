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

 Date: 07/04/2026 09:41:33
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for action_roles
-- ----------------------------
DROP TABLE IF EXISTS `action_roles`;
CREATE TABLE `action_roles`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `action_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '功能id',
  `role_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '角色id',
  `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '角色名称',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 50 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '功能角色表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of action_roles
-- ----------------------------
INSERT INTO `action_roles` VALUES (1, 1, 1, 'admin', '2025-05-26 17:48:41', '2025-05-26 17:48:43', '2025-06-13 16:26:48');
INSERT INTO `action_roles` VALUES (2, 2, 1, 'admin', '2025-05-26 17:48:41', '2025-05-26 17:48:41', '2025-06-13 16:27:31');
INSERT INTO `action_roles` VALUES (3, 3, 1, 'admin', '2025-05-26 17:48:41', '2025-05-26 17:48:41', '2025-06-13 16:27:43');
INSERT INTO `action_roles` VALUES (4, 4, 1, 'admin', '2025-05-26 17:48:41', '2025-05-26 17:48:41', '2025-06-13 16:27:51');
INSERT INTO `action_roles` VALUES (5, 6, 1, 'admin', '2025-06-11 11:46:16', '2025-06-11 11:46:16', '2025-06-11 11:51:53');
INSERT INTO `action_roles` VALUES (6, 6, 1, 'admin', '2025-06-11 11:46:16', '2025-06-11 11:46:16', '2025-06-11 11:51:53');
INSERT INTO `action_roles` VALUES (7, 6, 1, 'admin', '2025-06-11 14:13:56', '2025-06-11 14:13:56', '2025-06-11 14:13:56');
INSERT INTO `action_roles` VALUES (8, 6, 1, 'admin', '2025-06-11 14:13:56', '2025-06-11 14:13:56', '2025-06-11 14:14:18');
INSERT INTO `action_roles` VALUES (9, 6, 2, 'test', '2025-06-11 14:14:18', '2025-06-11 14:14:18', '2025-06-11 14:14:18');
INSERT INTO `action_roles` VALUES (10, 6, 2, 'test', '2025-06-11 14:14:18', '2025-06-11 14:14:18', '2025-06-11 14:14:25');
INSERT INTO `action_roles` VALUES (11, 6, 1, 'admin', '2025-06-11 14:14:25', '2025-06-11 14:14:25', '2025-06-11 14:14:25');
INSERT INTO `action_roles` VALUES (12, 6, 1, 'admin', '2025-06-11 14:14:25', '2025-06-11 14:14:25', '2025-06-13 16:28:07');
INSERT INTO `action_roles` VALUES (13, 1, 1, 'admin', '2025-06-13 16:26:48', '2025-06-13 16:26:48', '2025-06-13 16:26:48');
INSERT INTO `action_roles` VALUES (14, 1, 1, 'admin', '2025-06-13 16:26:48', '2025-06-13 16:26:48', NULL);
INSERT INTO `action_roles` VALUES (15, 2, 1, 'admin', '2025-06-13 16:27:31', '2025-06-13 16:27:31', '2025-06-13 16:27:31');
INSERT INTO `action_roles` VALUES (16, 2, 1, 'admin', '2025-06-13 16:27:31', '2025-06-13 16:27:31', '2025-06-13 16:50:04');
INSERT INTO `action_roles` VALUES (17, 3, 1, 'admin', '2025-06-13 16:27:43', '2025-06-13 16:27:43', '2025-06-13 16:27:43');
INSERT INTO `action_roles` VALUES (18, 3, 1, 'admin', '2025-06-13 16:27:43', '2025-06-13 16:27:43', NULL);
INSERT INTO `action_roles` VALUES (19, 4, 1, 'admin', '2025-06-13 16:27:51', '2025-06-13 16:27:51', '2025-06-13 16:27:51');
INSERT INTO `action_roles` VALUES (20, 4, 1, 'admin', '2025-06-13 16:27:51', '2025-06-13 16:27:51', NULL);
INSERT INTO `action_roles` VALUES (21, 6, 1, 'admin', '2025-06-13 16:28:07', '2025-06-13 16:28:07', '2025-06-13 16:28:07');
INSERT INTO `action_roles` VALUES (22, 6, 1, 'admin', '2025-06-13 16:28:07', '2025-06-13 16:28:07', NULL);
INSERT INTO `action_roles` VALUES (23, 2, 1, 'admin', '2025-06-13 16:50:04', '2025-06-13 16:50:04', '2025-06-13 16:50:04');
INSERT INTO `action_roles` VALUES (24, 2, 1, 'admin', '2025-06-13 16:50:04', '2025-06-13 16:50:04', '2025-06-13 16:56:30');
INSERT INTO `action_roles` VALUES (25, 2, 1, 'admin', '2025-06-13 16:56:30', '2025-06-13 16:56:30', '2025-06-13 16:56:30');
INSERT INTO `action_roles` VALUES (26, 2, 1, 'admin', '2025-06-13 16:56:30', '2025-06-13 16:56:30', '2025-06-13 17:01:21');
INSERT INTO `action_roles` VALUES (27, 2, 1, 'admin', '2025-06-13 17:01:21', '2025-06-13 17:01:21', '2025-06-13 17:01:21');
INSERT INTO `action_roles` VALUES (28, 2, 1, 'admin', '2025-06-13 17:01:21', '2025-06-13 17:01:21', '2025-06-13 17:02:12');
INSERT INTO `action_roles` VALUES (29, 2, 1, 'admin', '2025-06-13 17:02:12', '2025-06-13 17:02:12', '2025-06-13 17:02:12');
INSERT INTO `action_roles` VALUES (30, 2, 1, 'admin', '2025-06-13 17:02:12', '2025-06-13 17:02:12', '2025-06-13 17:07:48');
INSERT INTO `action_roles` VALUES (31, 2, 1, 'admin', '2025-06-13 17:07:48', '2025-06-13 17:07:48', '2025-06-13 17:07:48');
INSERT INTO `action_roles` VALUES (32, 2, 1, 'admin', '2025-06-13 17:07:48', '2025-06-13 17:07:48', '2025-06-13 17:17:32');
INSERT INTO `action_roles` VALUES (33, 2, 1, 'admin', '2025-06-13 17:17:32', '2025-06-13 17:17:32', '2025-06-13 17:17:32');
INSERT INTO `action_roles` VALUES (34, 2, 1, 'admin', '2025-06-13 17:17:32', '2025-06-13 17:17:32', NULL);
INSERT INTO `action_roles` VALUES (35, 7, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (36, 8, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (37, 9, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (38, 10, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (39, 11, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (40, 12, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (41, 13, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (42, 14, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (43, 15, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (44, 16, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (45, 17, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (46, 18, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (47, 19, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (48, 20, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (49, 21, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);
INSERT INTO `action_roles` VALUES (50, 22, 1, 'admin', '2025-06-16 08:53:37', '2025-06-16 08:53:37', NULL);

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
  `is_publish` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '是否发布 1=已发布 2=未发布',
  `tag` json NULL COMMENT '标签',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 14 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文章表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of article
-- ----------------------------
INSERT INTO `article` VALUES (1, 1, '标题1', '<p>测试1</p>', 1, 2, 1, '[\"测试标签1\", \"测试标签2\"]', '2023-09-19 11:43:58', '2025-07-14 15:10:18', NULL);
INSERT INTO `article` VALUES (13, 1, '标题1', '<p>内容13</p>', 0, 2, 1, '[\"测试标签11\", \"测试标签22\"]', '2024-07-22 11:21:18', '2025-06-17 10:27:27', NULL);
INSERT INTO `article` VALUES (14, 1, 'Go语言', '内容', 0, 0, 0, 'null', '2025-07-03 17:32:08', '2025-07-03 17:32:08', NULL);

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '分类表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of category
-- ----------------------------
INSERT INTO `category` VALUES (1, '分类名称', '2023-09-19 11:43:43', '2023-09-19 11:43:43', NULL);

-- ----------------------------
-- Table structure for dict
-- ----------------------------
DROP TABLE IF EXISTS `dict`;
CREATE TABLE `dict`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `pid` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父级id',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '字段名称(英文)',
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '字段名称(中文)',
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
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '字典表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of dict
-- ----------------------------
INSERT INTO `dict` VALUES (1, 0, 'gender', '性别', '', 1, 0, '{\"test\": 111, \"test2\": \"test222\"}', '', '2025-06-06 21:48:17', '2025-06-06 21:48:17', NULL);
INSERT INTO `dict` VALUES (2, 1, 'gender', '男', '1', 1, 0, '{\"a\": 11, \"b\": 22}', '测试', '2025-06-06 21:49:00', '2025-06-08 16:05:39', NULL);
INSERT INTO `dict` VALUES (3, 1, 'gender', '女', '2', 1, 0, '{\"test\": 111, \"test2\": \"test222\"}', '性别女', '2025-06-06 21:49:10', '2025-06-08 20:39:03', NULL);

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
  `meta` json NULL COMMENT 'meta',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_pid`(`pid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of menu
-- ----------------------------
INSERT INTO `menu` VALUES (1, 0, 'home', '/home', '', 'home/index', 2, 1, 0, '{\"icon\": \"iconfont icon-shouye\", \"roles\": [1], \"title\": \"message.router.home\", \"isHide\": false, \"isLink\": \"\", \"isAffix\": true, \"isIframe\": false, \"isKeepAlive\": true}', '2025-05-23 15:37:03', '2025-06-13 11:10:18', NULL);
INSERT INTO `menu` VALUES (2, 0, 'system', '/system', '/system/menu', 'layouts/routerView/parent', 2, 1, 0, '{\"icon\": \"iconfont icon-xitongshezhi\", \"roles\": [1], \"title\": \"message.router.system\", \"isHide\": false, \"isLink\": \"\", \"isAffix\": false, \"isIframe\": false, \"isKeepAlive\": true}', '2025-05-23 15:39:37', '2025-05-27 16:49:52', NULL);
INSERT INTO `menu` VALUES (3, 2, 'systemMenu', '/system/menu', '', 'system/menu/index', 2, 1, 0, '{\"icon\": \"iconfont icon-caidan\", \"roles\": [1], \"title\": \"message.router.systemMenu\", \"isHide\": false, \"isLink\": \"\", \"isAffix\": false, \"isIframe\": false, \"isKeepAlive\": true}', '2025-05-23 15:41:38', '2025-06-11 17:17:14', NULL);
INSERT INTO `menu` VALUES (4, 2, 'systemUser', '/system/user', '', 'system/user/index', 2, 1, 0, '{\"icon\": \"iconfont icon-icon-\", \"roles\": [1], \"title\": \"message.router.systemUser\", \"isHide\": false, \"isLink\": \"\", \"isAffix\": false, \"isIframe\": false, \"isKeepAlive\": true}', '2025-05-23 23:26:38', '2025-06-11 17:17:29', NULL);
INSERT INTO `menu` VALUES (5, 2, 'systemRole', '/system/role', '', 'system/role/index', 2, 1, 0, '{\"icon\": \"fa fa-user-circle-o\", \"roles\": [1], \"title\": \"message.router.systemRole\", \"isHide\": false, \"isLink\": \"\", \"isAffix\": false, \"isIframe\": false, \"isKeepAlive\": true}', '2025-05-25 14:37:04', '2025-06-11 17:17:36', NULL);
INSERT INTO `menu` VALUES (6, 2, 'systemDic', '/system/dic', '', 'system/dic/index', 2, 1, 0, '{\"icon\": \"ele-Collection\", \"roles\": [1], \"title\": \"message.router.systemDic\", \"isHide\": false, \"isLink\": \"\", \"isAffix\": false, \"isIframe\": false, \"isKeepAlive\": true}', '2025-05-25 14:54:04', '2025-06-11 17:17:42', NULL);
INSERT INTO `menu` VALUES (10, 0, 'article', '/article', '', 'article/index', 2, 1, 0, '{\"icon\": \"ele-Collection\", \"roles\": [1], \"title\": \"message.article.title\", \"isHide\": false, \"isLink\": \"\", \"isAffix\": false, \"isIframe\": false, \"authBtnList\": null, \"isKeepAlive\": true}', '2025-06-16 15:34:11', '2025-06-16 15:34:11', NULL);

-- ----------------------------
-- Table structure for menu_action
-- ----------------------------
DROP TABLE IF EXISTS `menu_action`;
CREATE TABLE `menu_action`  (
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
) ENGINE = InnoDB AUTO_INCREMENT = 22 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单功能表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of menu_action
-- ----------------------------
INSERT INTO `menu_action` VALUES (1, 0, 3, 1, 'btn', 'primary', 'small', 2, '新增菜单', 'sys.menu.add', 2, 0, '2025-05-21 10:24:14', '2025-06-13 16:26:48', NULL);
INSERT INTO `menu_action` VALUES (2, 0, 3, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.menu.edit', 2, 0, '2025-05-21 10:30:24', '2025-06-13 17:17:32', NULL);
INSERT INTO `menu_action` VALUES (3, 0, 3, 2, 'btn', 'primary', 'small', 2, '功能', 'sys.menu.action', 2, 0, '2025-05-21 10:30:37', '2025-06-13 16:27:43', NULL);
INSERT INTO `menu_action` VALUES (4, 0, 3, 2, 'btn', 'danger', 'small', 1, '删除', 'sys.menu.del', 2, 0, '2025-05-21 10:30:49', '2025-06-13 16:27:51', NULL);
INSERT INTO `menu_action` VALUES (5, 0, 3, 1, 'btn', 'primary', 'small', 2, '测试', '', 2, 0, '2025-06-03 16:48:56', '2025-06-03 16:48:56', '2025-06-03 16:55:38');
INSERT INTO `menu_action` VALUES (6, 3, 3, 1, 'btn', 'primary', 'small', 2, '新增功能', 'sys.menu.action.add', 2, 0, '2025-06-11 11:46:16', '2025-06-13 16:28:07', NULL);
INSERT INTO `menu_action` VALUES (7, 3, 3, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.menu.action.edit', 2, 0, '2025-06-13 11:36:56', '2025-06-13 16:28:14', NULL);
INSERT INTO `menu_action` VALUES (8, 3, 3, 2, 'btn', 'danger', 'small', 2, '删除', 'sys.menu.action.del', 2, 0, '2025-06-13 11:37:07', '2025-06-13 16:28:21', NULL);
INSERT INTO `menu_action` VALUES (9, 0, 4, 1, 'btn', 'primary', 'small', 2, '新增用户', 'sys.user.add', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_action` VALUES (10, 0, 4, 1, 'btn', 'danger', 'small', 2, '批量删除', 'sys.user.batchDel', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_action` VALUES (11, 0, 4, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.user.edit', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_action` VALUES (12, 0, 4, 2, 'btn', 'danger', 'small', 2, '删除', 'sys.user.del', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_action` VALUES (13, 0, 5, 1, 'btn', 'danger', 'small', 2, '批量删除', 'sys.role.batchDel', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_action` VALUES (14, 0, 5, 1, 'btn', 'primary', 'small', 2, '新增角色', 'sys.role.add', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_action` VALUES (15, 0, 5, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.role.edit', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_action` VALUES (16, 0, 5, 2, 'btn', 'danger', 'small', 2, '删除', 'sys.role.del', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_action` VALUES (17, 0, 6, 1, 'btn', 'primary', 'small', 2, '新增字典', 'sys.dic.add', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_action` VALUES (18, 0, 6, 2, 'btn', 'primary', 'small', 2, '编辑', 'sys.dic.edit', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_action` VALUES (19, 0, 6, 2, 'btn', 'danger', 'small', 2, '删除', 'sys.dic.del', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_action` VALUES (20, 0, 10, 1, 'btn', 'primary', 'small', 2, '新增文章', 'article.add', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_action` VALUES (21, 0, 10, 2, 'btn', 'primary', 'small', 2, '编辑', 'article.edit', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);
INSERT INTO `menu_action` VALUES (22, 0, 10, 2, 'btn', 'danger', 'small', 2, '删除', 'article.del', 2, 0, '2025-06-16 08:57:04', '2025-06-16 08:57:04', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of migrations
-- ----------------------------
INSERT INTO `migrations` VALUES (1, '20251212_create_user_table', '2025-12-12 17:04:27.313');

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
) ENGINE = InnoDB AUTO_INCREMENT = 37 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单角色表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role_menus
-- ----------------------------
INSERT INTO `role_menus` VALUES (1, 1, 1, 'admin', '2025-05-26 17:50:37', '2025-05-26 17:50:37', '2025-05-27 16:26:52');
INSERT INTO `role_menus` VALUES (2, 1, 2, 'admin', '2025-05-26 17:50:37', '2025-05-26 17:50:37', '2025-05-27 16:27:48');
INSERT INTO `role_menus` VALUES (3, 1, 3, 'admin', '2025-05-26 17:50:37', '2025-05-26 17:50:37', '2025-05-29 10:39:41');
INSERT INTO `role_menus` VALUES (4, 1, 4, 'admin', '2025-05-26 17:50:37', '2025-05-26 17:50:37', '2025-06-11 17:17:29');
INSERT INTO `role_menus` VALUES (5, 1, 5, 'admin', '2025-05-26 17:50:37', '2025-05-26 17:50:37', '2025-06-11 17:17:36');
INSERT INTO `role_menus` VALUES (6, 1, 6, 'admin', '2025-05-26 17:50:37', '2025-05-26 17:50:37', '2025-06-11 17:17:42');
INSERT INTO `role_menus` VALUES (17, 1, 3, 'admin', '2025-05-29 10:39:41', '2025-05-29 10:39:41', '2025-05-29 10:39:41');
INSERT INTO `role_menus` VALUES (18, 1, 3, 'admin', '2025-05-29 10:39:41', '2025-05-29 10:39:41', '2025-06-11 17:17:14');
INSERT INTO `role_menus` VALUES (19, 2, 7, 'test', '2025-06-06 10:21:50', '2025-06-06 10:21:50', '2025-06-06 10:21:50');
INSERT INTO `role_menus` VALUES (20, 2, 7, 'test', '2025-06-06 10:21:50', '2025-06-06 10:21:50', '2025-06-06 10:55:59');
INSERT INTO `role_menus` VALUES (21, 2, 8, 'test', '2025-06-06 10:35:00', '2025-06-06 10:35:00', '2025-06-06 10:55:33');
INSERT INTO `role_menus` VALUES (22, 2, 8, 'test', '2025-06-06 10:35:00', '2025-06-06 10:35:00', '2025-06-06 10:55:33');
INSERT INTO `role_menus` VALUES (23, 2, 9, 'test', '2025-06-06 11:01:30', '2025-06-06 11:01:30', '2025-06-06 11:01:42');
INSERT INTO `role_menus` VALUES (24, 2, 9, 'test', '2025-06-06 11:01:30', '2025-06-06 11:01:30', '2025-06-06 11:01:42');
INSERT INTO `role_menus` VALUES (25, 1, 3, 'admin', '2025-06-11 17:17:14', '2025-06-11 17:17:14', '2025-06-11 17:17:14');
INSERT INTO `role_menus` VALUES (26, 1, 3, 'admin', '2025-06-11 17:17:14', '2025-06-11 17:17:14', NULL);
INSERT INTO `role_menus` VALUES (27, 1, 4, 'admin', '2025-06-11 17:17:29', '2025-06-11 17:17:29', '2025-06-11 17:17:29');
INSERT INTO `role_menus` VALUES (28, 1, 4, 'admin', '2025-06-11 17:17:29', '2025-06-11 17:17:29', NULL);
INSERT INTO `role_menus` VALUES (29, 1, 5, 'admin', '2025-06-11 17:17:36', '2025-06-11 17:17:36', '2025-06-11 17:17:36');
INSERT INTO `role_menus` VALUES (30, 1, 5, 'admin', '2025-06-11 17:17:36', '2025-06-11 17:17:36', NULL);
INSERT INTO `role_menus` VALUES (31, 1, 6, 'admin', '2025-06-11 17:17:42', '2025-06-11 17:17:42', '2025-06-11 17:17:42');
INSERT INTO `role_menus` VALUES (32, 1, 6, 'admin', '2025-06-11 17:17:42', '2025-06-11 17:17:42', NULL);
INSERT INTO `role_menus` VALUES (33, 1, 1, 'admin', '2025-06-13 11:10:18', '2025-06-13 11:10:18', '2025-06-13 11:10:18');
INSERT INTO `role_menus` VALUES (34, 1, 1, 'admin', '2025-06-13 11:10:18', '2025-06-13 11:10:18', NULL);
INSERT INTO `role_menus` VALUES (35, 1, 10, 'admin', '2025-06-16 15:34:11', '2025-06-16 15:34:11', NULL);
INSERT INTO `role_menus` VALUES (36, 1, 10, 'admin', '2025-06-16 15:34:11', '2025-06-16 15:34:11', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of roles
-- ----------------------------
INSERT INTO `roles` VALUES (1, 'admin', '超级管理员', 1, '2025-05-26 16:52:43', '2025-05-26 16:52:50', NULL);
INSERT INTO `roles` VALUES (2, 'test', '测试', 1, '2025-05-28 10:47:22', '2025-05-28 10:47:22', NULL);
INSERT INTO `roles` VALUES (3, '1111', '1111', 1, '2025-06-06 09:51:06', '2025-06-06 09:51:06', '2025-06-06 09:51:10');

-- ----------------------------
-- Table structure for system_config
-- ----------------------------
DROP TABLE IF EXISTS `system_config`;
CREATE TABLE `system_config`  (
  `id` smallint(5) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
  `cn_name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '中文名称',
  `en_name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '英文名称',
  `default_value` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '默认值',
  `option_value` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '可选值',
  `type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '配置类型 1=输入框 2=单选 3=复选 4=下拉菜单 5=文本域 6=附件',
  `category` tinyint(1) NOT NULL DEFAULT 1 COMMENT '配置分类 1=基本信息 2=联系方式 3=seo设置 ',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uni_en_name`(`en_name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统配置表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of system_config
-- ----------------------------
INSERT INTO `system_config` VALUES (1, '网站域名', 'web_domain', 'www.a.com', '', 1, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (2, '关闭站点', 'is_open_site', '开启', '关闭,开启', 2, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (3, '网站Logo', 'site_logo', '', '', 6, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (4, '邮件端口', 'email_port', '465', '', 1, 2, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (5, '邮件标题', 'email_title', '【xxx】验证码', '', 1, 2, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (6, '发件人信息', 'send_user_info', '【管理员】', '', 1, 2, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (7, '发送内容', 'email_content', '【xxx】你的验证码是：', '', 5, 2, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (8, '关键词', 'web_keyword', '关键词...', '', 5, 3, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (9, '邮箱账号', 'email', 'xxx@email.com', '', 1, 2, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (10, '备案编号', 'record_number', 'Copyright© 2014-2019 | Powered by ***1.1 | 粤ICP备****号', '', 1, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (11, '网站描述', 'web_description', 'web', '', 1, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (12, '下拉选项', 'select', '下拉3', '下拉1,下拉2,下拉3', 4, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (13, '复选框', 'checkbox', 'HTML,CSS', 'AJAX,HTML,JS,CSS', 3, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (14, '文本域', 'textarea', '文本域', '0', 5, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (15, '默认头像', 'default_head_img', '', '', 6, 1, NULL, NULL, NULL);
INSERT INTO `system_config` VALUES (16, '描述', 'seo_description', '', '', 5, 3, NULL, NULL, NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 47 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, 'https://cdn.qitx.net/local/myblog/user_header_image/20230517/577a53d123bc4c4f19db0cb2c6c980a8.jpg', 'admin', '超级管理员', 'dsx.emil@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', '大师兄', 1, 31, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (2, '', 'test2', '李四1', 'ls@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2025-06-16 14:32:14', NULL);
INSERT INTO `user` VALUES (4, '', 'test1', '测试1', 'test1@qq.com', '$2a$10$Ww.IvYhlDpNt6Uq07X5W0OswksocMpae9dmaE2TaHclINQoBUF3Fq', '昵称', 1, 22, 1, '2023-09-07 17:48:39', '2023-09-12 09:52:47', '2023-09-12 09:53:12');
INSERT INTO `user` VALUES (10, '', 'dsx', '大师兄111', 'dsx@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄', 1, 0, 1, '2024-07-22 17:34:36', '2006-01-02 15:04:05', NULL);
INSERT INTO `user` VALUES (11, '', 'admin1', '张三1', 'zs1@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (12, '', 'test3', '李四1', 'ls3@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', NULL);
INSERT INTO `user` VALUES (13, '', 'test4', '测试2', 'test4@qq.com', '$2a$10$Ww.IvYhlDpNt6Uq07X5W0OswksocMpae9dmaE2TaHclINQoBUF3Fq', '昵称', 1, 22, 1, '2023-09-07 17:48:39', '2023-09-12 09:52:47', '2023-09-12 09:53:12');
INSERT INTO `user` VALUES (14, '', 'dsx1', '大师兄1', 'dsx1@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄1', 1, 0, 1, '2024-07-22 17:34:36', '2024-07-22 17:34:36', NULL);
INSERT INTO `user` VALUES (15, '', 'admin2', '张三2', 'zs2@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (16, '', 'test5', '李四5', 'ls5@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', NULL);
INSERT INTO `user` VALUES (17, '', 'test6', '测试6', 'test6@qq.com', '$2a$10$Ww.IvYhlDpNt6Uq07X5W0OswksocMpae9dmaE2TaHclINQoBUF3Fq', '昵称', 1, 22, 1, '2023-09-07 17:48:39', '2023-09-12 09:52:47', '2023-09-12 09:53:12');
INSERT INTO `user` VALUES (18, '', 'dsx2', '大师兄2', 'dsx2@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄1', 1, 0, 1, '2024-07-22 17:34:36', '2024-07-22 17:34:36', NULL);
INSERT INTO `user` VALUES (19, '', 'admin3', '张三3', 'zs3@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (20, '', 'test7', '李四7', 'ls7@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', NULL);
INSERT INTO `user` VALUES (21, '', 'test8', '测试8', 'test8@qq.com', '$2a$10$Ww.IvYhlDpNt6Uq07X5W0OswksocMpae9dmaE2TaHclINQoBUF3Fq', '昵称', 1, 22, 1, '2023-09-07 17:48:39', '2023-09-12 09:52:47', '2023-09-12 09:53:12');
INSERT INTO `user` VALUES (22, '', 'admin4', '张三4', 'zs4@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (23, '', 'test9', '李四9', 'ls9@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', NULL);
INSERT INTO `user` VALUES (24, '', 'test10', '测试10', 'test10@qq.com', '$2a$10$Ww.IvYhlDpNt6Uq07X5W0OswksocMpae9dmaE2TaHclINQoBUF3Fq', '昵称', 1, 22, 1, '2023-09-07 17:48:39', '2023-09-12 09:52:47', '2023-09-12 09:53:12');
INSERT INTO `user` VALUES (25, '', 'dsx3', '大师兄3', 'dsx3@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄1', 1, 0, 1, '2024-07-22 17:34:36', '2024-07-22 17:34:36', NULL);
INSERT INTO `user` VALUES (26, '', 'admin5', '张三5', 'zs5@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (27, '', 'test11', '李四11', 'ls11@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', NULL);
INSERT INTO `user` VALUES (28, '', 'test12', '测试12', 'test12@qq.com', '$2a$10$Ww.IvYhlDpNt6Uq07X5W0OswksocMpae9dmaE2TaHclINQoBUF3Fq', '昵称', 1, 22, 1, '2023-09-07 17:48:39', '2023-09-12 09:52:47', '2023-09-12 09:53:12');
INSERT INTO `user` VALUES (29, '', 'dsx4', '大师兄4', 'dsx4@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄1', 1, 0, 1, '2024-07-22 17:34:36', '2024-07-22 17:34:36', NULL);
INSERT INTO `user` VALUES (30, '', 'admin6', '张三6', 'zs6@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (31, '', 'test13', '李四13', 'ls13@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', NULL);
INSERT INTO `user` VALUES (32, '', 'test14', '测试14', 'test14@qq.com', '$2a$10$Ww.IvYhlDpNt6Uq07X5W0OswksocMpae9dmaE2TaHclINQoBUF3Fq', '昵称', 1, 22, 1, '2023-09-07 17:48:39', '2023-09-12 09:52:47', '2023-09-12 09:53:12');
INSERT INTO `user` VALUES (33, '', 'dsx5', '大师兄5', 'dsx5@qq.com', '$2a$10$Y2FUvgUMpMlJ5h/oooH7OOdInCZgheFQaiVkKu0Wx6YcXhiylAT3a', '大师兄1', 1, 0, 1, '2024-07-22 17:34:36', '2024-07-22 17:34:36', NULL);
INSERT INTO `user` VALUES (34, '', 'admin7', '张三7', 'zs7@qq.com', '$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6', 'dsx', 1, 28, 1, '2023-09-05 17:29:36', '2023-09-12 14:47:48', NULL);
INSERT INTO `user` VALUES (35, '', 'test15', '李四15', 'ls15@qq.com', '$2a$10$kycb2DM8CnubeoWABNPA1O2b0MrQQDqGsEZg8EuqK4G0a63EYDr.2', '昵称22', 1, 1, 1, '2023-09-06 11:38:50', '2023-09-13 09:29:27', '2025-06-06 09:50:57');
INSERT INTO `user` VALUES (36, '', 'test111', '测试16', 'test16@qq.com', '$2a$10$Ww.IvYhlDpNt6Uq07X5W0OswksocMpae9dmaE2TaHclINQoBUF3Fq', '昵称111', 1, 22, 1, '2023-09-07 17:48:39', '2023-09-12 09:52:47', '2023-09-12 09:53:12');
INSERT INTO `user` VALUES (37, '', 'testlisi', '测试李四', 'testlisi@qq.com', '$2a$10$1ED8aAx2IyYXUczXNh/2h.lAWcZdjXZShKBse6/0UBRDmUG1Y8j5G', '测试李四', 1, 31, 1, '2025-05-16 10:43:32', '2025-05-16 16:34:07', '2025-05-16 16:38:16');

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
) ENGINE = InnoDB AUTO_INCREMENT = 22 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户角色表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user_roles
-- ----------------------------
INSERT INTO `user_roles` VALUES (1, 1, 1, 'admin', '2025-05-26 17:53:10', '2025-05-26 17:53:10', NULL);
INSERT INTO `user_roles` VALUES (2, 2, 2, 'test', '2025-05-29 11:02:23', '2025-05-29 11:02:23', '2025-05-29 11:02:23');
INSERT INTO `user_roles` VALUES (3, 2, 2, 'test', '2025-05-29 11:02:23', '2025-05-29 11:02:23', '2025-05-29 11:18:03');
INSERT INTO `user_roles` VALUES (4, 10, 2, 'test', '2025-05-29 11:05:07', '2025-05-29 11:05:07', '2025-05-29 11:05:07');
INSERT INTO `user_roles` VALUES (5, 10, 2, 'test', '2025-05-29 11:05:07', '2025-05-29 11:05:07', '2025-05-29 11:18:12');
INSERT INTO `user_roles` VALUES (6, 2, 2, 'test', '2025-05-29 11:18:03', '2025-05-29 11:18:03', '2025-05-29 11:18:03');
INSERT INTO `user_roles` VALUES (7, 2, 2, 'test', '2025-05-29 11:18:03', '2025-05-29 11:18:03', '2025-05-30 15:59:51');
INSERT INTO `user_roles` VALUES (8, 10, 2, 'test', '2025-05-29 11:18:12', '2025-05-29 11:18:12', '2025-05-29 11:18:12');
INSERT INTO `user_roles` VALUES (9, 10, 2, 'test', '2025-05-29 11:18:12', '2025-05-29 11:18:12', '2025-05-29 14:37:18');
INSERT INTO `user_roles` VALUES (10, 10, 2, 'test', '2025-05-29 14:37:18', '2025-05-29 14:37:18', '2025-05-29 14:37:18');
INSERT INTO `user_roles` VALUES (11, 10, 2, 'test', '2025-05-29 14:37:18', '2025-05-29 14:37:18', NULL);
INSERT INTO `user_roles` VALUES (12, 2, 2, 'test', '2025-05-30 15:59:51', '2025-05-30 15:59:51', '2025-05-30 15:59:51');
INSERT INTO `user_roles` VALUES (13, 2, 2, 'test', '2025-05-30 15:59:51', '2025-05-30 15:59:51', '2025-06-16 10:02:05');
INSERT INTO `user_roles` VALUES (14, 2, 2, 'test', '2025-06-16 10:02:05', '2025-06-16 10:02:05', '2025-06-16 10:02:05');
INSERT INTO `user_roles` VALUES (15, 2, 2, 'test', '2025-06-16 10:02:05', '2025-06-16 10:02:05', '2025-06-16 14:31:39');
INSERT INTO `user_roles` VALUES (16, 2, 2, 'test', '2025-06-16 14:31:39', '2025-06-16 14:31:39', '2025-06-16 14:31:39');
INSERT INTO `user_roles` VALUES (17, 2, 2, 'test', '2025-06-16 14:31:39', '2025-06-16 14:31:39', '2025-06-16 14:32:08');
INSERT INTO `user_roles` VALUES (18, 2, 2, 'test', '2025-06-16 14:32:08', '2025-06-16 14:32:08', '2025-06-16 14:32:08');
INSERT INTO `user_roles` VALUES (19, 2, 2, 'test', '2025-06-16 14:32:08', '2025-06-16 14:32:08', '2025-06-16 14:32:14');
INSERT INTO `user_roles` VALUES (20, 2, 2, 'test', '2025-06-16 14:32:14', '2025-06-16 14:32:14', '2025-06-16 14:32:14');
INSERT INTO `user_roles` VALUES (21, 2, 2, 'test', '2025-06-16 14:32:14', '2025-06-16 14:32:14', NULL);

SET FOREIGN_KEY_CHECKS = 1;
