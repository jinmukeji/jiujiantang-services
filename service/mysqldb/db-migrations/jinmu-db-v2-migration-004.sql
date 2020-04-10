-- 初始化组织的数据
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

SET @now = NOW( );

INSERT INTO `organization`
(
	`legacy_account`,
	`name`,
	`phone`,
	`email`,
	`contact`,
	`type`,
	`country`,
	`state`,
	`city`,
	`district`,
	`street`,
	`address`,
	`postal_code`,
	`remark`,
	`customized_code`,
	`is_valid`,
	`created_at`,
	`updated_at`,
	`deleted_at`
)
SELECT
	P.account as legacy_account,
	IFNULL(C.company_name, '') as name,
	IFNULL(C.company_phone, '') as phone,
	IFNULL(C.company_email, '') as email,
	IFNULL(C.company_representative, '') as contact,
	IFNULL(C.company_type, '') as type,
	'' as country,
	IFNULL(C.company_state, '') as state,
	IFNULL(C.company_city, '') as city,
	IFNULL(C.company_block, '') as district,
	IFNULL(C.company_address, '') as street,
	IFNULL(C.company_address, '') as address,
	'' as postal_code,
	'迁移数据初始化导入的组织' as remark,
	CASE
			WHEN P.account_type <> 'normal' THEN P.account_type 
			ELSE '' 
	END AS customized_code,
	'1' as is_valid,
	@now as created_at,
	@now as updated_at,
	NULL as deleted_at
FROM
	jinmu_product AS P 
INNER JOIN jinmu_company AS C
ON P.account = C.account AND C.is_valid = 1
WHERE
	P.is_valid = 1;


SET FOREIGN_KEY_CHECKS = 1;
