package com.iflytek.ccr.polaris.cynosure.response.track;

/**
 * 轨迹配置
 *
 * @author sctang2
 * @create 2018-02-02 16:33
 **/
public class TrackConfig {
    //服务配置名称
    private String name;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    @Override
    public String toString() {
        return "TrackConfig{" +
                "name='" + name + '\'' +
                '}';
    }
}
