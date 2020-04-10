# 数据库迁移流程(生产环境)

## 1 停服

1. 停止依赖数据库的服务，特别是写操作相关的。

## 2 备份

### 2.1 dump 生产环境数据库

执行Shell命令：

```sh
# mysqldump Ver 8.0.11
# Dump Production data
mysqldump \
--host=jinmu-prod.cjzrjn31gtsw.rds.cn-north-1.amazonaws.com.cn \
--port=53306 \
--user=dbadmin \
--password \
--compress \
--databases jinmu \
--set-charset \
--default-character-set=utf8mb4 \
--opt \
--no-create-db \
--column-statistics=0 \
--skip-lock-tables \
> jinmu_production_data.sql
```

**注意事项：**

- 如果`mysqldump`是用的5.x版本的，则不需要 `--column-statistics=0` 参数

###2.2 建立数据库快照

在AWS控制台上操作:

- https://console.amazonaws.cn/rds/home?region=cn-north-1#dbinstance:id=jinmu-prod

## 3 迁移数据库

### 3.1 连接到`jinmu`数据库

连接生产环境数据库：

```sh
mysql \
--host jinmu-prod.cjzrjn31gtsw.rds.cn-north-1.amazonaws.com.cn \
--port 53306 \
--database jinmu \
--default-character-set=utf8mb4 \
--user dbadmin \
--password
```

### 3.2 变更原数据库中错误的数据库引擎

> 注意使用 `USE` 语句目标数据库

