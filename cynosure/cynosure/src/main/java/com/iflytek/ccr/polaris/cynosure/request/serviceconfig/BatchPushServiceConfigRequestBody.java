package com.iflytek.ccr.polaris.cynosure.request.serviceconfig;

import java.io.Serializable;
import java.util.List;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

/**
 * 服务配置批量推送-请求
 *
 * @author sctang2
 * @create 2018-01-05 10:13
 **/
@ApiModel("服务配置批量推送请求参数")
public class BatchPushServiceConfigRequestBody implements Serializable {
    private static final long serialVersionUID = -7445649664488660270L;

    //配置id列表
    @ApiModelProperty("配置id列表")
    private List<String> ids;

    //区域id列表
    @ApiModelProperty("区域id列表")
    private List<String> regionIds;

    public List<String> getIds() {
        return ids;
    }

    public void setIds(List<String> ids) {
        this.ids = ids;
    }

    public List<String> getRegionIds() {
        return regionIds;
    }

    public void setRegionIds(List<String> regionIds) {
        this.regionIds = regionIds;
    }

    @Override
    public String toString() {
        return "BatchPushServiceConfigRequestBody{" +
                "ids=" + ids +
                ", regionIds=" + regionIds +
                '}';
    }
}
