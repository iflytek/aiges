package com.iflytek.ccr.polaris.companion.utils;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.zkutil.ZkHelper;

import java.util.concurrent.TimeUnit;

/**
 * zk工具类
 */
public class ZkInstanceUtil {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ZkInstanceUtil.class);
    private static ZkHelper zkHelper = new ZkHelper(ConfigManager.getStringConfigByKey(ConfigManager.KEY_ZKSTR));

    public static ZkHelper getInstance() {
        return zkHelper;
    }


    /**
     * 检查zookeeper连接是否可用
     *
     * @return
     */
    public static boolean checkZkHelper() {
        boolean canConnect = false;
        try {
            canConnect = zkHelper.canConnect(3, TimeUnit.SECONDS);
            if (!canConnect) {
                zkHelper.closeClient();
                zkHelper = new ZkHelper(ConfigManager.getStringConfigByKey(ConfigManager.KEY_ZKSTR));
            }
            canConnect = zkHelper.canConnect(5, TimeUnit.SECONDS);
        } catch (InterruptedException e) {
            logger.error(e);
            Thread.currentThread().interrupt();
        }
        return canConnect;
    }
}
