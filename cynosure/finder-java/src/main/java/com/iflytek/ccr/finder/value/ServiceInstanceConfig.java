package com.iflytek.ccr.finder.value;

import java.io.Serializable;

/**
 * sdk 服务实例 属性配置类
 */
public class ServiceInstanceConfig implements Serializable {

    private int weight;
    private boolean isValid;

    public int getWeight() {
        return weight;
    }

    public void setWeight(int weight) {
        this.weight = weight;
    }

    public boolean isValid() {
        return isValid;
    }

    public void setValid(boolean valid) {
        isValid = valid;
    }

    @Override
    public String toString() {
        return "ServiceItemConfig{" +
                "weight=" + weight +
                ", isValid=" + isValid +
                '}';
    }
}
