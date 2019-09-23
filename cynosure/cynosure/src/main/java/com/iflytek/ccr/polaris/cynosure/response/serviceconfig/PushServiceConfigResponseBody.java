package com.iflytek.ccr.polaris.cynosure.response.serviceconfig;

import com.iflytek.ccr.polaris.cynosure.response.track.TrackConfigResponseBody;

import java.io.Serializable;
import java.util.List;

/**
 * 推送服务配置-响应
 *
 * @author sctang2
 * @create 2018-02-07 9:54
 **/
public class PushServiceConfigResponseBody extends TrackConfigResponseBody implements Serializable {
    private static final long serialVersionUID = 3666359628044640640L;

    //配置历史列表
    private List<ServiceConfigHistoryResponseBody> histories;

    public List<ServiceConfigHistoryResponseBody> getHistories() {
        return histories;
    }

    public void setHistories(List<ServiceConfigHistoryResponseBody> histories) {
        this.histories = histories;
    }

    @Override
    public String toString() {
        return "PushServiceConfigResponseBody{" +
                "histories=" + histories +
                '}';
    }
}
