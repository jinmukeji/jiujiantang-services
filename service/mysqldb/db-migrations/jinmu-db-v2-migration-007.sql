-- 初始化 user_preferences 数据
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

SET @now = NOW( );
INSERT INTO `user_preferences`
(
	`user_id`,
	`enable_heart_rate_chart`,
	`enable_pulse_wave_chart`,
	`enable_warm_prompt`,
	`enable_choose_status`,
	`enable_constitution_differentiation`,
	`enable_syndrome_differentiation`,
	`enable_western_medicine_analysis`,
	`enable_meridian_bar_graph`,
	`enable_comment`,
	`enable_health_trending`,
	`created_at`,
	`updated_at`,
	`deleted_at`
)
SELECT 
	OU.user_id,
	AP.enable_heart_rate_chart,
	AP.enable_pulse_wave_chart,
	AP.enable_warm_prompt,
	AP.enable_choose_status,
	AP.enable_constitution_differentiation,
	AP.enable_syndrome_differentiation,
	AP.enable_western_medicine_analysis,
	AP.enable_meridian_bar_graph,
	AP.enable_comment,
	1 as enable_health_trending,
	@now as created_at,
	@now as updated_at,
	NULL as deleted_at
FROM organization_user AS OU
INNER JOIN organization AS O
ON OU.organization_id = O.organization_id
INNER JOIN jinmu_account_preferences AS AP
ON AP.account = O.legacy_account;

SET FOREIGN_KEY_CHECKS = 1;
