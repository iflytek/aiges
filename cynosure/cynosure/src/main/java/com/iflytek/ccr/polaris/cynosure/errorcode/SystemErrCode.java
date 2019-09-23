package com.iflytek.ccr.polaris.cynosure.errorcode;

/**
 * 系统级别错误码
 *
 * @author sctang2
 * @create 2017-11-09 15:54
 **/
public class SystemErrCode {
    //基础错误码
    public static final int    ERRCODE_REQUEST_FAIL          = 100001;
    public static final String ERRMSG_REQUEST_FAIL           = "请求失败";
    public static final int    ERRCODE_NOT_FOUND_API         = 100002;
    public static final String ERRMSG_NOT_FOUND_API          = "不存在该业务";
    public static final int    ERRCODE_INVALID_PARAMETER     = 100003;
    public static final int    ERRCODE_NOT_AUTH              = 100004;
    public static final String ERRMSG_NOT_AUTH               = "没有权限执行此操作";
    public static final int    ERRCODE_NOT_FILE              = 100005;
    public static final String ERRMSG_NOT_FILE               = "请上传文件";
    public static final int    ERRCODE_FILE_CONTENT_NOT_NULL = 100006;
    public static final String ERRMSG_FILE_CONTENT_NOT_NULL  = "文件内容不能为空";
    public static final int    ERRCODE_FILE_NAME_BATCH_SAME  = 100007;
    public static final String ERRMSG_FILE_NAME_BATCH_SAME   = "文件名不能相同";
    public static final int    ERRCODE_FILE_TOO_BIG          = 100008;
    public static final String ERRMSG_FILE_TOO_BIG           = "文件太大";

    public static final String ERRMSG_ID_NOT_NULL = "id不能为空";

    //用户模块错误码
    public static final int    ERRCODE_USER_NOT_EXISTS         = 110001;
    public static final String ERRMSG_USER_NOT_EXISTS          = "不存在该用户";
    public static final int    ERRCODE_USER_PASSWORD_INCORRECT = 110002;
    public static final String ERRMSG_USER_PASSWORD_INCORRECT  = "密码错误";
    public static final int    ERRCODE_USER_NOT_LOGIN          = 110003;
    public static final String ERRMSG_USER_NOT_LOGIN           = "用户未登录";
    public static final String ERRMSG_APP_ID_ERROR             = "app_id不正确";
    public static final String ERRMSG_TOKEN_ERROR              ="token不正确";
    public static final String ERRMSG_TOKEN_DECRPT_ERROR       ="token解签不正确";
    public static final int    ERRCODE_USER_EXISTS             = 110004;
    public static final String ERRMSG_USER_EXISTS              = "已存在该用户";
    public static final int    ERRCODE_USER_PASSWORD_CONFIRM   = 110005;
    public static final String ERRMSG_USER_PASSWORD_CONFIRM    = "密码和确认密码不相同";

    public static final String ERRMSG_USER_ACCOUNT_NOT_NULL   = "账号不能为空";
    public static final String ERRMSG_USRE_ACCOUNT_MAX_LENGTH = "账号最大支持50个字符";
    public static final String ERRMSG_USER_PASSWORD_NOT_NULL  = "密码不能为空";
    public static final String ERRMSG_USER_NAME_NOT_NULL      = "用户名不能为空";
    public static final String ERRMSG_USRE_NAME_MAX_LENGTH    = "用户名最大支持50个字符";
    public static final String ERRMSG_USER_TELEPHONE_NOT_NULL = "手机号不能为空";
    public static final String ERRMSG_USRE_TELEPHONE_INVALID  = "手机号无效";
    public static final String ERRMSG_USER_EMAIL_NOT_NULL     = "邮箱不能为空";
    public static final String ERRMSG_USRE_EMAIL_MAX_LENGTH   = "邮箱最大支持100个字符";
    public static final String ERRMSG_USRE_EMAIL_INVALID      = "邮箱无效";
    public static final String ERRMSG_USER_ID_NOT_NULL        = "用户id不能为空";

