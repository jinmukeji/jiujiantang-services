SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
-- ----------------------------
-- Table structure for user_profile
-- ----------------------------
CREATE TABLE `user_profile` (
    `user_id` INT(10) NOT NULL COMMENT '用户ID',
    `nickname` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '用户昵称',
    `gender` ENUM('M', 'F') NOT NULL DEFAULT 'M' COMMENT '用户性别 M 男 F 女',
    `birthday` DATE NOT NULL COMMENT '用户生日',
    `height` SMALLINT(6) NOT NULL DEFAULT '0' COMMENT '用户身高 单位厘米',
    `weight` SMALLINT(6) NOT NULL DEFAULT '0' COMMENT '用户体重 单位千克',
    `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
    `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`user_id`)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8 COMMENT='用户档案';

-- ----------------------------
-- Table structure for audit_user_credential_update
-- ----------------------------
CREATE TABLE `audit_user_credential_update` ( 
    `record_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `user_id` INT(10) NOT NULL COMMENT '用户ID',
    `client_id` VARCHAR(255) NOT NULL COMMENT '客户端ID',
    `updated_record_type` ENUM('username', 'phone', 'email','password') NOT NULL DEFAULT 'email' COMMENT '更新记录的类型',
    `old_value` VARCHAR(255) NOT NULL COMMENT '旧的数值',
    `new_value` VARCHAR(255) NOT NULL COMMENT '新的数值',
    `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
    `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`record_id`),
    KEY `idx_user_id` (`user_id`) USING BTREE
)  ENGINE=INNODB DEFAULT CHARSET=UTF8 COMMENT='用户审计记录';

-- ----------------------------
-- Table structure for audit_user_signin_signout
-- ----------------------------
CREATE TABLE `audit_user_signin_signout` (
    `record_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `user_id` INT(10) NOT NULL COMMENT '用户ID',
    `client_id` VARCHAR(255) NOT NULL COMMENT '客户端ID',
    `ip` VARCHAR(20) NOT NULL COMMENT 'ip地址',
    `extra_params` VARCHAR(255) NOT NULL COMMENT '登录/登出参数',
    `record_type` ENUM('signin', 'signout') NOT NULL DEFAULT 'signin' COMMENT '登录或登出记录的类型',
    `sign_in_machine` VARCHAR(255) NOT NULL COMMENT '登录机器',
    `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
    `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`record_id`),
    KEY `idx_user_id` (`user_id`) USING BTREE
)  ENGINE=INNODB DEFAULT CHARSET=UTF8 COMMENT='用户登录/登出记录';

--  添加 enable_location_notification
ALTER table user_preferences ADD COLUMN enable_location_notification tinyint(4) NOT NULL DEFAULT 1 COMMENT '是否开启本地通知' AFTER `enable_health_trending`;

ALTER TABLE user RENAME legacy_user;

