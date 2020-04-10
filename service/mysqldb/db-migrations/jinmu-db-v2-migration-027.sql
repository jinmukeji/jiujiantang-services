-- 添加索引
ALTER TABLE `record` ADD  KEY `idx_record_token` (`record_token`) USING BTREE;

-- 添加analyze_status
ALTER table `record` ADD COLUMN `analyze_status` tinyint(4)  COMMENT '分析的状态，0 pending,1 in_progress,2 completed,3 error' AFTER `analyze_body`;
