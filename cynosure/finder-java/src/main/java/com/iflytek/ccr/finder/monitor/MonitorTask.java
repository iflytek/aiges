package com.iflytek.ccr.finder.monitor;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.constants.TaskType;
import com.iflytek.ccr.finder.handler.ConfigChangedHandler;
import com.iflytek.ccr.finder.handler.ServiceHandle;
import com.iflytek.ccr.finder.utils.ConfigManager;
import com.iflytek.ccr.finder.utils.RemoteUtil;
import com.iflytek.ccr.finder.utils.StringUtils;
import com.iflytek.ccr.finder.utils.ZkInstanceUtil;
import com.iflytek.ccr.finder.value.BootConfig;
import com.iflytek.ccr.finder.value.CommonResult;
import com.iflytek.ccr.finder.value.SubscribeRequestValue;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.Map;

public class MonitorTask implements Runnable {


    private final Logger logger = LoggerFactory.getLogger(MonitorTask.class);

    private FinderManager finderManager;

    public MonitorTask(FinderManager finderManager) {
        this.finderManager = finderManager;
    }


    @Override
    public void run() {
        while (true) {
            try {
                Thread.sleep(1000 * 60);
                TaskType taskType = finderManager.getGlobalCache().taskQueue.take();
                if (TaskType.INIT.equals(taskType)) {
                    reInit(finderManager.getBootConfig());
                } else if (TaskType.CONFIG.equals(taskType)) {
                    logger.info("MonitorTask CONFIG start");
                    reInitConfig();
                } else if (TaskType.SERVICE.equals(taskType)) {
                    reInitService();
                }
            } catch (Exception e) {
                logger.error("", e);
            }
        }

    }

    /**
     * 重新订阅服务
     */
    private void reInitService() {
        List<SubscribeRequestValue> requestValueList = finderManager.getServiceCache().getRequestValueList();
        ServiceHandle serviceHandle = finderManager.getServiceCache().getServiceHandle();
        if (null != requestValueList && null != serviceHandle) {
            finderManager.useAndSubscribeService(requestValueList, serviceHandle);
        }
    }

    /**
     * 重新订阅配置
     */
    private void reInitConfig() {
        List<String> configNameList = (List<String>) finderManager.getGlobalCache().getConfigNameList();
        ConfigChangedHandler configChangedHandler = finderManager.getGlobalCache().getConfigChangedHandler();
        if (null != configNameList && null != configChangedHandler) {
            finderManager.getConfigFinder().useAndSubscribeConfig(finderManager, configNameList, configChangedHandler, true);
        }
    }

    private void reInit(BootConfig bootConfig) {

        CommonResult response = RemoteUtil.queryZkInfo(bootConfig);
        if (null != response && Constants.SUCCESS == response.getRet()) {
            Map map = (Map) response.getData();
            String configPath = (String) map.get(Constants.CONFIG_PATH);
            String servicePath = (String) map.get(Constants.SERVICE_PATH);
            String zkNodePath = (String) map.get(Constants.KEY_ZK_NODE_PATH);
            List zkAddrList = (List) map.get(Constants.ZK_ADDR);
            ConfigManager.getInstance().put(Constants.CONFIG_PATH, configPath);
            ConfigManager.getInstance().put(Constants.SERVICE_PATH, servicePath);
            ConfigManager.getInstance().put(Constants.KEY_ZK_NODE_PATH, zkNodePath);

            String oriZkAddr = ConfigManager.getInstance().getStringConfigByKey(Constants.ZK_ADDR);
            String newZkAddr = StringUtils.join(zkAddrList, ",");
            if (null != oriZkAddr && oriZkAddr.equals(newZkAddr)) {
                return;
            }

            //如果zk 地址为空，则不需要关闭
            if (null != oriZkAddr) {
                ZkInstanceUtil.getInstance().closeClient();
            }
            ZkInstanceUtil.setZkAddr(newZkAddr);

            if (bootConfig.getZkSessionTimeout() > 0 && bootConfig.getZkConnectTimeout() > 0) {
                ZkInstanceUtil.init(StringUtils.join(zkAddrList, ","), bootConfig.getZkConnectTimeout(), bootConfig.getZkSessionTimeout());
            } else {
                ZkInstanceUtil.init(StringUtils.join(zkAddrList, ","));
            }

            reInitService();
            reInitConfig();
        } else {
            finderManager.getGlobalCache().taskQueue.add(TaskType.INIT);
        }
    }
}
