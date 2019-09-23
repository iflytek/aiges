package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.customdomain.SearchCondition;

/**
 * 最近搜索条件接口
 *
 * @author sctang2
 * @create 2018-01-26 10:21
 **/
public interface ILastestSearchCondition {
    /**
     * 新增
     *
     * @param url
     * @param condition
     * @return
     */
    int insert(String url, String condition);

    /**
     * 更新
     *
     * @param url
     * @param condition
     * @return
     */
    int update(String url, String condition);

    /**
     * 查询搜索条件
     *
     * @param url
     * @return
     */
    SearchCondition find(String url);

    /**
     * 同步搜索条件
     *
     * @param url
     * @param condition
     * @return
     */
    int syncSearchCondition(String url, String condition);
}
