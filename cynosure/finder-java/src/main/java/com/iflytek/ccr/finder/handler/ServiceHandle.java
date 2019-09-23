package com.iflytek.ccr.finder.handler;

import com.iflytek.ccr.finder.value.InstanceChangedEvent;

import java.util.List;

/**
 * 服务handler,服务发现回调函数
 */
public interface ServiceHandle {

    /**
     * 服务实例上面的配置信息发生变化
     *
     * @param serviceName 服务名称
     * @param instance    一般为addr：port
     * @param jsonConfig  配置信息
     */
    boolean onServiceInstanceConfigChanged(String serviceName, String instance, String jsonConfig);

    /**
     * 服务整体配置信息发生变化
     *
     * @param serviceName 服务名称
     * @param jsonConfig 配置信息
     */
    boolean onServiceConfigChanged(String serviceName, String jsonConfig);

    /**
     * 服务实例数量发生变化
     *
     * @param serviceName
     * @param eventList
     */
    boolean onServiceInstanceChanged(String serviceName, List<InstanceChangedEvent> eventList);
}
