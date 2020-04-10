SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
-- 修改原先subscription_transfer_record错误的COMMENT
ALTER  TABLE `subscription_transfer_record` COMMENT '订阅转让记录表';
-- 修改枚举类型
ALTER  TABLE `user_profile` MODIFY COLUMN `gender` varchar(10) COMMENT '用户性别 M 男 F 女';
ALTER  TABLE `local_notifications` MODIFY COLUMN `frequency` varchar(30) COMMENT '推送时间间隔基本单位,FrequencyDaily,FrequencyWeekly,FrequencyMonthly';
ALTER  TABLE `audit_user_credential_update` MODIFY COLUMN `updated_record_type` varchar(20) COMMENT '更新记录的类型,username,phone,email,password';
ALTER  TABLE `audit_user_signin_signout` MODIFY COLUMN `record_type` varchar(10) NOT NULL COMMENT '登录或登出记录的类型,signin,signout';
ALTER  TABLE `verification_code` MODIFY COLUMN `send_via` varchar(10) NOT NULL  COMMENT '发送方式,phone,email';
-- 修改user_id类型
ALTER  TABLE `audit_user_signin_signout` MODIFY COLUMN `user_id` INT(10) UNSIGNED COMMENT '用户ID';
ALTER  TABLE `user_profile` MODIFY COLUMN `user_id` INT(10) UNSIGNED COMMENT '用户ID';
ALTER  TABLE `phone_or_email_verfication` MODIFY COLUMN `user_id` INT(10) UNSIGNED COMMENT '用户ID';
ALTER  TABLE `pn_record` MODIFY COLUMN `user_id` INT(10) UNSIGNED COMMENT '用户ID';
--  删除2个索引
ALTER TABLE `sem_record` DROP  KEY `idx_to_address`;
ALTER TABLE `sms_record` DROP  KEY `idx_phone`;
-- 添加索引
ALTER TABLE `user_access_token` ADD  KEY `idx_expired_at` (`expired_at`) USING BTREE;
ALTER TABLE `verification_code` 
ADD  KEY `idx_sn` (`sn`) USING BTREE,
ADD  KEY `idx_code` (`code`) USING BTREE,
ADD  KEY `idx_send_to` (`send_to`) USING BTREE;
ALTER TABLE `wechat_user` ADD  KEY `idx_user_id` (`user_id`) USING BTREE;
--  废弃一个表
ALTER TABLE `organization_admin` RENAME `legacy_organization_admin`;

-- 修改phone_or_email_verfication的枚举类型
ALTER  TABLE `phone_or_email_verfication` MODIFY COLUMN `verification_type` varchar(10) NOT NULL  COMMENT '发送方式,phone,email';
-- 删除sms_record多余字段
ALTER  TABLE `sms_record` DROP COLUMN `expired_at`;
ALTER  TABLE `sms_record` DROP COLUMN `is_valid`;
-- 增加默认值
ALTER  TABLE `feedback` ALTER COLUMN `content` SET default '';
ALTER  TABLE `local_notifications` ALTER COLUMN `has_weekdays` SET default 0;
ALTER  TABLE `local_notifications` ALTER COLUMN `has_month_days` SET default 0;
-- 限制字段不可为空
ALTER TABLE `phone_or_email_verfication` MODIFY  `expired_at` datetime NOT NULL COMMENT '到期时间';
ALTER TABLE `subscription_transfer_record` MODIFY  `subscription_id` int(10) unsigned NOT NULL COMMENT '订阅ID';
ALTER TABLE `subscription_transfer_record` MODIFY  `user_id` int(10) unsigned NOT NULL COMMENT '用户ID';
ALTER TABLE `subscription_transfer_record` MODIFY  `old_user_id` int(10) unsigned NOT NULL COMMENT '原用户ID';
ALTER TABLE `verification_code` MODIFY  `expired_at` datetime NOT NULL COMMENT '到期时间';
ALTER TABLE `wechat_user` MODIFY  `user_id` int(10) unsigned NOT NULL COMMENT '喜马把脉用户ID';
