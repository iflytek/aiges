package fluent.httpclient;

import org.apache.http.HttpResponse;
import org.apache.http.HttpStatus;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.DefaultHttpClient;
import org.apache.http.params.CoreConnectionPNames;
import org.apache.http.util.EntityUtils;

import java.util.ArrayList;
import java.util.List;

public class TestHttpClientGet {
    public static void main(String[] args){
        String url1 = "http://10.1.87.71:9999/cynosure/refresh_service?path=/polaris/service/05127d76c3a6fe7c3375562921560a20";
        String url2 = "http://10.1.87.70:6868/cynosure/refresh_consumer?path=/polaris/service/05127d76c3a6fe7c3375562921560a20";
        List<String> urls = new ArrayList<>();
//        urls.add(url1);
        urls.add(url2);
        List<String> unParseResult = getApiVersion(urls);
        for (String resultStr: unParseResult) {
            System.out.println(resultStr);
        }

    }

    public static List<String> getApiVersion(List<String> urls){
        List<String> unParseResult = new ArrayList<>();
        for (String url: urls) {
            String resultStr = getRequestFromCompansion(url);
            unParseResult.add(resultStr);
        }
        return unParseResult;
    }

    public static String getRequestFromCompansion(String url){
        HttpClient httpClient = new DefaultHttpClient();
        httpClient.getParams().setParameter(CoreConnectionPNames.CONNECTION_TIMEOUT, 6000);
        httpClient.getParams().setParameter(CoreConnectionPNames.SO_TIMEOUT, 6000);
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
