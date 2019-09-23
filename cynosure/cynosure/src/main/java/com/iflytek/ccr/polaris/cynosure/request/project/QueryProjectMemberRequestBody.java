package com.iflytek.ccr.polaris.cynosure.request.project;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 查询项目成员列表-请求
 *
 * @author sctang2
 * @create 2018-01-16 15:26
 **/
public class QueryProjectMemberRequestBody extends BaseRequestBody implements Serializable {
    private static final long serialVersionUID = -6987774150125297428L;

    //项目id
    @NotBlank(message = SystemErrCode.ERRMSG_PROJECT_ID_NOT_NULL)
    private String id;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    @Override
    public String toString() {
        return "QueryProjectMemberRequestBody{" +
                "id='" + id + '\'' +
                '}';
    }
}
