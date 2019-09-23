package com.iflytek.ccr.polaris.cynosure.companionservice.domain.response;

import java.util.List;

/**
 * 服务数据明细-响应
 *
 * @author sctang2
 * @create 2017-12-07 15:13
 **/
public class ServiceDataNewDetailResponse {
    private String path;
    private List<String> versionList;

    public String getPath() {
        return path;
    }

    public void setPath(String path) {
        this.path = path;
    }

    public List<String> getVersionList() {
        return versionList;
    }

    public void setVersionList(List<String> versionList) {
        this.versionList = versionList;
    }
}
