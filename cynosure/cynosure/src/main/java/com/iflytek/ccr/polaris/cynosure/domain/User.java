package com.iflytek.ccr.polaris.cynosure.domain;

import java.io.Serializable;
import java.util.Date;
import java.util.List;

/**
 * 用户模型
 *
 * @author sctang2
 * @create 2017-11-10 9:13
 **/
public class User implements Serializable {
    private static final long serialVersionUID = -6697122248628232817L;

    //用户id
    private String id;

    //账号
    private String account;

    //密码
    private String password;

    //用户名
    private String userName;

    //手机号
    private String phone;

    //邮箱
    private String email;

    //创建时间
    private Date createTime;

    //更新时间
    private Date updateTime;

    //角色类型
    private Byte roleType;

    //服务列表
    private List<Service> serviceList;

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

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getUserName() {
        return userName;
    }

    public void setUserName(String userName) {
        this.userName = userName;
    }

    public String getPhone() {
        return phone;
    }

    public void setPhone(String phone) {
        this.phone = phone;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
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

    public Byte getRoleType() {
        return roleType;
    }

    public void setRoleType(Byte roleType) {
        this.roleType = roleType;
    }

    public List<Service> getServiceList() {
        return serviceList;
    }

    public void setServiceList(List<Service> serviceList) {
        this.serviceList = serviceList;
    }

    @Override
    public String toString() {
        return "User{" +
                "id='" + id + '\'' +
                ", account='" + account + '\'' +
                ", password='" + password + '\'' +
                ", userName='" + userName + '\'' +
                ", phone='" + phone + '\'' +
                ", email='" + email + '\'' +
                ", createTime=" + createTime +
                ", updateTime=" + updateTime +
                ", roleType=" + roleType +
                ", serviceList=" + serviceList +
                '}';
    }
}
