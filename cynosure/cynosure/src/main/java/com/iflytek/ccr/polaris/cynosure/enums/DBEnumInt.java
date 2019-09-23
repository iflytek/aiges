package com.iflytek.ccr.polaris.cynosure.enums;

/**
 * 数据库整型枚举
 *
 * @author sctang2
 * @create 2017-11-14 14:26
 **/
public enum DBEnumInt {
    //用户角色
    ROLE_TYPE_ADMIN("管理员", 1), ROLE_TYPE_USER("普通用户", 2),

    //是否为创建者
    CREATOR_N("否", 0), CREATOR_Y("是", 1);

    private String name;
    private int    index;

    private DBEnumInt(String name, int index) {
        this.name = name;
        this.index = index;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public int getIndex() {
        return index;
    }

    public void setIndex(int index) {
        this.index = index;
    }
}
