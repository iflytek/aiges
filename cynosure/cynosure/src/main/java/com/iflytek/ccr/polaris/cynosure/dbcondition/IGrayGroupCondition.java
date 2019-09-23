package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.domain.GrayGroup;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.AddGrayGroupRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.EditGrayGroupRequestBody;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Created by DELL-5490 on 2018/7/4.
 */
public interface IGrayGroupCondition {
    /**
     * 新增灰度组
     *
     * @param body
     * @return
     */
    GrayGroup add(AddGrayGroupRequestBody body);

    /**
     * 通过灰度组名称，版本id查询灰度组
     *
     * @param name
     * @param versionId
     * @return
     */
    GrayGroup find(String name, String versionId);

    /**
     * 查询灰度组列表
     *
     * @param map
     * @return
     */
    List<GrayGroup> findList(HashMap<String, Object> map);

    /**
     * 查询灰度组总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 根据id查询灰度配置
     *
     * @param id
     * @return
     */
    GrayGroup findById(String id);

    /**
     * 根据id删除灰度组
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 根据id更新灰度组
     *
     * @param id
     * @param body
     * @return
     */
    GrayGroup updateById(String id, EditGrayGroupRequestBody body);

    /**
     * 复制灰度组
     * @param grayGroupList
     * @param versionId
     * @return
     */
    Map<String, String> copy1(List<GrayGroup> grayGroupList, String versionId);
}