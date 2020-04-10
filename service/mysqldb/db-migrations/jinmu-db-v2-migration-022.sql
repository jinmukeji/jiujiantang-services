SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 添加s3-key
ALTER table `record` ADD COLUMN `s3_key` VARCHAR(50)  COMMENT 's3的key' AFTER `transaction_number`;
-- 迁移老的数据
UPDATE `record` as R set `s3_key` = concat('spec-v1/',R.record_id,'.txt') where R.record_id >= '100000';
