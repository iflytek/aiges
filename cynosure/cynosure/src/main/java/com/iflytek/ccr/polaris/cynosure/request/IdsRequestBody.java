package com.iflytek.ccr.polaris.cynosure.request;

import java.io.Serializable;
import java.util.List;

/**
 * ids-请求
 *
 * @author sctang2
 * @create 2018-01-31 10:38
 **/
public class IdsRequestBody implements Serializable {
    private static final long serialVersionUID = 5326331893813978065L;

    //id列表
    private List<String> ids;

    public List<String> getIds() {
        return ids;
    }

    public void setIds(List<String> ids) {
        this.ids = ids;
    }

    @Override
    public String toString() {
        return "IdsRequestBody{" +
                "ids=" + ids +
                '}';
    }
}
