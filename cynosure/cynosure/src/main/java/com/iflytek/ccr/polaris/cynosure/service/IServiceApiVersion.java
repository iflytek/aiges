package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.domain.Service;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceApiVersion;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.AddServiceApiVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.EditServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.QueryServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.servicediscovery.ServiceApiVersionDetailResponseBody;

import java.util.List;

/**
 * 服务版本业务接口
 *
 * @author sctang2
 * @create 2017-11-17 16:02
 **/
public interface IServiceApiVersion {
    /**
     * 新增版本
     *
     * @param body
     * @return
     */
    Response<ServiceApiVersionDetailResponseBody> add(AddServiceApiVersionRequestBody body);
    /**
     * 查询最近的版本列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findLastestList(QueryServiceVersionRequestBody body);

    /**
     * 查询版本列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findList(QueryServiceVersionRequestBody body);

    /**
     * 查询服务列表
     *
     * @param serviceIds
     * @return
     */
    List<ServiceApiVersion> findList1(List<String> serviceIds);

    /**
     * 查询版本明细
     *
     * @param id
     * @return
     */
    Response<ServiceApiVersionDetailResponseBody> find(String id);

    /**
     * 编辑版本
     *
     * @param body
     * @return
     */
    Response<ServiceApiVersionDetailResponseBody> edit(EditServiceVersionRequestBody body);

    /**
     * 删除版本
     *
     * @param body
     * @return
     */
    Response<String> delete(IdRequestBody body);


}