    //区域模块错误码
    public static final int    ERRCODE_REGION_NOT_EXISTS = 120001;
    public static final String ERRMSG_REGION_NOT_EXISTS  = "不存在该区域";
    public static final int    ERRCODE_REGION_EXISTS     = 120002;
    public static final String ERRMSG_REGION_EXISTS      = "已存在该区域";
    public static final int    ERRCODE_COMPANION_EXISTS  = 120003;
    public static final String ERRMSG_COMPANION_EXISTS   = "已存在该companion";

    public static final String ERRMSG_REGION_NAME_NOT_NULL       = "区域名称不能为空";
    public static final String ERRMSG_REGION_NAME_MAX_LENGTH     = "区域名称最大支持100个字符";
    public static final String ERRMSG_REGION_PUSH_URL_NOT_NULL   = "推送地址不能为空";
    public static final String ERRMSG_REGION_PUSH_URL_MAX_LENGTH = "推送地址最大支持500个字符";
    public static final String ERRMSG_REGION_PUSH_URL_INVALID    = "推送地址无效";
    public static final String ERRMSG_REGION_ID_NOT_NULL         = "区域id不能为空";
    public static final String ERRMSG_REGION_IDS_NOT_NULL        = "区域ids不能为null";
    public static final String ERRMSG_REGION_IDS_IS_NOT_EMPTY        = "区域ids不能为空";

    //项目模块错误码
    public static final int    ERRCODE_PROJECT_NOT_EXISTS        = 130001;
    public static final String ERRMSG_PROJECT_NOT_EXISTS         = "不存在该项目";
    public static final int    ERRCODE_PROJECT_EXISTS            = 130002;
    public static final String ERRMSG_PROJECT_EXISTS             = "已存在该项目";
    public static final int    ERRCODE_PROJECT_MEMBER_EXIEST     = 130003;
    public static final String ERRMSG_PROJECT_MEMBER_EXIEST      = "用户已在该项目中";
    public static final int    ERRCODE_PROJECT_MEMBER_NOT_EXIEST = 130004;
    public static final String ERRMSG_PROJECT_MEMBER_NOT_EXIEST  = "用户不在该项目中";

    public static final String ERRMSG_PROJECT_NAME_NOT_NULL   = "项目名称不能为空";
    public static final String ERRMSG_PROJECT_NAME_MAX_LENGTH = "项目名称最大支持100个字符";
    public static final String ERRMSG_PROJECT_DESC_MAX_LENGTH = "项目描述最大支持500个字符";
    public static final String ERRMSG_PROJECT_ID_NOT_NULL     = "项目id不能为空";
    public static final String ERRMSG_PROJECT_ID_MAX_LENGTH   = "项目id最大支持50个字符";

    //集群模块错误码
    public static final int    ERRCODE_CLUSTER_NOT_EXISTS      = 140001;
    public static final String ERRMSG_CLUSTER_NOT_EXISTS       = "不存在该集群";
    public static final int    ERRCODE_CLUSTER_EXISTS          = 140002;
    public static final String ERRMSG_CLUSTER_EXISTS           = "已存在该集群";
    public static final int    ERRCODE_CLUSTER_CREATE          = 140003;
    public static final String ERRMSG_CLUSTER_CREATE           = "该用户已创建集群，无法进行删除";
    public static final int    ERRCODE_CLUSTER_COPY_NOT_EXISTS = 140004;
    public static final String ERRMSG_CLUSTER_COPY_NOT_EXISTS  = "不存在该复制集群";

    public static final String ERRMSG_CLUSTER_NAME_NOT_NULL   = "集群名称不能为空";
    public static final String ERRMSG_CLUSTER_NAME_MAX_LENGTH = "集群名称最大支持100个字符";
    public static final String ERRMSG_CLUSTER_ID_NOT_NULL     = "集群id不能为空";
    public static final String ERRMSG_CLUSTER_ID_MAX_LENGTH   = "集群id最大支持50个字符";
    public static final String ERRMSG_CLUSTER_DESC_MAX_LENGTH = "集群描述最大支持500个字符";

    //服务模块错误码
    public static final int    ERRCODE_SERVICE_NOT_EXISTS      = 150001;
    public static final String ERRMSG_SERVICE_NOT_EXISTS       = "不存在该服务";
    public static final int    ERRCODE_SERVICE_EXISTS          = 150002;
    public static final String ERRMSG_SERVICE_EXISTS           = "已存在该服务";
    public static final int    ERRCODE_SERVICE_CREATE          = 150003;
    public static final String ERRMSG_SERVICE_CREATE           = "该用户已创建服务，无法进行删除";
    public static final int    ERRCODE_SERVICE_COPY_NOT_EXISTS = 150004;
    public static final String ERRMSG_SERVICE_COPY_NOT_EXISTS  = "不存在该复制服务";

