package com.iflytek.ccr.finder.listenr;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.service.CommonService;
import com.iflytek.ccr.finder.service.GrayConfigService;
import com.iflytek.ccr.finder.service.impl.CommonServiceImpl;
import com.iflytek.ccr.finder.service.impl.GrayConfigServiceImpl;
import com.iflytek.ccr.finder.utils.ConfigManager;
import com.iflytek.ccr.finder.utils.ListenerUtil;
import com.iflytek.ccr.finder.utils.StringUtils;
import com.iflytek.ccr.finder.utils.ZkInstanceUtil;
import com.iflytek.ccr.finder.value.GrayConfigValue;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.framework.recipes.cache.NodeCache;
import org.apache.curator.framework.recipes.cache.NodeCacheListener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;

/**
 * 支持灰度配置的节点监听
 */
public class GrayNodeCacheListener implements NodeCacheListener {

    private static final Logger logger = LoggerFactory.getLogger(GrayNodeCacheListener.class);
    GrayConfigService grayConfigService = new GrayConfigServiceImpl();
    CommonService commonService = new CommonServiceImpl();
    private FinderManager finderManager;
    private List<String> configNameList;

    public GrayNodeCacheListener(FinderManager finderManager, List<String> configNameList) {
        this.finderManager = finderManager;
        this.configNameList = configNameList;
    }

    @Override
    public void nodeChanged() throws Exception {
        try {
            String rootConfigPath = ConfigManager.getInstance().getStringConfigByKey(Constants.CONFIG_PATH);
            String grayConfigPath = rootConfigPath + Constants.GRAY_NODE_PATH;
            List<GrayConfigValue> grayValueList = grayConfigService.parseGrayData(grayConfigPath);
            GrayConfigValue grayConfigValue = grayConfigService.getGrayServer(finderManager, grayValueList);

            //如果是灰度服务，则监听并加载灰度配置，否则监听并加载正常配置
            //null 说明不在灰度组中
            //得到当前正在使用中的灰度配置对象
            GrayConfigValue usedGrayConfigValue = finderManager.getGlobalCache().getGrayConfigValue();
            boolean flag = needUpdateMonitorPath(usedGrayConfigValue, grayConfigValue);
            finderManager.getGlobalCache().setGrayConfigValue(grayConfigValue);
            if (flag) {
                if (null != grayConfigValue) {
                    String monitorRootPath = rootConfigPath + Constants.GRAY_NODE_PATH + "/" + grayConfigValue.getGroupId();
                    monitorConfig(monitorRootPath, this.configNameList);
                } else {
                    monitorConfig(rootConfigPath, this.configNameList);
                }

                //如果nodeCachesList是空，说明没有监听任何配置，这时候监听正常配置
            } else if (finderManager.getGlobalCache().nodeCachesList.isEmpty()) {
                monitorConfig(rootConfigPath, this.configNameList);
            }

            //如果没有监听任何节点，不再注册
            if (finderManager.getGlobalCache().configListenerMap.isEmpty()) {
                return;
            }

            String consumerPath = "";
            String basePath = "";
            String ipAddr = finderManager.getBootConfig().getMeteData().getAddress();
            if (null != grayConfigValue) {
                basePath = rootConfigPath + Constants.GRAY_CONSUMER_NODE_PATH + "/" + grayConfigValue.getGroupId();
                consumerPath = basePath + "/" + ipAddr;
            } else {
                basePath = rootConfigPath;
                consumerPath = rootConfigPath + Constants.NORMAL_CONSUMER_NODE_PATH + "/" + ipAddr;
            }

            if (StringUtils.isNOtNullOrEmpty(finderManager.getGlobalCache().getConfigConsumerPath())) {
                commonService.unRegisterConsumer(finderManager.getGlobalCache().getConfigConsumerPath());
                //删除对无用消费者路径的监控
                finderManager.getGlobalCache().monitorPathList.remove(finderManager.getGlobalCache().getConfigConsumerPath());
            }

            finderManager.getGlobalCache().setConfigConsumerPath(consumerPath);
            finderManager.getGlobalCache().setBasePath(basePath);
            commonService.registerConsumer(consumerPath);
            //增加对消费者路径的监控
            finderManager.getGlobalCache().monitorPathList.add(consumerPath);

        } catch (Exception e) {
            logger.error(String.format("GrayNodeCacheListener error:%s", e.getMessage()), e);
        }

    }

    /**
     * 监听配置的变化
     *
     * @param monitorRootPath
     */
    private void monitorConfig(String monitorRootPath, List<String> configNameList) {
        removeOldConfigMonitor();
        for (String fileName : configNameList) {
            //如果fileName 被取消订阅，则不需要监控
            if (finderManager.getGlobalCache().configListenerMap.containsKey(fileName)) {
                String fullPath = monitorRootPath + "/" + fileName;
                monitorSingleConfig(fullPath, fileName);
            }
        }
    }

    private void removeOldConfigMonitor() {
        List<NodeCache> nodeCachesList = finderManager.getGlobalCache().nodeCachesList;
        ListenerUtil.closeNodeCache(nodeCachesList);
        nodeCachesList.clear();
    }

    /**
     * 监控单个配置的变更
     *
     * @param fullPath
     */
    private void monitorSingleConfig(String fullPath, String fileName) {
        ZkHelper zkHelper = ZkInstanceUtil.getInstance();
        if (zkHelper.checkExists(fullPath)) {
            NodeCache nodeCache = zkHelper.addListener(new ConfigNodeCacheListener(this.finderManager, fullPath), fullPath, false);
            finderManager.getGlobalCache().nodeCachesList.add(nodeCache);
            finderManager.getGlobalCache().configListenerMap.put(fileName, nodeCache);
        } else {
            logger.error(String.format("path : s% does not exists!", fullPath));
        }
    }


    /**
     * 比较灰度组的变化情况：
     * 1、原来在灰度组A，现在依旧在A  忽略
     * 2、原来在灰度组A，现在在灰度组B 变更监听路径
     * 3、原来在灰度组A，现在不在灰度组  变更监听路径
     * 4、原来不在灰度组，现在在灰度组A 变更监听路径
     * 5、原来不在灰度组，现在依旧不在灰度组 忽略
     *
     * @param usedGrayConfigValue
     * @param grayConfigValue
     * @return
     */
    private boolean needUpdateMonitorPath(GrayConfigValue usedGrayConfigValue, GrayConfigValue grayConfigValue) {
        boolean flag = false;
        // 原来在灰度组A，现在在灰度组B 变更监听路径
        if (null != usedGrayConfigValue && null != grayConfigValue) {
            if (!usedGrayConfigValue.getGroupId().equals(grayConfigValue.getGroupId())) {
                flag = true;
            }

            // 原来在灰度组A，现在不在灰度组
        } else if (null != usedGrayConfigValue && null == grayConfigValue) {
            flag = true;

            //原来不在灰度组，现在在灰度组A
        } else if (null == usedGrayConfigValue && null != grayConfigValue) {
            flag = true;
        }

        return flag;
    }
}
