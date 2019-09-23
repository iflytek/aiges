package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.Region;

import java.util.HashMap;
import java.util.List;

/**
 * 区域持久层接口
 *
 * @author sctang2
 * @create 2017-11-14 20:06
 **/
public interface RegionMapper {
    /**
     * 新增区域
     *
     * @param region
     * @return
     */
    int insert(Region region);

    /**
     * 删除区域
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 更新区域
     *
     * @param region
     * @return
     */
    int updateById(Region region);

    /**
     * 根据id查询区域信息
     *
     * @param id
     * @return
     */
    Region findById(String id);

    /**
     * 通过区域名称查询区域信息
     *
     * @param name
     * @return
     */
    Region findByName(String name);

    /**
     * 通过区域companion查询区域信息
     *
     * @param pushUrl
     * @return
     */
    Region findByPushUrl(String pushUrl);

    /**
     * 查询区域总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询区域列表
     *
     * @param map
     * @return
     */
    List<Region> findList(HashMap<String, Object> map);
}
