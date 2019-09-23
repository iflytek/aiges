package com.iflytek.ccr.polaris.cynosure.request.project;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 删除项目成员-请求
 *
 * @author sctang2
 * @create 2018-01-16 15:30
 **/
public class DeleteProjectMemberRequestBody implements Serializable {
    private static final long serialVersionUID = -544162227283621201L;

    //项目id
    @NotBlank(message = SystemErrCode.ERRMSG_PROJECT_ID_NOT_NULL)
    private String id;

    //用户id
    @NotBlank(message = SystemErrCode.ERRMSG_USER_ID_NOT_NULL)
    private String userId;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getUserId() {
        return userId;
    }

    public void setUserId(String userId) {
        this.userId = userId;
    }

    @Override
    public String toString() {
        return "DeleteProjectMemberRequestBody{" +
                "id='" + id + '\'' +
                ", userId='" + userId + '\'' +
                '}';
    }
}
