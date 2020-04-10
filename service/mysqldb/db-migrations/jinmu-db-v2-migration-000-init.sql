SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for client
-- ----------------------------
DROP TABLE IF EXISTS `client`;
CREATE TABLE `client` (
  `client_id` varchar(255) NOT NULL COMMENT '客户端ID',
  `secret_key` varchar(255) NOT NULL COMMENT '客户端授权密钥',
  `name` varchar(255) NOT NULL COMMENT '客户端名称',
  `zone` varchar(255) NOT NULL COMMENT '客户端所在区域',
  `customized_code` varchar(255) NOT NULL DEFAULT '' COMMENT '定制化代码',
  `remark` varchar(255) NOT NULL COMMENT '客户端所备注',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='客户端';

-- ----------------------------
-- Table structure for device
-- ----------------------------
DROP TABLE IF EXISTS `device`;
CREATE TABLE `device` (
  `device_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '设备ID',
  `mac` bigint(20) unsigned NOT NULL COMMENT '设备MAC地址',
  `sn` varchar(255) NOT NULL COMMENT 'SN号',
  `pin` varchar(255) NOT NULL DEFAULT '' COMMENT '验证码',
  `zone` varchar(255) NOT NULL DEFAULT '' COMMENT '设备所在区域',
  `model` varchar(255) NOT NULL DEFAULT '' COMMENT '设备型号',
  `customized_code` varchar(255) NOT NULL DEFAULT '' COMMENT '定制化代码',
  `tags` varchar(255) NOT NULL DEFAULT '' COMMENT '标签列表',
  `remarks` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`device_id`),
  UNIQUE KEY `idx_sn` (`sn`) USING BTREE,
  KEY `idx_mac` (`mac`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='设备';

-- ----------------------------
-- Table structure for device_organization_binding
-- ----------------------------
DROP TABLE IF EXISTS `device_organization_binding`;
CREATE TABLE `device_organization_binding` (
  `device_id` int(10) unsigned NOT NULL COMMENT '设备ID',
  `organization_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '与设备关联的组织ID',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  PRIMARY KEY (`device_id`,`organization_id`),
  KEY `idx_organization_id` (`organization_id`) USING BTREE,
  KEY `idx_device_id` (`device_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='设备与组织绑定关系';

-- ----------------------------
-- Table structure for feedback
-- ----------------------------
DROP TABLE IF EXISTS `feedback`;
CREATE TABLE `feedback` (
  `feedback_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '意见反馈ID',
  `user_id` int(10) unsigned NOT NULL COMMENT '用户ID',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '反馈内容',
  `contact_way` varchar(255) NOT NULL DEFAULT '' COMMENT '联系方式',
  `is_valid` tinyint(4) NOT NULL DEFAULT '1' COMMENT '是否有效',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`feedback_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户反馈';

-- ----------------------------
-- Table structure for organization
-- ----------------------------
DROP TABLE IF EXISTS `organization`;
CREATE TABLE `organization` (
  `organization_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '组织ID',
  `legacy_account` varchar(50) NOT NULL DEFAULT '' COMMENT 'Legacy的acount信息',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '组织名称',
  `phone` varchar(50) NOT NULL DEFAULT '' COMMENT '联系电话',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '联系Email',
  `contact` varchar(255) NOT NULL DEFAULT '' COMMENT '联系人',
  `type` varchar(255) NOT NULL DEFAULT '' COMMENT '组织机构类型',
  `country` varchar(50) DEFAULT NULL COMMENT '组织机构所在国家',
  `state` varchar(255) NOT NULL DEFAULT '' COMMENT '组织所在省份',
  `city` varchar(255) NOT NULL COMMENT '组织所在城市',
  `district` varchar(255) DEFAULT NULL COMMENT '组织所在区',
  `street` varchar(255) NOT NULL COMMENT '组织所在街道',
  `address` varchar(255) NOT NULL COMMENT '组织地址',
  `postal_code` varchar(20) NOT NULL DEFAULT '' COMMENT '邮政编码',
  `remark` varchar(255) DEFAULT '' COMMENT '备注信息',
  `customized_code` varchar(255) NOT NULL DEFAULT '' COMMENT '定制化代码 例如 custom_sanshui custom_dengyun',
  `is_valid` tinyint(4) NOT NULL COMMENT '是否有效',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`organization_id`),
  KEY `idx_is_valid` (`is_valid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='组织';

-- ----------------------------
-- Table structure for organization_admin
-- ----------------------------
DROP TABLE IF EXISTS `organization_admin`;
CREATE TABLE `organization_admin` (
  `organization_id` int(10) unsigned NOT NULL COMMENT '组织ID',
  `admin_id` int(10) unsigned NOT NULL COMMENT '组织的管理员的用户ID',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`organization_id`,`admin_id`),
  KEY `idx_organization_id` (`organization_id`) USING BTREE,
  KEY `idx_admin_id` (`admin_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='组织与管理员的关系';

-- ----------------------------
-- Table structure for organization_owner
-- ----------------------------
DROP TABLE IF EXISTS `organization_owner`;
CREATE TABLE `organization_owner` (
  `organization_id` int(10) unsigned NOT NULL COMMENT '组织ID',
  `owner_id` int(10) unsigned NOT NULL COMMENT '组织的拥有者的用户ID',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`organization_id`,`owner_id`),
  KEY `idx_organization_id` (`organization_id`) USING BTREE,
  KEY `idx_owner_id` (`owner_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='组织与拥有者的关系';

-- ----------------------------
-- Table structure for organization_user
-- ----------------------------
DROP TABLE IF EXISTS `organization_user`;
CREATE TABLE `organization_user` (
  `organization_id` int(10) unsigned NOT NULL COMMENT '组织ID',
  `user_id` int(10) unsigned NOT NULL COMMENT '组织的用户ID',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  KEY `idx_organization_id` (`organization_id`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='组织与用户（常客）的关系';

-- ----------------------------
-- Table structure for record
-- ----------------------------
DROP TABLE IF EXISTS `record`;
CREATE TABLE `record` (
  `record_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '测量结果记录ID',
  `client_id` varchar(255) NOT NULL COMMENT '客户端ID',
  `user_id` int(10) unsigned NOT NULL COMMENT '用户ID',
  `device_id` int(10) unsigned DEFAULT NULL COMMENT '设备ID',
  -- 用户选择相关字段
  `finger` tinyint(4) NOT NULL COMMENT '1左小指-5左大姆 6右大姆-10右小姆',
  `app_heart_rate` double(15,8) NOT NULL DEFAULT '0' COMMENT 'App测的心率',
  `is_sport_or_drunk` int(4) DEFAULT '-1' COMMENT '运动或饮酒',
  `cold` int(4) DEFAULT '-1' COMMENT '感冒或病毒感染期',
  `menstrual_cycle` int(4) DEFAULT '-1' COMMENT '生理周期',
  `oviposit_period` int(4) DEFAULT '-1' COMMENT '排卵期',
  `lactation` int(4) DEFAULT '-1' COMMENT '哺乳期',
  `pregnancy` int(4) DEFAULT '-1' COMMENT '怀孕',
  `cm_app_status_a` int(2) DEFAULT '-1' COMMENT '口苦口黏，皮肤瘙痒，大便不成形，头重身痛',
  `cm_app_status_b` int(2) DEFAULT '-1' COMMENT '急躁易怒，头晕胀痛',
  `cm_app_status_c` int(2) DEFAULT '-1' COMMENT '口苦听力下降女性带下异味小便黄短',
  `cm_app_status_d` int(2) DEFAULT '-1' COMMENT '口中异味反酸便秘喉咙干痒牙龈出血',
  `cm_app_status_e` int(2) DEFAULT '-1' COMMENT '胃部冷痛，得温缓解',
  `cm_app_status_f` int(2) DEFAULT '-1' COMMENT '失眠多梦健忘眩晕',

  -- 算法服务返回结果字段
  `c0` double(15,8) COMMENT '心包经测量指标',
  `c1` double(15,8) COMMENT '肝经测量指标',
  `c2` double(15,8) COMMENT '肾经测量指标',
  `c3` double(15,8) COMMENT '脾经测量指标',
  `c4` double(15,8) COMMENT '肺经测量指标',
  `c5` double(15,8) COMMENT '胃经测量指标',
  `c6` double(15,8) COMMENT '胆经测量指标',
  `c7` double(15,8) COMMENT '膀胱经测量指标',
  `c0cv` double(15,8) COMMENT '心包经_变异 ',
  `c1cv` double(15,8) COMMENT '肝经_变异',
  `c2cv` double(15,8) COMMENT '肾经_变异',
  `c3cv` double(15,8) COMMENT '脾经_变异',
  `c4cv` double(15,8) COMMENT '肺经_变异',
  `c5cv` double(15,8) COMMENT '胃经_变异',
  `c6cv` double(15,8) COMMENT '胆经_变异',
  `c7cv` double(15,8) COMMENT '膀胱经_变异',
  `g0` tinyint(4) COMMENT '风_心包经',
  `g1` tinyint(4) COMMENT '风_肝经',
  `g2` tinyint(4) COMMENT '风_肾经',
  `g3` tinyint(4) COMMENT '风_脾经',
  `g4` tinyint(4) COMMENT '风_肺经',
  `g5` tinyint(4) COMMENT '风_胃经',
  `g6` tinyint(4) COMMENT '风_胆经',
  `g7` tinyint(4) COMMENT '风_膀胱经',
  `heart_rate` double(15,8) COMMENT '算法服务器测的心率',
  `heart_rate_cv` float COMMENT '心率变异',
  `health_code` int(10) COMMENT '测试结果健康代码',
  `snr` float COMMENT '信噪比',
  `dc_drift` float COMMENT '直流漂移',
  `elapsed` int(11) COMMENT '算法服务器计算所消耗时间',
  -- 其它字段
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `sys_remark` varchar(1000) NOT NULL DEFAULT '' COMMENT '系统备注',
  `record_type` int(10) DEFAULT '5' COMMENT '5 代表 1.5，6 代表1.6，7 代表 1.7 以此类推',
  `answers` varchar(4000) NOT NULL DEFAULT '' COMMENT '智能分析的回答，JSON格式',
  `is_valid` tinyint(4) NOT NULL DEFAULT '1' COMMENT '是否有效',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`record_id`),
  KEY `idx_user_id` (`user_id`) USING BTREE,
  KEY `idx_created_at` (`created_at`) USING BTREE,
  KEY `idx_is_valid` (`is_valid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='脉诊测量记录';

-- ----------------------------
-- Table structure for subscription
-- ----------------------------
DROP TABLE IF EXISTS `subscription`;
CREATE TABLE `subscription` (
  `subscription_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '订阅ID',
  `organization_id` int(10) unsigned DEFAULT NULL COMMENT '组织ID',
  `subscription_type` tinyint(4) DEFAULT '0' COMMENT '0 定制化 1 试用版 2 黄喜马把脉 3 白喜马把脉 4 钻石姆 5 礼品版',
  `active` tinyint(4) DEFAULT '0' COMMENT '是否激活 0未激活 1已激活',
  `activated_at` DATETIME NULL DEFAULT NULL COMMENT '激活时间',
  `expired_at` DATETIME NULL DEFAULT NULL COMMENT '到期时间',
  `customized_code` varchar(255) DEFAULT '' COMMENT '自定义代码',
  `max_user_limits` int(11) DEFAULT NULL COMMENT '最大用户（常客）数量',
  `contract_year` smallint(6) DEFAULT NULL COMMENT '订阅合同年限',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`subscription_id`),
  KEY `idx_organization_id` (`organization_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='产品服务订阅';


-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `user_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `legacy_subject_id` int(20) unsigned NULL COMMENT 'Legacy 迁移前 subject_id，即老的用户ID',
  `username` varchar(255) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '用户登录密码',
  `zone` varchar(255) NOT NULL COMMENT '用户选择的区域',
  `register_type` varchar(20) NOT NULL COMMENT '用户注册方式',
  `register_time` timestamp NOT NULL COMMENT '注册时间',
  `register_source_client_id` varchar(255) NOT NULL COMMENT '发起注册的客户端ID',
  `customized_code` varchar(255) NOT NULL DEFAULT '' COMMENT '定制化代码',
  `nickname` varchar(255) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `gender` enum('M','F') NOT NULL DEFAULT 'M' COMMENT '用户性别 M 男 F 女',
  `birthday` date NOT NULL COMMENT '用户生日',
  `height` smallint(6) NOT NULL DEFAULT '0' COMMENT '用户身高 单位厘米',
  `weight` smallint(6) NOT NULL DEFAULT '0' COMMENT '用户体重 单位千克',
  `phone` varchar(50) NOT NULL DEFAULT '' COMMENT '用户电话号码',
  `email` varchar(255) NOT NULL DEFAULT '' COMMENT '用户联系邮箱',
  `country` varchar(255) DEFAULT NULL COMMENT '国家',
  `state` varchar(255) NOT NULL DEFAULT '' COMMENT '用户所在省份',
  `city` varchar(255) NOT NULL DEFAULT '' COMMENT '用户所在城市',
  `district` varchar(255) DEFAULT NULL COMMENT '地区',
  `street` varchar(255) NOT NULL DEFAULT '' COMMENT '用户所在街道',
  `user_defined_code` varchar(255) DEFAULT NULL COMMENT '用户自定义代码',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户备注',
  `is_profile_completed` tinyint(4) DEFAULT '0' COMMENT '是否开启心率扇形图',
  `is_activated` tinyint(4) NOT NULL DEFAULT '1' COMMENT '是否激活',
  `activated_at` timestamp NULL DEFAULT NULL COMMENT '激活时间',
  `deactivated_at` timestamp NULL DEFAULT NULL COMMENT '取消激活（禁用）时间',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`user_id`),
  KEY `idx_register_type` (`register_type`) USING BTREE,
  KEY `idx_username` (`username`) USING BTREE,
  KEY `idx_email` (`email`) USING BTREE,
  KEY `idx_phone` (`phone`) USING BTREE,
  KEY `idx_legacy_subject_id` (`legacy_subject_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户';

-- ----------------------------
-- Table structure for user_access_token
-- ----------------------------
DROP TABLE IF EXISTS `user_access_token`;
CREATE TABLE `user_access_token` (
  `user_id` int(10) unsigned NOT NULL COMMENT '用户ID',
  `token` varchar(36) NOT NULL COMMENT '当前会话的令牌',
  `expired_at` datetime NOT NULL COMMENT '会话过期时间',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  PRIMARY KEY (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户访问令牌';

-- ----------------------------
-- Table structure for user_preferences
-- ----------------------------
DROP TABLE IF EXISTS `user_preferences`;
CREATE TABLE `user_preferences` (
  `user_id` int(10) unsigned NOT NULL COMMENT '用户ID',
  `enable_heart_rate_chart` tinyint(4) DEFAULT '1' COMMENT '是否开启心率扇形图',
  `enable_pulse_wave_chart` tinyint(4) DEFAULT '1' COMMENT '是否开启波形图',
  `enable_warm_prompt` tinyint(4) DEFAULT '1' COMMENT '是否开启温馨提示',
  `enable_choose_status` tinyint(4) DEFAULT '1' COMMENT '是否开启选择状态',
  `enable_constitution_differentiation` tinyint(4) DEFAULT '1' COMMENT '是否开启中医体质判读',
  `enable_syndrome_differentiation` tinyint(4) DEFAULT '1' COMMENT '是否开启中医脏腑判读',
  `enable_western_medicine_analysis` tinyint(4) DEFAULT '1' COMMENT '是否开启西医判读',
  `enable_meridian_bar_graph` tinyint(4) DEFAULT '1' COMMENT '是否开启柱状图',
  `enable_comment` tinyint(4) DEFAULT '1' COMMENT '是否开启备注',
  `enable_health_trending` tinyint(4) DEFAULT '1' COMMENT '是否健康趋势入口',
  `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
  `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户首选项信息';

SET FOREIGN_KEY_CHECKS = 1;

-- ----------------------------
-- Table structure for transaction_number
-- ----------------------------
CREATE TABLE `transaction_number` (
 `transaction_date` date NOT NULL COMMENT '日期',
 `transaction_number` int(32) unsigned NOT NULL COMMENT '流水号',
 `created_at` timestamp NOT NULL COMMENT '数据记录创建时间',
 `updated_at` timestamp NOT NULL COMMENT '数据记录更新时间',
 PRIMARY KEY (`transaction_date`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='流水号';
