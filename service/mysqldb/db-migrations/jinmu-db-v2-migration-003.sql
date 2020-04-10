-- 初始化组织的 Owner 用户数据
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;


SET @now = NOW( );
INSERT INTO `user` (
	`username`,
	`password`,
	`zone`,
	`register_type`,
	`register_time`,
	`register_source_client_id`,
	`customized_code`,
	`nickname`,
	`gender`,
	`birthday`,
	`height`,
	`weight`,
	`remark`,
	`is_profile_completed`,
	`is_activated`,
	`activated_at`,
	`deactivated_at`,
	`created_at`,
	`updated_at`,
	`deleted_at` 
) SELECT
	P.account,
	P.password,
	'CN' AS zone,
	'username' AS register_type,
	@now AS register_time,
	'jm-10001' AS register_source_client_id,
	CASE
			WHEN P.account_type <> 'normal' THEN P.account_type 
			ELSE '' 
	END AS customized_code,
	CONCAT( 'JM', P.account ) AS nickname,
	'M' AS gender,
	'1998-01-01' AS birthday,
	'170' AS height,
	'60' AS weight,
	'迁移数据导入的Owner用户' AS remark,
	0 AS is_profile_completed,
	'1' AS is_activated,
	@now AS activated_at,
	NULL AS deactivated_at,
	@now AS created_at,
	@now AS updated_at,
	NULL AS deleted_at 
FROM
	jinmu_product AS P 
WHERE
	P.is_valid = 1;

SET FOREIGN_KEY_CHECKS = 1;
