package com.iflytek.ccr.polaris.companion.common;

import java.util.HashMap;
import java.util.Map;

/**
 * Created by eric on 2017/11/21.
 */
public class JsonResult {
    public JsonResult() {
        ret = ErrorCode.SUCCESS;
        msg = "success";
        data = new HashMap<String, Object>();
    }

    public int getRet() {
        return ret;
    }

    public void setRet(int ret) {
        this.ret = ret;
    }

    public String getMsg() {
        return msg;
    }

    public void setMsg(String msg) {
        this.msg = msg;
    }

    public Map<String, Object> getData() {
        return data;
    }

    public void setData(Map<String, Object> data) {
        this.data = data;
    }

    private int ret;
    private String msg;
    private Map<String, Object> data;
}
