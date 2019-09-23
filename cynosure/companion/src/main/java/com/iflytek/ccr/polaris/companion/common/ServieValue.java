package com.iflytek.ccr.polaris.companion.common;

import java.util.List;

/**
 * 服务对象
 */
public class ServieValue {

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



    @Override
    public String toString() {
        return "ServieValue{" +
                "path='" + path + '\'' +
                ", versionList=" + versionList +
                '}';
    }
}
