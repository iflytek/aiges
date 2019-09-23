package com.iflytek.ccr.polaris.cynosure.apicompatible;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * 接口版本标识注解
 *
 * @author sctang2
 * @create 2017-11-09 15:47
 **/
@Target({ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
public @interface ApiVersion {
    /**
     * 版本号
     *
     * @return
     */
    int value();
}
