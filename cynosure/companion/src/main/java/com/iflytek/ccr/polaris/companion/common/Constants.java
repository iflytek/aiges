package com.iflytek.ccr.polaris.companion.common;

public interface Constants {

    String PROVIDER = "provider";

    String CONSUMER = "consumer";

    String DATE_PATTERN = "yyyyMMddHH";

    String BACKUP_PRE = "/opt/server/backup/";

    String QUEUE_PATH_ZK_NODE = "/polaris/zkNode";

    int QUEUE_MAX_NUM = 3000;
    /**
     * 配置反馈路径
     */
    String QUEUE_PATH_CONFIG = "/polaris/queue/config";

    /**
     * 服务反馈路径
     */
    String QUEUE_PATH_SERVICE = "/polaris/queue/service";

    String CONFIG_PATH_PREFIX = "/polaris/config";

    String SERVICE_PATH_PREFIX = "/polaris/service";

    String FEED_BACK_PATH_PRE = "/feedback/data/";

    String PUSH_CONFIG_FEEDBACK_SITE_URI = "/api/v1/service/config/feedback";

    String PUSH_CONFIG_SERVICE_SITE_URI = "/api/v1/service/discovery/feedback";

    String DISCOVERY_SERVICE_SITE_URI = "/api/v1/service/discovery/add";

    String CONTENT_TYPE_OCTET_STREAM = "application/octet-stream";
    String ZK_NODE_PUSHID = "pushId";
    String ZK_NODE_GRAY = "gray";
    String ZK_NODE_FILE_NAME = "fileName";
    String GRAY_SERVERS = "grayServers";
    String GRAY_GROUP = "grayGroupId";
    String GROUP_ID = "group_id";
    String SERVER_LIST = "server_list";
    String ZK_NODE_PATH = "path";
    String USER_DATA = "user";
    String SDK_DATA = "sdk";
    String ZK_PATH_CONF = "/conf";
    String ZK_PATH_ROUTE = "/route";
    String ZK_NODE_DATA = "data";
    String ZK_NODE_DATA_PATH = "/data";

    String DEFAULT_CHARSET = "UTF-8";

    String SUCCESS = "0";
    String DEFAULT_WEIGHT = "100";
    String DEFAULT_VALID = "true";
    String WEIGHT = "weight";
    String IS_VALID = "is_valid";

    int NUM_3 = 3;
    int NUM_5 = 5;
    int NUM_6 = 6;

    int DEFAULT_MONITOR_TIME = 5000;
    /**
     * 1分钟对应的毫秒数
     */
    long MILLIS_1_MINUTES = 1000 * 60;
}
