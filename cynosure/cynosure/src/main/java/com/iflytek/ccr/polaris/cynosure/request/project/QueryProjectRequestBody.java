package com.iflytek.ccr.polaris.cynosure.request.project;

import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;

import java.io.Serializable;

/**
 * 查询项目列表-请求
 *
 * @author sctang2
 * @create 2018-01-03 9:44
 **/
public class QueryProjectRequestBody extends BaseRequestBody implements Serializable {
    private static final long serialVersionUID = -6712040302891312193L;

    //项目名称
    private String name;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }
}
