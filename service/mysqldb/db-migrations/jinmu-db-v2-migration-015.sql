--  添加 record_token
ALTER table record ADD COLUMN record_token varchar(255) AFTER `has_sent_wx_view_report_notification`;


-- 初始化 push_notification 数据
DROP TABLE IF EXISTS `push_notification`;
CREATE TABLE `push_notification` (
  `pn_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Message ID',
  `pn_display_time` varchar(255) NOT NULL COMMENT '消息上的显示的时间',
  `pn_title` varchar(255) NOT NULL COMMENT '消息的标题',
  `pn_image_url` varchar(255)  NOT NULL COMMENT '消息图片的URL',
  `pn_type` tinyint(4) NOT NULL DEFAULT 0 COMMENT '消息推送方式 0广播所有人 1设备标签 2设备别名 3Registration ID 4用户分群推送',
  `pn_content_url` varchar(255) NOT NULL COMMENT '消息内容的URL',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`pn_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='推送通知';

-- 初始化 pn_record 阅读消息的数据
DROP TABLE IF EXISTS `pn_record`;
CREATE TABLE `pn_record` (
  `pn_id` int(10) unsigned NOT NULL COMMENT 'Push Notification ID',
  `user_id` int(10)  NOT NULL COMMENT 'User ID',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  KEY `idx_pn_id` (`pn_id`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE,
  PRIMARY KEY (`pn_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='User阅读推送通知的记录';

--  添加 has_ae_error
ALTER table record ADD COLUMN has_ae_error tinyint(4) NOT NULL DEFAULT 0 COMMENT 'ae的结果是否异常' AFTER `record_token`;

--  添加 measurement_posture
ALTER table record ADD COLUMN measurement_posture tinyint(4) NOT NULL DEFAULT 0 COMMENT '测量姿态 0 坐姿 1 站姿 2 躺姿' AFTER `has_ae_error`;