```sql
ALTER TABLE jinmu_company ENGINE = InnoDB;
ALTER TABLE jinmu_contract ENGINE = InnoDB;
ALTER TABLE jinmu_feedback ENGINE = InnoDB;
ALTER TABLE jinmu_product ENGINE = InnoDB;
ALTER TABLE jinmu_record ENGINE = InnoDB;
ALTER TABLE jinmu_status ENGINE = InnoDB;
ALTER TABLE jinmu_subject ENGINE = InnoDB;

CREATE INDEX `idx_subject_id` ON `jinmu_status` (`subject_id`) USING BTREE;

ALTER TABLE `jinmu_record` 
MODIFY COLUMN `comment` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL AFTER `other_status`;
ALTER TABLE `jinmu_token` COLLATE = utf8mb4_general_ci;

ALTER TABLE `jinmu_token` 
MODIFY COLUMN `account_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL FIRST,
MODIFY COLUMN `token` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL AFTER `expired_at`;
```

### 3.3 变更数据结构与迁移数据

#### 3.3.1 建立新的数据表结构 

1. 运行 **[jinmu-db-v2-migration-000-init.sql](jinmu-db-v2-migration-000-init.sql)** 脚本文件内容。

#### 3.3.2 初始化`client`表 

1. 运行 **[jinmu-db-v2-migration-001.sql](jinmu-db-v2-migration-001.sql)** 脚本文件内容。

#### 3.3.2 初始化`device`表数据 

1. 基于工厂生产信息初始化`device`数据：

   运行 **[jinmu-db-v2-migration-002.sql](jinmu-db-v2-migration-002.sql)** 脚本文件内容。

2. 验证是否有遗漏MAC地址没有进行初始化：

   ```sql
   -- 查询数据库记录中存在MAC地址，但不存在SN号的记录
   SELECT
   	T.m AS mac,
   	CONV( T.m, 16, 10 ) AS mac_number,
   	CONCAT('JMCNUK', T.m) AS unknown_sn
   FROM
   	( SELECT DISTINCT mac AS m FROM jinmu_mac UNION SELECT DISTINCT appratus_mac_address AS m FROM jinmu_product ) AS T
   	LEFT JOIN device AS D ON CONV( T.m, 16, 10 ) = D.mac 
   WHERE
   	D.mac IS NULL
   	AND T.m <> ''
   	AND T.m IS NOT NULL;
   ```

3. 如果上述第2步存在返回结果，则执行以下语句将修正后的记录保存到`device`数据表中：

   ```sql
   -- 清理脏数据
   -- 将只有MAC地址却没有SN号的设备信息初始化插入到device数据表之中
   
   SET @now = NOW();
   INSERT INTO `device` ( `mac`, `sn`, `zone`, `model`, `tags`, `remarks`, `created_at`, `updated_at` ) SELECT
   mac_number,
   unknown_sn,
   'CN',
   'XMW23',
   'missing_sn',
   '历史记录中缺少SN信息的产品',
   @now,
   @now 
   FROM
   	(
   	SELECT
   		T.m AS mac,
   		CONV( T.m, 16, 10 ) AS mac_number,
   		CONCAT( 'JMCNUK', T.m ) AS UNKNOWN_SN 
   	FROM
   		( SELECT DISTINCT mac AS m FROM jinmu_mac UNION SELECT DISTINCT appratus_mac_address AS m FROM jinmu_product ) AS T
   		LEFT JOIN device AS D ON CONV( T.m, 16, 10 ) = D.mac 
   	WHERE
   		D.mac IS NULL 
   		AND T.m <> '' 
   	AND T.m IS NOT NULL 
   	) AS P;
   ```

   

#### 3.3.3 初始化用户与组织机构相关数据

> 注意事项：
>
> - 已经失效（is_valid = 1) 的数据不导入迁移。包括之前台湾停用账号部分。

1. 初始化 Owner 数据

   运行 **[jinmu-db-v2-migration-003.sql](jinmu-db-v2-migration-003.sql)** 脚本文件内容。

2. 初始化`organization`数据

   运行 **[jinmu-db-v2-migration-004.sql](jinmu-db-v2-migration-004.sql)** 脚本文件内容。

3. 初始化`organization_owner`数据

   运行 **[jinmu-db-v2-migration-005.sql](jinmu-db-v2-migration-005.sql)** 脚本文件内容。

4. 初始化常客 User 数据以及组织与常客关联关系数据

   运行 **[jinmu-db-v2-migration-006.sql](jinmu-db-v2-migration-006.sql)** 脚本文件内容。

#### 3.3.4 初始化 `user_preferences` 数据

1. 运行 **[jinmu-db-v2-migration-007.sql](jinmu-db-v2-migration-007.sql)** 脚本文件内容。

#### 3.3.5 初始化 `feedback` 数据

1. 运行 **[jinmu-db-v2-migration-008.sql](jinmu-db-v2-migration-008.sql)** 脚本文件内容。

#### 3.3.6 初始化 `device_organization_binding` 数据

1. 运行 **[jinmu-db-v2-migration-009.sql](jinmu-db-v2-migration-009.sql)** 脚本文件内容。

#### 3.3.7 初始化 `subscription` 数据

1. 运行 **[jinmu-db-v2-migration-010.sql](jinmu-db-v2-migration-010.sql)** 脚本文件内容。

#### 3.3.8 初始化 `record` 数据

1. 运行 **[jinmu-db-v2-migration-011.sql](jinmu-db-v2-migration-011.sql)** 脚本文件内容。

## 4 后续维护计划

### 4.1 将老数据迁移到别的地方备份，并建库保留以便查询

### 4.2 重命名老的数据库表，避免应用程序继续访问旧的数据，破获历史记录

执行以下SQL脚本：

```sql
RENAME TABLE `jinmu_account_preferences` TO `legacy_jinmu_account_preferences`;
RENAME TABLE `jinmu_company` TO `legacy_jinmu_company`;
RENAME TABLE `jinmu_contract` TO `legacy_jinmu_contract`;
RENAME TABLE `jinmu_feedback` TO `legacy_jinmu_feedback`;
RENAME TABLE `jinmu_log` TO `legacy_jinmu_log`;
RENAME TABLE `jinmu_mac` TO `legacy_jinmu_mac`;
RENAME TABLE `jinmu_product` TO `legacy_jinmu_product`;
RENAME TABLE `jinmu_record` TO `legacy_jinmu_record`;
RENAME TABLE `jinmu_status` TO `legacy_jinmu_status`;
RENAME TABLE `jinmu_subject` TO `legacy_jinmu_subject`;
RENAME TABLE `jinmu_token` TO `legacy_jinmu_token`;
```

## 5 补充登云缺少的那部分的用户和组织
执行以下SQL脚本：
```sql
START TRANSACTION;
SET @ACTIVATED_AT = now();
SET @CREATED_AT = @ACTIVATED_AT;
SET @UPDATED_AT = @ACTIVATED_AT;
SET @REGISTER_TIME = @ACTIVATED_AT;
SET autocommit = 0;
-- 创建user
INSERT INTO `user` (`username`, `password`, `zone`, `register_type`, `register_time`, `register_source_client_id`,`customized_code`, `nickname`, `gender`, `birthday`, `height`, `weight`, `is_activated`, `activated_at`, `created_at`, `updated_at`) VALUES ('dengyun-10001', 'Ac-XJUOar7vgQ5O0', 'CN-X', 'username', @REGISTER_TIME,'dengyun-10001','custom_dengyun', 'dengyun', 'M', '1998-01-01', '170', '60', '1', @ACTIVATED_AT, @CREATED_AT, @UPDATED_AT);
SET @USER_ID = LAST_INSERT_ID();
SELECT @USER_ID;
-- 创建organization
INSERT INTO `organization` (`name`, `phone`, `email`, `type`, `state`, `city`, `district`, `street`, `address`,`customized_code`, `is_valid`, `created_at`, `updated_at`) VALUES ('深圳市每天美耶科技有限公司', '1366004663', 'caijianping@idengyun.com', '养生', '深圳', '深圳市', '前海深港合作区', '深圳市前海深港合作区前湾一路1号A栋201室', '深圳市前海深港合作区前湾一路1号A栋201室','custom_dengyun', '1', @CREATED_AT, @UPDATED_AT);
SET @ORG_ID = LAST_INSERT_ID();
SELECT concat('@USER_ID =  ', @USER_ID,'@ORG_ID = ' ,@ORG_ID);
-- 创建subscription
INSERT INTO `subscription` (`organization_id`, `subscription_type`, `active`, `activated_at`, `expired_at`, `max_user_limits`, `contract_year`, `created_at`, `updated_at`) VALUES (@ORG_ID, '0', '1', @ACTIVATED_AT, date_add(@ACTIVATED_AT, interval 10 year), '500000', '10', @CREATED_AT, @UPDATED_AT);
-- 创建组织owner
INSERT INTO `organization_owner` (`organization_id`, `owner_id`, `created_at`, `updated_at`) VALUES (@ORG_ID, @USER_ID,  @CREATED_AT, @UPDATED_AT);
-- 创建组织user
INSERT INTO `organization_user` (`organization_id`, `user_id`, `created_at`, `updated_at`) VALUES (@ORG_ID, @USER_ID,  @CREATED_AT, @UPDATED_AT);
-- 创建userPreferences
insert into user_preferences (`user_id`,`enable_constitution_differentiation`,`enable_syndrome_differentiation`,`created_at`, `updated_at`) values (@USER_ID, 0, 0, @CREATED_AT, @UPDATED_AT);
COMMIT;
```



