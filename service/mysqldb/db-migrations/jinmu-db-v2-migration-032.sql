-- 增加算法服务返回的最高最低心率到数据库
ALTER TABLE `record` 
ADD COLUMN `algorithm_highest_heart_rate` INT(10) NULL DEFAULT 0 COMMENT '算法服务计算得到的最高心率' AFTER `heart_rate`,
ADD COLUMN `algorithm_lowest_heart_rate` INT(10) NULL DEFAULT 0 COMMENT '算法服务计算得到的最低心率' AFTER `algorithm_highest_heart_rate`;
