-- 初始化 二维码 数据
DROP TABLE IF EXISTS `wxmp_tmp_qrcode`;
CREATE TABLE `wxmp_tmp_qrcode` (
  `scene_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '场景ID',
  `raw_url` varchar(255) NOT NULL DEFAULT '' COMMENT '二维码原始URL',
  `ticket` varchar(255) NOT NULL DEFAULT '' COMMENT '二维码的ticket',
  `account` varchar(255)  NOT NULL COMMENT '账户',
  `machine_uuid` varchar(255)  NOT NULL COMMENT '机器的标识',
  `origin_id` varchar(255)  NOT NULL COMMENT '公众号原始ID',
  `expired_at` timestamp NOT NULL COMMENT '二维码的过期时间',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`scene_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='二维码';

-- 初始化 扫码记录 数据
DROP TABLE IF EXISTS `scanned_qrcode_record`;
CREATE TABLE `scanned_qrcode_record` (
  `record_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '扫码记录ID',	
  `scene_id` int(10) unsigned NOT NULL  COMMENT '场景ID',
  `user_id` int(10) unsigned NOT NULL COMMENT '用户ID',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`record_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='二维码扫码记录';

-- 初始化 一体机账户 数据
DROP TABLE IF EXISTS `jinmu_l_account`;
CREATE TABLE `jinmu_l_account` (
  `account` varchar(255)  NOT NULL COMMENT '账户',
  `password` varchar(255) NOT NULL COMMENT '用户登录密码',
  `organization_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '与设备关联的组织ID',
  `remark` varchar(255) NOT NULL COMMENT '账户备注',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
	KEY `idx_organization_id` (`organization_id`) USING BTREE,
  PRIMARY KEY (`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='一体机账户';

-- 初始化 微信用户 数据
DROP TABLE IF EXISTS `wechat_user`;
CREATE TABLE `wechat_user` (
  `open_id` varchar(255)  NOT NULL COMMENT '微信OpenID',
  `union_id` varchar(255)  NOT NULL COMMENT '微信UnionID',
  `user_id` int(10) unsigned  COMMENT '喜马把脉用户ID',
  `origin_id` varchar(255)  NOT NULL COMMENT '公众号原始ID',
  `nickname` varchar(255) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `avatar_image_url` varchar(255) NOT NULL DEFAULT '' COMMENT '头像图片URL',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  KEY `idx_user_id` (`user_id`) USING BTREE,
  PRIMARY KEY (`union_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='微信用户';

-- 初始化 token 数据
DROP TABLE IF EXISTS `jinmu_l_access_token`;
CREATE TABLE `jinmu_l_access_token` (
  `account` varchar(255)  NOT NULL COMMENT '账户',
  `token` varchar(36) NOT NULL COMMENT '当前会话的令牌',
  `expired_at` datetime NOT NULL COMMENT '会话过期时间',
  `machine_uuid` varchar(255)  NOT NULL COMMENT '机器的标识',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  PRIMARY KEY (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户访问令牌';

ALTER table record
ADD COLUMN app_highest_heart_rate int(10) AFTER `is_valid` , 
ADD COLUMN app_lowest_heart_rate int(10) AFTER `app_highest_heart_rate` ,
ADD COLUMN has_paid tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否完成支付' AFTER `app_lowest_heart_rate`,
ADD COLUMN show_full_report tinyint(4) NOT NULL DEFAULT '1' COMMENT '是否显示完成测量报告' AFTER `has_paid`,
ADD COLUMN has_sent_wx_view_report_notification tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否已经发送过微信查看报告通知' AFTER `show_full_report`;

-- 初始化 session 数据
DROP TABLE IF EXISTS `wechat_h5_session`;
CREATE TABLE `wechat_h5_session` (
  `session_id` varchar(255)  NOT NULL COMMENT 'Session ID',
  `state` varchar(255) NOT NULL COMMENT '微信 OAuth 验证的 state',
  `open_id` varchar(255) NOT NULL COMMENT '微信OpenID',
  `union_id` varchar(255)  NOT NULL COMMENT '微信UnionID',
  `user_id` int(10) NOT NULL COMMENT '用户ID',
  `authorized` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否已经验证通过',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `expired_at` timestamp NOT NULL COMMENT '到期时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`session_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户访问令牌';

