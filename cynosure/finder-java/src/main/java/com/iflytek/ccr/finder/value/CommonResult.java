package com.iflytek.ccr.finder.value;

import java.io.Serializable;

/**
 * ret：
 * 0->成功
 * 1->来自缓存
 *
 * @param <T>
 */
public class CommonResult<T> implements Serializable {
    private int ret;
    private T data;
    private String msg;

    public int getRet() {
        return ret;
    }

    public void setRet(int ret) {
        this.ret = ret;
    }

    public T getData() {
        return data;
    }

    public void setData(T data) {
        this.data = data;
    }

    public String getMsg() {
        return msg;
    }

    public void setMsg(String msg) {
        this.msg = msg;
    }

    @Override
    public String toString() {
        return "CommonResult{" +
                "ret='" + ret + '\'' +
                ", data=" + data +
                ", msg='" + msg + '\'' +
                '}';
    }
}
