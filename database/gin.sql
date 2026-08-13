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

 Date: 13/08/2026 14:22:34
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for agent_message
-- ----------------------------
DROP TABLE IF EXISTS `agent_message`;
CREATE TABLE `agent_message`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `session_id` int(10) UNSIGNED NOT NULL COMMENT 'AI会话ID',
  `role` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'user=用户问题 assistant=模型回答 tool=工具执行结果 system=系统提示语',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '消息正文,user存问题原文,assistant存回答或空(有工具调用时content为空),tool存工具返回的JSON',
  `tool_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '工具名称,assistant消息存要调用的工具名(如route:list),tool消息存已执行的工具名',
  `tool_call_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'DeepSeek返回的工具调用ID,assistant和对应的tool消息共用一个ID,用于关联\"调了什么→返回了什么\"',
  `tool_args` json NULL COMMENT '工具参数,仅assistant消息,存模型传过来的参数JSON(如{\"command\":\"route:list\"})',
  `tokens` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '本条消息消耗的token数（API返回的usage）',
  `cost_ms` decimal(10, 2) UNSIGNED NOT NULL DEFAULT 0.00 COMMENT '本轮对话耗时(毫秒),包含API调用+工具执行的总时间',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_session_id`(`session_id`) USING BTREE,
  INDEX `idx_tool_call_id`(`tool_call_id`) USING BTREE,
  INDEX `idx_created_at`(`created_at`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'AI消息详情表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of agent_message
-- ----------------------------
INSERT INTO `agent_message` VALUES (5, 3, 'user', '你能做什么', '', '', NULL, 0, 0.00, '2026-08-13 09:48:40', '2026-08-13 09:48:40', NULL);
INSERT INTO `agent_message` VALUES (6, 3, 'assistant', '你好！我是 Gin-Admin 后台 AI 助手，可以帮你处理以下范围的事情：\n\n## 📊 数据查询\n- **操作日志统计**：查询今日/指定日期的 PV、UV、请求方法分布、状态码统计等\n- **系统配置查询**：查询站点名称、LOGO 等系统配置项\n- **用户搜索**：根据姓名或用户名搜索系统用户\n- **部门查询**：查询部门树结构、指定上级的子部门、部门成员数量\n- **字典查询**：查询性别、状态等枚举字典项\n\n## 🛠️ CLI 命令执行\n- **代码生成**：\n  - 创建控制器（`make:controller`）\n  - 创建模型（`make:model`）\n  - 创建服务（`make:service`）\n  - 创建验证器（`make:request`）\n  - 创建中间件（`make:middleware`）\n  - 创建路由（`make:router`）\n  - 生成 Swagger 文档（`make:docs`）\n- **权限管理**：\n  - 同步用户权限到 Redis（`permission:sync`）\n  - 授予 admin 全部权限（`permission:grant-admin`）\n- **查看信息**：\n  - 查看路由列表（`route:list`）\n  - 查看 Job 列表（`job:list`）\n\n有什么需要帮忙的，尽管告诉我～ 😊', '', '', NULL, 0, 6815.00, '2026-08-13 09:48:47', '2026-08-13 09:48:47', NULL);
INSERT INTO `agent_message` VALUES (7, 3, 'user', '给我重新生成swagger文档', '', '', NULL, 0, 0.00, '2026-08-13 09:57:46', '2026-08-13 09:57:46', NULL);
INSERT INTO `agent_message` VALUES (8, 3, 'assistant', '好的，我来为你重新生成 Swagger 文档。', 'execute_cli', 'call_00_kw6tQZJkVLir15isHcsw2949', '{\"command\": \"make:docs\"}', 0, 1779.00, '2026-08-13 09:57:48', '2026-08-13 09:57:48', NULL);
INSERT INTO `agent_message` VALUES (9, 3, 'tool', '{\"args\":{},\"command\":\"make:docs\",\"message\":\"命令 make:docs 执行完成\",\"output\":\"2026/08/13 09:57:48 Generate swagger docs....\\n2026/08/13 09:57:48 Generate general API Info, search dir:./\\n2026/08/13 09:57:53 Generating request.AgentAsk\\n2026/08/13 09:57:53 Generating base.BaseRequest\\n2026/08/13 09:57:53 Generating base.Context\\n2026/08/13 09:57:53 Generating errcode.SuccessResponse\\n2026/08/13 09:57:53 Generating model.AgentMessage\\n2026/08/13 09:57:53 Generating model.JsonValue\\n2026/08/13 09:57:53 Generating model.DateTime\\n2026/08/13 09:57:53 Generating errcode.ArgsErrorResponse\\n2026/08/13 09:57:53 Generating errcode.SystemErrorResponse\\n2026/08/13 09:57:53 Generating request.PageData\\n2026/08/13 09:57:53 Generating model.Article\\n2026/08/13 09:57:53 Generating model.User\\n2026/08/13 09:57:53 Generating model.UserRoles\\n2026/08/13 09:57:53 Skipping \'model.User\', recursion detected.\\n2026/08/13 09:57:53 Generating model.UserDepartments\\n2026/08/13 09:57:53 Generating model.Department\\n2026/08/13 09:57:53 Generating model.DepartmentLeaders\\n2026/08/13 09:57:53 Skipping \'model.User\', recursion detected.\\n2026/08/13 09:57:53 Generating pkg.TreeNode\\n2026/08/13 09:57:53 Generating request.ArticleCreate\\n2026/08/13 09:57:53 Generating request.ArticleUpdate\\n2026/08/13 09:57:53 Generating model.ConfigCategory\\n2026/08/13 09:57:53 Generating request.ConfigCategoryCreate\\n2026/08/13 09:57:53 Generating request.ConfigCategoryUpdate\\n2026/08/13 09:57:53 Generating service.DashboardCard\\n2026/08/13 09:57:53 Generating service.OperatorLogStatistics\\n2026/08/13 09:57:53 Generating service.MethodStat\\n2026/08/13 09:57:53 Generating service.CostDist\\n2026/08/13 09:57:53 Generating service.HourlyStat\\n2026/08/13 09:57:53 Generating service.StatusCodeStat\\n2026/08/13 09:57:53 Generating service.DashboardSummary\\n2026/08/13 09:57:53 Generating service.SystemResource\\n2026/08/13 09:57:53 Generating service.CPUInfo\\n2026/08/13 09:57:53 Generating service.MemoryInfo\\n2026/08/13 09:57:53 Generating service.DiskInfo\\n2026/08/13 09:57:53 Generating service.NetInfo\\n2026/08/13 09:57:53 Generating request.DepartmentCreate\\n2026/08/13 09:57:53 Generating request.DepartmentUpdate\\n2026/08/13 09:57:53 Generating model.Dict\\n2026/08/13 09:57:53 Generating request.DictCreate\\n2026/08/13 09:57:53 Generating request.DictUpdate\\n2026/08/13 09:57:53 Generating model.ImportRecords\\n2026/08/13 09:57:53 Generating request.UserLogin\\n2026/08/13 09:57:53 Generating v1.LoginResponse\\n2026/08/13 09:57:53 Generating v1.Token\\n2026/08/13 09:57:53 Generating v1.CaptchaResponse\\n2026/08/13 09:57:53 Generating request.CheckCaptcha\\n2026/08/13 09:57:53 Generating model.Menu\\n2026/08/13 09:57:53 Generating model.MenuActions\\n2026/08/13 09:57:53 Generating model.RoleMenus\\n2026/08/13 09:57:53 Generating model.MenuMeta\\n2026/08/13 09:57:53 Generating request.MenuCreate\\n2026/08/13 09:57:53 Generating request.Meta\\n2026/08/13 09:57:53 Generating request.ActionCreate\\n2026/08/13 09:57:53 Generating request.RoleMenu\\n2026/08/13 09:57:53 Generating request.MenuUpdate\\n2026/08/13 09:57:53 Generating model.OperatorLog\\n2026/08/13 09:57:53 Generating request.OperatorLog\\n2026/08/13 09:57:53 Generating request.PageListValidate\\n2026/08/13 09:57:53 Generating model.Permission\\n2026/08/13 09:57:53 Generating model.Roles\\n2026/08/13 09:57:53 Generating model.RolePermissions\\n2026/08/13 09:57:53 Generating request.RoleCreate\\n2026/08/13 09:57:53 Generating request.UserRoles\\n2026/08/13 09:57:53 Generating request.RolePermission\\n2026/08/13 09:57:53 Generating request.RoleUpdate\\n2026/08/13 09:57:53 Generating model.SystemConfig\\n2026/08/13 09:57:53 Generating request.SystemConfigUpdates\\n2026/08/13 09:57:53 Generating request.SystemConfigValueUpdate\\n2026/08/13 09:57:53 Generating request.SystemConfigCreate\\n2026/08/13 09:57:53 Generating request.SystemConfigUpdate\\n2026/08/13 09:57:53 Generating request.UserCreate\\n2026/08/13 09:57:53 Generating request.Dept\\n2026/08/13 09:57:53 Generating request.UserUpdate\\n2026/08/13 09:57:53 Generating request.UserImport\\n2026/08/13 09:57:53 Generating request.UserImportItem\\n2026/08/13 09:57:53 Generating request.UserBatchDelete\\n2026/08/13 09:57:53 Generating request.UserPassword\\n2026/08/13 09:57:53 create docs.go at docs/docs.go\\n2026/08/13 09:57:53 create swagger.json at docs/swagger.json\\n2026/08/13 09:57:53 create swagger.yaml at docs/swagger.yaml\\n\\n\"}', 'execute_cli', 'call_00_kw6tQZJkVLir15isHcsw2949', NULL, 0, 0.00, '2026-08-13 09:57:54', '2026-08-13 09:57:54', NULL);
INSERT INTO `agent_message` VALUES (10, 3, 'user', '给admin用户赋予全部权限', '', '', NULL, 0, 0.00, '2026-08-13 09:58:58', '2026-08-13 09:58:58', NULL);
INSERT INTO `agent_message` VALUES (11, 3, 'assistant', '好的，我来为 admin 用户授予全部权限。', 'execute_cli', 'call_00_ZWGevBpdELwKr0n8CXsG8751', '{\"command\": \"permission:grant-admin\"}', 0, 1680.00, '2026-08-13 09:58:59', '2026-08-13 09:58:59', NULL);
INSERT INTO `agent_message` VALUES (12, 3, 'tool', '{\"args\":{},\"command\":\"permission:grant-admin\",\"message\":\"命令 permission:grant-admin 执行完成\",\"output\":\"\"}', 'execute_cli', 'call_00_ZWGevBpdELwKr0n8CXsG8751', NULL, 0, 0.00, '2026-08-13 09:59:00', '2026-08-13 09:59:00', NULL);
INSERT INTO `agent_message` VALUES (13, 3, 'assistant', '已成功为 admin 用户授予全部权限 ✅\n\n如需确认权限是否同步到 Redis，可以再执行 `permission:sync` 命令。', '', '', NULL, 0, 1247.00, '2026-08-13 09:59:01', '2026-08-13 09:59:01', NULL);
INSERT INTO `agent_message` VALUES (14, 3, 'user', '可以再执行 `permission:sync` 命令', '', '', NULL, 0, 0.00, '2026-08-13 09:59:26', '2026-08-13 09:59:26', NULL);
INSERT INTO `agent_message` VALUES (15, 3, 'assistant', '好的，我来执行权限同步命令。', 'execute_cli', 'call_00_4t6OhJzoEru8Hmi9g0ed0449', '{\"command\": \"permission:sync\"}', 0, 2182.00, '2026-08-13 09:59:28', '2026-08-13 09:59:28', NULL);
INSERT INTO `agent_message` VALUES (16, 3, 'tool', '{\"args\":{},\"command\":\"permission:sync\",\"message\":\"命令 permission:sync 执行完成\",\"output\":\"\"}', 'execute_cli', 'call_00_4t6OhJzoEru8Hmi9g0ed0449', NULL, 0, 0.00, '2026-08-13 09:59:28', '2026-08-13 09:59:28', NULL);
INSERT INTO `agent_message` VALUES (17, 3, 'assistant', '权限同步命令已执行完成 ✅\n\n现在 admin 用户的权限已成功授予并同步到 Redis 了。', '', '', NULL, 0, 1879.00, '2026-08-13 09:59:30', '2026-08-13 09:59:30', NULL);
INSERT INTO `agent_message` VALUES (20, 4, 'user', '当前系统配置有哪些配置', '', '', NULL, 0, 0.00, '2026-08-13 10:45:58', '2026-08-13 10:45:58', NULL);
INSERT INTO `agent_message` VALUES (21, 4, 'assistant', '', 'query_system_config', 'call_00_CexaC0YCXk4qjdKy545Z8425', '{}', 0, 1379.00, '2026-08-13 10:46:00', '2026-08-13 10:46:00', NULL);
INSERT INTO `agent_message` VALUES (22, 4, 'tool', '{\"error\":\"Error 1054 (42S22): Unknown column \'sort\' in \'order clause\'\"}', 'query_system_config', 'call_00_CexaC0YCXk4qjdKy545Z8425', NULL, 0, 0.00, '2026-08-13 10:46:00', '2026-08-13 10:46:00', NULL);
INSERT INTO `agent_message` VALUES (23, 4, 'assistant', '抱歉，查询系统配置时出现了数据库错误：\n\n**Error 1054 (42S22): Unknown column \'sort\' in \'order clause\'**\n\n这表示系统配置表在排序时引用了一个不存在的 `sort` 字段，属于后端程序或数据库结构问题，我无法获取到配置数据。\n\n建议联系开发人员检查系统配置表（如 `sys_config`）是否存在 `sort` 字段，或修正相关查询语句。', '', '', NULL, 0, 2496.00, '2026-08-13 10:46:02', '2026-08-13 10:46:02', NULL);
INSERT INTO `agent_message` VALUES (24, 4, 'user', '重新查询', '', '', NULL, 0, 0.00, '2026-08-13 10:51:39', '2026-08-13 10:51:39', NULL);
INSERT INTO `agent_message` VALUES (25, 4, 'assistant', '', 'query_system_config', 'call_00_4oHKJxxvbEGZh40iCNIM6608', '{}', 0, 1315.00, '2026-08-13 10:51:40', '2026-08-13 10:51:40', NULL);
INSERT INTO `agent_message` VALUES (26, 4, 'tool', '{\"list\":[{\"key\":\"web_domain\",\"name\":\"网站域名\",\"value\":\"www.a.com\",\"type\":1},{\"key\":\"is_open_site\",\"name\":\"关闭站点\",\"value\":\"开启\",\"type\":2},{\"key\":\"site_logo\",\"name\":\"网站Logo\",\"value\":\"\",\"type\":6},{\"key\":\"email_port\",\"name\":\"邮件端口\",\"value\":\"465\",\"type\":1},{\"key\":\"email_title\",\"name\":\"邮件标题\",\"value\":\"【xxx】验证码\",\"type\":1},{\"key\":\"send_user_info\",\"name\":\"发件人信息\",\"value\":\"【管理员】\",\"type\":1},{\"key\":\"email_content\",\"name\":\"发送内容\",\"value\":\"【xxx】你的验证码是：\",\"type\":5},{\"key\":\"web_keyword\",\"name\":\"关键词\",\"value\":\"关键词...\",\"type\":5},{\"key\":\"email\",\"name\":\"邮箱账号\",\"value\":\"xxx@email.com\",\"type\":1},{\"key\":\"record_number\",\"name\":\"备案编号\",\"value\":\"Copyright© 2014-2019 | Powered by ***1.1 | 粤ICP备****号\",\"type\":1},{\"key\":\"web_description\",\"name\":\"网站描述\",\"value\":\"web\",\"type\":1},{\"key\":\"select\",\"name\":\"下拉选项\",\"value\":\"下拉3\",\"type\":4},{\"key\":\"checkbox\",\"name\":\"复选框\",\"value\":\"HTML,CSS\",\"type\":3},{\"key\":\"textarea\",\"name\":\"文本域\",\"value\":\"文本域\",\"type\":5},{\"key\":\"default_head_img\",\"name\":\"默认头像\",\"value\":\"\",\"type\":6},{\"key\":\"seo_description\",\"name\":\"描述\",\"value\":\"11\",\"type\":5}],\"total\":16}', 'query_system_config', 'call_00_4oHKJxxvbEGZh40iCNIM6608', NULL, 0, 0.00, '2026-08-13 10:51:40', '2026-08-13 10:51:40', NULL);
INSERT INTO `agent_message` VALUES (27, 4, 'assistant', '当前系统配置共 **16 项**，详情如下：\n\n| 配置键 | 名称 | 值 | 类型 |\n|---|---|---|---|\n| web_domain | 网站域名 | www.a.com | 文本 |\n| is_open_site | 关闭站点 | 开启 | 开关 |\n| site_logo | 网站Logo | （空） | 图片 |\n| email_port | 邮件端口 | 465 | 文本 |\n| email_title | 邮件标题 | 【xxx】验证码 | 文本 |\n| send_user_info | 发件人信息 | 【管理员】 | 文本 |\n| email_content | 发送内容 | 【xxx】你的验证码是： | 文本域 |\n| web_keyword | 关键词 | 关键词... | 文本域 |\n| email | 邮箱账号 | xxx@email.com | 文本 |\n| record_number | 备案编号 | Copyright© 2014-2019 | Powered by ***1.1 | 粤ICP备****号 | 文本 |\n| web_description | 网站描述 | web | 文本 |\n| select | 下拉选项 | 下拉3 | 下拉框 |\n| checkbox | 复选框 | HTML,CSS | 复选框 |\n| textarea | 文本域 | 文本域 | 文本域 |\n| default_head_img | 默认头像 | （空） | 图片 |\n| seo_description | 描述 | 11 | 文本域 |\n\n**主要配置概览：**\n- **站点信息**：域名 `www.a.com`，站点处于开启状态，Logo 未设置\n- **邮件配置**：邮箱账号 `xxx@email.com`，端口 465，发件人为【管理员】\n- **SEO/备案**：关键词、描述、备案编号均已配置\n\n如需查询某项具体配置，可告诉我配置键名。', '', '', NULL, 0, 4759.00, '2026-08-13 10:51:45', '2026-08-13 10:51:45', NULL);
INSERT INTO `agent_message` VALUES (28, 4, 'user', '帮我重新生成swagger文档', '', '', NULL, 0, 0.00, '2026-08-13 10:53:08', '2026-08-13 10:53:08', NULL);
INSERT INTO `agent_message` VALUES (29, 4, 'assistant', '', 'execute_cli', 'call_00_ESvviGKFDutFyNiVnr7N0190', '{\"command\": \"make:docs\"}', 0, 1681.00, '2026-08-13 10:53:10', '2026-08-13 10:53:10', NULL);
INSERT INTO `agent_message` VALUES (30, 4, 'tool', '{\"args\":{},\"command\":\"make:docs\",\"message\":\"命令 make:docs 执行完成\",\"output\":\"2026/08/13 10:53:10 Generate swagger docs....\\n2026/08/13 10:53:10 Generate general API Info, search dir:./\\n2026/08/13 10:53:12 Generating request.AgentAsk\\n2026/08/13 10:53:12 Generating base.BaseRequest\\n2026/08/13 10:53:12 Generating base.Context\\n2026/08/13 10:53:12 Generating errcode.SuccessResponse\\n2026/08/13 10:53:12 Generating model.AgentSession\\n2026/08/13 10:53:12 Generating model.DateTime\\n2026/08/13 10:53:12 Generating errcode.SystemErrorResponse\\n2026/08/13 10:53:12 Generating model.AgentMessage\\n2026/08/13 10:53:12 Generating model.JsonValue\\n2026/08/13 10:53:12 Generating errcode.ArgsErrorResponse\\n2026/08/13 10:53:12 Generating request.PageData\\n2026/08/13 10:53:12 Generating model.Article\\n2026/08/13 10:53:12 Generating model.User\\n2026/08/13 10:53:12 Generating model.UserRoles\\n2026/08/13 10:53:12 Skipping \'model.User\', recursion detected.\\n2026/08/13 10:53:12 Generating model.UserDepartments\\n2026/08/13 10:53:12 Generating model.Department\\n2026/08/13 10:53:12 Generating model.DepartmentLeaders\\n2026/08/13 10:53:12 Skipping \'model.User\', recursion detected.\\n2026/08/13 10:53:12 Generating pkg.TreeNode\\n2026/08/13 10:53:12 Generating request.ArticleCreate\\n2026/08/13 10:53:12 Generating request.ArticleUpdate\\n2026/08/13 10:53:12 Generating model.ConfigCategory\\n2026/08/13 10:53:12 Generating request.ConfigCategoryCreate\\n2026/08/13 10:53:12 Generating request.ConfigCategoryUpdate\\n2026/08/13 10:53:12 Generating service.DashboardCard\\n2026/08/13 10:53:12 Generating service.OperatorLogStatistics\\n2026/08/13 10:53:12 Generating service.MethodStat\\n2026/08/13 10:53:12 Generating service.CostDist\\n2026/08/13 10:53:12 Generating service.HourlyStat\\n2026/08/13 10:53:12 Generating service.StatusCodeStat\\n2026/08/13 10:53:12 Generating service.DashboardSummary\\n2026/08/13 10:53:12 Generating service.SystemResource\\n2026/08/13 10:53:12 Generating service.CPUInfo\\n2026/08/13 10:53:12 Generating service.MemoryInfo\\n2026/08/13 10:53:12 Generating service.DiskInfo\\n2026/08/13 10:53:12 Generating service.NetInfo\\n2026/08/13 10:53:12 Generating request.DepartmentCreate\\n2026/08/13 10:53:12 Generating request.DepartmentUpdate\\n2026/08/13 10:53:12 Generating model.Dict\\n2026/08/13 10:53:12 Generating request.DictCreate\\n2026/08/13 10:53:12 Generating request.DictUpdate\\n2026/08/13 10:53:12 Generating model.ImportRecords\\n2026/08/13 10:53:12 Generating request.UserLogin\\n2026/08/13 10:53:12 Generating v1.LoginResponse\\n2026/08/13 10:53:12 Generating v1.Token\\n2026/08/13 10:53:12 Generating v1.CaptchaResponse\\n2026/08/13 10:53:12 Generating request.CheckCaptcha\\n2026/08/13 10:53:12 Generating model.Menu\\n2026/08/13 10:53:12 Generating model.MenuActions\\n2026/08/13 10:53:12 Generating model.RoleMenus\\n2026/08/13 10:53:12 Generating model.MenuMeta\\n2026/08/13 10:53:12 Generating request.MenuCreate\\n2026/08/13 10:53:12 Generating request.Meta\\n2026/08/13 10:53:12 Generating request.ActionCreate\\n2026/08/13 10:53:12 Generating request.RoleMenu\\n2026/08/13 10:53:12 Generating request.MenuUpdate\\n2026/08/13 10:53:12 Generating model.OperatorLog\\n2026/08/13 10:53:12 Generating request.OperatorLog\\n2026/08/13 10:53:12 Generating request.PageListValidate\\n2026/08/13 10:53:12 Generating model.Permission\\n2026/08/13 10:53:12 Generating model.Roles\\n2026/08/13 10:53:12 Generating model.RolePermissions\\n2026/08/13 10:53:12 Generating request.RoleCreate\\n2026/08/13 10:53:12 Generating request.UserRoles\\n2026/08/13 10:53:12 Generating request.RolePermission\\n2026/08/13 10:53:12 Generating request.RoleUpdate\\n2026/08/13 10:53:12 Generating model.SystemConfig\\n2026/08/13 10:53:12 Generating request.SystemConfigUpdates\\n2026/08/13 10:53:12 Generating request.SystemConfigValueUpdate\\n2026/08/13 10:53:12 Generating request.SystemConfigCreate\\n2026/08/13 10:53:12 Generating request.SystemConfigUpdate\\n2026/08/13 10:53:12 Generating request.UserCreate\\n2026/08/13 10:53:12 Generating request.Dept\\n2026/08/13 10:53:12 Generating request.UserUpdate\\n2026/08/13 10:53:12 Generating request.UserImport\\n2026/08/13 10:53:12 Generating request.UserImportItem\\n2026/08/13 10:53:12 Generating request.UserBatchDelete\\n2026/08/13 10:53:12 Generating request.UserPassword\\n2026/08/13 10:53:12 create docs.go at docs/docs.go\\n2026/08/13 10:53:12 create swagger.json at docs/swagger.json\\n2026/08/13 10:53:12 create swagger.yaml at docs/swagger.yaml\\n\\n\"}', 'execute_cli', 'call_00_ESvviGKFDutFyNiVnr7N0190', NULL, 0, 0.00, '2026-08-13 10:53:12', '2026-08-13 10:53:12', NULL);
INSERT INTO `agent_message` VALUES (31, 4, 'assistant', 'Swagger 文档已成功重新生成 ✅\n\n**生成结果：**\n- 📄 `docs/docs.go`\n- 📄 `docs/swagger.json`\n- 📄 `docs/swagger.yaml`\n\n**生成时间：** 2026/08/13 10:53:12\n\n本次共扫描并生成了所有模型、请求体、响应体等结构定义，涵盖用户、部门、菜单、角色、权限、系统配置、操作日志、字典、文章等模块，以及登录、验证码等接口相关结构。生成过程正常完成，无报错。', '', '', NULL, 0, 1966.00, '2026-08-13 10:53:14', '2026-08-13 10:53:14', NULL);
INSERT INTO `agent_message` VALUES (34, 6, 'user', '你能查询天气吗', '', '', NULL, 0, 0.00, '2026-08-13 11:07:50', '2026-08-13 11:07:50', NULL);
INSERT INTO `agent_message` VALUES (35, 6, 'assistant', '抱歉,我当前只能处理以下范围:查询操作日志统计、查询系统配置项、根据姓名或用户名搜索系统用户、执行CLI命令、查询部门列表或指定上级的子部门、根据字典类型查询字典项', '', '', NULL, 0, 2064.00, '2026-08-13 11:07:53', '2026-08-13 11:07:53', NULL);
INSERT INTO `agent_message` VALUES (39, 8, 'user', '给test角色新增ai助手权限', '', '', NULL, 0, 0.00, '2026-08-13 14:06:43', '2026-08-13 14:06:43', NULL);
INSERT INTO `agent_message` VALUES (40, 8, 'assistant', '', 'execute_cli', 'call_00_jlWKIiR7BCMuBiJ5N2YC3616', '{\"command\": \"route:list\"}', 0, 3208.00, '2026-08-13 14:06:46', '2026-08-13 14:06:46', NULL);
INSERT INTO `agent_message` VALUES (41, 8, 'tool', '{\"args\":{},\"command\":\"route:list\",\"message\":\"命令 route:list 执行完成\",\"output\":\"Method   Path                                Handler                                 \\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/agent/ask                  \\u001b[0m \\u001b[37mgin/app/controller/v1.(*AgentController).Ask\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/agent/history              \\u001b[0m \\u001b[37mgin/app/controller/v1.(*AgentController).History\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/agent/sessions             \\u001b[0m \\u001b[37mgin/app/controller/v1.(*AgentController).Sessions\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/article                    \\u001b[0m \\u001b[37mgin/app/controller/v1.(*ArticleController).Create\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/article                    \\u001b[0m \\u001b[37mgin/app/controller/v1.(*ArticleController).List\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/article/:id                \\u001b[0m \\u001b[37mgin/app/controller/v1.(*ArticleController).Detail\\u001b[0m\\n\\u001b[31mDELETE  \\u001b[0m \\u001b[36m/api/v1/article/:id                \\u001b[0m \\u001b[37mgin/app/controller/v1.(*ArticleController).Delete\\u001b[0m\\n\\u001b[34mPUT     \\u001b[0m \\u001b[36m/api/v1/article/:id                \\u001b[0m \\u001b[37mgin/app/controller/v1.(*ArticleController).Update\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/captcha                    \\u001b[0m \\u001b[37mgin/app/controller/v1.(*LoginController).CheckCaptcha\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/captcha                    \\u001b[0m \\u001b[37mgin/app/controller/v1.(*LoginController).GetCaptcha\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/config-category            \\u001b[0m \\u001b[37mgin/app/controller/v1.(*ConfigCategoryController).Create\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/config-category            \\u001b[0m \\u001b[37mgin/app/controller/v1.(*ConfigCategoryController).List\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/config-category/:id        \\u001b[0m \\u001b[37mgin/app/controller/v1.(*ConfigCategoryController).Detail\\u001b[0m\\n\\u001b[31mDELETE  \\u001b[0m \\u001b[36m/api/v1/config-category/:id        \\u001b[0m \\u001b[37mgin/app/controller/v1.(*ConfigCategoryController).Delete\\u001b[0m\\n\\u001b[34mPUT     \\u001b[0m \\u001b[36m/api/v1/config-category/:id        \\u001b[0m \\u001b[37mgin/app/controller/v1.(*ConfigCategoryController).Update\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/dashboard/cards            \\u001b[0m \\u001b[37mgin/app/controller/v1.(*DashboardController).Cards\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/dashboard/statistics       \\u001b[0m \\u001b[37mgin/app/controller/v1.(*DashboardController).Statistics\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/dashboard/system-resource  \\u001b[0m \\u001b[37mgin/app/controller/v1.(*DashboardController).SystemResource\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/department                 \\u001b[0m \\u001b[37mgin/app/controller/v1.(*DepartmentController).Create\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/department                 \\u001b[0m \\u001b[37mgin/app/controller/v1.(*DepartmentController).List\\u001b[0m\\n\\u001b[34mPUT     \\u001b[0m \\u001b[36m/api/v1/department/:id             \\u001b[0m \\u001b[37mgin/app/controller/v1.(*DepartmentController).Update\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/department/:id             \\u001b[0m \\u001b[37mgin/app/controller/v1.(*DepartmentController).Detail\\u001b[0m\\n\\u001b[31mDELETE  \\u001b[0m \\u001b[36m/api/v1/department/:id             \\u001b[0m \\u001b[37mgin/app/controller/v1.(*DepartmentController).Delete\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/dict                       \\u001b[0m \\u001b[37mgin/app/controller/v1.(*DictController).List\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/dict                       \\u001b[0m \\u001b[37mgin/app/controller/v1.(*DictController).Create\\u001b[0m\\n\\u001b[34mPUT     \\u001b[0m \\u001b[36m/api/v1/dict/:id                   \\u001b[0m \\u001b[37mgin/app/controller/v1.(*DictController).Update\\u001b[0m\\n\\u001b[31mDELETE  \\u001b[0m \\u001b[36m/api/v1/dict/:id                   \\u001b[0m \\u001b[37mgin/app/controller/v1.(*DictController).Delete\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/dict/:id                   \\u001b[0m \\u001b[37mgin/app/controller/v1.(*DictController).Detail\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/import-records             \\u001b[0m \\u001b[37mgin/app/controller/v1.(*ImportRecordsController).List\\u001b[0m\\n\\u001b[31mDELETE  \\u001b[0m \\u001b[36m/api/v1/import-records/:id         \\u001b[0m \\u001b[37mgin/app/controller/v1.(*ImportRecordsController).Delete\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/login                      \\u001b[0m \\u001b[37mgin/app/controller/v1.(*LoginController).Login\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/menu                       \\u001b[0m \\u001b[37mgin/app/controller/v1.(*MenuController).List\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/menu                       \\u001b[0m \\u001b[37mgin/app/controller/v1.(*MenuController).Create\\u001b[0m\\n\\u001b[31mDELETE  \\u001b[0m \\u001b[36m/api/v1/menu/:id                   \\u001b[0m \\u001b[37mgin/app/controller/v1.(*MenuController).Delete\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/menu/:id                   \\u001b[0m \\u001b[37mgin/app/controller/v1.(*MenuController).Detail\\u001b[0m\\n\\u001b[34mPUT     \\u001b[0m \\u001b[36m/api/v1/menu/:id                   \\u001b[0m \\u001b[37mgin/app/controller/v1.(*MenuController).Update\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/operator-log               \\u001b[0m \\u001b[37mgin/app/controller/v1.(*OperatorLogController).List\\u001b[0m\\n\\u001b[31mDELETE  \\u001b[0m \\u001b[36m/api/v1/operator-log/:id           \\u001b[0m \\u001b[37mgin/app/controller/v1.(*OperatorLogController).Delete\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/operator-log/:id           \\u001b[0m \\u001b[37mgin/app/controller/v1.(*OperatorLogController).Detail\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/operator-log/batch-delete  \\u001b[0m \\u001b[37mgin/app/controller/v1.(*OperatorLogController).BatchDelete\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/permission                 \\u001b[0m \\u001b[37mgin/app/controller/v1.(*PermissionController).List\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/refresh-token              \\u001b[0m \\u001b[37mgin/app/controller/v1.(*LoginController).RefreshToken\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/role                       \\u001b[0m \\u001b[37mgin/app/controller/v1.(*RoleController).List\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/role                       \\u001b[0m \\u001b[37mgin/app/controller/v1.(*RoleController).Create\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/role/:id                   \\u001b[0m \\u001b[37mgin/app/controller/v1.(*RoleController).Detail\\u001b[0m\\n\\u001b[34mPUT     \\u001b[0m \\u001b[36m/api/v1/role/:id                   \\u001b[0m \\u001b[37mgin/app/controller/v1.(*RoleController).Update\\u001b[0m\\n\\u001b[31mDELETE  \\u001b[0m \\u001b[36m/api/v1/role/:id                   \\u001b[0m \\u001b[37mgin/app/controller/v1.(*RoleController).Delete\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/role/:id/menu              \\u001b[0m \\u001b[37mgin/app/controller/v1.(*MenuController).RoleMenu\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/system-config              \\u001b[0m \\u001b[37mgin/app/controller/v1.(*SystemConfigController).Create\\u001b[0m\\n\\u001b[34mPUT     \\u001b[0m \\u001b[36m/api/v1/system-config              \\u001b[0m \\u001b[37mgin/app/controller/v1.(*SystemConfigController).UpdateConfig\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/system-config              \\u001b[0m \\u001b[37mgin/app/controller/v1.(*SystemConfigController).List\\u001b[0m\\n\\u001b[31mDELETE  \\u001b[0m \\u001b[36m/api/v1/system-config/:id          \\u001b[0m \\u001b[37mgin/app/controller/v1.(*SystemConfigController).Delete\\u001b[0m\\n\\u001b[34mPUT     \\u001b[0m \\u001b[36m/api/v1/system-config/:id          \\u001b[0m \\u001b[37mgin/app/controller/v1.(*SystemConfigController).Update\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/system-config/:id          \\u001b[0m \\u001b[37mgin/app/controller/v1.(*SystemConfigController).Detail\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/test                       \\u001b[0m \\u001b[37mgin/app/controller/v1.(*LoginController).Test\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/user                       \\u001b[0m \\u001b[37mgin/app/controller/v1.(*UserController).List\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/user                       \\u001b[0m \\u001b[37mgin/app/controller/v1.(*UserController).Create\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/api/v1/user/:id                   \\u001b[0m \\u001b[37mgin/app/controller/v1.(*UserController).Detail\\u001b[0m\\n\\u001b[34mPUT     \\u001b[0m \\u001b[36m/api/v1/user/:id                   \\u001b[0m \\u001b[37mgin/app/controller/v1.(*UserController).Update\\u001b[0m\\n\\u001b[31mDELETE  \\u001b[0m \\u001b[36m/api/v1/user/:id                   \\u001b[0m \\u001b[37mgin/app/controller/v1.(*UserController).Delete\\u001b[0m\\n\\u001b[34mPUT     \\u001b[0m \\u001b[36m/api/v1/user/:id/password          \\u001b[0m \\u001b[37mgin/app/controller/v1.(*UserController).Password\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/user/batch-delete          \\u001b[0m \\u001b[37mgin/app/controller/v1.(*UserController).BatchDelete\\u001b[0m\\n\\u001b[33mPOST    \\u001b[0m \\u001b[36m/api/v1/user/import                \\u001b[0m \\u001b[37mgin/app/controller/v1.(*UserController).Import\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/ping                              \\u001b[0m \\u001b[37mgin/router.LoadRouters                  \\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/public/*filepath                  \\u001b[0m \\u001b[37mgithub.com/gin-gonic/gin.(*RouterGroup).createStaticHandler\\u001b[0m\\n\\u001b[37mHEAD    \\u001b[0m \\u001b[36m/public/*filepath                  \\u001b[0m \\u001b[37mgithub.com/gin-gonic/gin.(*RouterGroup).createStaticHandler\\u001b[0m\\n\\u001b[32mGET     \\u001b[0m \\u001b[36m/swagger/*any                      \\u001b[0m \\u001b[37mgithub.com/swaggo/gin-swagger.CustomWrapHandler\\u001b[0m\\n\"}', 'execute_cli', 'call_00_jlWKIiR7BCMuBiJ5N2YC3616', NULL, 0, 0.00, '2026-08-13 14:06:46', '2026-08-13 14:06:46', NULL);
INSERT INTO `agent_message` VALUES (42, 8, 'assistant', '我已找到与\"AI助手\"相关的路由（AgentController 相关）。现在为 test 角色授予 AI 助手相关权限。', 'set_role_permissions', 'call_00_15ooOP4khTaevZYBSken9981', '{\"role\": \"test\", \"permissions\": [\"POST:/api/v1/agent/ask\", \"GET:/api/v1/agent/history\", \"GET:/api/v1/agent/sessions\"]}', 0, 3684.00, '2026-08-13 14:06:50', '2026-08-13 14:06:50', NULL);
INSERT INTO `agent_message` VALUES (43, 8, 'tool', '{\"granted\":3,\"message\":\"权限设置完成\",\"permissions\":[\"POST:/api/v1/agent/ask\",\"GET:/api/v1/agent/history\",\"GET:/api/v1/agent/sessions\"],\"role\":\"test\"}', 'set_role_permissions', 'call_00_15ooOP4khTaevZYBSken9981', NULL, 0, 0.00, '2026-08-13 14:06:50', '2026-08-13 14:06:50', NULL);
INSERT INTO `agent_message` VALUES (44, 8, 'assistant', '已为 **test** 角色成功授予 AI 助手相关权限，共 3 项：\n\n| 方法 | 路径 | 说明 |\n|------|------|------|\n| POST | /api/v1/agent/ask | AI 助手问答 |\n| GET | /api/v1/agent/history | AI 助手历史记录 |\n| GET | /api/v1/agent/sessions | AI 助手会话列表 |', '', '', NULL, 0, 1794.00, '2026-08-13 14:06:52', '2026-08-13 14:06:52', NULL);

-- ----------------------------
-- Table structure for agent_session
-- ----------------------------
DROP TABLE IF EXISTS `agent_session`;
CREATE TABLE `agent_session`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` int(10) UNSIGNED NOT NULL COMMENT '用户ID',
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '会话标题',
  `provider` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模型提供商,如deepseek、openai、moonshot',
  `model` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模型名,如deepseek-v4-pro',
  `message_count` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '该会话总消息条数(含user、assistant、tool)',
  `total_tokens` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '该会话累计消耗的token总数',
  `trace_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '请求追踪ID',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id`) USING BTREE,
  INDEX `idx_created_at`(`created_at`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'AI会话主表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of agent_session
-- ----------------------------
INSERT INTO `agent_session` VALUES (3, 1, '你能做什么', 'deepseek', 'deepseek-v4-pro', 15, 0, '2506d2e9-198e-4fef-9b6a-073111c05d7b', '2026-08-13 09:48:40', '2026-08-13 10:40:07', NULL);
INSERT INTO `agent_session` VALUES (4, 1, '当前系统配置有哪些配置', 'deepseek', 'deepseek-v4-pro', 12, 0, 'b084afb5-5c5d-4f9a-a818-f068d2337124', '2026-08-13 10:45:58', '2026-08-13 10:53:14', NULL);
INSERT INTO `agent_session` VALUES (6, 1, '你能查询天气吗', 'deepseek', 'deepseek-v4-pro', 2, 0, 'd5ef56bd-692a-44ed-92de-5a64c0690c47', '2026-08-13 11:07:50', '2026-08-13 11:07:53', NULL);
INSERT INTO `agent_session` VALUES (8, 1, '给test角色新增ai助手权限', 'deepseek', 'deepseek-v4-pro', 6, 0, '3345a3af-6356-4e41-996a-89bdbb83d98a', '2026-08-13 14:06:43', '2026-08-13 14:06:52', NULL);

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
INSERT INTO `menu` VALUES (72, 0, 1, 'agent', 1, 0, '2026-08-13 10:38:41', '2026-08-13 10:38:41', NULL);

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
INSERT INTO `menu_meta` VALUES (29, 72, 'AI助手', '', 'fa fa-window-restore', '/assistant', '', 'agent/index', 2, 1, 2, '', 2, '2026-08-13 10:38:41', '2026-08-13 10:38:41', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 61 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '权限表' ROW_FORMAT = DYNAMIC;

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
INSERT INTO `permission` VALUES (60, 'POST:/api/v1/agent/ask', 'POST', '/api/v1/agent/ask', '2026-08-12 15:04:41', '2026-08-12 15:04:41', NULL);
INSERT INTO `permission` VALUES (62, 'GET:/api/v1/agent/history', 'GET', '/api/v1/agent/history', '2026-08-13 10:10:15', '2026-08-13 10:10:15', NULL);
INSERT INTO `permission` VALUES (63, 'GET:/api/v1/agent/sessions', 'GET', '/api/v1/agent/sessions', '2026-08-13 10:23:51', '2026-08-13 10:23:51', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 768 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色菜单表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role_menus
-- ----------------------------
INSERT INTO `role_menus` VALUES (303, 2, 24, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `role_menus` VALUES (304, 2, 27, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `role_menus` VALUES (305, 2, 2, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `role_menus` VALUES (306, 2, 3, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `role_menus` VALUES (666, 2, 1, 'test', '2026-08-04 17:15:38', '2026-08-04 17:15:38', NULL);
INSERT INTO `role_menus` VALUES (818, 1, 1, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (819, 1, 2, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (820, 1, 67, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (821, 1, 68, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (822, 1, 69, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (823, 1, 70, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (824, 1, 71, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (825, 1, 3, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (826, 1, 58, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (827, 1, 24, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (828, 1, 25, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (829, 1, 27, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (830, 1, 4, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (831, 1, 54, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (832, 1, 56, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (833, 1, 61, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (834, 1, 62, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (835, 1, 32, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (836, 1, 33, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (837, 1, 34, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (838, 1, 60, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (839, 1, 35, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (840, 1, 5, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (841, 1, 37, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (842, 1, 38, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (843, 1, 39, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (844, 1, 6, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (845, 1, 40, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (846, 1, 41, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (847, 1, 59, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (848, 1, 42, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (849, 1, 20, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (850, 1, 23, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (851, 1, 51, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (852, 1, 52, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (853, 1, 53, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (854, 1, 21, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (855, 1, 48, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (856, 1, 49, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (857, 1, 50, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (858, 1, 22, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (859, 1, 10, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (860, 1, 43, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (861, 1, 44, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (862, 1, 45, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (863, 1, 63, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (864, 1, 65, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (865, 1, 64, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (866, 1, 66, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `role_menus` VALUES (867, 1, 72, 'admin', '2026-08-13 10:38:41', '2026-08-13 10:38:41', NULL);

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
) ENGINE = InnoDB AUTO_INCREMENT = 447 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色权限表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role_permissions
-- ----------------------------
INSERT INTO `role_permissions` VALUES (562, 1, 60, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (563, 1, 62, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (564, 1, 1, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (565, 1, 3, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (566, 1, 38, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (567, 1, 19, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (568, 1, 8, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (569, 1, 6, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (570, 1, 15, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (571, 1, 39, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (572, 1, 20, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (573, 1, 17, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (574, 1, 54, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (575, 1, 52, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (576, 1, 53, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (577, 1, 59, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (578, 1, 55, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (579, 1, 58, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (580, 1, 57, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (581, 1, 56, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (582, 1, 12, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (583, 1, 16, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (584, 1, 18, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (585, 1, 21, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (586, 1, 10, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (587, 1, 28, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (588, 1, 40, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (589, 1, 2, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (590, 1, 23, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (591, 1, 9, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (592, 1, 33, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (593, 1, 25, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (594, 1, 46, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (595, 1, 47, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (596, 1, 44, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (597, 1, 45, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (598, 1, 43, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (599, 1, 42, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (600, 1, 24, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (601, 1, 31, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (602, 1, 32, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (603, 1, 34, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (604, 1, 35, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (605, 1, 26, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (606, 1, 7, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (607, 1, 4, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (608, 1, 11, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (609, 1, 27, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (610, 1, 37, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (611, 1, 13, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (612, 1, 29, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (613, 1, 41, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (614, 1, 14, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (615, 1, 30, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (616, 1, 5, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (617, 1, 36, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (618, 1, 22, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (619, 1, 63, '2026-08-13 10:24:45', NULL);
INSERT INTO `role_permissions` VALUES (623, 2, 62, '2026-08-13 14:06:50', NULL);
INSERT INTO `role_permissions` VALUES (624, 2, 63, '2026-08-13 14:06:50', NULL);
INSERT INTO `role_permissions` VALUES (625, 2, 60, '2026-08-13 14:06:50', NULL);

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
INSERT INTO `roles` VALUES (1, 'admin', '超级管理员', 1, '2025-05-26 16:52:43', '2026-08-13 10:24:45', NULL);
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
) ENGINE = InnoDB AUTO_INCREMENT = 61 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户角色表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user_roles
-- ----------------------------
INSERT INTO `user_roles` VALUES (35, 2, 2, 'test', '2026-07-20 14:19:04', '2026-07-20 14:19:04', NULL);
INSERT INTO `user_roles` VALUES (55, 61, 2, 'test', '2026-08-04 16:28:39', '2026-08-04 16:28:39', NULL);
INSERT INTO `user_roles` VALUES (58, 10, 2, 'test', '2026-08-07 16:06:05', '2026-08-07 16:06:05', NULL);
INSERT INTO `user_roles` VALUES (64, 1, 1, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);
INSERT INTO `user_roles` VALUES (65, 10, 1, 'admin', '2026-08-13 10:24:45', '2026-08-13 10:24:45', NULL);

SET FOREIGN_KEY_CHECKS = 1;
