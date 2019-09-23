package com.iflytek.ccr.polaris.cynosure.request.track;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 查询轨迹明细-请求
 *
 * @author sctang2
 * @create 2017-11-25 15:41
 **/
public class QueryTrackDetailRequestBody extends BaseRequestBody implements Serializable {
    private static final long serialVersionUID = 9017475745236605152L;

    //推送id
    @NotBlank(message = SystemErrCode.ERRMSG_TRACK_ID_NOT_NULL)
    private String pushId;

    //@NotBlank(message = SystemErrCode.ERRMSG_TRACK_ISGRAY_NOT_NULL)
    private String isGray;

    public String getPushId() {
        return pushId;
    }

    public void setPushId(String pushId) {
        this.pushId = pushId;
    }

    public String getIsGray() {
        return isGray;
    }

    public void setIsGray(String isGray) {
        this.isGray = isGray;
    }

    @Override
    public String toString() {
        return "QueryTrackDetailRequestBody{" +
                "pushId='" + pushId + '\'' +
                ", isGray='" + isGray + '\'' +
                '}';
    }
}
