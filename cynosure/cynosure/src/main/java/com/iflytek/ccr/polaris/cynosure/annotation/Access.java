package com.iflytek.ccr.polaris.cynosure.annotation;

import java.lang.annotation.*;

/**
 * 自定义权限 注解
 *
 * @author sctang2
 * @create 2017-11-15 10:41
 **/
@Target(ElementType.METHOD)
@Retention(RetentionPolicy.RUNTIME)
@Documented
public @interface Access {
    String[] value() default {};

    String[] authorities() default {};

    String[] roles() default {};
}
