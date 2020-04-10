SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
START TRANSACTION;
-- 清除`user_subscription_sharing`数据
delete from `user_subscription_sharing`  
-- 添加主键
Alter table `user_subscription_sharing` add primary key(`subscription_id`,`user_id`);

-- 重新初始化user_subscription_sharing数据
INSERT INTO `user_subscription_sharing` (
    `subscription_id`,
    `user_id`, 
    `created_at`,
    `updated_at`,
    `deleted_at`
     ) SELECT
     S.subscription_id,
     OU.user_id,
     OU.created_at,
     OU.updated_at,
     OU.deleted_at
	FROM organization_user AS OU
     inner join organization_owner as OO on OU.organization_id = OO.organization_id 
    inner join subscription as S on S.owner_id = OO.owner_id and S.deleted_at IS NULL
WHERE        
	OU.deleted_at IS NULL;
COMMIT;