-- ----------------------------
-- Table structure for user
-- ----------------------------
CREATE TABLE `user` (
    `user_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `register_type` VARCHAR(20) NOT NULL COMMENT '用户注册方式',
    `register_time` TIMESTAMP NOT NULL COMMENT '注册时间',
    `zone` VARCHAR(255) NOT NULL COMMENT '用户选择的区域',
    `customized_code` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '定制化代码',
    `user_defined_code` varchar(255) DEFAULT NULL COMMENT '用户自定义代码',
    `remark` VARCHAR(255)CHARACTER SET UTF8MB4 COLLATE UTF8MB4_GENERAL_CI NOT NULL DEFAULT '' COMMENT '用户备注',
    `encrypted_password` VARCHAR(255) NOT NULL COMMENT '用户encrypt后的登录密码',
    `seed` VARCHAR(20) NOT NULL COMMENT '随机种子',
    `secure_email` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '安全邮箱',
    `signin_phone` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '登录Phone',
    `signin_username` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '登录用户名',
    `nation_code` VARCHAR(50) NOT NULL COMMENT '国家代码',
    `language` VARCHAR(50) NOT NULL DEFAULT 'zh-Hans' COMMENT '语言',
    `has_set_email` TINYINT(4) NOT NULL COMMENT '是否设置邮箱',
    `has_set_phone` TINYINT(4) NOT NULL COMMENT '是否设置电话',
    `has_set_username` TINYINT(4) NOT NULL COMMENT '是否设置用户',
    `has_set_password` TINYINT(4) NOT NULL COMMENT '是否设置密码',
    `has_set_user_profile` TINYINT(4) NOT NULL COMMENT '是否设置用户详情',
    `has_set_secure_questions` TINYINT(4) NOT NULL COMMENT '是否设置密保问题',
    `has_set_language` TINYINT(4) NOT NULL COMMENT '是否设置语言',
    `is_profile_completed` TINYINT(4) NOT NULL COMMENT '是否完成初始化',
    `register_source` VARCHAR(20) NOT NULL COMMENT '用户来源',
    `latest_login_time` TIMESTAMP NULL DEFAULT NULL  COMMENT '最近登录时间',
    `secure_question_1` VARCHAR(255)  COMMENT '密保问题1',
    `secure_question_2` VARCHAR(255)  COMMENT '密保问题2',
    `secure_question_3` VARCHAR(255)  COMMENT '密保问题3',
    `secure_answer_1` VARCHAR(255)  COMMENT '密保答案1',
    `secure_answer_2` VARCHAR(255)  COMMENT '密保答案2',
    `secure_answer_3` VARCHAR(255)  COMMENT '密保答案3',
    `latest_updated_email_at` TIMESTAMP NULL DEFAULT NULL COMMENT '最新更新邮箱的时间',
    `latest_updated_phone_at` TIMESTAMP NULL DEFAULT NULL COMMENT '最新更新电话的时间',
    `latest_updated_username_at` TIMESTAMP NULL DEFAULT NULL COMMENT '最新更新用户名的时间',
    `latest_updated_password_at` TIMESTAMP NULL DEFAULT NULL COMMENT '最新更新密码的时间',
    `latest_updated_secure_questions_at` TIMESTAMP NULL DEFAULT NULL COMMENT '最新更新密保问题的时间',
    `region` VARCHAR(255) DEFAULT 'mainland_china' COMMENT '区域',
    `has_set_region` TINYINT(4) NOT NULL COMMENT '是否设置区域',
    `is_activated` tinyint(4) NOT NULL DEFAULT '1' COMMENT '是否激活',
    `activated_at` timestamp NULL DEFAULT NULL COMMENT '激活时间',
    `deactivated_at` timestamp NULL DEFAULT NULL COMMENT '取消激活（禁用）时间',
    `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
    `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`user_id`),
    KEY `idx_signin_username` (`signin_username`) USING BTREE,
    KEY `idx_secure_email` (`secure_email`) USING BTREE,
    KEY `idx_signin_phone` (`signin_phone`) USING BTREE
)  ENGINE=INNODB DEFAULT CHARSET=UTF8 COMMENT='用户登录/登出记录';

SET @now = NOW();

-- 初始化用户数据
INSERT INTO `user` (
    `user_id`,
    `register_type`,
    `register_time`,
    `zone`,
    `customized_code`,
    `user_defined_code`,
    `remark`,
    `signin_username`,
    `nation_code`,
    `has_set_email`,
    `has_set_phone`,
    `has_set_username`,
    `has_set_password`,
    `has_set_user_profile`,
    `is_profile_completed`,
    `language`,
    `register_source`,
    `latest_login_time`,
    `secure_question_1`, 
    `secure_question_2`,
    `secure_question_3`,
    `secure_answer_1`,
    `secure_answer_2`,
    `secure_answer_3`,
    `latest_updated_email_at`, 
    `latest_updated_phone_at`,
    `latest_updated_username_at`,
    `latest_updated_password_at`,
    `is_activated`,
    `activated_at`,
    `deactivated_at`,
    `created_at`,
	`updated_at`,
	`deleted_at`
) SELECT
    LU.user_id,
    LU.register_type,
    LU.register_time,
    LU.zone,
    LU.customized_code,
    LU.user_defined_code,
    LU.remark,
    LU.username,
    '' as nation_code,
    0 AS has_set_email,
    0 AS has_set_phone,
	CASE
			WHEN LU.username <> '' THEN 1
			ELSE 0 
	END AS has_set_username,
	CASE
			WHEN LU.password <> '' THEN 1
			ELSE 0         
	END AS has_set_password,
    0 AS has_set_user_profile,
    LU.is_profile_completed,
    'zh-Hans' AS language,
    '数据库迁移' AS register_source, 
    NULL AS  latest_login_time,
    NULL AS  secure_question_1,
    NULL AS  secure_question_2, 
    NULL AS  secure_question_3,
    NULL AS  secure_answer_1,
    NULL AS  secure_answer_2, 
    NULL AS  secure_answer_3,
    NULL AS  latest_updated_email_at,
    NULL AS  latest_updated_phone_at, 
    NULL AS  latest_updated_username_at, 
    NULL AS  latest_updated_password_at,
    LU.is_activated,
    LU.activated_at,
    LU.deactivated_at,
    @now AS created_at,
	@now AS updated_at,
	NULL AS deleted_at                                
FROM
	legacy_user AS LU
WHERE        
	LU.deleted_at IS NULL;   

INSERT INTO `user_profile` (
    `user_id`,
    `nickname`,
    `gender`,
    `birthday`,
    `height`,
    `weight`,
    `created_at`,
    `updated_at`,
    `deleted_at`
) SELECT
    LU.user_id,
    LU.nickname,
    LU.gender,
    LU.birthday,
    LU.height,
    LU.weight,
    @now AS created_at,
	@now AS updated_at,
	NULL AS deleted_at                                
FROM
	legacy_user AS LU
WHERE        
	LU.deleted_at IS NULL;

-- 初始化 verification_code 数据
CREATE TABLE `verification_code` (
    `record_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `usage` VARCHAR(20) DEFAULT '' NOT NULL COMMENT '使用用途',
    `sn` VARCHAR(50)   DEFAULT '' COMMENT '序列号',
    `code` VARCHAR(6)   DEFAULT '' COMMENT '验证码',
    `send_via` ENUM('phone', 'email') NOT NULL DEFAULT 'phone' COMMENT '发送方式',
    `nation_code` VARCHAR(50) NOT NULL COMMENT '国家代码',
    `send_to` VARCHAR(255) DEFAULT '' NOT NULL COMMENT '接收人',
    `expired_at` DATETIME NULL DEFAULT NULL COMMENT '到期时间',
    `has_used` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否使用过',
    `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
    `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`record_id`)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8 COMMENT='验证码记录';


