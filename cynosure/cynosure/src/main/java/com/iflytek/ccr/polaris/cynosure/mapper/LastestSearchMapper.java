package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.LastestSearch;

/**
 * 最近搜索持久化接口
 *
 * @author sctang2
 * @create 2018-01-26 10:18
 **/
public interface LastestSearchMapper {
    /**
     * 新增
     *
     * @param lastestSearch
     * @return
     */
    int insert(LastestSearch lastestSearch);

    /**
     * 更新
     *
     * @param lastestSearch
     * @return
     */
    int update(LastestSearch lastestSearch);

    /**
     * 查询
     *
     * @param lastestSearch
     * @return
     */
    LastestSearch find(LastestSearch lastestSearch);
}
