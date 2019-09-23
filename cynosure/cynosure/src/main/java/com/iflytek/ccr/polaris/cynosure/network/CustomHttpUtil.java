package com.iflytek.ccr.polaris.cynosure.network;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import org.apache.http.HttpEntity;
import org.apache.http.NameValuePair;
import org.apache.http.client.fluent.Request;
import org.apache.http.entity.ContentType;
import org.apache.http.message.BasicNameValuePair;
import org.apache.http.util.EntityUtils;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Map;

/**
 * 自定义http
 *
 * @author sctang2
 * @create 2017-12-18 11:29
 **/
@Component
public class CustomHttpUtil {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(CustomHttpUtil.class);
    /**
     * multipart/mixed post请求
     *
     * @param url
     * @param paramsList
     * @return
     * @throws Exception
     */
    public static String doPostByMultipartMixed(String url, List<CustomHttpParams> paramsList) throws Exception {
        String BOUNDARY = java.util.UUID.randomUUID().toString();
        String PREFIX = "--";
        String LINEND = "\r\n";
        if (null == paramsList || paramsList.isEmpty()) {
            return null;
        }
        byte[] totalByte = null;
        String data = null;
        int size = paramsList.size();
        for (int i = 0; i < size; i++) {
            CustomHttpParams httpParam = paramsList.get(i);
            Map<String, Object> map = httpParam.getMap();
            byte[] bt = httpParam.getBt();
            String joinMapStr = joinParams(map);
            byte[] mapBt = joinMapStr.getBytes("UTF-8");
            if (0 == i) {
                data = PREFIX + BOUNDARY + LINEND + "Content-Type:application/octet-stream" + LINEND + "Content-Disposition:param;" + joinMapStr + LINEND + "Content-Length:" + mapBt.length + LINEND + LINEND;
                totalByte = byteMerge(data.getBytes(), bt);
            } else {
                data = LINEND + PREFIX + BOUNDARY + LINEND + "Content-Type:application/octet-stream" + LINEND + "Content-Disposition:param;" + joinMapStr + LINEND + "Content-Length:" + mapBt.length + LINEND + LINEND;
                totalByte = byteMerge(totalByte, data.getBytes());
                totalByte = byteMerge(totalByte, bt);
            }
        }
        data = LINEND + PREFIX + BOUNDARY + PREFIX;
        totalByte = byteMerge(totalByte, data.getBytes());
        NameValuePair nameValuePair = new BasicNameValuePair("boundary", BOUNDARY);
        HttpEntity entity = Request.Post(url).connectTimeout(5000)
                .bodyByteArray(totalByte, ContentType.create("multipart/mixed", nameValuePair))
                .execute().returnResponse().getEntity();
        return EntityUtils.toString(entity, "utf-8");
    }

    /**
     * 拼接参数
     *
     * @param params
     * @return
     */
    public static String joinParams(Map<String, Object> params) {
        if (params == null || params.size() == 0) {
            return "";
        }
        StringBuilder pas = new StringBuilder();
        for (Map.Entry<String, Object> entry : params.entrySet()) {
            pas.append(entry.getKey());
            pas.append("=");
            pas.append(entry.getValue().toString());
            pas.append(";");
        }
        return pas.substring(0, pas.length() - 1);
    }

    /**
     * 字节流merge
     *
     * @param byte_1
     * @param byte_2
     * @return
     */
    private static byte[] byteMerge(byte[] byte_1, byte[] byte_2) {
        byte[] byte_3 = new byte[byte_1.length + byte_2.length];
        System.arraycopy(byte_1, 0, byte_3, 0, byte_1.length);
        System.arraycopy(byte_2, 0, byte_3, byte_1.length, byte_2.length);
        return byte_3;
    }
}
