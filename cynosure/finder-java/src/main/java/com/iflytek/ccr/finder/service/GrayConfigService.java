package com.iflytek.ccr.finder.service;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.value.GrayConfigValue;

import java.util.List;

/**
 * 灰度服务
 */
public interface GrayConfigService {

    /**
     * 解析灰度配置数据
     * @param grayConfigPath
     * @return
     */
    List<GrayConfigValue> parseGrayData(String grayConfigPath);

    /**
     * 获取灰度组配置：如果在灰度组中则返回灰度组，否则返回空
     * @param finderManager
     * @param grayValueList
     * @return
     */
    GrayConfigValue getGrayServer(FinderManager finderManager, List<GrayConfigValue> grayValueList);
}
