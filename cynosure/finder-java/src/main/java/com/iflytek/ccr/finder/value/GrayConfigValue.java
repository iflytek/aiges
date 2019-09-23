package com.iflytek.ccr.finder.value;

import org.codehaus.jackson.annotate.JsonProperty;

import java.util.List;

public class GrayConfigValue {
    @JsonProperty("group_id")
    private String groupId;

    @JsonProperty("server_list")
    private List<String> serverList;
    
    public String getGroupId() {
        return groupId;
    }

    public void setGroupId(String groupId) {
        this.groupId = groupId;
    }

    public List<String> getServerList() {
        return serverList;
    }

    public void setServerList(List<String> serverList) {
        this.serverList = serverList;
    }

    @Override
    public String toString() {
        return "GrayConfigValue{" +
                "groupId='" + groupId + '\'' +
                ", serverList=" + serverList +
                '}';
    }
}
