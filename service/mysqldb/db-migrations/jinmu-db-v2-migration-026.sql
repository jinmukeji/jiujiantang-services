-- 添加昵称首字母
ALTER table user_profile ADD COLUMN nickname_initial VARCHAR(4)  COMMENT '昵称首字母' AFTER `nickname`;
