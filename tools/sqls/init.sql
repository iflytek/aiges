CREATE DATABASE /*!32312 IF NOT EXISTS*/`ifly_cynosure_3` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_bin */;

USE `ifly_cynosure_3`;


-- ----------------------------
-- Table structure for tb_cluster
-- ----------------------------
DROP TABLE IF EXISTS `tb_cluster`;
CREATE TABLE `tb_cluster` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '集群唯一标识',
  `name` varchar(100) COLLATE utf8_bin NOT NULL COMMENT '集群名称',
  `push_url` varchar(500) COLLATE utf8_bin NOT NULL COMMENT '集群推送地址',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `user_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '用户id',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_name` (`name`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_dic
-- ----------------------------
DROP TABLE IF EXISTS `tb_dic`;
CREATE TABLE `tb_dic` (
  `type` varchar(1) NOT NULL,
  `type_name` varchar(20) NOT NULL,
  PRIMARY KEY (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tb_gray_group
-- ----------------------------
DROP TABLE IF EXISTS `tb_gray_group`;
CREATE TABLE `tb_gray_group` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '版本id',
  `version_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '版本id',
  `user_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '用户id',
  `name` varchar(100) COLLATE utf8_bin DEFAULT NULL COMMENT '灰度组名称',
  `content` text COLLATE utf8_bin COMMENT '推送实例内容',
  `description` varchar(500) COLLATE utf8_bin DEFAULT NULL COMMENT '版本描述',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_name_version_id` (`version_id`,`name`) USING BTREE,
  KEY `idx_service_id` (`version_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_lastest_search
-- ----------------------------
DROP TABLE IF EXISTS `tb_lastest_search`;
CREATE TABLE `tb_lastest_search` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '唯一标识',
  `user_id` varchar(50) COLLATE utf8_bin DEFAULT NULL COMMENT '用户id',
  `url` varchar(500) COLLATE utf8_bin DEFAULT NULL COMMENT '请求地址',
  `pre_condition` varchar(500) COLLATE utf8_bin DEFAULT NULL COMMENT '前置条件',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_user_id_url` (`user_id`,`url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_loadbalance
-- ----------------------------
DROP TABLE IF EXISTS `tb_loadbalance`;
CREATE TABLE `tb_loadbalance` (
  `name` varchar(50) COLLATE utf8_bin NOT NULL,
  `abbr` varchar(100) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`name`,`abbr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_project
-- ----------------------------
DROP TABLE IF EXISTS `tb_project`;
CREATE TABLE `tb_project` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '项目id',
  `name` varchar(100) COLLATE utf8_bin NOT NULL COMMENT '项目名称',
  `description` varchar(500) COLLATE utf8_bin DEFAULT NULL COMMENT '项目描述',
  `user_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '用户id',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_name` (`name`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_project_member
-- ----------------------------
DROP TABLE IF EXISTS `tb_project_member`;
CREATE TABLE `tb_project_member` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '唯一标识',
  `user_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '用户id',
  `project_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '项目id',
  `creator` tinyint(4) NOT NULL COMMENT '是否为创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_user_id_project_id` (`user_id`,`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_role
-- ----------------------------
DROP TABLE IF EXISTS `tb_role`;
CREATE TABLE `tb_role` (
  `role_name` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '角色名称',
  `role_type` tinyint(4) NOT NULL COMMENT '角色类型',
  PRIMARY KEY (`role_name`,`role_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_service
-- ----------------------------
DROP TABLE IF EXISTS `tb_service`;
CREATE TABLE `tb_service` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '服务id',
  `group_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '服务组id',
  `name` varchar(100) COLLATE utf8_bin NOT NULL COMMENT '服务名称',
  `description` varchar(500) COLLATE utf8_bin DEFAULT NULL COMMENT '服务描述',
  `user_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '用户id',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_group_id_name` (`group_id`,`name`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_service_api_version
-- ----------------------------
DROP TABLE IF EXISTS `tb_service_api_version`;
CREATE TABLE `tb_service_api_version` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '版本id',
  `service_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '服务id',
  `api_version` varchar(20) COLLATE utf8_bin NOT NULL COMMENT '版本号',
  `user_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '用户id',
  `description` varchar(500) COLLATE utf8_bin DEFAULT NULL COMMENT '版本描述',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_service_id_version` (`api_version`,`service_id`),
  KEY `idx_service_id` (`service_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_service_config
-- ----------------------------
DROP TABLE IF EXISTS `tb_service_config`;
CREATE TABLE `tb_service_config` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '配置id',
  `user_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '用户id',
  `version_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '版本id',
  `gray_group_id` varchar(50) COLLATE utf8_bin NOT NULL DEFAULT '0' COMMENT '灰度组id，正常组默认为0',
  `name` varchar(100) COLLATE utf8_bin NOT NULL COMMENT '配置名称',
  `path` varchar(500) COLLATE utf8_bin NOT NULL COMMENT '配置路径',
  `content` longblob COMMENT '配置内容',
  `md5` varchar(50) COLLATE utf8_bin DEFAULT NULL COMMENT '配置内容md5值',
  `description` varchar(500) COLLATE utf8_bin DEFAULT NULL COMMENT '配置描述',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_version_id_group_id_name` (`version_id`,`gray_group_id`,`name`) USING BTREE,
  KEY `idx_user_id` (`user_id`),
  KEY `idx_version_id` (`version_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_service_config_history
-- ----------------------------
DROP TABLE IF EXISTS `tb_service_config_history`;
CREATE TABLE `tb_service_config_history` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '配置id',
  `user_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '用户id',
  `config_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '配置id',
  `content` longblob COMMENT '配置内容',
  `md5` varchar(50) COLLATE utf8_bin DEFAULT NULL COMMENT '配置内容md5值',
  `description` varchar(500) COLLATE utf8_bin DEFAULT NULL COMMENT '配置描述',
  `push_version` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '推送版本号',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_config_id` (`config_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_service_config_push_feedback
-- ----------------------------
DROP TABLE IF EXISTS `tb_service_config_push_feedback`;
CREATE TABLE `tb_service_config_push_feedback` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '推送反馈id',
  `push_id` varchar(50) COLLATE utf8_bin DEFAULT NULL COMMENT '推送id',
  `project` varchar(100) COLLATE utf8_bin DEFAULT NULL COMMENT '项目名称',
  `service_group` varchar(100) COLLATE utf8_bin DEFAULT NULL COMMENT '服务组名称',
  `service` varchar(100) COLLATE utf8_bin DEFAULT NULL COMMENT '服务名称',
  `version` varchar(20) COLLATE utf8_bin DEFAULT NULL COMMENT '服务版本号',
  `gray_group_id` varchar(50) COLLATE utf8_bin NOT NULL,
  `config` varchar(100) COLLATE utf8_bin DEFAULT NULL COMMENT '配置名称',
  `addr` varchar(100) COLLATE utf8_bin DEFAULT NULL COMMENT '地址',
  `update_status` tinyint(4) DEFAULT NULL COMMENT '更新状态',
  `load_status` tinyint(4) DEFAULT NULL COMMENT '加载状态',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `load_time` datetime DEFAULT NULL COMMENT '加载时间',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_push_id` (`push_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_service_config_push_history
-- ----------------------------
DROP TABLE IF EXISTS `tb_service_config_push_history`;
CREATE TABLE `tb_service_config_push_history` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '推送id',
  `user_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '用户id',
  `gray_group_id` varchar(50) COLLATE utf8_bin NOT NULL,
  `project` varchar(100) COLLATE utf8_bin NOT NULL COMMENT '项目名称',
  `service_group` varchar(100) COLLATE utf8_bin NOT NULL COMMENT '服务组名称',
  `service` varchar(100) COLLATE utf8_bin NOT NULL COMMENT '服务名称',
  `version` varchar(20) COLLATE utf8_bin NOT NULL COMMENT '服务版本号',
  `cluster_text` text COLLATE utf8_bin NOT NULL COMMENT '集群',
  `service_config_text` text COLLATE utf8_bin NOT NULL COMMENT '服务配置',
  `push_time` datetime DEFAULT NULL COMMENT '推送时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_service_discovery_push_feedback
-- ----------------------------
DROP TABLE IF EXISTS `tb_service_discovery_push_feedback`;
CREATE TABLE `tb_service_discovery_push_feedback` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '推送反馈id',
  `push_id` varchar(50) COLLATE utf8_bin DEFAULT NULL COMMENT '推送id',
  `project` varchar(100) COLLATE utf8_bin DEFAULT NULL COMMENT '项目名称',
  `service_group` varchar(100) COLLATE utf8_bin DEFAULT NULL COMMENT '服务组名称',
  `consumer_service` varchar(100) COLLATE utf8_bin DEFAULT NULL COMMENT '消费端服务名称',
  `consumer_version` varchar(20) COLLATE utf8_bin DEFAULT NULL COMMENT '消费端版本',
  `provider_service` varchar(100) COLLATE utf8_bin DEFAULT NULL COMMENT '提供端服务名称',
  `provider_version` varchar(20) COLLATE utf8_bin DEFAULT NULL COMMENT '提供端版本',
  `addr` varchar(100) COLLATE utf8_bin DEFAULT NULL COMMENT '地址',
  `update_status` tinyint(4) DEFAULT NULL COMMENT '更新状态',
  `load_status` tinyint(4) DEFAULT NULL COMMENT '加载状态',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `load_time` datetime DEFAULT NULL COMMENT '加载时间',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `api_version` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `type` varchar(1) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_push_id` (`push_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_service_discovery_push_history
-- ----------------------------
DROP TABLE IF EXISTS `tb_service_discovery_push_history`;
CREATE TABLE `tb_service_discovery_push_history` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '服务发现推送id',
  `user_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '用户id',
  `project` varchar(100) COLLATE utf8_bin NOT NULL COMMENT '项目名称',
  `service_group` varchar(100) COLLATE utf8_bin NOT NULL COMMENT '服务组名称',
  `service` varchar(100) COLLATE utf8_bin NOT NULL COMMENT '服务名称',
  `version` varchar(100) COLLATE utf8_bin NOT NULL COMMENT '服务版本',
  `cluster_text` text COLLATE utf8_bin NOT NULL COMMENT '集群',
  `push_time` datetime DEFAULT NULL COMMENT '推送时间',
  `api_version` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_service_group
-- ----------------------------
DROP TABLE IF EXISTS `tb_service_group`;
CREATE TABLE `tb_service_group` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '唯一标识',
  `project_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '项目id',
  `name` varchar(100) COLLATE utf8_bin NOT NULL COMMENT '服务组名称',
  `description` varchar(500) COLLATE utf8_bin DEFAULT NULL COMMENT '服务组描述',
  `user_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '用户id',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_project_id_name` (`project_id`,`name`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_project_id` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_service_version
-- ----------------------------
DROP TABLE IF EXISTS `tb_service_version`;
CREATE TABLE `tb_service_version` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '版本id',
  `service_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '服务id',
  `version` varchar(20) COLLATE utf8_bin NOT NULL COMMENT '版本号',
  `user_id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '用户id',
  `description` varchar(500) COLLATE utf8_bin DEFAULT NULL COMMENT '版本描述',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_service_id_version` (`version`,`service_id`),
  KEY `idx_service_id` (`service_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tb_user
-- ----------------------------
DROP TABLE IF EXISTS `tb_user`;
CREATE TABLE `tb_user` (
  `id` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '唯一标识',
  `account` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '账号',
  `password` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '密码',
  `user_name` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '用户姓名',
  `phone` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '联系方式',
  `email` varchar(100) COLLATE utf8_bin NOT NULL COMMENT '邮箱',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `role_type` tinyint(4) NOT NULL DEFAULT '2' COMMENT '角色类型 1.管理员 2.普通用户',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_account` (`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


/*Data for the table `tb_loadbalance` */

insert  into `tb_loadbalance`(`name`,`abbr`) values 
('一致性Hash','ConsistentHash'),
('最少活跃调用数','LeastActive'),
('轮循','RoundRobin'),
('随机','Random');

/*Data for the table `tb_role` */

insert  into `tb_role`(`role_name`,`role_type`) values 
('admin',1),
('user',2);

/*Data for the table `tb_user` */

insert  into `tb_user`(`id`,`account`,`password`,`user_name`,`phone`,`email`,`create_time`,`update_time`,`role_type`) values 
('3688919024172793856','admin','b301eba23dc7ab9df6b7315de2ac222a','admin','13739263609','sctang2@iflytek.com','2017-11-10 10:12:07','2018-02-05 05:20:27',1);

insert  into `tb_dic` values
 ('0', '服务配置'),
 ('1', '路由配置'),
 ('2', '实例配置');
