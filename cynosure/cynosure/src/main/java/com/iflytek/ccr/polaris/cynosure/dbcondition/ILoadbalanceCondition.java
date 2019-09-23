package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.domain.LoadBalance;

import java.util.HashMap;
import java.util.List;

/**
 * 负载均衡条件接口
 *
 * @author sctang2
 * @create 2017-12-11 10:46
 **/
public interface ILoadbalanceCondition {
    /**
     * 查询负载均衡总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询负载均衡列表
     *
     * @param map
     * @return
     */
    List<LoadBalance> findList(HashMap<String, Object> map);

    /**
     * 英文名查询负载均衡名称
     * @param abbr
     * @return
     */
    LoadBalance findByEnglishName(String abbr);
}
