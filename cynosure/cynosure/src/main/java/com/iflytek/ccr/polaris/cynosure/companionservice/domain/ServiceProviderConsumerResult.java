package com.iflytek.ccr.polaris.cynosure.companionservice.domain;

/**
 * 服务提供者、消费者结果
 *
 * @author sctang2
 * @create 2017-12-11 17:20
 **/
public class ServiceProviderConsumerResult {
    //地址
    private String addr;

    //状态
    private boolean is_valid;

    //权重
    private int weight;

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr;
    }

    public boolean isIs_valid() {
        return is_valid;
    }

    public void setIs_valid(boolean is_valid) {
        this.is_valid = is_valid;
    }

    public int getWeight() {
        return weight;
    }

    public void setWeight(int weight) {
        this.weight = weight;
    }
}
