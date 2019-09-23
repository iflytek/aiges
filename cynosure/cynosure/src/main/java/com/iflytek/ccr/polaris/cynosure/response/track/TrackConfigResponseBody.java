package com.iflytek.ccr.polaris.cynosure.response.track;

import java.io.Serializable;
import java.util.Date;
import java.util.List;

/**
 * 轨迹配置列表-响应
 *
 * @author sctang2
 * @create 2017-11-24 11:53
 **/
public class TrackConfigResponseBody implements Serializable {
    private static final long serialVersionUID = -3732083827373434360L;

    //推送id
    private String id;

    //配置列表
    private List<TrackConfig> configs;

    //区域列表
    private List<TrackRegion> regions;

    //推送时间
    private Date pushTime;

    //灰度配置组字段，若为0，则表示非灰度配置，其他值表示灰度配置
    private String grayGroupId;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public List<TrackConfig> getConfigs() {
        return configs;
    }

    public void setConfigs(List<TrackConfig> configs) {
        this.configs = configs;
    }

    public List<TrackRegion> getRegions() {
        return regions;
    }

    public void setRegions(List<TrackRegion> regions) {
        this.regions = regions;
    }

    public Date getPushTime() {
        return pushTime;
    }

    public void setPushTime(Date pushTime) {
        this.pushTime = pushTime;
    }

    public String getGrayGroupId() {
        return grayGroupId;
    }

    public void setGrayGroupId(String grayGroupId) {
        this.grayGroupId = grayGroupId;
    }

    @Override
    public String toString() {
        return "TrackConfigResponseBody{" +
                "id='" + id + '\'' +
                ", configs=" + configs +
                ", regions=" + regions +
                ", pushTime=" + pushTime +
                '}';
    }
}
