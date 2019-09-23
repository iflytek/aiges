package com.iflytek.ccr.finder.value;

import org.apache.curator.framework.recipes.cache.NodeCache;
import org.apache.curator.framework.recipes.cache.TreeCache;

/**
 * 监控对象
 */
public class MonitorValue {
    private String confPath;
    private String routePath;
    private String providerPath;
    private NodeCache confNodeCache;
    private NodeCache routeNodeCache;
    private TreeCache providerTreeCache;

    public String getConfPath() {
        return confPath;
    }

    public void setConfPath(String confPath) {
        this.confPath = confPath;
    }

    public String getRoutePath() {
        return routePath;
    }

    public void setRoutePath(String routePath) {
        this.routePath = routePath;
    }

    public String getProviderPath() {
        return providerPath;
    }

    public void setProviderPath(String providerPath) {
        this.providerPath = providerPath;
    }

    public NodeCache getConfNodeCache() {
        return confNodeCache;
    }

    public void setConfNodeCache(NodeCache confNodeCache) {
        this.confNodeCache = confNodeCache;
    }

    public NodeCache getRouteNodeCache() {
        return routeNodeCache;
    }

    public void setRouteNodeCache(NodeCache routeNodeCache) {
        this.routeNodeCache = routeNodeCache;
    }

    public TreeCache getProviderTreeCache() {
        return providerTreeCache;
    }

    public void setProviderTreeCache(TreeCache providerTreeCache) {
        this.providerTreeCache = providerTreeCache;
    }

    @Override
    public String toString() {
        return "MonitorValue{" +
                "confPath='" + confPath + '\'' +
                ", routePath='" + routePath + '\'' +
                ", providerPath='" + providerPath + '\'' +
                ", confNodeCache=" + confNodeCache +
                ", routeNodeCache=" + routeNodeCache +
                ", providerTreeCache=" + providerTreeCache +
                '}';
    }
}