    public static final String ERRMSG_SERVICE_NAME_NOT_NULL   = "服务名称不能为空";
    public static final String ERRMSG_SERVICE_NAME_MAX_LENGTH = "服务名称最大支持100个字符";
    public static final String ERRMSG_SERVICE_DESC_MAX_LENGTH = "服务描述最大支持500个字符";
    public static final String ERRMSG_SERVICE_ID_NOT_NULL     = "服务id不能为空";
    public static final String ERRMSG_SERVICE_ID_MAX_LENGTH   = "服务id最大支持50个字符";

    //版本模块错误码
    public static final int    ERRCODE_SERVICE_VERSION_NOT_EXISTS = 160001;
    public static final String ERRMSG_SERVICE_VERSION_NOT_EXISTS  = "不存在该版本";
    public static final int    ERRCODE_SERVICE_VERSION_EXISTS     = 160002;
    public static final String ERRMSG_SERVICE_VERSION_EXISTS      = "已存在该版本";
    public static final int    ERRCODE_SERVICE_VERSION_CREATE     = 160003;
    public static final String ERRMSG_SERVICE_VERSION_CREATE      = "该用户已创建版本，无法进行删除";

    public static final String ERRMSG_SERVICE_VERSION_NOT_NULL        = "服务版本号不能为空";
    public static final String ERRMSG_SERVICE_VERSION_MAX_LENGTH      = "服务版本号最大支持20个字符";
    public static final String ERRMSG_SERVICE_VERSION_ID_NOT_NULL     = "服务版本id不能为空";
    public static final String ERRMSG_SERVICE_VERSION_ID_MAX_LENGTH   = "服务版本id最大支持50个字符";
    public static final String ERRMSG_SERVICE_VERSION_DESC_MAX_LENGTH = "服务版本描述最大支持500个字符";

    public static final String ERRMSG_SERVICEAPI_VERSION_NOT_NULL     = "服务API版本号不能为空";

