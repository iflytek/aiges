package com.iflytek.ccr.polaris.cynosure.request.cluster;

import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;

import java.io.Serializable;

/**
 * 查询集群列表-请求
 *
 * @author sctang2
 * @create 2017-12-05 8:58
 **/
public class QueryClusterRequestBody extends BaseRequestBody implements Serializable {
    private static final long serialVersionUID = 5013430883596490683L;

    //项目名称
    private String project;

    public String getProject() {
        return project;
    }

    public void setProject(String project) {
        this.project = project;
    }

    @Override
    public String toString() {
        return "QueryClusterRequestBody{" +
                "project='" + project + '\'' +
                '}';
    }
}
