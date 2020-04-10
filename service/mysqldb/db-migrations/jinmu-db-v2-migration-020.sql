SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 添加锁
ALTER table subscription_activation_code ADD COLUMN `activation_lock` VARCHAR(50)  COMMENT '锁' AFTER `activated_at`;

