DROP TABLE
IF EXISTS `tb_dic`;

CREATE TABLE `tb_dic` (
	`type` VARCHAR (1) NOT NULL,
	`type_name` VARCHAR (20) NOT NULL,
	PRIMARY KEY (`type`)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;

INSERT INTO `tb_dic`
VALUES
	('0', '服务配置');

INSERT INTO `tb_dic`
VALUES
	('1', '路由配置');

INSERT INTO `tb_dic`
VALUES
	('2', '实例配置');

DROP TABLE
IF EXISTS `tb_gray_group`;

CREATE TABLE `tb_gray_group` (
	`id` VARCHAR (50) COLLATE utf8_bin NOT NULL COMMENT '版本id',
	`version_id` VARCHAR (50) COLLATE utf8_bin NOT NULL COMMENT '版本id',
	`user_id` VARCHAR (50) COLLATE utf8_bin NOT NULL COMMENT '用户id',
	`name` VARCHAR (100) COLLATE utf8_bin DEFAULT NULL COMMENT '灰度组名称',
	`content` text COLLATE utf8_bin COMMENT '推送实例内容',
	`description` VARCHAR (500) COLLATE utf8_bin DEFAULT NULL COMMENT '版本描述',
	`create_time` datetime DEFAULT NULL COMMENT '创建时间',
	`update_time` datetime DEFAULT NULL COMMENT '更新时间',
	PRIMARY KEY (`id`),
	UNIQUE KEY `uq_name_version_id` (`version_id`, `name`) USING BTREE,
	KEY `idx_service_id` (`version_id`),
	KEY `idx_user_id` (`user_id`)
) ENGINE = INNODB DEFAULT CHARSET = utf8 COLLATE = utf8_bin;

DROP TABLE
IF EXISTS `tb_service_api_version`;

CREATE TABLE `tb_service_api_version` (
	`id` VARCHAR (50) COLLATE utf8_bin NOT NULL COMMENT '版本id',
	`service_id` VARCHAR (50) COLLATE utf8_bin NOT NULL COMMENT '服务id',
	`api_version` VARCHAR (20) COLLATE utf8_bin NOT NULL COMMENT '版本号',
	`user_id` VARCHAR (50) COLLATE utf8_bin NOT NULL COMMENT '用户id',
	`description` VARCHAR (500) COLLATE utf8_bin DEFAULT NULL COMMENT '版本描述',
	`create_time` datetime DEFAULT NULL COMMENT '创建时间',
	`update_time` datetime DEFAULT NULL COMMENT '更新时间',
	PRIMARY KEY (`id`),
	UNIQUE KEY `uq_service_id_version` (`api_version`, `service_id`),
	KEY `idx_service_id` (`service_id`),
	KEY `idx_user_id` (`user_id`)
) ENGINE = INNODB DEFAULT CHARSET = utf8 COLLATE = utf8_bin;

ALTER TABLE tb_service_config ADD COLUMN gray_group_id VARCHAR (50) NOT NULL DEFAULT "0";

ALTER TABLE tb_service_config_push_feedback ADD COLUMN gray_group_id VARCHAR (50 DEFAULT "0";

ALTER TABLE tb_service_config_push_history ADD COLUMN gray_group_id VARCHAR (50) NOT NULL DEFAULT "0";

ALTER TABLE tb_service_discovery_push_history ADD COLUMN version VARCHAR (100) NOT NULL DEFAULT "";

ALTER TABLE tb_service_config_push_feedback ADD COLUMN api_version VARCHAR (20);

ALTER TABLE tb_service_discovery_push_feedback ADD COLUMN type VARCHAR (1);

ALTER TABLE tb_service_config DROP INDEX uq_version_id_name;

ALTER TABLE tb_service_config ADD CONSTRAINT uq_version_id_group_id_name UNIQUE (
	`version_id`,
	`gray_group_id`,
	`name`
) USING BTREE;

ALTER TABLE tb_lastest_search ADD PRIMARY KEY (id);

ALTER TABLE tb_lastest_search ADD UNIQUE KEY uq_user_id_url (user_id, url);

DROP INDEX uq_version_id_name ON tb_service_config;

ALTER TABLE tb_service_config ADD UNIQUE KEY `uq_version_id_group_id_name` (
	version_id,
	gray_group_id,
	`name`
);

ALTER TABLE tb_service_discovery_push_feedback ADD COLUMN `api_version` varchar(20);


ALTER TABLE tb_lastest_search ADD PRIMARY KEY (id);

ALTER TABLE tb_lastest_search ADD UNIQUE KEY uq_user_id_url (user_id, url);

DROP INDEX uq_version_id_name ON tb_service_config;

ALTER TABLE tb_service_config ADD UNIQUE KEY `uq_version_id_group_id_name` (
	version_id,
	gray_group_id,
	`name`
);

ALTER TABLE tb_service_discovery_push_feedback ADD COLUMN `api_version` varchar(20);

