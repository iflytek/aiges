package com.iflytek.ccr.polaris.cynosure.aop;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.serializer.SerializerFeature;
import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.network.HttpUtil;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.springframework.stereotype.Component;
import org.springframework.web.context.request.RequestAttributes;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;
import org.springframework.web.multipart.MaxUploadSizeExceededException;

import javax.servlet.http.HttpServletRequest;
import java.util.Map;

/**
 * 控制器时间监控AOP
 *
 * @author scta ng2
 * @create 2017-11-09 16:40
 **/
@Aspect
@Component
public class ControllerTimeAdvice extends BaseTimeAdvice {
    //private static final Logger logger = LoggerFactory.getLogger(ControllerAspect.class);
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ControllerTimeAdvice.class);

    /**
     * 环绕通知
     *
     * @param pjp
     * @return
     * @throws Throwable
     */
    @Around("execution(* com.iflytek.ccr.polaris.cynosure.controller..*.*(..)))")
    public Object doAround(ProceedingJoinPoint pjp) throws Throwable {
        //记录开始时间
        long startTimeMillis = System.currentTimeMillis();

        RequestAttributes ra = RequestContextHolder.getRequestAttributes();
        ServletRequestAttributes sra = (ServletRequestAttributes) ra;
        HttpServletRequest request = sra.getRequest();
        //获取请求头
        Map<String, String> headerNameMap = HttpUtil.findHeaderNameMap(request);
        String requestPath = request.getRequestURI();

        //获取请求参数
        String params = this.getParams(pjp);

        // result的值就是被拦截方法的返回值
        Object result;
        try {
            result = pjp.proceed();
            long endTimeMillis = System.currentTimeMillis();
            logger.info(requestPath + " cost time " + (endTimeMillis - startTimeMillis) + "ms" + "\nclient header " + headerNameMap + "\nclient request " + params + "\nserver response " + JSON.toJSONString(result, SerializerFeature.WriteMapNullValue));
        } catch (Exception ex) {
            if (ex instanceof MaxUploadSizeExceededException) {
                result = new Response<String>(SystemErrCode.ERRCODE_FILE_TOO_BIG, SystemErrCode.ERRMSG_FILE_TOO_BIG);
            } else {
                result = new Response<String>(SystemErrCode.ERRCODE_REQUEST_FAIL, SystemErrCode.ERRMSG_REQUEST_FAIL);
            }
            String error = "";
            StackTraceElement[] trace = ex.getStackTrace();
            for (StackTraceElement s : trace) {
                error += "\tat " + s + "\r\n";
            }
            //记录错误日志
            logger.error(requestPath + "\nclient header " + headerNameMap + "\nclient request " + params + "\nserver response " + JSON.toJSONString(result, SerializerFeature.WriteMapNullValue) + "\nerror " + error + "\n" + ex);
        }
        return result;
    }
}
