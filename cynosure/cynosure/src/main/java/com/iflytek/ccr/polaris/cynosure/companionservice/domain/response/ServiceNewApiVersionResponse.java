package com.iflytek.ccr.polaris.cynosure.companionservice.domain.response;

/**
 * 服务-响应
 *
 * @author sctang2
 * @create 2017-11-24 10:23
 **/
public class ServiceNewApiVersionResponse {
    //消息
    private String msg;

    //错误码
    private int ret;

    //数据
    private ServiceDataApiVersionResponse data;

    public String getMsg() {
        return msg;
    }

    public void setMsg(String msg) {
        this.msg = msg;
    }

    public int getRet() {
        return ret;
    }

    public void setRet(int ret) {
        this.ret = ret;
    }

    public ServiceDataApiVersionResponse getData() {
        return data;
    }

    public void setData(ServiceDataApiVersionResponse data) {
        this.data = data;
    }
}
