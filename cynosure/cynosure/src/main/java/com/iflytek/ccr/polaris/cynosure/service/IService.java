package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.AddServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.CopyServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.EditServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.QueryServiceListRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.service.ServiceDetailResponseBody;

/**
 * 服务业务逻辑接口
 *
 * @author sctang2
 * @create 2017-11-16 19:54
 **/
public interface IService {
    /**
     * 查询最近的服务列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findLastestList(QueryServiceListRequestBody body);

    /**
     * 查询服务列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findList(QueryServiceListRequestBody body);

    /**
     * 查询服务明细
     *
     * @param id
     * @return
     */
    Response<ServiceDetailResponseBody> find(String id);

    /**
     * 新增服务
     *
     * @param body
     * @return
     */
    Response<ServiceDetailResponseBody> add(AddServiceRequestBody body);

    /**
     * 编辑服务
     *
     * @param body
     * @return
     */
    Response<ServiceDetailResponseBody> edit(EditServiceRequestBody body);

    /**
     * 删除服务
     *
     * @param body
     * @return
     */
    Response<String> delete(IdRequestBody body);

    /**
     * 复制服务
     *
     * @param body
     * @return
     */
    Response<ServiceDetailResponseBody> copy(CopyServiceRequestBody body);
}
