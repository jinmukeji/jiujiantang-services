SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 添加批次
ALTER table subscription_activation_code ADD COLUMN batch VARCHAR(50)  COMMENT '批次' AFTER `expired_at`;

