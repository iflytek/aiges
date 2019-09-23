package com.iflytek.ccr.polaris.cynosure.exception;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;

/**
 * @author sctang
 * @version 1.0
 * @create 2017年5月10日 下午2:39:13
 * @description 应用程序异常
 */
public class GlobalExceptionUtil {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(GlobalExceptionUtil.class);

    /**
     * @param ex
     * @description 记录异常
     * @author sctang
     * @create 2017年5月10日 下午2:43:10
     * @version 1.0
     */
    public static void log(Exception ex) {
        logger.error(ex + " " + ex.getMessage());
    }
}
