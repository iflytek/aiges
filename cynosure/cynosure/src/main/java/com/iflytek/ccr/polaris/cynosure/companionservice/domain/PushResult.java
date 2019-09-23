package com.iflytek.ccr.polaris.cynosure.companionservice.domain;

import java.util.List;

/**
 * 推送结果
 *
 * @author sctang2
 * @create 2017-12-10 20:11
 **/
public class PushResult {
    //推送id
    private String pushId;

    //结果
    private String result;

    private List<PushDetailResult> data;

    public PushResult() {
    }

    public PushResult(String pushId, String result, List<PushDetailResult> data) {
        this.pushId = pushId;
        this.result = result;
        this.data = data;
    }

    public List<PushDetailResult> getData() {
        return data;
    }

    public void setData(List<PushDetailResult> data) {
        this.data = data;
    }

    public String getPushId() {
        return pushId;
    }

    public void setPushId(String pushId) {
        this.pushId = pushId;
    }

    public String getResult() {
        return result;
    }

    public void setResult(String result) {
        this.result = result;
    }
}
