package com.iflytek.ccr.polaris.cynosure.request.project;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 添加项目成员-请求
 *
 * @author sctang2
 * @create 2018-01-24 21:23
 **/
public class AddProjectMemberRequestBody implements Serializable {
    private static final long serialVersionUID = -544162227283621201L;

    //项目id
    @NotBlank(message = SystemErrCode.ERRMSG_PROJECT_ID_NOT_NULL)
    private String id;

    //账号
    @NotBlank(message = SystemErrCode.ERRMSG_USER_ACCOUNT_NOT_NULL)
    @Length(min = 1, max = 50, message = SystemErrCode.ERRMSG_USRE_ACCOUNT_MAX_LENGTH)
    private String account;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getAccount() {
        return account;
    }

    public void setAccount(String account) {
        this.account = account;
    }

    @Override
    public String toString() {
        return "AddProjectMemberRequestBody{" +
                "id='" + id + '\'' +
                ", account='" + account + '\'' +
                '}';
    }
}
