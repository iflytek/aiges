## 部署文档

### 简介

部署配置中心&服务发现Polaris套件，包括Cynosure、Companion、zk集群，采用容器化方案。

### 基础环境

docker版本1.12.0，请自行安装。

MySQL资源请协调DBA处理。

### ZK集群

请参考zkinstall-new.md

### Cynosure

创建目录:mkdir -p /opt/polaris/cynosure

在cynosure目录下创建config目录，在config目录中上传application.yml。

application.yml文件在https://git.xfyun.cn/AIaaS/polaris/src/master/cynosure/src/main/resources/config/application.yml中，请自行下载。

cynosure/start.sh，参数中172.16.59.153/develop/cynosure:2.0.3为镜像地址
```
docker run -d --name cynosure --net="host" -v /opt/polaris/cynosure/config/application.yml:/opt/server/cynosure/config/application.yml -v /opt/polaris/cynosure/logs:/log/server -v /etc/localtime:/etc/localtime 172.16.59.153/develop/cynosure:2.0.3 sh watchdog.sh
```

**注:application.yml为cynosure的配置文件，包含mysql配置、redis配置。支持高可用，若部署集群，请自行部署redis，修改application.yml配置，并修改sessionShare为1**。
如：
```
  #redis配置
  redis:
      # Redis数据库索引（默认为0）
      database: 0
      # Redis服务器地址
      host: 10.1.86.212
      # Redis服务器连接端口
      port: 6379
      # Redis服务器连接密码（默认为空）
      password:
      # 连接超时时间（毫秒）
      timeout: 0
      pool:
          # 连接池最大连接数（使用负值表示没有限制）
          max-active: 100
          # 连接池最大阻塞等待时间（使用负值表示没有限制）
          max-wait: -1
          # 连接池中的最大空闲连接
          max-idle: 8
          # 连接池中的最小空闲连接
          min-idle: 0
  #会话共享
  sessionShare: 1
  #对外开发接口认证配置
    private_key: MIIBVQIBADANBgkqhkiG9w0BAQEFAASCAT8wggE7AgEAAkEAnDci65joIaWLoUC7GjKsKWEK9425W7vQlIxnTqpd2KJsvrI1mcI_1L9RXGqkrovrtfsCZt4x6p4mT7zjTViIyQIDAQABAkAE_mI4Y8_v22nmQrp4cOw9-mMuXLJzM0LMrNxUkG-lkCbIsaIsMu43wj5WPhI-yaUP_oItgR0NTIGoOBF7ZnSBAiEA4Iyrvwcj3Ihs8WKdJb_il6IRv4Q7hhLAt43r0SnUaTkCIQCyGFAp2rOZWLq35bZ2o-QfdJ7StUPm0Dh5TKoggCXsEQIgfnfb9yAfW4Le0Oj4lx1Gkp5uHo5sM-v17KubCFfl0UkCIQCeUTS58Dvl3tWlcqRAVTMOr2ocn5ysC3-YfQljeOe9MQIhAIma_c07u8s1Xd6UpljZ8Ui1frkxGfpvTuXA24QiP0Dv
    valid_interval: 5000
```
**注：基于k8s部署的时候，如果使用configmap，可以参考 https://git.xfyun.cn/sctang2/cynosure
### Companion

创建目录:mkdir -p /opt/polaris/companion

companion/start.sh
```
docker run -d --name companion --net="host" -v /opt/polaris/companion/logs:/log/server -v /etc/localtime:/etc/localtime 172.16.59.153/develop/companion:2.0.2 sh watchdog.sh -h10.1.86.211 -p6868 -z10.1.86.211:2181,10.1.86.70:2181,10.1.86.212:2181 -whttps://10.1.87.69:8095
```

**注：启动companion参数分别代表companion的ip地址-h10.1.86.211，端口-p6868，zk集群地址-z10.1.86.211:2181,10.1.86.70:2181,10.1.86.212:2181，cynosure地址-whttps://10.1.87.69:8095**

启动完成后，可以输入以下url验证是否正常启动了。如果正常启动了，会有返回结果。<br/>
http://ip:port/finder/query_zk_info?project=AIaaS&group=aitest&service=iatExecutor&version=2.0.0<br/>
(ip\port根据实际情况替换)

### Mysql脚本

```
CREATE DATABASE /*!32312 IF NOT EXISTS*/`ifly_cynosure_2` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_bin */;

USE `ifly_cynosure_2`;


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
 ('0', '服务配置')，
 ('1', '路由配置')，
 ('2', '实例配置');

```
