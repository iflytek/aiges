package com.iflytek.ccr.polaris.companion.common;

import com.iflytek.ccr.polaris.companion.utils.JacksonUtils;
import org.codehaus.jackson.annotate.JsonProperty;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class GrayConfigValue {
    @JsonProperty("group_id")
    private String groupId;

    @JsonProperty("server_list")
    private List<String> serverList;

    public static void main(String[] args) {
        GrayConfigValue value = new GrayConfigValue();
        value.setGroupId("groupa");
        List<String> ss = new ArrayList<>();
        ss.add("123.22.2.2.2:111");
        ss.add("124.22.2.2.2:111");
        value.setServerList(ss);

        List<GrayConfigValue> list = new ArrayList<>();
        list.add(value);
        list.add(value);
        String result = JacksonUtils.toJson(list);
        System.out.println(result);
        List map = JacksonUtils.toObject(result, List.class);
        System.out.println(map.size());
        for (Object v : map) {
            Map<String, Object> vMap = (Map<String, Object>) v;
            System.out.println(vMap.get("group_id"));
            System.out.println(vMap.get("server_list"));

            GrayConfigValue vv = new GrayConfigValue();
            vv.setGroupId(vMap.get("group_id").toString());
            vv.setServerList((List<String>) vMap.get("server_list"));
            System.out.println(vv.serverList.get(1));
        }
    }

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
