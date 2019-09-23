package com.iflytek.ccr.polaris.cynosure.request.user;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Email;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import javax.validation.constraints.Pattern;
import java.io.Serializable;

/**
 * 编辑用户-请求
 *
 * @author sctang2
 * @create 2017-11-14 15:16
 **/
public class EditUserRequestBody implements Serializable {
    private static final long serialVersionUID = 3892608139381950399L;

    //用户id
    @NotBlank(message = SystemErrCode.ERRMSG_USER_ID_NOT_NULL)
    private String id;

    //用户名
    @Length(min = 1, max = 50, message = SystemErrCode.ERRMSG_USRE_NAME_MAX_LENGTH)
    private String userName;

    //手机号
    @Length(min = 1, max = 50, message = SystemErrCode.ERRMSG_USRE_TELEPHONE_INVALID)
    @Pattern(regexp = "^1[0-9]{10}$", message = SystemErrCode.ERRMSG_USRE_TELEPHONE_INVALID)
    private String phone;

    //邮箱
    @Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_USRE_EMAIL_MAX_LENGTH)
    @Email(message = SystemErrCode.ERRMSG_USRE_EMAIL_INVALID)
    private String email;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getUserName() {
        return userName;
    }

    public void setUserName(String userName) {
        this.userName = userName;
    }

    public String getPhone() {
        return phone;
    }

    public void setPhone(String phone) {
        this.phone = phone;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    @Override
    public String toString() {
        return "EditUserRequestBody{" +
                "id='" + id + '\'' +
                ", userName='" + userName + '\'' +
                ", phone='" + phone + '\'' +
                ", email='" + email + '\'' +
                '}';
    }
}
