package com.iflytek.ccr.polaris.cynosure.network;

import com.iflytek.ccr.polaris.cynosure.exception.GlobalExceptionUtil;

import javax.servlet.http.HttpServletRequest;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.nio.charset.Charset;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.Map;

/**
 * http工具包
 *
 * @author sctang2
 * @create 2017-12-18 11:50
 **/
public class HttpUtil {
    /**
     * 获取http header信息
     *
     * @param request
     * @return
     */
    public static Map<String, String> findHeaderNameMap(HttpServletRequest request) {
        Map<String, String> headerNameMap = new HashMap<String, String>();
        Enumeration<String> headerNames = request.getHeaderNames();
        while (headerNames.hasMoreElements()) {
            String key = (String) headerNames.nextElement();
            String value = request.getHeader(key);
            headerNameMap.put(key, value);
        }
        return headerNameMap;
    }

    /**
     * 获取请求参数
     *
     * @param request
     * @return
     */
    public static String findParameters(HttpServletRequest request) {
        String contentType = request.getContentType();
        String result = "";
        if (null == contentType) {
            return result;
        }
        if (contentType.toLowerCase().contains("application/x-www-form-urlencoded")) {
            //当为application/x-www-form-urlencoded协议
            result = getBodyString(request);
        } else if (contentType.toLowerCase().contains("application/json")) {
            //当为application/json协议
        } else if (contentType.toLowerCase().contains("multipart/form-data")) {
            //当为multipart/form-data协议
            result = getBodyString(request);
        }
        return result;
    }

    /**
     * 当为application/x-www-form-urlencoded、multipart/form-data协议
     *
     * @param request
     * @return
     */
    private static String getBodyString(HttpServletRequest request) {
        String result = "";
        Map<String, String[]> params = request.getParameterMap();
        if (!params.isEmpty()) {
            for (String key : params.keySet()) {
                String[] values = params.get(key);
                for (int i = 0; i < values.length; i++) {
                    String value = values[i];
                    result += key + "=" + value + "&";
                }
            }
            // 去掉最后一个空格
            result = result.substring(0, result.length() - 1);
        }
        return result;
    }

    /**
     * 当为application/json协议
     *
     * @param request
     * @return
     */
    private static String getBodyStringForJson(HttpServletRequest request) {
        StringBuilder sb = new StringBuilder();
        InputStream inputStream = null;
        BufferedReader reader = null;
        try {
            inputStream = request.getInputStream();
            reader = new BufferedReader(new InputStreamReader(inputStream, Charset.forName("UTF-8")));
            String line = "";
            while ((line = reader.readLine()) != null) {
                sb.append(line);
            }
        } catch (IOException e) {
            GlobalExceptionUtil.log(e);
        } finally {
            if (inputStream != null) {
                try {
                    inputStream.close();
                } catch (IOException e) {
                    GlobalExceptionUtil.log(e);
                }
            }
            if (reader != null) {
                try {
                    reader.close();
                } catch (IOException e) {
                    GlobalExceptionUtil.log(e);
                }
            }
        }
        return sb.toString();
    }
}
