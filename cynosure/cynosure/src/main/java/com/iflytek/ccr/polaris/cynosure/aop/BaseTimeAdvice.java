package com.iflytek.ccr.polaris.cynosure.aop;

import org.aspectj.lang.ProceedingJoinPoint;
import org.springframework.stereotype.Component;
import org.springframework.validation.BindingResult;

/**
 * 时间监控基类
 *
 * @author sctang2
 * @create 2018-03-05 11:21
 **/
@Component
public class BaseTimeAdvice {
    /**
     * 获取请 求参数
     *
     * @param pjp
     * @return
     */
    protected String getParams(ProceedingJoinPoint pjp) {
        String params = "";
        Object[] args = pjp.getArgs();
        if (null != args && args.length > 0) {
            for (Object arg : args) {
                if (arg instanceof BindingResult) {
                    continue;
                }
                if (null == arg) {
                    continue;
                }
                params = params + arg.toString();
            }
        }
        return params;
    }
}
