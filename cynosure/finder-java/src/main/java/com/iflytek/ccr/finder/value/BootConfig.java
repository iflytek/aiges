package com.iflytek.ccr.finder.value;

public class BootConfig {

    /**
     * companion url
     */
    private String companionUrl;

    /**
     * 缓存路径
     */
    private String cachePath;

    /**
     *
     */
    private long tickerDuration;

    /**
     * 配置默认使用本地缓存
     */
    private boolean configCache = true;

    /**
     * 服务默认使用本地缓存
     */
    private boolean serviceCache = true;

    /**
     * 单位秒
     */
    private int zkSessionTimeout;

    /**
     * 单位秒
     */
    private int zkConnectTimeout;

    private int zkMaxSleepTime;

    private int zkMaxRetryNum;

    private ServiceMeteData meteData;

    public BootConfig(String companionUrl, ServiceMeteData meteData) {
        this.companionUrl = companionUrl;
        this.meteData = meteData;
    }

    public BootConfig(String companionUrl, String cachePath, ServiceMeteData meteData) {
        this.companionUrl = companionUrl;
        this.cachePath = cachePath;
        this.meteData = meteData;
    }

    public BootConfig(String companionUrl, String cachePath, int tickerDuration, int zkSessionTimeout, int zkConnectTimeout, int zkMaxSleepTime, int zkMaxRetryNum, ServiceMeteData meteData) {
        this.companionUrl = companionUrl;
        this.cachePath = cachePath;
        this.tickerDuration = tickerDuration;
        this.zkSessionTimeout = zkSessionTimeout;
        this.zkConnectTimeout = zkConnectTimeout;
        this.zkMaxSleepTime = zkMaxSleepTime;
        this.zkMaxRetryNum = zkMaxRetryNum;
        this.meteData = meteData;
    }

    public boolean isConfigCache() {
        return configCache;
    }

    public void setConfigCache(boolean configCache) {
        this.configCache = configCache;
    }

    public boolean isServiceCache() {
        return serviceCache;
    }

    public void setServiceCache(boolean serviceCache) {
        this.serviceCache = serviceCache;
    }

    public String getCompanionUrl() {
        return companionUrl;
    }

    public void setCompanionUrl(String companionUrl) {
        this.companionUrl = companionUrl;
    }

    public String getCachePath() {
        return cachePath;
    }

    public void setCachePath(String cachePath) {
        this.cachePath = cachePath;
    }

    public long getTickerDuration() {
        return tickerDuration;
    }

    public void setTickerDuration(long tickerDuration) {
        this.tickerDuration = tickerDuration;
    }

    public int getZkSessionTimeout() {
        return zkSessionTimeout;
    }

    public void setZkSessionTimeout(int zkSessionTimeout) {
        this.zkSessionTimeout = zkSessionTimeout;
    }

    public int getZkConnectTimeout() {
        return zkConnectTimeout;
    }

    public void setZkConnectTimeout(int zkConnectTimeout) {
        this.zkConnectTimeout = zkConnectTimeout;
    }

    public int getZkMaxSleepTime() {
        return zkMaxSleepTime;
    }

    public void setZkMaxSleepTime(int zkMaxSleepTime) {
        this.zkMaxSleepTime = zkMaxSleepTime;
    }

    public int getZkMaxRetryNum() {
        return zkMaxRetryNum;
    }

    public void setZkMaxRetryNum(int zkMaxRetryNum) {
        this.zkMaxRetryNum = zkMaxRetryNum;
    }

    public ServiceMeteData getMeteData() {
        return meteData;
    }

    public void setMeteData(ServiceMeteData meteData) {
        this.meteData = meteData;
    }

    @Override
    public String toString() {
        return "BootConfig{" +
                "companionUrl='" + companionUrl + '\'' +
                ", cachePath='" + cachePath + '\'' +
                ", tickerDuration=" + tickerDuration +
                ", zkSessionTimeout=" + zkSessionTimeout +
                ", zkConnectTimeout=" + zkConnectTimeout +
                ", zkMaxSleepTime=" + zkMaxSleepTime +
                ", zkMaxRetryNum=" + zkMaxRetryNum +
                ", meteData=" + meteData +
                '}';
    }
}
