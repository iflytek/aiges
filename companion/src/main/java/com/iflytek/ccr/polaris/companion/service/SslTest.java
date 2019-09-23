package com.iflytek.ccr.polaris.companion.service;

import com.iflytek.ccr.polaris.companion.common.ConfigFeedBackValue;
import com.iflytek.ccr.polaris.companion.utils.JacksonUtils;
import org.apache.commons.io.IOUtils;

import java.io.OutputStreamWriter;
import java.net.URL;
import java.net.URLConnection;

public class SslTest {

    public String getRequest(String url,int timeOut) throws Exception{
        URL u = new URL(url);
        if("https".equalsIgnoreCase(u.getProtocol())){
            SslUtils.ignoreSsl();
        }
        URLConnection conn = u.openConnection();
        conn.setConnectTimeout(timeOut);
        conn.setReadTimeout(timeOut);
        return IOUtils.toString(conn.getInputStream());
    }

    public String postRequest(String urlAddress,String args,int timeOut) throws Exception{
        URL url = new URL(urlAddress);
        if("https".equalsIgnoreCase(url.getProtocol())){
            SslUtils.ignoreSsl();
        }
        URLConnection u = url.openConnection();
        u.setDoInput(true);
        u.setDoOutput(true);
        u.setConnectTimeout(timeOut);
        u.setReadTimeout(timeOut);
        OutputStreamWriter osw = new OutputStreamWriter(u.getOutputStream(), "UTF-8");
        osw.write(args);
        osw.flush();
        osw.close();
        u.getOutputStream();
        return IOUtils.toString(u.getInputStream());
    }

    public static void main(String[] args) {
        try {
            SslTest st = new SslTest();
//            String a = st.getRequest("https://10.1.87.69:8095/index.html#/main/discovery", 3000);
//            System.out.println(a);

            String url = "https://10.1.87.69:8095/api/v1/service/config/feedback";
            ConfigFeedBackValue feedBackValue = new ConfigFeedBackValue();
            feedBackValue.setGrayGroupId("aaa");
            String reqJson = JacksonUtils.toJson(feedBackValue);
            System.out.println(st.postRequest(url,reqJson,5000));
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

}

