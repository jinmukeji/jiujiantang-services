/*
 Navicat Premium Data Transfer

 Source Server         : jinmu-test
 Source Server Type    : MySQL
 Source Server Version : 50717
 Source Host           : jinmu-test.cjzrjn31gtsw.rds.cn-north-1.amazonaws.com.cn
 Source Database       : jinmudev

 Target Server Type    : MySQL
 Target Server Version : 50717
 File Encoding         : utf-8

 Date: 01/24/2018 13:51:03 PM
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `jinmu_company`
-- ----------------------------
DROP TABLE IF EXISTS `jinmu_company`;
CREATE TABLE `jinmu_company` (
  `company_id` int(20) unsigned NOT NULL AUTO_INCREMENT,
  `company_code` varchar(50) NOT NULL COMMENT '公司统编',
  `account` varchar(50) NOT NULL COMMENT '产品账号',
  `company_name` varchar(100) NOT NULL COMMENT '公司名称',
  `company_address` varchar(100) DEFAULT NULL COMMENT '公司地址',
  `company_phone` varchar(50) DEFAULT NULL COMMENT '公司移动电话',
  `company_line` varchar(50) DEFAULT NULL COMMENT '公司固定电话',
  `company_email` varchar(50) DEFAULT NULL COMMENT '公司联系邮箱',
  `company_state` varchar(50) DEFAULT NULL COMMENT '公司所在省份',
  `company_city` varchar(50) DEFAULT NULL COMMENT '公司所在城市',
  `company_block` varchar(50) DEFAULT NULL COMMENT '公司所在区域',
  `legal_person` varchar(50) DEFAULT NULL COMMENT '公司法人',
  `company_representative` varchar(50) DEFAULT NULL COMMENT '公司联系人',
  `company_bank_account` varchar(50) DEFAULT NULL COMMENT '公司银行账号',
  `company_type` varchar(50) DEFAULT NULL COMMENT '公司类型',
  `create_date` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `update_date` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
  `end_date` datetime DEFAULT NULL COMMENT '结束日期',
  `is_valid` int(2) DEFAULT NULL COMMENT '是否有效',
  PRIMARY KEY (`company_id`),
  UNIQUE KEY `company_code` (`company_code`),
  KEY `account` (`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
--  Table structure for `jinmu_contract`
-- ----------------------------
DROP TABLE IF EXISTS `jinmu_contract`;
CREATE TABLE `jinmu_contract` (
  `contract_id` int(20) unsigned NOT NULL AUTO_INCREMENT,
  `contract_code` varchar(100) NOT NULL COMMENT '合同编号',
  `account` varchar(50) NOT NULL COMMENT '产品账号',
  `amount` int(20) DEFAULT NULL COMMENT '合同金额',
  `monetary_unit` varchar(10) DEFAULT NULL COMMENT '合同金额单位',
  `max_client` int(10) NOT NULL COMMENT '最大客户数',
  `contract_content` varchar(300) DEFAULT NULL COMMENT '合同内容',
  `contract_sign_date` datetime DEFAULT NULL COMMENT '合同签到日期',
  `contract_type` varchar(4) DEFAULT NULL COMMENT '合同类型T|F',
  `contract_start_date` datetime DEFAULT NULL COMMENT '合同开始日期',
  `contract_end_date` datetime DEFAULT NULL COMMENT '合同结束日期',
  `contract_year` int(10) DEFAULT NULL COMMENT '产品允许使用的年限',
  `total_times` int(10) DEFAULT NULL COMMENT '总共可用次数',
  `remaining_times` int(10) DEFAULT NULL COMMENT '剩余可用次数',
  `buyer_code` varchar(30) NOT NULL COMMENT '购买产品公司',
  `seller_code` varchar(30) NOT NULL COMMENT '销售产品公司',
  `create_date` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `update_date` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
  `end_date` datetime DEFAULT NULL COMMENT '结束日期',
  `is_valid` int(2) DEFAULT NULL COMMENT '是否有效',
  PRIMARY KEY (`contract_id`),
  UNIQUE KEY `contract_code` (`contract_code`),
  KEY `account` (`account`),
  KEY `is_valid` (`is_valid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
--  Table structure for `jinmu_feedback`
-- ----------------------------
DROP TABLE IF EXISTS `jinmu_feedback`;
CREATE TABLE `jinmu_feedback` (
  `feedback_id` int(20) unsigned NOT NULL AUTO_INCREMENT,
  `account` varchar(50) NOT NULL COMMENT '产品账号',
  `contact_way` varchar(150) DEFAULT NULL COMMENT '联系方式',
  `content` text COMMENT '意见内容',
  `create_date` datetime DEFAULT NULL COMMENT '创建日期',
  `end_date` datetime DEFAULT NULL COMMENT '结束日期',
  `is_valid` int(2) DEFAULT NULL COMMENT '是否有效',
  PRIMARY KEY (`feedback_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
--  Table structure for `jinmu_log`
-- ----------------------------
DROP TABLE IF EXISTS `jinmu_log`;
CREATE TABLE `jinmu_log` (
  `log_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'LOG記錄ID',
  `log_key` varchar(64) DEFAULT NULL COMMENT 'Log唯一key',
  `request_uri` varchar(128) DEFAULT NULL COMMENT '請求的UTL',
  `request_time` timestamp NULL DEFAULT NULL COMMENT '網絡請求時間',
  `server_addr` varchar(32) DEFAULT NULL COMMENT '服務器IP地址',
  `remote_addr` varchar(32) DEFAULT NULL COMMENT '客戶端IP地址',
  `request_method` varchar(32) DEFAULT NULL COMMENT '網絡請求方法',
  `info_lang` varchar(32) DEFAULT NULL COMMENT 'LANG資訊',
  `info_route` varchar(64) DEFAULT NULL COMMENT '伺服器內路由資訊',
  `info_header` int(11) DEFAULT NULL COMMENT '網絡請求Header資訊',
  `request_param` longtext COMMENT '請求的參數',
  `response_status` int(11) DEFAULT NULL COMMENT '返回内容',
  `response_description` varchar(64) DEFAULT NULL,
  `response_data` longtext,
  `run_time` double(15,6) DEFAULT NULL COMMENT '運行時間：S',
  `throughput_rate` double(15,6) DEFAULT NULL COMMENT '吞吐率：req/s',
  `memory_use` varchar(32) DEFAULT NULL COMMENT '內存使用情況',
  `file_load` int(11) DEFAULT NULL COMMENT '文件加載',
  `error_msg` text COMMENT '错误资讯',
  `is_valid` int(11) DEFAULT '1' COMMENT '記錄是否失效',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '創建時間',
  `update_time` timestamp NULL DEFAULT NULL COMMENT '更新時間',
  PRIMARY KEY (`log_id`),
  KEY `log_key` (`log_key`),
  KEY `response_status` (`response_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
--  Table structure for `jinmu_mac`
-- ----------------------------
DROP TABLE IF EXISTS `jinmu_mac`;
CREATE TABLE `jinmu_mac` (
  `mac_id` int(20) unsigned NOT NULL AUTO_INCREMENT,
  `mac` varchar(100) DEFAULT NULL COMMENT '仪器MAC',
  `account` int(20) DEFAULT NULL,
  PRIMARY KEY (`mac_id`),
  UNIQUE KEY `mac` (`mac`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
--  Table structure for `jinmu_product`
-- ----------------------------
DROP TABLE IF EXISTS `jinmu_product`;
CREATE TABLE `jinmu_product` (
  `product_id` int(20) unsigned NOT NULL AUTO_INCREMENT,
  `account` varchar(50) NOT NULL COMMENT '产品账号',
  `password` varchar(150) DEFAULT NULL COMMENT '账号密码',
  `pad_app_version_name` varchar(100) DEFAULT NULL COMMENT 'APP 版本名称',
  `pad_app_version_code` varchar(10) DEFAULT NULL COMMENT 'APP 版本代码',
  `pad_app_description` text COMMENT 'APP 版本说明',
  `pad_app_downloadurl` varchar(320) DEFAULT NULL COMMENT 'APP 下载地址',
  `pad_mac_address` varchar(100) DEFAULT NULL COMMENT '平板MAC',
  `pad_device_code` varchar(100) DEFAULT NULL COMMENT '平板编码',
  `appratus_model` varchar(100) DEFAULT NULL COMMENT '仪器型号',
  `appratus_mac_address` varchar(100) DEFAULT NULL COMMENT '仪器MAC',
  `appratus_device_code` varchar(100) DEFAULT NULL COMMENT '仪器编码',
  `appratus_manu_date` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '仪器生产日期',
  `appratus_end_date` datetime DEFAULT NULL COMMENT '仪器预估寿命',
  `create_date` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `update_date` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
  `end_date` datetime DEFAULT NULL COMMENT '结束日期',
  `is_valid` int(2) DEFAULT NULL COMMENT '是否有效',
  PRIMARY KEY (`product_id`),
  UNIQUE KEY `account` (`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
--  Table structure for `jinmu_record`
-- ----------------------------
DROP TABLE IF EXISTS `jinmu_record`;
CREATE TABLE `jinmu_record` (
  `record_id` int(20) unsigned NOT NULL AUTO_INCREMENT,
  `subject_id` int(20) NOT NULL COMMENT '测试者ID',
  `state_id` int(20) NOT NULL COMMENT '测试者生理状态ID',
  `appratus_mac` varchar(100) DEFAULT NULL COMMENT '仪器MAC',
  `meat_status` int(4) DEFAULT NULL COMMENT '空腹状态',
  `drink_status` int(4) DEFAULT NULL COMMENT '饮酒状态',
  `flu_status` int(4) DEFAULT NULL COMMENT '感冒状态',
  `period_status` int(4) DEFAULT NULL COMMENT '生理状态',
  `other_status` varchar(20) DEFAULT NULL COMMENT '其他状态',
  `comment` text COMMENT '测量备注',
  `C0` double(15,8) DEFAULT NULL COMMENT '心测量指标',
  `C1` double(15,8) DEFAULT NULL COMMENT '肝测量指标',
  `C2` double(15,8) DEFAULT NULL COMMENT '肾测量指标',
  `C3` double(15,8) DEFAULT NULL COMMENT '脾测量指标',
  `C4` double(15,8) DEFAULT NULL COMMENT '肺测量指标',
  `C5` double(15,8) DEFAULT NULL COMMENT '胃测量指标',
  `C6` double(15,8) DEFAULT NULL COMMENT '胆测量指标',
  `C7` double(15,8) DEFAULT NULL COMMENT '膀胱测量指标',
  `C8` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C9` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C10` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C11` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C0CV` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C1CV` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C2CV` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C3CV` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C4CV` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C5CV` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C6CV` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C7CV` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C8CV` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C9CV` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C10CV` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `C11CV` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P0` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P1` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P2` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P3` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P4` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P5` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P6` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P7` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P8` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P9` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P10` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P11` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P0SD` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P1SD` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P2SD` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P3SD` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P4SD` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P5SD` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P6SD` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P7SD` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P8SD` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P9SD` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P10SD` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `P11SD` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `HR` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `HRCV` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `MBP` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `DBP` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `SBP` double(15,8) DEFAULT NULL COMMENT '测量指标',
  `PP` double(15,8) DEFAULT NULL COMMENT '腰围',
  `hand_opt` int(15) DEFAULT NULL COMMENT '测量指标',
  `G0` tinyint(4) DEFAULT NULL COMMENT '测量指标',
  `G1` tinyint(4) DEFAULT NULL COMMENT '测量指标',
  `G2` tinyint(4) DEFAULT NULL COMMENT '测量指标',
  `G3` tinyint(4) DEFAULT NULL COMMENT '测量指标',
  `G4` tinyint(4) DEFAULT NULL COMMENT '测量指标',
  `G5` tinyint(4) DEFAULT NULL COMMENT '测量指标',
  `G6` tinyint(4) DEFAULT NULL COMMENT '测量指标',
  `G7` tinyint(4) DEFAULT NULL COMMENT '测量指标',
  `health_code` int(10) NOT NULL COMMENT '测试结果健康代码',
  `inspect_mode` int(20) NOT NULL COMMENT '后台运算服务器',
  `bits` int(4) DEFAULT NULL COMMENT '后台运算服务器',
  `c_code` varchar(30) DEFAULT NULL COMMENT '其他',
  `condition_flag` varchar(200) DEFAULT NULL COMMENT 'condition_flag',
  `create_date` datetime DEFAULT NULL COMMENT '创建日期',
  `update_date` datetime DEFAULT NULL COMMENT '更新备注日期',
  `mobile_type` char(20) DEFAULT NULL COMMENT '手机版本',
  `snr` float DEFAULT NULL,
  `dc_drift` float DEFAULT NULL,
  `status` text,
  `is_sport_or_drunk` int(4) DEFAULT '-1',
  `cold` int(4) DEFAULT '-1',
  `menstrual_cycle` int(4) DEFAULT '-1',
  `oviposit_period` int(4) DEFAULT '-1',
  `lactation` int(4) DEFAULT '-1',
  `pregnancy` int(4) DEFAULT '-1',
  `finger` int(4) DEFAULT NULL,
  `is_valid` int(4) DEFAULT '0' COMMENT '是否有效',
  `APP_HR` int(10) DEFAULT '0',
  `record_type` int(10) DEFAULT '5' COMMENT '5 代表 1.5，6 代表1.6，7 代表 1.7 以此类推',
  `cm_app_status_a` int(2) DEFAULT '-1' COMMENT '口苦口黏，皮肤瘙痒，大便不成形，头重身痛',
  `cm_app_status_b` int(2) DEFAULT '-1' COMMENT '急躁易怒，头晕胀痛',
  `cm_app_status_c` int(2) DEFAULT '-1' COMMENT '口苦听力下降女性带下异味小便黄短',
  `cm_app_status_d` int(2) DEFAULT '-1' COMMENT '口中异味反酸便秘喉咙干痒牙龈出血',
  `cm_app_status_e` int(2) DEFAULT '-1' COMMENT '胃部冷痛，得温缓解',
  `cm_app_status_f` int(2) DEFAULT '-1' COMMENT '失眠多梦健忘眩晕',
  PRIMARY KEY (`record_id`),
  KEY `subject_id` (`subject_id`),
  KEY `create_date` (`create_date`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;

-- ----------------------------
--  Table structure for `jinmu_status`
-- ----------------------------
DROP TABLE IF EXISTS `jinmu_status`;
CREATE TABLE `jinmu_status` (
  `status_id` int(20) unsigned NOT NULL AUTO_INCREMENT,
  `subject_id` int(20) NOT NULL COMMENT '测试者ID',
  `height` int(10) DEFAULT NULL COMMENT '身高',
  `weight` int(10) DEFAULT NULL COMMENT '体重',
  `waistline` int(10) DEFAULT NULL COMMENT '腰围',
  `create_date` datetime DEFAULT NULL COMMENT '创建日期',
  `end_date` datetime DEFAULT NULL COMMENT '结束日期',
  `is_valid` int(2) DEFAULT NULL COMMENT '是否有效',
  PRIMARY KEY (`status_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;

-- ----------------------------
--  Table structure for `jinmu_subject`
-- ----------------------------
DROP TABLE IF EXISTS `jinmu_subject`;
CREATE TABLE `jinmu_subject` (
  `subject_id` int(20) unsigned NOT NULL AUTO_INCREMENT,
  `account` varchar(50) NOT NULL COMMENT '产品账号',
  `pin` varchar(20) DEFAULT NULL COMMENT '身份证号码',
  `name` varchar(256) DEFAULT NULL COMMENT '姓名',
  `nickname` varchar(150) DEFAULT NULL COMMENT '账号昵称',
  `birthdate` date DEFAULT NULL COMMENT '生日',
  `gender` varchar(5) DEFAULT NULL COMMENT '性别M|F',
  `phone` varchar(50) DEFAULT NULL COMMENT '电话',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `nationality` varchar(50) DEFAULT NULL COMMENT '国籍或地区',
  `state` varchar(50) DEFAULT NULL COMMENT '所在省份',
  `city` varchar(50) DEFAULT NULL COMMENT '所在城市',
  `block` varchar(50) DEFAULT NULL COMMENT '所在区域',
  `create_date` datetime DEFAULT NULL COMMENT '创建日期',
  `update_date` datetime DEFAULT NULL COMMENT '更新日期',
  `end_date` datetime DEFAULT NULL COMMENT '结束日期',
  `is_valid` int(2) DEFAULT NULL COMMENT '是否有效',
  `subject_type` int(2) DEFAULT '1',
  PRIMARY KEY (`subject_id`),
  KEY `account` (`account`),
  KEY `is_valid` (`is_valid`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;

-- ----------------------------
--  Table structure for `jinmu_token`
-- ----------------------------
DROP TABLE IF EXISTS `jinmu_token`;
CREATE TABLE `jinmu_token` (
  `account_id` varchar(50) NOT NULL,
  `created_at` datetime NOT NULL,
  `expired_at` datetime NOT NULL,
  `token` varchar(36) NOT NULL,
  PRIMARY KEY (`token`),
  UNIQUE KEY `token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;

--
-- 转utc 时间
-- 加字段
ALTER TABLE jinmu_record ADD create_date_utc datetime DEFAULT NOW()
COMMENT 'utc时区的测量时间';

-- 批量转
update jinmu_record as t
set t.create_date_utc = CONVERT_TZ
(t.create_date,'+08:00','+00:00');

-- 插入测试帐号
insert into `jinmu_contract` ( `contract_code`, `account`, `amount`, `monetary_unit`, `max_client`, `contract_content`, `contract_sign_date`, `contract_type`, `contract_start_date`, `contract_end_date`, `contract_year`, `total_times`, `remaining_times`, `buyer_code`, `seller_code`, `create_date`, `update_date`, `end_date`, `is_valid`) values ( '1_JMZ26C-1026', '1', '-1', 'RMB', '300', 'JM100150', null, 'T', '2018-01-23 15:01:05', '2021-01-23 15:01:05', '3', null, null, '1', '91320400MA1MP76Q53', '2017-05-08 11:17:38', '2018-01-23 15:01:18', null, '1');

insert into `jinmu_company` ( `company_code`, `account`, `company_name`, `company_address`, `company_phone`, `company_line`, `company_email`, `company_state`, `company_city`, `company_block`, `legal_person`, `company_representative`, `company_bank_account`, `company_type`, `create_date`, `update_date`, `end_date`, `is_valid`) values ( '1', '1', '李连杰', '江苏常州天宁区北塘河路恒生科技园1栋503', '13530410518', 'null', 'zhuyingjie@jinmuhealth.com', '广东省', '广州市', '番禺区', '', '', null, '养生', '2017-05-08 11:17:38', '2018-01-23 15:01:17', null, '1');

insert into `jinmu_product` ( `account`, `password`, `pad_app_version_name`, `pad_app_version_code`, `pad_app_description`, `pad_app_downloadurl`, `pad_mac_address`, `pad_device_code`, `appratus_model`, `appratus_mac_address`, `appratus_device_code`, `appratus_manu_date`, `appratus_end_date`, `create_date`, `update_date`, `end_date`, `is_valid`) values ( '1', 'release1', '金姆健康大陆版1.3.0', '10', '新增与改进\r\n1.使用记录左右滑动\r\n2.备注按钮的修改\r\n3.测量结果界面增加切换示意图按钮\r\n4.完善了用户搜索弹窗的模糊搜索\r\n5.增加了解绑设置的页面\r\n6.智能分析解析更新\r\n7.未选择用户，提示选择用户\r\n8.新增，设置，解绑，返回主页面时，拉出侧滑菜单\r\n9.新增点击系统或者第三方应用，在测量过程中，中断蓝牙的提示\r\n10.App导航页的图片更换\r\n\r\n修复以下问题\r\n1.企业信息中区域显示错误（null）\r\n2.个人信息中邮箱输入中文（中文+@qq.com），点击保存报错\r\n3.历史记录中修改查询日期后，查询的数据有误\r\n4.新浪微博分享的修复\r\n', 'http://jinmu.oss-cn-shanghai.aliyuncs.com/com.jinmu.healthdlb_1.3.0.apk', null, null, 'XMW23', '38D269ED8184', null, '2017-05-02 00:00:00', '2022-05-01 00:00:00', '2017-05-08 11:17:38', '2018-01-23 15:01:17', null, '1');

insert into `jinmu_subject` ( `account`, `subject_id`,`pin`, `name`, `nickname`, `birthdate`, `gender`, `phone`, `email`, `nationality`, `state`, `city`, `block`, `create_date`, `update_date`, `end_date`, `is_valid`) values 
( '1', '7824' ,'', '贾跃亭', '', '2018-01-23', 'M', '025-111111', 'jiayueting@jinmuhealth.com', '', '江苏', '常州', '天宁区', '2018-01-23 15:01:08', '2018-01-23 15:01:08', '2050-01-01 00:00:00', '1');

insert into `jinmu_status` ( `subject_id`, `height`, `weight`, `waistline`, `create_date`, `end_date`, `is_valid`) values ( '7824', '195', '78', '0', '2018-01-23 15:01:08', '2018-01-30 15:01:08', '1');

-- 用于测试空时间
insert into `jinmu_contract` ( `contract_code`, `account`, `amount`, `monetary_unit`, `max_client`, `contract_content`, `contract_sign_date`, `contract_type`, `contract_start_date`, `contract_end_date`, `contract_year`, `total_times`, `remaining_times`, `buyer_code`, `seller_code`, `create_date`, `update_date`, `end_date`, `is_valid`) values ( '14_JMZ26C-1026', '14', '-1', 'RMB', '300', 'JM100150', null, 'T', null, null, '100', null, null, '14', '91320400MA1MP76Q53', '2017-05-08 14:54:19', '2017-05-08 14:54:19', null, '1');

insert into `jinmu_record` (`record_id`, `subject_id`, `state_id`, `appratus_mac`, `meat_status`, `drink_status`, `flu_status`, `period_status`, `other_status`, `comment`, `C0`, `C1`, `C2`, `C3`, `C4`, `C5`, `C6`, `C7`, `C8`, `C9`, `C10`, `C11`, `C0CV`, `C1CV`, `C2CV`, `C3CV`, `C4CV`, `C5CV`, `C6CV`, `C7CV`, `C8CV`, `C9CV`, `C10CV`, `C11CV`, `P0`, `P1`, `P2`, `P3`, `P4`, `P5`, `P6`, `P7`, `P8`, `P9`, `P10`, `P11`, `P0SD`, `P1SD`, `P2SD`, `P3SD`, `P4SD`, `P5SD`, `P6SD`, `P7SD`, `P8SD`, `P9SD`, `P10SD`, `P11SD`, `HR`, `HRCV`, `MBP`, `DBP`, `SBP`, `PP`, `hand_opt`, `G0`, `G1`, `G2`, `G3`, `G4`, `G5`, `G6`, `G7`, `health_code`, `inspect_mode`, `bits`, `c_code`, `condition_flag`, `create_date`, `update_date`, `mobile_type`, `snr`, `dc_drift`, `status`, `is_sport_or_drunk`, `cold`, `menstrual_cycle`, `oviposit_period`, `lactation`, `pregnancy`, `finger`, `is_valid`) values ( '1','1', '1', '38:D2:69:ED:74:0D', '-1', '0', '0', '0', null, null, '-2.00000000', '0.00000000', '10.00000000', '-5.00000000', '0.00000000', '-3.00000000', '-1.00000000', '0.00000000', null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, '100.00000000', '4.60000000', null, null, null, null, null, '0', '0', '0', '1', '0', '0', '0', '0', '8', '0', '22', 'CL', '//8=', '2017-05-08 12:31:55', null, 'ANDROID', null, null, null, '-1', '-1', '-1', '-1', '-1', '-1', '-1', '1');

insert into `jinmu_subject` ( `account`, `subject_id`,`pin`, `name`, `nickname`, `birthdate`, `gender`, `phone`, `email`, `nationality`, `state`, `city`, `block`, `create_date`, `update_date`, `end_date`, `is_valid`) values 
( '1', '1' ,'', '刘强东', '', '2018-01-23', 'M', '025-111111', 'liuqiangdong@jinmuhealth.com', '', '江苏', '常州', '天宁区', '2018-01-23 15:01:08', '2018-01-23 15:01:08', '2050-01-01 00:00:00', '1');

insert into `jinmu_status` ( `subject_id`, `height`, `weight`, `waistline`, `create_date`, `end_date`, `is_valid`) values ( 1 , '195', '78', '0', '2018-01-23 15:01:08', '2018-01-30 15:01:08', '1');


insert into `jinmu_mac` ( `mac`, `account`) values ( '39D269E877A1', '1');

-- 测第一次登录

insert into `jinmu_contract` ( `contract_code`,
`account`, `amount`, `monetary_unit`, `max_client`, `contract_content`, `contract_sign_date`, `contract_type`, `contract_start_date`, `contract_end_date`, `contract_year`, `total_times`, `remaining_times`, `buyer_code`, `seller_code`, `create_date`, `update_date`, `end_date`, `is_valid`) values
( '1_JMZ26C-1027', '11', '-1', 'RMB', '300', 'JM100150', null, 'T', null , null, '3', null, null, '1', '91320400MA1MP76Q53', '2017-05-08 11:17:38', '2018-01-23 15:01:18', null, '1');

insert into `jinmu_product` ( `account`,
`password`, `pad_app_version_name`, `pad_app_version_code`, `pad_app_description`, `pad_app_downloadurl`, `pad_mac_address`, `pad_device_code`, `appratus_model`, `appratus_mac_address`, `appratus_device_code`, `appratus_manu_date`, `appratus_end_date`, `create_date`, `update_date`, `end_date`, `is_valid`) values
( '11', 'release1', '金姆健康大陆版1.3.0', '10', '新增与改进\r\n1.使用记录左右滑动\r\n2.备注按钮的修改\r\n3.测量结果界面增加切换示意图按钮\r\n4.完善了用户搜索弹窗的模糊搜索\r\n5.增加了解绑设置的页面\r\n6.智能分析解析更新\r\n7.未选择用户，提示选择用户\r\n8.新增，设置，解绑，返回主页面时，拉出侧滑菜单\r\n9.新增点击系统或者第三方应用，在测量过程中，中断蓝牙的提示\r\n10.App导航页的图片更换\r\n\r\n修复以下问题\r\n1.企业信息中区域显示错误（null）\r\n2.个人信息中邮箱输入中文（中文+@qq.com），点击保存报错\r\n3.历史记录中修改查询日期后，查询的数据有误\r\n4.新浪微博分享的修复\r\n', 'http://jinmu.oss-cn-shanghai.aliyuncs.com/com.jinmu.healthdlb_1.3.0.apk', null, null, 'XMW23', '38D269ED8184', null, '2017-05-02 00:00:00', '2022-05-01 00:00:00', '2017-05-08 11:17:38', '2018-01-23 15:01:17', null, '1');

-- 给subject表name字段添加普通索引(字段可重复)
ALTER TABLE jinmu_subject ADD INDEX subject_name_index(name);

-- 给product表添加account_type字段
ALTER TABLE jinmu_product ADD account_type varchar(256) DEFAULT 'normal' COMMENT '区分项目与普通账号',

-- 给product表account_type添加索引
ALTER TABLE jinmu_product ADD INDEX account_type_index(account_type);

-- 添加测量过程开关需要建表
CREATE TABLE `jinmu_account_preferences`
(
    `account` varchar(50) NOT NULL COMMENT '产品账号',
    `enable_heart_rate_chart` int(2) DEFAULT 1 COMMENT '是否开启心率扇形图',
    `enable_pulse_wave_chart` int(2) DEFAULT 1 COMMENT '是否开启波形图',
    `enable_warm_prompt` int(2) DEFAULT 1 COMMENT '是否开启温馨提示',
    `enable_choose_status` int(2) DEFAULT 1 COMMENT '是否开启选择状态',
    `enable_constitution_differentiation` int(2) DEFAULT 1 COMMENT '是否开启中医体质判读',
    `enable_syndrome_differentiation` int(2) DEFAULT 1 COMMENT '是否开启中医脏腑判读',
    `enable_western_medicine_analysis` int(2) DEFAULT 1 COMMENT '是否开启西医判读',
    `enable_meridian_bar_graph` int(2) DEFAULT 1 COMMENT '是否开启柱状图',
    `enable_comment` int(2) DEFAULT 1 COMMENT '是否开启备注',
    PRIMARY KEY `account` (`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 给jinmu_account_preferences插入account=1做测试
INSERT INTO jinmu_account_preferences (account) VALUES (1);
