DROP TABLE IF EXISTS `local_notifications`;
CREATE TABLE `local_notifications` (
    `ln_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '本地通知ID',
    `title` VARCHAR(36) COMMENT '推送标题',
    `content` VARCHAR(255) COMMENT '推送内容',
    `event_happen_at` DATETIME NOT NULL COMMENT '推送时间',
    `timezone` VARCHAR(36) NOT NULL COMMENT '时区信息',
    `frequency` ENUM('FrequencyDaily', 'FrequencyWeekly', 'FrequencyMonthly') NOT NULL DEFAULT 'FrequencyDaily' COMMENT '推送时间间隔基本单位',
    `interval` INT(32) COMMENT '推送时间间隔',
    `has_weekdays` TINYINT(1) COMMENT '是否需要以周为基本单位',
    `weekdays` VARCHAR(100) COMMENT '一周内有哪些天推送,这里保存的是json类型',
    `has_month_days` TINYINT(1) COMMENT '是否需要以月为基本单位',
    `month_days` VARCHAR(300) COMMENT '一个月内有哪些天推送,这里保存的是json类型',
    `max_notification_times` INT(32) COMMENT '最大推送次数',
    `end_at` DATETIME COMMENT '推送结束时间',
    `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
    `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`ln_id`)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8 COMMENT='本地通知';
