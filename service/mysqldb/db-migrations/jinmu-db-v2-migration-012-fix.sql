-- 1 修正 Record 创建日期问题
UPDATE `record`
INNER JOIN `legacy_jinmu_record`
SET record.created_at=legacy_jinmu_record.create_date_utc 
WHERE record.record_id = legacy_jinmu_record.record_id;

-- 2 修正 User 创建日期问题

-- 2.1 Owner 用户
-- 局部老账号没有日期，初始化一个比较早的值
SET @long_long_ago = '2016-07-01 00:00:00';

UPDATE `user` AS U
INNER JOIN `legacy_jinmu_product` AS P
ON U.username = P.account
SET U.created_at = IFNULL(CONVERT_TZ(P.create_date,'+08:00','+00:00'), @long_long_ago);

-- 2.2 一般用户
UPDATE `user` AS U
INNER JOIN `legacy_jinmu_subject` AS S
ON U.legacy_subject_id = S.subject_id
SET U.created_at = CONVERT_TZ(S.create_date,'+08:00','+00:00');
