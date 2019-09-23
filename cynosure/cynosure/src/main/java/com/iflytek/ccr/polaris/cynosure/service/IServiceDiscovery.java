package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.domain.LoadBalance;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.*;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.servicediscovery.AddApiVersionResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.servicediscovery.ServiceDiscoveryResponseBody;

import java.util.List;

/**
 * 服务发现接口
 *
 * @author sctang2
 * @create 2017-12-05 15:36
 **/
public interface IServiceDiscovery {
    /**
     * 新增服务
     *
     * @param body
     * @return
     */
    Response<AddApiVersionResponseBody> add(AddServiceDiscoveryRequestBody body);

    /**
     * 查询最近的服务发现列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findLastestList(QueryServiceDiscoveryListRequestBody body);

    /**
     * 根据项目、集群、服务，查询服务发现列表
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findLastestList1(QueryServiceDiscoveryListRequestBody body);

    /**
     * 查询服务发现明细
     *
     * @param body
     * @return
     */
    Response<ServiceDiscoveryResponseBody> find(QueryServiceDiscoveryDetailRequestBody body);

    /**
     * 编辑服务发现
     *
     * @param body
     * @return
     */
    Response<String> edit(EditServiceDiscoveryRequestBody body);


    /**
     * 查询服务提供者列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> provider(QueryServiceDiscoveryDetailRequestBody body);

    /**
     * 编辑提供端
     *
     * @param body
     * @return
     */
    Response<String> editProvider(EditServiceDiscoveryProviderRequestBody body);

    /**
     * 查询服务消费端列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> consumer(QueryServiceDiscoveryDetailRequestBody body);

    /**
     * 更新反馈
     *
     * @param body
     * @return
     */
    Response<String> feedback(ServiceDiscoveryFeedBackRequestBody body);

    /**
     * 查询负载均衡列表
     *
     * @return
     */
    Response<List<LoadBalance>> findBalanceList();
}
