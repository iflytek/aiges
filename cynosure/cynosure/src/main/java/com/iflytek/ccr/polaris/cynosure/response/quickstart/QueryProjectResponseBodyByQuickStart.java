package com.iflytek.ccr.polaris.cynosure.response.quickstart;

import java.io.Serializable;
import java.util.List;

/**
 * 查询项目-响应
 *
 * @author sctang2
 * @create 2018-02-09 11:08
 **/
public class QueryProjectResponseBodyByQuickStart implements Serializable {
    private static final long serialVersionUID = 6201204205347158453L;

    //唯一标识
    private String id;

    //名称
    private String name;

    //集群列表
    private List<QueryClusterResponseBodyByQuickStart> children;

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

    public List<QueryClusterResponseBodyByQuickStart> getChildren() {
        return children;
    }

    public void setChildren(List<QueryClusterResponseBodyByQuickStart> children) {
        this.children = children;
    }

    @Override
    public String toString() {
        return "QueryProjectResponseBodyByQuickStart{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", children=" + children +
                '}';
    }
}
