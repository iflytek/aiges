package com.iflytek.ccr.polaris.cynosure.companionservice.domain.response;

/**
 * 服务-响应
 *
 * @author sctang2
 * @create 2017-11-24 10:23
 **/
public class ServiceNewResponse {
    //消息
    private String msg;

    //错误码
    private int ret;

    //数据
    private ServiceDataNewResponse data;

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

    public ServiceDataNewResponse getData() {
        return data;
    }

    public void setData(ServiceDataNewResponse data) {
        this.data = data;
    }
}

