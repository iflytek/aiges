package com.iflytek.ccr.polaris.companion.common;

public class ConfigValue {
    private String ipAddr;
    private String zkStr;
    private int port;
    private String websiteUrl;

    public String getIpAddr() {
        return ipAddr;
    }

    public void setIpAddr(String ipAddr) {
        this.ipAddr = ipAddr;
    }

    public String getZkStr() {
        return zkStr;
    }

    public void setZkStr(String zkStr) {
        this.zkStr = zkStr;
    }

    public int getPort() {
        return port;
    }

    public void setPort(int port) {
        this.port = port;
    }

    public String getWebsiteUrl() {
        return websiteUrl;
    }

    public void setWebsiteUrl(String websiteUrl) {
        this.websiteUrl = websiteUrl;
    }

    @Override
    public String toString() {
        return "ConfigValue{" +
                "ipAddr='" + ipAddr + '\'' +
                ", zkStr='" + zkStr + '\'' +
                ", port=" + port +
                ", websiteUrl='" + websiteUrl + '\'' +
                '}';
    }
}
