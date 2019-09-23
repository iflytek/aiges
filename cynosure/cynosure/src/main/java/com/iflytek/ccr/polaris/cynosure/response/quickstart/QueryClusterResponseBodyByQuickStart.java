package com.iflytek.ccr.polaris.cynosure.response.quickstart;

import java.io.Serializable;
import java.util.List;

/**
 * 查询集群-响应
 *
 * @author sctang2
 * @create 2018-02-09 11:09
 **/
public class QueryClusterResponseBodyByQuickStart implements Serializable {
    private static final long serialVersionUID = 102335604843153283L;

    //唯一标识
    private String id;

    //名称
    private String name;

    //服务列表
    private List<QueryServiceResponseBodyByQuickStart> children;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public List<QueryServiceResponseBodyByQuickStart> getChildren() {
        return children;
    }

    public void setChildren(List<QueryServiceResponseBodyByQuickStart> children) {
        this.children = children;
    }

    @Override
    public String toString() {
        return "QueryClusterResponseBodyByQuickStart{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", children=" + children +
                '}';
    }
}
