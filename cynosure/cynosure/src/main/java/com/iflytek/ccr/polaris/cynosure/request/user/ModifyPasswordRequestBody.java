package com.iflytek.ccr.polaris.cynosure.request.user;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 修改密码-请求
 *
 * @author sctang2
 * @create 2017-11-14 10:46
 **/
public class ModifyPasswordRequestBody implements Serializable {
    private static final long serialVersionUID = 5240275462716246055L;

    //旧密码
    @NotBlank(message = SystemErrCode.ERRMSG_USER_PASSWORD_NOT_NULL)
    private String oldPassword;

    //新密码
    @NotBlank(message = SystemErrCode.ERRMSG_USER_PASSWORD_NOT_NULL)
    private String newPassword;

    //确认新密码
    private String confirmPassword;

    public String getOldPassword() {
        return oldPassword;
    }

    public void setOldPassword(String oldPassword) {
        this.oldPassword = oldPassword;
    }

    public String getNewPassword() {
        return newPassword;
    }

    public void setNewPassword(String newPassword) {
        this.newPassword = newPassword;
    }

    public String getConfirmPassword() {
        return confirmPassword;
    }

    public void setConfirmPassword(String confirmPassword) {
        this.confirmPassword = confirmPassword;
    }

    @Override
    public String toString() {
        return "ModifyPasswordRequestBody{" +
                "oldPassword='" + oldPassword + '\'' +
                ", newPassword='" + newPassword + '\'' +
                ", confirmPassword='" + confirmPassword + '\'' +
                '}';
    }
}
