package com.iflytek.ccr.polaris.cynosure.response.servicediscovery;

import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.ServiceParam;

import java.io.Serializable;
import java.util.List;
import java.util.Map;

/**
 * 查询服务发现提供者-响应
 *
 * @author sctang2
 * @create 2017-12-05 19:01
 **/
public class QueryServiceDiscoveryProviderResponseBody implements Serializable {
    private static final long serialVersionUID = 8705609676406826317L;

    //地址
    private String addr;

    //状态
    private boolean valid;

    //权重
    private int weight;

    //存放自定义K_V参数
   private List<ServiceParam>  user;

    public List<ServiceParam> getUser() {
        return user;
    }

    public void setUser(List<ServiceParam> user) {
        this.user = user;
    }

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr;
    }

    public boolean isValid() {
        return valid;
    }

    public void setValid(boolean valid) {
        this.valid = valid;
    }

    public int getWeight() {
        return weight;
    }

    public void setWeight(int weight) {
        this.weight = weight;
    }

    @Override
    public String toString() {
        return "QueryServiceDiscoveryProviderResponseBody{" +
                "addr='" + addr + '\'' +
                ", valid=" + valid +
                ", weight=" + weight +
                '}';
    }
}
