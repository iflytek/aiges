package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.AddServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.CopyServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.EditServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.QueryServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.serviceversion.ServiceVersionDetailResponseBody;

/**
 * 服务版本业务接口
 *
 * @author sctang2
 * @create 2017-11-17 16:02
 **/
public interface IServiceVersion {
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
     * 查询版本明细
     *
     * @param id
     * @return
     */
    Response<ServiceVersionDetailResponseBody> find(String id);

    /**
     * 新增版本
     *
     * @param body
     * @return
     */
    Response<ServiceVersionDetailResponseBody> add(AddServiceVersionRequestBody body);

    /**
     * 编辑版本
     *
     * @param body
     * @return
     */
    Response<ServiceVersionDetailResponseBody> edit(EditServiceVersionRequestBody body);

    /**
     * 删除版本
     *
     * @param body
     * @return
     */
    Response<String> delete(IdRequestBody body);

    /**
     * 复制版本
     *
     * @param body
     * @return
     */
    Response<ServiceVersionDetailResponseBody> copy(CopyServiceVersionRequestBody body);
}
