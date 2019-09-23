package com.iflytek.ccr.finder.monitor;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.constants.TaskType;
import com.iflytek.ccr.finder.utils.ConfigManager;
import com.iflytek.ccr.finder.utils.ZkInstanceUtil;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Iterator;

/**
 * 监控临时节点是否存在
 */
public class PathMonitorTask implements Runnable {

    private static final Logger logger = LoggerFactory.getLogger(PathMonitorTask.class);

    private FinderManager finderManager;

    public PathMonitorTask(FinderManager finderManager) {
        this.finderManager = finderManager;
    }

    @Override
    public void run() {
        while (true) {
            survivalCheck();
        }
    }

    /**
     * 节点存活检查
     */
    private void survivalCheck() {
        try {
            //10秒检查一次
            Thread.sleep(10000);
            Iterator<String> iterator = finderManager.getGlobalCache().monitorPathList.iterator();
            while (iterator.hasNext()) {
                String path = iterator.next();
                ZkHelper zkHelper = ZkInstanceUtil.getInstance();
                if (!zkHelper.checkExists(path)) {
                    zkHelper.addEphemeral(path, "");
                }
            }
        } catch (Exception e) {
            logger.error("", e);
        }

        try {
            //检查companion的zookeeper节点，如果节点发生变更，说明zookeeper的集群信息可能发生了变化，需要重新初始化
            String zkNodePath = ConfigManager.getInstance().getStringConfigByKey(Constants.KEY_ZK_NODE_PATH);
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            if (!zkHelper.checkExists(zkNodePath)) {
                finderManager.getGlobalCache().taskQueue.add(TaskType.INIT);
                //由于monitortask的执行周期是60秒，所以这里要sleep 60秒后再执行
                Thread.sleep(1000 * 60);
            }
        } catch (Exception e) {
            logger.error("", e);
        }
    }
}
