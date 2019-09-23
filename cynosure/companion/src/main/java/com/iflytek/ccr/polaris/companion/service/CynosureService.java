package com.iflytek.ccr.polaris.companion.service;

import com.iflytek.ccr.polaris.companion.common.JsonResult;

public interface CynosureService {

    JsonResult queryServiceList(String path);

    /**
     * 刷新提供者状态
     *
     * @param path
     * @return
     */
    JsonResult queryProviderOrConsumerList(String path);

    /**
     * 刷新conf 状态
     *
     * @param path
     * @return
     */
    JsonResult refreshConfStatus(String path);


    /**
     * 刷新service  conf 状态
     *
     * @param path
     * @return
     */
    JsonResult refreshServiceConfStatus(String path);

    /**
     * 删除zk数据
     *
     * @param path
     * @return
     */
    JsonResult delZkData(String path);
}
