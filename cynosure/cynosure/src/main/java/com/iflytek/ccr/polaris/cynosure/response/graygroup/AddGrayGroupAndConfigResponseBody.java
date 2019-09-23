package com.iflytek.ccr.polaris.cynosure.response.graygroup;

import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.ServiceConfigDetailResponseBody;

import java.io.Serializable;
import java.util.List;

/**
 * 新增灰度组和灰度配置-响应
 *
 * @author sctang2
 * @create 2018-01-30 9:13
 **/
public class AddGrayGroupAndConfigResponseBody implements Serializable {
    private static final long serialVersionUID = -5769922887741353548L;

    //灰度组
    private GrayGroupDetailResponseBody grayGroup;

    //灰度服务配置
    private List<ServiceConfigDetailResponseBody> configs;

    public GrayGroupDetailResponseBody getGrayGroup() {
        return grayGroup;
    }

    public void setGrayGroup(GrayGroupDetailResponseBody grayGroup) {
        this.grayGroup = grayGroup;
    }

    public List<ServiceConfigDetailResponseBody> getConfigs() {
        return configs;
    }

    public void setConfigs(List<ServiceConfigDetailResponseBody> configs) {
        this.configs = configs;
    }

    @Override
    public String toString() {
        return "AddGrayGroupAndConfigResponseBody{" +
                "grayGroup=" + grayGroup +
                ", configs=" + configs +
                '}';
    }
}
