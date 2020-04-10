-- 初始化 client 表

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

SET @now = NOW();
INSERT INTO `client` ( `client_id`, `secret_key`, `name`, `zone`,`customized_code`, `remark`, `created_at`, `updated_at`, `deleted_at` )
VALUES
	( 'jm-10001', '3+o6Q2y7+g.tzF,U=4qpy7orzd9@(}X8', 'JinmuHealth-app', 'CN', '' , '喜马把脉健康APP', @now, @now, NULL ),
	( 'jm-10002', 'v9n7e/B%EgvD7^P9%UV37^fs7T3*z^(g', 'JinmuHealth-android-tablet','CN','', '喜马把脉健康APP Android一体机版', @now, @now, NULL ),
	( 'jm-10003', 'mJk#W&y7?F^8A7+V6f6.7]9qg4t)nsmo', 'JinmuHealth-android-enginneringtool','CN','', '喜马把脉健康APP Android工程测试工具', @now, @now, NULL ),
	( 'dengyun-10001', '2=827Uor/^L76szN%.zTpba#8(?3stMg', 'dengyun-server', 'CN-X','custom_dengyun', '登云服务端Application', @now, @now, NULL );

SET FOREIGN_KEY_CHECKS = 1;
