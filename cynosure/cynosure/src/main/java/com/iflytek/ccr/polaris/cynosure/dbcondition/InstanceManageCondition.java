package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.domain.GrayGroup;
import com.iflytek.ccr.polaris.cynosure.request.InstanceManageRequestBody.EditInstanceRequestBody;

import java.util.List;

/**
 * Created by DELL-5490 on 2018/7/4.
 */
public interface InstanceManageCondition {
    /**
     * 根据id更新灰度组
     *
     * @param id
     * @param body
     * @return
     */
    GrayGroup updateById(String id, EditInstanceRequestBody body);

    List<String> findTotal(String versionId, String grayGroupId);
}