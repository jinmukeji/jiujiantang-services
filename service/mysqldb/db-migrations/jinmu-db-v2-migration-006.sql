-- 初始化常客 User 数据
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

SET @now = NOW( );


INSERT INTO `user` (
	`legacy_subject_id`,
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
	`phone`,
	`email`,
	`country`,
	`state`,
	`city`,
	`district`,
	`street`,
	`user_defined_code`,
	`remark`,
	`is_profile_completed`,
	`is_activated`,
	`activated_at`,
	`deactivated_at`,
	`created_at`,
	`updated_at`,
	`deleted_at` 
) 
SELECT
	S.subject_id AS legacy_subject_id,
	CONCAT(S.account,'_', S.subject_id) as username,
	'' as password,
	'CN' AS zone,
	'legacy' AS register_type,
	@now AS register_time,
	'jm-10001' AS register_source_client_id,
	CASE
			WHEN P.account_type <> 'normal' THEN P.account_type 
			ELSE '' 
	END AS customized_code,
	S.name AS nickname,
	S.gender AS gender,
	S.birthdate AS birthday,
	IFNULL(S.phone, '') AS phone,
	IFNULL(S.email, '') AS email,
	IFNULL(S.nationality, '') AS country,
	IFNULL(S.state, '') AS state,
	IFNULL(S.city, '') AS city,
	IFNULL(S.block, '') AS district,
	'' AS street,
	S.customized_code AS user_defined_code,
	'迁移数据导入的常客用户' AS remark,
	'1' AS is_profile_completed,
	'1' AS is_activated,
	@now AS activated_at,
	NULL AS deactivated_at,
	@now AS created_at,
	@now AS updated_at,
	NULL AS deleted_at 
FROM jinmu_product AS P 
INNER JOIN jinmu_subject AS S
ON P.account = S.account AND S.is_valid = 1
WHERE
	P.is_valid = 1;

-- 更新身高和体重信息
UPDATE `user` AS U
LEFT JOIN jinmu_status AS ST
ON U.legacy_subject_id = ST.subject_id AND ST.is_valid = 1
SET 
	U.height = IFNULL(ST.height, '165'),
	U.weight = IFNULL(ST.weight, '65')
WHERE 
	U.legacy_subject_id IS NOT NULL;

-- 初始化组织与常客关联关系数据
-- 1/2 所有的 Owner 都作为一个常客
INSERT INTO `organization_user`
(
	`organization_id`,
	`user_id`,
	`created_at`,
	`updated_at`,
	`deleted_at`
)
SELECT
	organization_id,
	owner_id as user_id,
	@now as created_at,
	@now as updated_at,
	NULL as deleted_at
FROM organization_owner;

-- 2/2 导入原 jinmu_subject 中的用户
INSERT INTO `organization_user`
(
	`organization_id`,
	`user_id`,
	`created_at`,
	`updated_at`,
	`deleted_at`
)
SELECT 
	O.organization_id,
	U.user_id,
	@now AS created_at,
	@now AS updated_at,
	NULL AS deleted_at
FROM `user` AS U
INNER JOIN jinmu_subject AS S
ON U.legacy_subject_id = S.subject_id
INNER JOIN organization AS O
ON O.legacy_account = S.account
WHERE U.legacy_subject_id IS NOT NULL;

SET FOREIGN_KEY_CHECKS = 1;
