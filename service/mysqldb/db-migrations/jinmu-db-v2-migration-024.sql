SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

 -- 把三水项目涉及的账号全部设为已过期
UPDATE `user`
SET `is_activated` = 0
WHERE `customized_code` = 'custom_sanshui';
