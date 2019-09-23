package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.domain.Region;
import com.iflytek.ccr.polaris.cynosure.request.region.AddRegionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.region.EditRegionRequestBody;

import java.util.HashMap;
import java.util.List;

/**
 * 区域条件接口
 *
 * @author sctang2
 * @create 2017-12-09 15:08
 **/
public interface IRegionCondition {
    /**
     * 查询总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询列表
     *
     * @param map
     * @return
     */
    List<Region> findList(HashMap<String, Object> map);

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
     * 新增区域
     *
     * @param body
     * @return
     */
    Region add(AddRegionRequestBody body);

    /**
     * 根据id更新区域
     *
     * @param id
     * @param body
     * @return
     */
    Region updateById(String id, EditRegionRequestBody body);

    /**
     * 根据id删除区域
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 通过区域ids查询区域
     *
     * @param regionIds
     * @return
     */
    List<Region> findListByIds(List<String> regionIds);
}
