-- 初始化 sms_record 表
DROP TABLE IF EXISTS `sms_record`;
CREATE TABLE `sms_record` (
    `sms_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '短信ID',
    `phone` VARCHAR(50) NOT NULL COMMENT '联系电话',
    `sms_status` TINYINT(4) NOT NULL DEFAULT '0' COMMENT '0 待发送 1 发送中 2 发送成功 3 发送失败',
    `template_action` VARCHAR(20) NOT NULL COMMENT '模版行为',
    `nation_code` VARCHAR(50) NOT NULL COMMENT '国家代码',
    `platform_type` ENUM('Aliyun', 'Tencent') NOT NULL DEFAULT 'Aliyun' COMMENT '平台类型',
    `template_param` VARCHAR(255) NOT NULL COMMENT '模版参数',
    `language` ENUM('zh-Hans', 'zh-Hant', 'en') NOT NULL DEFAULT 'zh-Hans' COMMENT '语言',
    `sms_error_log` VARCHAR(1000) DEFAULT '' COMMENT '短信错误信息',
    `serial_number` VARCHAR(50)   DEFAULT '' COMMENT '序列号',
    `expired_at` DATETIME NULL DEFAULT NULL COMMENT '到期时间',
    `is_valid` tinyint(4) NOT NULL DEFAULT '1' COMMENT '是否有效',
    `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
    `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`sms_id`),
    KEY `idx_phone` (`phone`) USING BTREE
)  ENGINE=INNODB DEFAULT CHARSET=UTF8 COMMENT='短信记录';
