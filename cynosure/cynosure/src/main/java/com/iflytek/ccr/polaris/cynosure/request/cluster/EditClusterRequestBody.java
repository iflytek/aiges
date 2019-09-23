package com.iflytek.ccr.polaris.cynosure.request.cluster;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 编辑集群-请求
 *
 * @author sctang2
 * @create 2017-11-16 9:14
 **/
public class EditClusterRequestBody implements Serializable {
    private static final long serialVersionUID = 1156816781317418365L;

    //集群id
    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_ID_NOT_NULL)
    private String id;

    //集群描述
    @Length(max = 500, message = SystemErrCode.ERRMSG_CLUSTER_DESC_MAX_LENGTH)
    private String desc;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getDesc() {
        return desc;
    }

    public void setDesc(String desc) {
        this.desc = desc;
    }

    @Override
    public String toString() {
        return "EditClusterRequestBody{" +
                "id='" + id + '\'' +
                ", desc='" + desc + '\'' +
                '}';
    }
}
