package com.iflytek.ccr.polaris.cynosure.util;

import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;

import java.util.HashMap;

/**
 * 分页工具
 *
 * @author sctang2
 * @create 2017-11-14 17:04
 **/
public class PagingUtil {
    /**
     * 创建分页查询条件
     *
     * @param body
     * @return
     */
    public static HashMap<String, Object> createCondition(BaseRequestBody body) {
        int currentPage = body.getCurrentPage();
        int pagesize = body.getPagesize();
        int isPage = body.getIsPage();
        HashMap<String, Object> map = new HashMap<>();
        if (0 == isPage) {
            return map;
        }
        // 每次默认返回10条，最大支持100条
        if (pagesize <= 0) {
            pagesize = 10;
        } else if (pagesize > 100) {
            pagesize = 100;
        }
        map.put("pagesize", pagesize);
        if (currentPage <= 0) {
            currentPage = 1;
        }
        map.put("startIndex", (currentPage - 1) * pagesize);
        return map;
    }

    /**
     * 获取起始索引
     *
     * @param body
     * @return
     */
    public static int getStartIndex(BaseRequestBody body) {
        int currentPage = body.getCurrentPage();
        int pagesize = body.getPagesize();
        // 每次默认返回10条，最大支持100条
        if (pagesize <= 0) {
            pagesize = 10;
        } else if (pagesize > 100) {
            pagesize = 100;
        }
        if (currentPage <= 0) {
            currentPage = 1;
        }
        return (currentPage - 1) * pagesize;
    }

    /**
     * 获取结束索引
     *
     * @param body
     * @return
     */
    public static int getEndIndex(BaseRequestBody body) {
        int currentPage = body.getCurrentPage();
        int pagesize = body.getPagesize();
        // 每次默认返回10条，最大支持100条
        if (pagesize <= 0) {
            pagesize = 10;
        } else if (pagesize > 100) {
            pagesize = 100;
        }
        if (currentPage <= 0) {
            currentPage = 1;
        }
        return (currentPage - 1) * pagesize + pagesize;
    }

    /**
     * 创建分页结果
     *
     * @param body
     * @param totalCount
     * @return
     */
    public static QueryPagingListResponseBody createResult(BaseRequestBody body, int totalCount) {
        int currentPage = body.getCurrentPage();
        QueryPagingListResponseBody result = new QueryPagingListResponseBody();
        if (currentPage <= 0) {
            currentPage = 1;
        }
        result.setCurrentPage(currentPage);
        result.setTotalCount(totalCount);
        return result;
    }
}
