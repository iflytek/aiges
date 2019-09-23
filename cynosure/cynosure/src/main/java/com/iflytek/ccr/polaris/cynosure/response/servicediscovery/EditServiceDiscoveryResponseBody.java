package com.iflytek.ccr.polaris.cynosure.response.servicediscovery;

import java.io.Serializable;

/**
 * 编辑服务发现-响应
 *
 * @author sctang2
 * @create 2017-12-13 16:09
 **/
public class EditServiceDiscoveryResponseBody implements Serializable {
    private static final long serialVersionUID = 5770486045168599535L;

    //推送id
    private String pushId;

    public String getPushId() {
        return pushId;
    }

    public void setPushId(String pushId) {
        this.pushId = pushId;
    }

    @Override
    public String toString() {
        return "EditServiceDiscoveryResponseBody{" +
                "pushId='" + pushId + '\'' +
                '}';
    }
}
