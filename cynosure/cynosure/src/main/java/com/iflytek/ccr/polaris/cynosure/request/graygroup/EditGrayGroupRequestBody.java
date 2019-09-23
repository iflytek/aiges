package com.iflytek.ccr.polaris.cynosure.request.graygroup;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 编辑灰度组-请求
 *
 * @author sctang2
 * @create 2017-11-17 16:05
 **/
public class EditGrayGroupRequestBody implements Serializable {
    private static final long serialVersionUID = -1441222124991579884L;

    //灰度组id
    @NotBlank(message = SystemErrCode.ERRMSG_GRAY_GROUP_ID_NOT_NULL)
    private String id;

    //灰度组描述
    @Length(max = 500, message = SystemErrCode.ERRMSG_GRAY_GROUP_DESC_MAX_LENGTH)
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
        return "EditGrayGroupRequestBody{" +
                "id='" + id + '\'' +
                ", desc='" + desc + '\'' +
                '}';
    }
}
