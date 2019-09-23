package com.iflytek.ccr.polaris.cynosure.response.quickstart;

import java.io.Serializable;

/**
 * 查询版本-响应
 *
 * @author sctang2
 * @create 2018-02-09 11:13
 **/
public class QueryVersionResponseBodyByQuickStart implements Serializable {
    private static final long serialVersionUID = 8054008351863921899L;

    //唯一标识
    private String id;

    //名称
    private String name;

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

    @Override
    public String toString() {
        return "QueryVersionResponseBodyByQuickStart{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                '}';
    }
}
