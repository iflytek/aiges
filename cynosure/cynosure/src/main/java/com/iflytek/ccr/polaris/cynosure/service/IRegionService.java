package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.region.AddRegionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.region.EditRegionRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.region.RegionDetailResponseBody;

/**
 * 区域业务接口
 *
 * @author sctang2
 * @create 2017-11-14 20:52
 **/
public interface IRegionService {
    /**
     * 查询区域列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findList(BaseRequestBody body);

    /**
     * 查询区域详情
     *
     * @param id
     * @return
     */
    Response<RegionDetailResponseBody> find(String id);

    /**
     * 新增区域
     *
     * @param body
     * @return
     */
    Response<RegionDetailResponseBody> add(AddRegionRequestBody body);

    /**
     * 编辑区域
     *
     * @param body
     * @return
     */
    Response<RegionDetailResponseBody> edit(EditRegionRequestBody body);

    /**
     * 删除区域
     *
     * @param body
     * @return
     */
    Response<String> delete(IdRequestBody body);
}
