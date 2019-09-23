package com.iflytek.ccr.polaris.cynosure.companionservice.domain;

import java.io.Serializable;
import java.util.List;

/**
 * 服务明细结果
 *
 * @author sctang2
 * @create 2018-02-06 10:40
 **/
public class ServiceDetailNewResult implements Serializable {
    private static final long serialVersionUID = 5150271535369434190L;

    //集群名称(区域名称)
    private String name;

    //服务名
    private List<String> versionList;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public List<String> getVersionList() {
        return versionList;
    }

    public void setVersionList(List<String> versionList) {
        this.versionList = versionList;
    }

    @Override
    public String toString() {
        return "ServiceDetailNewResult{}";
    }
}
