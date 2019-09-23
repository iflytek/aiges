package com.iflytek.ccr.polaris.companion.common;

/**
 * 响应实体
 *
 * @author sctang2
 * @create 2017-11-09 15:59
 **/

public class WebsitResult{
    private String code;
    private String message;
    private String data;
    public String getData() {
        return data;
    }

    public void setData(String data) {
        this.data = data;
    }



    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    @Override
    public String toString() {
        return "WebsitResult{" +
                "code='" + code + '\'' +
                ", message='" + message + '\'' +
                '}';
    }
}
