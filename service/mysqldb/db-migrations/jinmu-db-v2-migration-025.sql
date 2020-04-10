SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 添加has_stress_state
ALTER table `record` ADD COLUMN `has_stress_state` tinyint(4) NOT NULL DEFAULT '0'  COMMENT '是否是应激态' AFTER `s3_key`;

-- 添加stress_state
ALTER table `record` ADD COLUMN `stress_state` varchar(1000)  COMMENT 'json格式的应激态状态，格式map[string]bool' AFTER `has_stress_state`;

-- 添加analyze_body
ALTER table `record` ADD COLUMN `analyze_body` varchar(8000)  COMMENT '新分析接口的body' AFTER `stress_state`;
