package com.iflytek.ccr.polaris.cynosure.companionservice.domain.response;

/**
 * 服务数据明细-响应
 *
 * @author sctang2
 * @create 2017-12-07 15:13
 **/
public class ServiceDataDetailResponse {
    private String childPath;
    private String pushId;
    private String data;

    public String getChildPath() {
        return childPath;
    }

    public void setChildPath(String childPath) {
        this.childPath = childPath;
    }

    public String getPushId() {
        return pushId;
    }

    public void setPushId(String pushId) {
        this.pushId = pushId;
    }

    public String getData() {
        return data;
    }

    public void setData(String data) {
        this.data = data;
    }
}
