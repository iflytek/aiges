package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.AddClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.CopyClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.EditClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.QueryClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.cluster.ClusterDetailResponseBody;

/**
 * 集群业务接口
 *
 * @author sctang2
 * @create 2017-11-15 17:39
 **/
public interface IClusterService {
    /**
     * 查询最近的集群列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findLastestList(QueryClusterRequestBody body);

    /**
     * 查询集群列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findList(QueryClusterRequestBody body);

    /**
     * 新增集群
     *
     * @param body
     * @return
     */
    Response<ClusterDetailResponseBody> add(AddClusterRequestBody body);

    /**
     * 编辑集群
     *
     * @param body
     * @return
     */
    Response<ClusterDetailResponseBody> edit(EditClusterRequestBody body);

    /**
     * 删除集群
     *
     * @param body
     * @return
     */
    Response<String> delete(IdRequestBody body);

    /**
     * 查询集群明细
     *
     * @param id
     * @return
     */
    Response<ClusterDetailResponseBody> find(String id);

    /**
     * 复制集群
     *
     * @param body
     * @return
     */
    Response<ClusterDetailResponseBody> copy(CopyClusterRequestBody body);
}
