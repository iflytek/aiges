package com.iflytek.ccr.polaris.companion.common;

import java.util.Arrays;

public class ZkDataValue {
    private int ret;
    private String desc;
    private String pushId;
    private byte[] realData;

    public int getRet() {
        return ret;
    }

    public void setRet(int ret) {
        this.ret = ret;
    }

    public String getDesc() {
        return desc;
    }

    public void setDesc(String desc) {
        this.desc = desc;
    }

    public String getPushId() {
        return pushId;
    }

    public void setPushId(String pushId) {
        this.pushId = pushId;
    }

    public byte[] getRealData() {
        return realData;
    }

    public void setRealData(byte[] realData) {
        this.realData = realData;
    }

    @Override
    public String toString() {
        return "ZkDataValue{" +
                "ret=" + ret +
                ", desc='" + desc + '\'' +
                ", pushId='" + pushId + '\'' +
                ", realData=" + Arrays.toString(realData) +
                '}';
    }
}
