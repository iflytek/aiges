package com.iflytek.ccr.polaris.cynosure.util;

import org.apache.http.HttpResponse;
import org.apache.http.HttpStatus;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.DefaultHttpClient;
import org.apache.http.params.CoreConnectionPNames;
import org.apache.http.util.EntityUtils;

public class HttpUtil {

    //该方法发送get请求
    public static String myHttpGet(String url){
        HttpClient httpClient = new DefaultHttpClient();
        //设置超时时间
        httpClient.getParams().setParameter(CoreConnectionPNames.CONNECTION_TIMEOUT, 3000);
        httpClient.getParams().setParameter(CoreConnectionPNames.SO_TIMEOUT, 3000);
        HttpGet get = new HttpGet(url);
        try{
            // 发送GET请求
            HttpResponse httpResponse = httpClient.execute(get);
            // 如果服务器成功地返回响应
            if (httpResponse.getStatusLine().getStatusCode() == HttpStatus.SC_OK)
            {
                // 获取服务器响应字符串
                String result = EntityUtils.toString(httpResponse.getEntity());
                return result;
            } else {
                // 如果服务器失败返回响应数据"error"
                return "error";
            }
        }catch(Exception e){
            // 捕获超时异常 并反馈给调用者
            e.printStackTrace();
            return "connection time out";
        }finally{
            httpClient.getConnectionManager().shutdown();
        }
    }
}
