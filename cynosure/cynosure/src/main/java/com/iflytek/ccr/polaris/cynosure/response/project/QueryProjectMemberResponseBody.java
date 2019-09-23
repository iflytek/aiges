package com.iflytek.ccr.polaris.cynosure.response.project;

import java.io.Serializable;
import java.util.Date;

/**
 * 查询项目成员列表-响应
 *
 * @author sctang2
 * @create 2018-01-16 17:14
 **/
public class QueryProjectMemberResponseBody implements Serializable {
    private static final long serialVersionUID = 4247296122181459710L;

    //用户id
    private String id;

    //账号
    private String account;

    //创建时间
    private Date createTime;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getAccount() {
        return account;
    }

    public void setAccount(String account) {
        this.account = account;
    }

    public Date getCreateTime() {
        return createTime;
    }

    public void setCreateTime(Date createTime) {
        this.createTime = createTime;
    }

    @Override
    public String toString() {
        return "QueryProjectMemberResponseBody{" +
                "id='" + id + '\'' +
                ", account='" + account + '\'' +
                ", createTime=" + createTime +
                '}';
    }
}
