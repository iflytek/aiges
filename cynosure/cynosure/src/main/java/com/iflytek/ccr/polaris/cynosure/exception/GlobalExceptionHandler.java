package com.iflytek.ccr.polaris.cynosure.exception;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.serializer.SerializerFeature;
import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.network.HttpUtil;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.multipart.MultipartException;

import javax.servlet.http.HttpServletRequest;
import java.util.Map;

/**
 * 全局异常捕获
 *
 * @author sctang2
 * @create 2017-11-09 15:55
 **/
@ControllerAdvice
@ResponseBody
public class GlobalExceptionHandler {
    //private static final Logger logger = LoggerFactory.getLogger(GlobalExceptionHandler.class);
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(GlobalExceptionHandler.class);

    /**
     * 获取全局异常
     *
     * @param request
     * @param exception
     * @return
     * @throws Exception
     */
    @ExceptionHandler(value = Exception.class)
    public Response allExceptionHandler(HttpServletRequest request, Exception exception) throws Exception {
        //获取请求头
        Map<String, String> headerNameMap = HttpUtil.findHeaderNameMap(request);
        String requestPath = request.getRequestURI();
        //创建返回结果
        boolean isPrint = true;
        Response<String> result;
        if (exception instanceof org.springframework.web.servlet.NoHandlerFoundException) {
            result = new Response<>(SystemErrCode.ERRCODE_NOT_FOUND_API, SystemErrCode.ERRMSG_NOT_FOUND_API);
            isPrint = false;
        } else if (exception instanceof MultipartException) {
            result = new Response<>(SystemErrCode.ERRCODE_FILE_TOO_BIG, SystemErrCode.ERRMSG_FILE_TOO_BIG);
        } else {
            result = new Response<>(SystemErrCode.ERRCODE_REQUEST_FAIL, SystemErrCode.ERRMSG_REQUEST_FAIL);
        }
        //记录错误日志
        if (isPrint) {
            String error = "";
            StackTraceElement[] trace = exception.getStackTrace();
            for (StackTraceElement s : trace) {
                error += "\tat " + s + "\r\n";
            }
            logger.error(requestPath + "\nclient header " + headerNameMap + "\nserver response " + JSON.toJSONString(result, SerializerFeature.WriteMapNullValue) + "\nerror " + error + "\n" + exception);
        } else {
            logger.error(requestPath + "\nclient header " + headerNameMap + "\nserver response " + JSON.toJSONString(result, SerializerFeature.WriteMapNullValue));
        }
        return result;
    }
}
