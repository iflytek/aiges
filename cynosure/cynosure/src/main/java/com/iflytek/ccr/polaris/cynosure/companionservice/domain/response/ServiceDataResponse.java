package com.iflytek.ccr.polaris.cynosure.companionservice.domain.response;

import java.util.List;

/**
 * 服务数据-响应
 *
 * @author sctang2
 * @create 2017-12-07 15:03
 **/
public class ServiceDataResponse {
    private List<ServiceDataDetailResponse> pathList;

    public List<ServiceDataDetailResponse> getPathList() {
        return pathList;
    }

    public void setPathList(List<ServiceDataDetailResponse> pathList) {
        this.pathList = pathList;
    }
}
