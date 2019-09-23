package com.iflytek.ccr.polaris.cynosure.response.project;

import java.io.Serializable;
import java.util.Date;

/**
 * 项目明细-响应
 *
 * @author sctang2
 * @create 2017-11-19 21:08
 **/
public class ProjectDetailResponseBody implements Serializable {
    private static final long serialVersionUID = 5530617498093517305L;

    //项目id
    private String id;

    //项目名称
    private String name;

    //项目描述
    private String desc;

    //创建时间
    private Date createTime;

    //更新时间
    private Date updateTime;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDesc() {
        return desc;
    }

    public void setDesc(String desc) {
        this.desc = desc;
    }

    public Date getCreateTime() {
        return createTime;
    }

    public void setCreateTime(Date createTime) {
        this.createTime = createTime;
    }

    public Date getUpdateTime() {
        return updateTime;
    }

    public void setUpdateTime(Date updateTime) {
        this.updateTime = updateTime;
    }

    @Override
    public String toString() {
        return "ProjectDetailResponseBody{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", desc='" + desc + '\'' +
                ", createTime=" + createTime +
                ", updateTime=" + updateTime +
                '}';
    }
}
