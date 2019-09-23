package com.iflytek.ccr.polaris.cynosure.aop;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.serializer.SerializerFeature;
import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.springframework.stereotype.Component;

import java.util.Arrays;

/**
 * 方法时间监控AOP
 *
 * @author sctang2
 * @create 2018-03-05 10:50
 **/
@Aspect
@Component
public class MethodTimeAdvice extends BaseTimeAdvice {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(MethodTimeAdvice.class);

    /**
     * 环绕通知
     *
     * @param pjp
     * @return
     * @throws Throwable
     */
    @Around("execution(* com.iflytek.ccr.polaris.cynosure.dbcondition..*.*(..)))")
    public Object doAround(ProceedingJoinPoint pjp) throws Throwable {
        long startTimeMillis = System.currentTimeMillis();

        try {
            Object result = pjp.proceed();
            long endTimeMillis = System.currentTimeMillis();
            logger.info(pjp.getTarget().getClass().getName() + "." + pjp.getSignature().getName() + " " + Arrays.toString(pjp.getArgs()) + " take time " + (endTimeMillis - startTimeMillis) + "ms" + "\nresult " + JSON.toJSONString(result, SerializerFeature.WriteMapNullValue));
            return result;
        } catch (Throwable e) {
            long endTimeMillis = System.currentTimeMillis();
            logger.error(pjp.getTarget().getClass().getName() + "." + pjp.getSignature().getName() + " " + Arrays.toString(pjp.getArgs()) + " take time " + (endTimeMillis - startTimeMillis) + "ms" + "\nresult " + e);
            throw e;
        }
    }
}
