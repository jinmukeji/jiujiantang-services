-- 初始化 organization_owner 数据
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

SET @now = NOW( );

INSERT INTO `organization_owner`
(
	`organization_id`,
	`owner_id`,
	`created_at`,
	`updated_at`,
	`deleted_at`
)
SELECT 
	o.organization_id as organization_id,
	u.user_id as owner_id, 
	@now as created_at,
	@now as updated_at,
	NULL as deleted_at
FROM user as u
LEFT JOIN organization AS o
ON u.username = o.legacy_account;

SET FOREIGN_KEY_CHECKS = 1;
