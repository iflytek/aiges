package com.iflytek.ccr.finder.value;

/**
 * Created by eric on 2017/11/21.
 */
public abstract class ErrorCode {

    /**
     * 成功
     */
    public static final int SUCCESS = 0;

    /**
     * 成功
     */
    public static final int READ_CACHE_SUCCESS = 1;

    /**
     * 内部异常
     */
    public static final int INTERNAL_EXCEPTION = -2;

    /**
     * 参数非法
     */
    public static final int PARAM_INVALID = 10000;

    /**
     * companion url is empty
     */
    public static final int MISS_COMPANION_URL = 10001;

    /**
     * 读取文件失败
     */
    public static final int READ_FILE_FAIL = 10002;

    /**
     * 请求companion失败
     */
    public static final int QUERY_ZK_INFO_FAIL = 10003;

    /**
     * 文件不存在
     */
    public static final int CONFIG_MISS_FILE = 10100;
}


//const (
//        Success             ReturnCode = 0
//        InvalidParam        ReturnCode = 10000
//        MissCompanionUrl 	ReturnCode = 10001
//        ConfigMissName      ReturnCode = 10100
//        ConfigMissCacheFile ReturnCode = 10101
//        ZkMissRootPath      ReturnCode = 10200 + iota
//        ZkMissAddr
//        ZkGetInfoError
//        )
