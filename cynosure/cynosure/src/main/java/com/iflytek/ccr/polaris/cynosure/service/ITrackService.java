package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.IdsRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.track.QueryTrackDetailRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.track.QueryTrackRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;

/**
 * 轨迹业务逻辑接口
 *
 * @author sctang2
 * @create 2017-11-24 11:56
 **/
public interface ITrackService {
    /**
     * 查询最近的配置推送轨迹列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findLastestConfigList(QueryTrackRequestBody body);

    /**
     * 查询最近的服务发现推送轨迹列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findLastestDiscoveryList(QueryTrackRequestBody body);

    /**
     * 查询配置推送轨迹明细
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findConfig(QueryTrackDetailRequestBody body);

    /**
     * 查询服务发现推送轨迹明细
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findDiscovery(QueryTrackDetailRequestBody body);

    /**
     * 删除配置推送轨迹
     *
     * @param body
     * @return
     */
    Response<String> deleteConfig(IdRequestBody body);

    /**
     * 批量删除配置推送轨迹
     *
     * @param body
     * @return
     */
    Response<String> batchDeleteConfig(IdsRequestBody body);

    /**
     * 删除服务发现推送轨迹
     *
     * @param body
     * @return
     */
    Response<String> deleteDiscovery(IdRequestBody body);

    /**
     * 批量删除服务发现推送轨迹
     *
     * @param body
     * @return
     */
    Response<String> batchDeleteDiscovery(IdsRequestBody body);

    /**
     * 快速查询
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findList(BaseRequestBody body);
}
