-- 初始化 subscription 数据

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

SET @now = NOW();

INSERT INTO `subscription`
(
	`organization_id`,
	`subscription_type`,
	`active`,
	`activated_at`,
	`expired_at`,
	`customized_code`,
	`max_user_limits`,
	`contract_year`,
	`created_at`,
	`updated_at`,
	`deleted_at`
)
SELECT
	O.organization_id,
	CASE
		WHEN C.max_client < 100 THEN '5' -- 礼品版
		WHEN C.max_client >= 100 AND C.contract_year = 1 THEN '2' -- 黄金姆
		WHEN C.max_client >= 100 AND C.contract_year = 2 THEN '3' -- 白金姆
		WHEN C.max_client >= 100 AND C.contract_year = 5 THEN '4' -- 钻石姆
		WHEN C.max_client >= 100 AND C.contract_year = 100 THEN '4' -- 钻石姆
		ELSE '0' -- 定制版
	END	AS subscription_type,
	CASE
		WHEN C.contract_start_date IS NOT NULL THEN '1'
		ELSE '0'
	END AS active,
	CONVERT_TZ(C.contract_start_date,'+08:00','+00:00') AS activated_at,
	CONVERT_TZ(C.contract_end_date,'+08:00','+00:00') AS expired_at,
	'' AS customized_code,
	C.max_client AS max_user_limits,
	C.contract_year AS contract_year,
	@now as created_at,
	@now as updated_at,
	NULL as deleted_at
FROM jinmu_contract AS C
INNER JOIN organization AS O
ON O.legacy_account = C.account
WHERE C.is_valid = 1;


SET FOREIGN_KEY_CHECKS = 1;
