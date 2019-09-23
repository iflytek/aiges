package com.iflytek.ccr.polaris.cynosure.companionservice.domain;

/**
 * 结果
 *
 * @author sctang2
 * @create 2017-12-10 21:28
 **/
public class Result {
    //名称(区域名称)
    private String name;

    //返回结果
    private String result;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getResult() {
        return result;
    }

    public void setResult(String result) {
        this.result = result;
    }

    @Override
    public String toString() {
        return "Result{" +
                "name='" + name + '\'' +
                ", result='" + result + '\'' +
                '}';
    }
}
