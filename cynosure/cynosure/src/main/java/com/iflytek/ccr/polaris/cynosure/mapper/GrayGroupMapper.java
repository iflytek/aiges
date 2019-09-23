package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.GrayGroup;

import java.util.HashMap;
import java.util.List;

/**
 * Created by DELL-5490 on 2018/7/4.
 */
public interface GrayGroupMapper {
    /**
     * 新增灰度组
     *
     * @param grayGroup
     * @return
     */
    int insert(GrayGroup grayGroup);

    /**
     * 查询灰度组
     *
     * @param grayGroup
     * @return
     */
    GrayGroup find(GrayGroup grayGroup);

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
     * 根据id查询灰度组
     *
     * @param id
     * @return
     */
    GrayGroup findById(String id);

    /**
     * 删除灰度组
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 更新灰度组
     *
     * @param grayGroup
     * @return
     */
    int updateById(GrayGroup grayGroup);
}
