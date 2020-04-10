-- 初始化 device_organization_binding 数据

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

SET @now = NOW();

INSERT INTO `device_organization_binding`
(
	`device_id`,
	`organization_id`,
	`created_at`,
	`updated_at`
)
SELECT 
	D.device_id,
	O.organization_id,
	@now,
	@now
FROM jinmu_mac AS M
INNER JOIN device AS D
ON CONV( M.mac, 16, 10 )  = D.mac
INNER JOIN organization AS O
ON O.legacy_account = M.account;

SET FOREIGN_KEY_CHECKS = 1;
