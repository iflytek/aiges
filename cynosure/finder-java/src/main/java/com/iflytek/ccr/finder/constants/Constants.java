package com.iflytek.ccr.finder.constants;

public interface Constants {


    String DEFAULT_CHARSET = "UTF-8";

    String GROUP_ID = "group_id";
    String SERVER_LIST = "server_list";

    String REGISTER_PRE = "register_";
    String UNREGISTER_PRE = "unregister_";
    String KEY_REGISTER_ADDR = "register_addr";
    String KEY_CONFIG_NAME_LIST = "configNameList";
    String KEY_SERVICE_NAME_LIST = "serviceNameList";
    String KEY_CONFIG_HANDLE = "configChangedHandler";
    String KEY_SERVICE_HANDLE = "serviceHandle";
    String KEY_BOOT_CONFIG = "bootConfig";
    String KEY_ROUTE_CONSUMER = "consumer";
    String KEY_ROUTE_PROVIDER = "provider";
    String KEY_ROUTE_ONLY = "only";
    String KEY_ROUTE_ROUTE_RULE_ID = "id";
    String KEY_PROJECT = "project";
    String KEY_SERVICE = "service";
    String KEY_CACHEP_PATH = "cachePath";
    String KEY_GROUP = "group";
    String KEY_ZK_NODE_PATH = "zk_node_path";
    int SUCCESS = 0;
    int INTERNAL_EXCEPTION = -2;
    String CONFIG_PATH = "config_path";
    String GRAY_NODE_PATH = "/gray";
    String GRAY_CONSUMER_NODE_PATH = "/consumer/gray";
    String NORMAL_CONSUMER_NODE_PATH = "/consumer/normal";
    String SERVICE_PATH = "service_path";
    String ZK_ADDR = "zk_addr";
    String UPDATE_STATUS_SUCCESS = "1";
    String UPDATE_STATUS_FAIL = "0";
    String LOAD_STATUS_SUCCESS = "1";
    String LOAD_STATUS_FAIL = "0";
    String IP_ADDR = "ip_addr";
    String POLARIS_PATH_ZK_NODE = "/polaris/zkNode";
    String COMPANION_URL_PRE = "/finder/query_zk_info";
    String CONFIG_FEEDBACK_URL = "/finder/push_config_feedback";
    String SERVICE_FEEDBACK_URL = "/finder/push_service_feedback";
    String REGISTER_SERVICE_INFO_URL = "/finder/register_service_info";
    String REGISTER_SERVICE_URL = "/finder/register_service_info";
    String UNREGISTER_SERVICE_URL = "/finder/unregister_service";
    String WEIGHT = "weight";
    String IS_VALID = "is_valid";

    int DEFAULT_WEIGHT = 100;
    boolean DEFAULT_VALID = true;
    String DEFAULT_CACHEP_PATH = "/polaris/finder/cache";


    String PROXY_MODE = "proxy_mode";
    String LOAD_BALANCE = "lb_mode";

    int DEFAULT_MONITOR_TIME = 5000;

    int NUM_1 = 1;

    long MILLIS_1_MINUTES = 1000 * 60;

    String KEY_ROUTE_ONLY_Y = "Y";

    String KEY_SERVICE_CONFIG_CHANGE = "0";
    String KEY_SERVICE_ROUTE_CHANGE = "1";
    String KEY_SERVICE_INSTANCE_CHANGE = "2";
}
