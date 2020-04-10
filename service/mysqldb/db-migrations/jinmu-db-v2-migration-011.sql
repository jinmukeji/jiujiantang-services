-- 初始化 record 表

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

SET @now = NOW();

INSERT INTO `record`
(
	`record_id`,
	`client_id`,
	`user_id`,
	`device_id`,
	`finger`,
	`app_heart_rate`,
	`is_sport_or_drunk`,
	`cold`,
	`menstrual_cycle`,
	`oviposit_period`,
	`lactation`,
	`pregnancy`,
	`cm_app_status_a`,
	`cm_app_status_b`,
	`cm_app_status_c`,
	`cm_app_status_d`,
	`cm_app_status_e`,
	`cm_app_status_f`,
	`c0`,
	`c1`,
	`c2`,
	`c3`,
	`c4`,
	`c5`,
	`c6`,
	`c7`,
	`c0cv`,
	`c1cv`,
	`c2cv`,
	`c3cv`,
	`c4cv`,
	`c5cv`,
	`c6cv`,
	`c7cv`,
	`g0`,
	`g1`,
	`g2`,
	`g3`,
	`g4`,
	`g5`,
	`g6`,
	`g7`,
	`heart_rate`,
	`heart_rate_cv`,
	`health_code`,
	`snr`,
	`dc_drift`,
	`elapsed`,
	`remark`,
	`sys_remark`,
	`record_type`,
	`answers`,
	`is_valid`,
	`created_at`,
	`updated_at`,
	`deleted_at`
)
SELECT 
	R.record_id,
	'jm-10001' AS client_id,
	U.user_id,
	D.device_id,
	CASE
		WHEN R.finger = 0 THEN 1 -- 左手食指
		WHEN R.finger = 1 THEN 6 -- 右手食指
		ELSE 4 -- 未知的，初始化为左手大拇指
	END AS finger,
	CASE 
		WHEN R.APP_HR < 0 OR  R.APP_HR > 250 THEN 0	-- 修正错误范围的心率为0
		ELSE R.APP_HR -- 正常范围
	END AS app_heart_rate,
	R.is_sport_or_drunk,
	R.cold,
	R.menstrual_cycle,
	R.oviposit_period,
	R.lactation,
	R.pregnancy,
	R.cm_app_status_a,
	R.cm_app_status_b,
	R.cm_app_status_c,
	R.cm_app_status_d,
	R.cm_app_status_e,
	R.cm_app_status_f,
	R.C0 as c0,
	R.C1 as c1,
	R.C2 as c2,
	R.C3 as c3,
	R.C4 as c4,
	R.C5 as c5,
	R.C6 as c6,
	R.C7 as c7,
	R.C0CV as c0cv,
	R.C1CV as c1cv,
	R.C2CV as c2cv,
	R.C3CV as c3cv,
	R.C4CV as c4cv,
	R.C5CV as c5cv,
	R.C6CV as c6cv,
	R.C7CV as c7cv,
	R.G0 as g0,
	R.G1 as g1,
	R.G2 as g2,
	R.G3 as g3,
	R.G4 as g4,
	R.G5 as g5,
	R.G6 as g6,
	R.G7 as g7,
	R.HR as heart_rate,
	R.HRCV as heart_rate_cv,
	R.health_code,
	R.snr,
	R.dc_drift,
	0 AS elapsed,
	CASE 
		WHEN R.comment IS NULL THEN ''	
		ELSE R.comment 
	END AS remark,
	'迁移旧版数据' as sys_remark,
	R.record_type,
	'' as answers,
	1 AS is_valid,
	IFNULL(CONVERT_TZ(R.create_date,'+08:00','+00:00'), @now) as created_at,
	@now as updated_at,
	NULL as deleted_at
FROM jinmu_record AS R
INNER JOIN user AS U
ON R.subject_id = U.legacy_subject_id
LEFT JOIN device AS D
ON D.mac =  CONV( REPLACE(R.appratus_mac, ':', ''), 16, 10 )
WHERE R.is_valid = 1;


SET FOREIGN_KEY_CHECKS = 1;
