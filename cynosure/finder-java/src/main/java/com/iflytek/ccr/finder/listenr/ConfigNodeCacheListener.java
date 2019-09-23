package com.iflytek.ccr.finder.listenr;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.handler.ConfigChangedHandler;
import com.iflytek.ccr.finder.utils.*;
import com.iflytek.ccr.finder.value.Config;
import com.iflytek.ccr.finder.value.ErrorCode;
import com.iflytek.ccr.finder.value.ZkDataValue;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.framework.recipes.cache.NodeCacheListener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * 配置文件节点监听者
 */
public class ConfigNodeCacheListener implements NodeCacheListener {

    private static final Logger logger = LoggerFactory.getLogger(ConfigNodeCacheListener.class);

    private FinderManager finderManager;
    private String path;
    private ConfigChangedHandler configChangedHandler;

    public ConfigNodeCacheListener(FinderManager finderManager, String path) {
        this.finderManager = finderManager;
        this.path = path;
        this.configChangedHandler = finderManager.getGlobalCache().getConfigChangedHandler();
    }

    @Override
    public void nodeChanged() throws Exception {
        try {
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            ZkDataValue zkDataValue = ByteUtil.parseZkData(zkHelper.getByteData(this.path));
            String updateStauts = null;
            if (ErrorCode.SUCCESS == zkDataValue.getRet()) {
                updateStauts = Constants.UPDATE_STATUS_SUCCESS;
            } else {
                updateStauts = Constants.UPDATE_STATUS_FAIL;
            }
            String pushId = zkDataValue.getPushId();
            Config config = new Config();
            config.setFile(zkDataValue.getRealData());
            config.setName(path.substring(path.lastIndexOf("/") + 1, path.length()));
            if (config.getName().endsWith(".toml") || config.getName().endsWith(".TOML")) {
                config.setConfigMap(FinderFileUtils.parseTomlFile(config.getFile()));
            }
            String groupId = getGroupId(path);
            long updateTime = System.currentTimeMillis();
            if (finderManager.getGlobalCache().initMap.get(config.getName())) {
                finderManager.getGlobalCache().initMap.put(config.getName(), false);
                return;
            }
            //更新最新的内容到缓存文件
            FinderFileUtils.writeByteArrayToFile(PathUtils.getCacheFilePath(finderManager, "config") + config.getName(), config.getFile());
            boolean isSuccess = configChangedHandler.onConfigFileChanged(config);
            //pushId不是正常值，则不进行反馈
            if (!("0".equals(pushId) || null == pushId || "".equals(pushId))) {
                if (isSuccess) {
                    RemoteUtil.pushConfigFeedback(finderManager,  groupId, config.getName(),pushId, updateStauts, Constants.LOAD_STATUS_SUCCESS, String.valueOf(updateTime), String.valueOf(System.currentTimeMillis()));
                } else {
                    RemoteUtil.pushConfigFeedback(finderManager, groupId, config.getName(), pushId, updateStauts, Constants.LOAD_STATUS_FAIL, String.valueOf(updateTime), String.valueOf(System.currentTimeMillis()));
                }
            }
        } catch (Exception e) {
            logger.error(String.format("ConfigNodeCacheListener error:%s", e.getMessage()), e);
        }
    }

    /**
     * 获取groupId
     *
     * @param path
     * @return
     */
    private String getGroupId(String path) {
        if (path.contains(Constants.GRAY_NODE_PATH)) {
            String pre = ConfigManager.getInstance().getStringConfigByKey(Constants.CONFIG_PATH) + Constants.GRAY_NODE_PATH;
            int end = path.lastIndexOf("/");
            return path.substring(pre.length() + 1, end);
        } else {
            return "0";
        }
    }
}
