package com.iflytek.ccr.polaris.cynosure.response.servicediscovery;

import java.io.Serializable;

/**
 * 查询服务发现消费者-响应
 *
 * @author sctang2
 * @create 2018-02-02 14:57
 **/
public class QueryServiceDiscoveryConsumerResponseBody implements Serializable {
    private static final long serialVersionUID = 1666312289438219963L;

    //地址
    private String addr;

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr;
    }

    @Override
    public String toString() {
        return "QueryServiceDiscoveryConsumerResponseBody{" +
                "addr='" + addr + '\'' +
                '}';
    }
}