    //配置模块错误码
    public static final int    ERRCODE_SERVICE_CONFIG_NOT_EXISTS = 170001;
    public static final int    ERRCODE_SERVICE_CONFIG_CREATE     = 170002;
    public static final int    ERRCODE_SERVICE_CONFIG_CONTENT_BYTE_MAX_LENGTH     = 170003;
    public static final int    ERRCODE_SERVICE_CONFIG_DOWNLOAD_FALSE = 170004;
    public static final String ERRMSG_SERVICE_CONFIG_NOT_EXISTS  = "不存在该配置";
    public static final String ERRMSG_SERVICE_CONFIG_CREATE      = "该用户已创建配置，无法进行删除";
    public static final String ERRMSG_SERVICE_CONFIG_ID_NOT_NULL        = "服务配置id不能为空";
    public static final String ERRMSG_SERVICE_CONFIG_IDS_NOT_NULL       = "服务配置ids不能为空";
    public static final String ERRMSG_SERVICE_CONFIG_NAME_NOT_NULL      = "服务配置名称不能为空";
    public static final String ERRMSG_SERVICE_CONFIG_NAME_MAX_LENGTH    = "服务配置名称最大支持100个字符";
    public static final String ERRMSG_SERVICE_CONFIG_DESC_MAX_LENGTH    = "服务配置描述最大支持500个字符";
    public static final String ERRMSG_SERVICE_CONFIG_CONTENT_NOT_NULL   = "服务配置内容不能为空";
    public static final String ERRMSG_SERVICE_CONFIG_CONTENT_MAX_LENGTH = "服务配置内容最大支持10240个字符";
    public static final String ERRMSG_SERVICE_CONFIG_ADDR_NOT_NULL      = "服务配置地址不能为空";
    public static final String ERRMSG_SERVICE_CONFIG_ADDR_MAX_LENGTH    = "服务配置地址最大支持100个字符";
    public static final String ERRMSG_SERVICE_CONFIG_CONTENT_BYTE_MAX_LENGTH = "服务配置文件内容最大支持1M";
    public static final String ERRCODE_SERVICE_CONFIG_DOWNLOAD_FALSE_MESSAGE = "下载配置文件出错";
    //灰度组模块错误码
    public static final int    ERRCODE_GRAY_GROUP_NAME_NOT_NULL   = 210006;
    public static final String ERRMSG_GRAY_GROUP_NAME_NOT_NULL    = "灰度组名称不能为空";
    public static final int    ERRCODE_GRAY_GROUP_NAME_MAX_LENGTH = 210007;
    public static final String ERRMSG_GRAY_GROUP_NAME_MAX_LENGTH  = "灰度组名称最大支持100个字符";
    public static final String ERRMSG_GRAY_GROUP_ID_NOT_NULL      = "灰度组id不能为空";
    public static final String ERRMSG_GRAY_GROUP_ID_MAX_LENGTH    = "灰度组id最大支持50个字符";
    public static final String ERRMSG_INSTANCE_NOT_NULL           = "推送实例不能为空";
    public static final int    ERRCODE_GRAY_GROUP_NOT_EXISTS      = 210001;
    public static final String ERRMSG_GRAY_GROUP_NOT_EXISTS       = "不存在该灰度组";
    public static final int    ERRCODE_GRAY_GROUP_EXISTS          = 210002;
    public static final String ERRMSG_GRAY_GROUP_EXISTS           = "已存在该灰度组";
    public static final int    ERRCODE_GRAY_CONFIGS_MAX_SIZE      = 210003;
    public static final String ERRMSG_GRAY_CONFIGS_MAX_SIZE       = "灰度配置文件上传超过10个";
    public static final String ERRMSG_GRAY_GROUP_DESC_MAX_LENGTH  = "灰度组描述最大支持500个字符";
    public static final int    ERRCODE_GRAY_INSTANCE_REPEAT       = 210004;
    public static final String ERRMSG_GRAY_INSTANCE_REPEAT        = "推送实例内容重复";
    public static final int    ERRCODE_GRAY_CONFIG_REPEAT         = 210005;
    public static final String ERRMSG_GRAY_CONFIG_REPEAT          = "上传的灰度配置文件已存在";
    public static final int    ERRCODE_GRAY_INSTANCE_ARE_USED     = 210006;
    public static final String ERRMSG_GRAY_INSTANCE_ARE_USED      = "推送实例已被版本下其它灰度组使用";


    //配置历史模块错误码
    public static final int    ERRCODE_SERVICE_CONFIG_HISTORY_NOT_EXISTS = 180001;
    public static final String ERRMSG_SERVICE_CONFIG_HISTORY_NOT_EXISTS  = "不存在该配置历史";

    //服务发现模块错误码
    public static final String ERRMSG_SERVICE_DISCOVERY_WEIGHT_INVALID       = "服务发现权重范围在0-100";
    public static final String ERRMSG_SERVICE_DISCOVERY_LOADBALANCE_NOT_NULL = "负载均衡不能为空";
    public static final int    ERRCODE_SERVICE_DISCOVERY_PARAMS_REPEAT       = 180002;
    public static final String ERRMSG_SERVICE_DISCOVERY_PARAMS_REPEAT        = "自定义规则key值重复";
    public static final String ERRMSG_SERVICE_DISCOVERY_ROULE_PROVIDER_REPEAT        = "一个服务提供者出现在多个路由规则中";
    public static final String ERRMSG_SERVICE_DISCOVERY_ROULE_CONSUMER_REPEAT        = "一个消费者实例出现在多个路由规则中";

    //轨迹模块错误码
    public static final int    ERRCODE_TRACK_NOT_EXISTS = 200001;
    public static final String ERRMSG_TRACK_NOT_EXISTS  = "不存在该轨迹";
    public static final String ERRMSG_TRACK_ID_NOT_NULL     = "推送id不能为空";
    public static final String ERRMSG_TRACK_IDS_NOT_NULL    = "推送ids不能为空";
    public static final String ERRMSG_TRACK_ID_MAX_LENGTH   = "推送id最大支持50个字符";
    public static final String ERRMSG_TRACK_ISGRAY_NOT_NULL = "灰度标识isGray不能为空";
}
