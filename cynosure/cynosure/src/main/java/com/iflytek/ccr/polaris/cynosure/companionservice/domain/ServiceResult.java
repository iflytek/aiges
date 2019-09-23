package com.iflytek.ccr.polaris.cynosure.companionservice.domain;

import java.io.Serializable;
import java.util.List;

/**
 * 服务结果
 *
 * @author sctang2
 * @create 2017-12-11 13:55
 **/
public class ServiceResult<T> implements Serializable {
    private static final long serialVersionUID = -1999943002807585083L;

    //总数
    private int totalCount;

    //服务明细结果
    private List<T> results;

    //统计请求失败的区域数目
    private int failureCount;

    //记录请求失败的区域名
    private List<String> failureArea;

    public int getFailureCount() {
        return failureCount;
    }

    public void setFailureCount(int failureCount) {
        this.failureCount = failureCount;
    }

    public List<String> getFailureArea() {
        return failureArea;
    }

    public void setFailureArea(List<String> failureArea) {
        this.failureArea = failureArea;
    }

    public int getTotalCount() {
        return totalCount;
    }

    public void setTotalCount(int totalCount) {
        this.totalCount = totalCount;
    }

    public List<T> getResults() {
        return results;
    }

    public void setResults(List<T> results) {
        this.results = results;
    }
}
