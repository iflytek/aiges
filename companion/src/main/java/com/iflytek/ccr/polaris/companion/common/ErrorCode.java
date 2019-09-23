package com.iflytek.ccr.polaris.companion.common;

/**
 * Created by eric on 2017/11/21.
 */
public abstract class ErrorCode {
    public static final int SUCCESS = 0;
    public static final int INTERNAL_EXCEPTION = -2;
    /**
     * 路径不存在
     */
    public static final int PATH_NOT_EXISTS = 10002;
    public static final int PARAM_INVALID = 10001;
}
