package com.iflytek.ccr.polaris.cynosure.companionservice.domain;

/**
 * 推送明细结果
 *
 * @author sctang2
 * @create 2017-12-10 22:00
 **/
public class PushDetailResult {
    //集群名称
    private String name;

    //结果
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
}
