-- 初始化 account_l_record 数据
DROP TABLE IF EXISTS `account_l_record`;
CREATE TABLE `account_l_record` (
    `record_id` int(10) unsigned NOT NULL COMMENT '记录ID',
    `account` varchar(255)  NOT NULL COMMENT '账户',
    `created_at` TIMESTAMP NOT NULL COMMENT '数据记录创建时间',
    `updated_at` TIMESTAMP NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`record_id`)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8 COMMENT='一体机账户与记录关联表';
