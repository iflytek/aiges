package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.LoadBalance;

import java.util.HashMap;
import java.util.List;

/**
 * 负载均衡持久层接口
 *
 * @author sctang2
 * @create 2017-12-07 11:05
 **/
public interface LoadBalanceMapper {
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

    LoadBalance findByEnglishName(String abbr);
}
