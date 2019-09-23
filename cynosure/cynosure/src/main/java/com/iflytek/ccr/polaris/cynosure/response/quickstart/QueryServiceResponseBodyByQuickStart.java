package com.iflytek.ccr.polaris.cynosure.response.quickstart;

import java.io.Serializable;
import java.util.List;

/**
 * 查询服务-响应
 *
 * @author sctang2
 * @create 2018-02-09 11:12
 **/
public class QueryServiceResponseBodyByQuickStart implements Serializable {
    private static final long serialVersionUID = -3615101364086366506L;

    //唯一标识
    private String id;

    //名称
    private String name;

    //版本列表
    private List<QueryVersionResponseBodyByQuickStart> children;

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

    public List<QueryVersionResponseBodyByQuickStart> getChildren() {
        return children;
    }

    public void setChildren(List<QueryVersionResponseBodyByQuickStart> children) {
        this.children = children;
    }

    @Override
    public String toString() {
        return "QueryServiceResponseBodyByQuickStart{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", children=" + children +
                '}';
    }
}
