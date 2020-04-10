DROP TABLE IF EXISTS `sem_record`;
CREATE TABLE `sem_record` (
   `sem_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '邮件ID',
   `to_address` VARCHAR(50) NOT NULL COMMENT '收件邮箱地址',
   `sem_status` TINYINT(4) NOT NULL DEFAULT '0' COMMENT '0 待发送 1 发送中 2 发送成功 3 发送失败',
   `template_action` VARCHAR(20) NOT NULL COMMENT '模版行为',
   `platform_type` ENUM('Aliyun', 'NetEase') NOT NULL DEFAULT 'Aliyun' COMMENT '平台类型',
   `template_param` VARCHAR(255) NOT NULL COMMENT '模版参数',
   `language` ENUM('zh-Hans', 'zh-Hant','en') NOT NULL DEFAULT 'zh-Hans' COMMENT  '语言',
   `sem_error_log` VARCHAR(1000) DEFAULT '' COMMENT  '邮箱错误信息',
   `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
   `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
   `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
   PRIMARY KEY (`sem_id`),
   KEY `idx_to_address` (`to_address`) USING BTREE
)  ENGINE=INNODB DEFAULT CHARSET=UTF8 COMMENT='邮箱记录';
