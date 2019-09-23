package com.iflytek.ccr.polaris.companion.task;

import com.iflytek.ccr.polaris.companion.common.Constants;
import com.iflytek.ccr.polaris.companion.utils.ConfigManager;
import com.iflytek.ccr.polaris.companion.utils.ZkInstanceUtil;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class MonitorZkTask implements Runnable {
    private final Logger logger = LoggerFactory.getLogger(MonitorZkTask.class);

    @Override
    public void run() {

        while (true) {
            try {
                Thread.sleep(Constants.DEFAULT_MONITOR_TIME);
                ZkHelper zkHelper = ZkInstanceUtil.getInstance();
                String path = ConfigManager.getStringConfigByKey(ConfigManager.KEY_ZK_NODE_PATH);
                if (!zkHelper.checkExists(path)) {
                    String zkNodePath = zkHelper.addEphemeralSequential(Constants.QUEUE_PATH_ZK_NODE + "/zk", ConfigManager.getStringConfigByKey(ConfigManager.KEY_ZKSTR));
                    ConfigManager.getInstance().put("zkNodePath", zkNodePath);
                }
            } catch (Exception e) {
                logger.error("", e);
            }
        }

    }
}