-- 初始化 verification_phone_and_email 验证手机与邮件
CREATE TABLE `phone_or_email_verfication` (
    `record_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `user_id` INT(10) NOT NULL COMMENT '用户ID',
    `verification_type` ENUM('phone', 'email') COMMENT '验证类型',
    `verification_number` VARCHAR(50) DEFAULT '' COMMENT '验证号',
    `expired_at` DATETIME NULL DEFAULT NULL COMMENT '到期时间',
    `has_used` TINYINT(4) NOT NULL DEFAULT '0' COMMENT '是否使用过',
    `send_to` VARCHAR(255) DEFAULT '' NOT NULL COMMENT '接收人',
    `nation_code` VARCHAR(50) NOT NULL COMMENT '国家代码',
    `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
    `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`record_id`),
    KEY `idx_verification_number` (`verification_number`) USING BTREE
)  ENGINE=INNODB DEFAULT CHARSET=UTF8 COMMENT='验证码记录';

-- 初始化user_subscription_sharing
CREATE TABLE `user_subscription_sharing` (
    `subscription_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '订阅ID',
    `user_id` int(10) unsigned NOT NULL COMMENT '用户ID',
    `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
    `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    KEY `idx_subscription_id` (`subscription_id`) USING BTREE,
    KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='拥有或被分享的订阅';

-- 初始化user_subscription_sharing数据
INSERT INTO `user_subscription_sharing` (
    `subscription_id`,
    `user_id`, 
    `created_at`,
    `updated_at`,
    `deleted_at`
     ) SELECT
     S.subscription_id,
     OU.user_id,
     OU.created_at,
     OU.updated_at,
     OU.deleted_at
	FROM organization_user AS OU
    inner join subscription as S on S.organization_id = OU.organization_id
WHERE        
	OU.deleted_at IS NULL;

-- 初始化subscription_renew_record
CREATE TABLE `subscription_renew_record` (
    `record_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '订阅修改记录ID',
    `subscription_id` int(10) unsigned NOT NULL  COMMENT '订阅ID',
    `user_id` int(10) unsigned NOT NULL COMMENT '用户ID',
    `client_id` varchar(255) NOT NULL COMMENT '客户端ID',
    `old_contract_year` int(11) NOT NULL COMMENT '修改前的订阅合同年限',
    `old_max_user_limits` smallint(6) NOT NULL COMMENT '修改前的最大用户（常客）数量',
    `max_user_limits` int(11) DEFAULT NULL COMMENT '最大用户（常客）数量',
    `contract_year` smallint(6) DEFAULT NULL COMMENT '订阅合同年限',
    `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
    `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`record_id`),
    KEY `idx_subscription_id` (`subscription_id`) USING BTREE,
    KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='订阅修改记录表';

--  初始化subscription_activation_code
CREATE TABLE `subscription_activation_code` (
    `code` VARCHAR(255) NOT NULL COMMENT '激活码',
    `user_id` INT(10) UNSIGNED DEFAULT NULL COMMENT '使用激活码的用户ID',
    `subscription_id` INT(10) UNSIGNED DEFAULT NULL COMMENT '激活生成的订阅ID',
    `subscription_type` tinyint(4) DEFAULT '0' COMMENT '0 定制化 1 试用版 2 黄金姆 3 白金姆 4 钻石姆 5 礼品版',
    `max_user_limits` INT(11) DEFAULT NULL COMMENT '最大用户（常客）数量',
    `contract_year` SMALLINT(6) DEFAULT NULL COMMENT '订阅合同年限',
    `checksum` VARCHAR(255) NOT NULL COMMENT '校验位',
    `activated` TINYINT(4) DEFAULT '0' COMMENT '是否激活 0未激活 1已激活',
    `activated_at` DATETIME NULL DEFAULT NULL COMMENT '激活时间',
    `sold` TINYINT(4) DEFAULT '0' COMMENT '是否售出 0未售出 1已售出',
    `sold_at` DATETIME NULL DEFAULT NULL COMMENT '售出时间',
    `expired_at` DATETIME NULL DEFAULT NULL COMMENT '到期时间',
    `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
    `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`code`),
    KEY `idx_subscription_id` (`subscription_id`) USING BTREE,
    KEY `idx_user_id` (`user_id`) USING BTREE
)  ENGINE=INNODB DEFAULT CHARSET=UTF8 COMMENT='激活码表';

-- 初始化subscription_transfer_record
CREATE TABLE `subscription_transfer_record` (
    `record_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '订阅转让记录ID',
    `subscription_id` INT(10) UNSIGNED DEFAULT NULL COMMENT '订阅ID',
    `user_id` INT(10) UNSIGNED DEFAULT NULL COMMENT '用户ID',
    `client_id` VARCHAR(255) NOT NULL COMMENT '客户端ID',
    `old_user_id` INT(10) UNSIGNED DEFAULT NULL COMMENT '原用户ID',
    `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
    `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`record_id`),
    KEY `idx_subscription_id` (`subscription_id`) USING BTREE,
    KEY `idx_user_id` (`user_id`) USING BTREE
)  ENGINE=INNODB DEFAULT CHARSET=UTF8 COMMENT='订阅修改记录表';

