package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.request.InstanceManageRequestBody.AddGrayGroupInstanceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.InstanceManageRequestBody.EditInstanceRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.graygroup.GrayGroupDetailResponseBody;

/**
 * Created by DELL-5490 on 2018/7/7.
 */
public interface InstanceService {
    /**
     * 通过id查询对应的推送实例内容
     *
     * @param id
     * @return
     */
    Response<String> findById(String id);

    /**
     * 查询推送实例列表
     *
     * @return
     */
    Response<String> findList(String versionId);

    /**
     * 编辑灰度组
     *
     * @param body
     * @return
     */
    Response<GrayGroupDetailResponseBody> edit(EditInstanceRequestBody body);

    /**
     * 查询推送实例列表
     *
     * @return
     */
    Response<String> appointList(AddGrayGroupInstanceRequestBody body);
}
