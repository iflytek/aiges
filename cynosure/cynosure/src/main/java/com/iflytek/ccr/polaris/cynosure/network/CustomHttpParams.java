package com.iflytek.ccr.polaris.cynosure.network;

import java.util.HashMap;

/**
 * 自定义http参数
 *
 * @author sctang2
 * @create 2017-12-10 20:31
 **/
public class CustomHttpParams {
    //参数map，K-V
    private HashMap<String, Object> map;

    //字节流
    private byte[] bt;

    public HashMap<String, Object> getMap() {
        return map;
    }

    public void setMap(HashMap<String, Object> map) {
        this.map = map;
    }

    public byte[] getBt() {
        return bt;
    }

    public void setBt(byte[] bt) {
        this.bt = bt;
    }
}