ALTER TABLE subscription RENAME legacy_subscription;

-- 初始化subscription
CREATE TABLE `subscription` (
  `subscription_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '订阅ID',
  `subscription_type` tinyint(4) DEFAULT '0' COMMENT '0 定制化 1 试用版 2 黄金姆 3 白金姆 4 钻石姆 5 礼品版',
  `customized_code` varchar(255) DEFAULT '' COMMENT '自定义代码',
  `activated` tinyint(4) DEFAULT '0' COMMENT '是否激活 0未激活 1已激活',
  `activated_at` DATETIME NULL DEFAULT NULL COMMENT '激活时间',
  `expired_at` DATETIME NULL DEFAULT NULL COMMENT '到期时间',
  `max_user_limits` int(11) DEFAULT NULL COMMENT '最大用户（常客）数量',
  `contract_year` smallint(6) DEFAULT NULL COMMENT '订阅合同年限',
  `owner_id` int(10) unsigned NOT NULL  COMMENT '拥有者ID',
  `is_selected` tinyint(4) DEFAULT '0' COMMENT '是否选择该订阅作为正在使用的订阅 0未选择 1已选择',
  `is_migrated_activated` tinyint(4) DEFAULT '1' COMMENT '是否是迁移前激活状态 0未激活 1已激活',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`subscription_id`),
  KEY `idx_owner_id` (`owner_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='产品服务订阅';

-- 初始化subscription 数据
INSERT INTO `subscription` (
    `subscription_id`,
    `subscription_type`, 
    `customized_code`,
    `activated`,
    `activated_at`,
    `expired_at`,
    `max_user_limits`,
    `contract_year`,
    `owner_id`,
    `is_selected`,
    `is_migrated_activated`,
    `created_at`,
    `updated_at`,
    `deleted_at`
     ) SELECT
     LS.subscription_id,
     LS.subscription_type,
     LS.customized_code,
     LS.active,
     LS.activated_at,
     LS.expired_at,
     LS.max_user_limits,
     LS.contract_year,
     OO.owner_id,
     '1' as is_selected,
     LS.active as is_migrated_activated,
     LS.created_at,
     LS.updated_at,
     LS.deleted_at
	FROM legacy_subscription AS LS
    inner join organization_owner as OO on OO.organization_id = LS.organization_id
    inner join user as U on U.user_id = OO.owner_id
WHERE        
	LS.deleted_at IS NULL;

-- user_used_device 数据
CREATE TABLE `user_used_device` (
  `user_id` int(10) unsigned NOT NULL  COMMENT '订阅ID',
  `device_id` int(10) unsigned NOT NULL  COMMENT '设备ID',
  `client_id` varchar(255) NOT NULL COMMENT '客户端ID',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`user_id`,`device_id`,`client_id`),
  KEY `idx_user_id` (`user_id`) USING BTREE,
  KEY `idx_device_id` (`device_id`) USING BTREE,
  KEY `idx_client_id` (`client_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='使用过的脉诊仪';

--  添加 usage
ALTER table client ADD COLUMN `usage` varchar(255) NOT NULL COMMENT '用途' AFTER `remark`;
INSERT INTO `client` ( `client_id`, `secret_key`, `name`, `zone`,`customized_code`, `remark`,`usage`,`created_at`, `updated_at`, `deleted_at` )
VALUES
	( 'jm-10006', 'CHSnHWkepLThkmPw8IUX', 'JinmuHealth-web', 'CN', '' , '金姆ID官方网站','随时随地管理金姆帐号', @now, @now, NULL );

-- 更新client内容
UPDATE `client`
    SET `remark` = CASE `client_id`
        WHEN 'jm-10001' THEN '金姆健康APP'
        WHEN 'jm-10002' THEN '金姆健康一体机'
        WHEN 'jm-10004' THEN '金姆健康APP'
        WHEN 'jm-10005' THEN '金姆健康APP'    
    END,
	`usage` = CASE `client_id`
        WHEN 'jm-10001' THEN '家庭和机构的健康管家'
        WHEN 'jm-10002' THEN '快速检测,微信查看分析报告'
        WHEN 'jm-10004' THEN '家庭和机构的健康管家'
        WHEN 'jm-10005' THEN '家庭和机构的健康管家'
    END    
WHERE `client_id` IN ('jm-10001','jm-10002','jm-10004','jm-10005');


-- ----------------------------
-- Table structure for notification
-- ----------------------------
CREATE TABLE `notification_preferences` (
    `user_id` INT(10) UNSIGNED NOT NULL COMMENT '用户ID',
    `phone_enabled` TINYINT(4) NOT NULL DEFAULT 1 COMMENT '是否允许手机通知',
    `phone_enabled_updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT '最新更新是否允许手机通知状态的时间',
    `wechat_enabled` TINYINT(4) NOT NULL DEFAULT 1 COMMENT '是否允许微信通知',   
    `wechat_enabled_updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT '最新更新是否允许微信通知状态的时间',
    `weibo_enabled` TINYINT(4) NOT NULL DEFAULT 1 COMMENT '是否允许微博通知',
    `weibo_enabled_updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT '最新更新是否允许微博通知状态的时间',
    `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
    `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`user_id`)
) ENGINE=INNODB DEFAULT CHARSET=UTF8 COMMENT='通知配置首选项';

 -- 初始化 notification_preferences 数据
INSERT INTO `notification_preferences` (
    `user_id`,
    `phone_enabled`, 
    `wechat_enabled`,
    `weibo_enabled`,
    `phone_enabled_updated_at`,
    `wechat_enabled_updated_at`,
    `weibo_enabled_updated_at`,
    `created_at`,
    `updated_at`
     ) SELECT
     U.user_id,
     '1' as phone_enabled,
     '1' as wechat_enabled,
     '1' as weibo_enabled,
     @now,
     @now,
     @now,
     @now,
     @now
 FROM user AS U
WHERE        
 U.deleted_at IS NULL;


ALTER table record ADD COLUMN transaction_number varchar(255) DEFAULT NULL COMMENT '流水号' AFTER `measurement_posture`;

-- 添加系统备注字段
ALTER table user ADD COLUMN sys_remark VARCHAR(255) NOT NULL DEFAULT '' COMMENT '系统备注' AFTER `deactivated_at`;


