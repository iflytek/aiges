package com.iflytek.ccr.polaris.cynosure.response.track;

import java.io.Serializable;
import java.util.Date;
import java.util.List;

/**
 * 轨迹发现列表-响应
 *
 * @author sctang2
 * @create 2017-12-14 10:17
 **/
public class TrackDiscoveryResponseBody implements Serializable {
    private static final long serialVersionUID = -7169110602323765215L;

    //推送id
    private String id;

    //区域列表
    private List<TrackRegion> regions;

    //推送时间
    private Date pushTime;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
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
}
