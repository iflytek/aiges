package com.iflytek.ccr.finder.utils;

import com.iflytek.ccr.zkutil.ZkHelper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ZkInstanceUtil {
    private static final Logger logger = LoggerFactory.getLogger(ZkInstanceUtil.class);
    private static ZkHelper zkHelper = null;
    private static String ZK_ADDR = null;

    public static void setZkAddr(String zkAddr) {
        ZK_ADDR = zkAddr;
    }

    public static void init(String zkAddr) {
        zkHelper = new ZkHelper(zkAddr);
    }

    public static void init(String connectionString, int connectionTimeoutMs, int sessionTimeoutMs) {
        zkHelper = new ZkHelper(connectionString, connectionTimeoutMs, sessionTimeoutMs);
    }

    public static ZkHelper getInstance() {
        if (null == zkHelper) {
            zkHelper = new ZkHelper(ZK_ADDR);
        }
        return zkHelper;
    }
}
