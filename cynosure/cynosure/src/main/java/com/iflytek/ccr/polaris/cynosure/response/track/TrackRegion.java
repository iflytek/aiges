package com.iflytek.ccr.polaris.cynosure.response.track;

/**
 * 轨迹区域
 *
 * @author sctang2
 * @create 2018-02-02 16:37
 **/
public class TrackRegion {
    //区域名称
    private String name;

    //推送结果
    private int successed;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public int getSuccessed() {
        return successed;
    }

    public void setSuccessed(int successed) {
        this.successed = successed;
    }

    @Override
    public String toString() {
        return "TrackRegion{" +
                "name='" + name + '\'' +
                ", successed=" + successed +
                '}';
    }
}
