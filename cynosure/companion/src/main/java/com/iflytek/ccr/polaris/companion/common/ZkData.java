package com.iflytek.ccr.polaris.companion.common;

public class ZkData {



    private String childPath;

    private String pushId;

    private String data;

    public String getData() {
        return data;
    }

    public void setData(String data) {
        this.data = data;
    }

    public String getPushId() {
        return pushId;
    }

    public void setPushId(String pushId) {
        this.pushId = pushId;
    }

    public String getChildPath() {
        return childPath;
    }

    public void setChildPath(String childPath) {
        this.childPath = childPath;
    }

    @Override
    public String toString() {
        return "ZkData{" +
                "childPath='" + childPath + '\'' +
                ", pushId='" + pushId + '\'' +
                ", data='" + data + '\'' +
                '}';
    }
}
