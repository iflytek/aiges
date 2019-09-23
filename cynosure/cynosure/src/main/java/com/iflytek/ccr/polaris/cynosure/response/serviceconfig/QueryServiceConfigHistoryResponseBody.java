package com.iflytek.ccr.polaris.cynosure.response.serviceconfig;

import java.io.Serializable;
import java.util.Date;
import java.util.List;

/**
 * 查询配置历史列表-响应
 *
 * @author sctang2
 * @create 2018-03-10 20:07
 **/
public class QueryServiceConfigHistoryResponseBody implements Serializable {
    private static final long serialVersionUID = 1672438611332020325L;

    //推送版本号
    private String pushVersion;

    //创建时间
    private Date createTime;

    //配置历史列表
    private List<ServiceConfigHistoryResponseBody> histories;

    public String getPushVersion() {
        return pushVersion;
    }

    public void setPushVersion(String pushVersion) {
        this.pushVersion = pushVersion;
    }

    public Date getCreateTime() {
        return createTime;
    }

    public void setCreateTime(Date createTime) {
        this.createTime = createTime;
    }

    public List<ServiceConfigHistoryResponseBody> getHistories() {
        return histories;
    }

    public void setHistories(List<ServiceConfigHistoryResponseBody> histories) {
        this.histories = histories;
    }

    @Override
    public String toString() {
        return "QueryServiceConfigHistoryResponseBody{" +
                "pushVersion='" + pushVersion + '\'' +
                ", createTime=" + createTime +
                ", histories=" + histories +
                '}';
    }
}
