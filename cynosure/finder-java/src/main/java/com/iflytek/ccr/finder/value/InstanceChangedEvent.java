package com.iflytek.ccr.finder.value;

import java.util.List;

/**
 * 实例变化事件
 */
public class InstanceChangedEvent {
    /**
     * 事件类型
     */
    private final Type type;

    /**
     * 服务实例列表
     */
    private List<ServiceInstance> serviceInstanceList;


    public InstanceChangedEvent(Type type, List<ServiceInstance> serviceInstanceList) {
        this.type = type;
        this.serviceInstanceList = serviceInstanceList;
    }

    public List<ServiceInstance> getServiceInstanceList() {
        return serviceInstanceList;
    }


    public Type getType() {
        return type;
    }

    @Override
    public String toString() {
        return InstanceChangedEvent.class.getSimpleName() + "{" +
                "type=" + type +
                '}';
    }

    public enum Type {
        ADD, REMVOE;
    }
}
