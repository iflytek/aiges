package com.iflytek.ccr.finder.utils;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.value.BootConfig;
import com.iflytek.ccr.finder.value.CommonResult;
import com.iflytek.ccr.finder.value.ErrorCode;
import org.apache.http.client.fluent.Request;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.net.URL;
import java.net.URLConnection;
import java.util.List;
import java.util.Map;

public class RemoteUtil {

    private static final Logger logger = LoggerFactory.getLogger(RemoteUtil.class);

    /**
     * 向指定URL发送GET方法的请求
     *
     * @param url 发送请求的URL
     *            请求参数，请求参数应该是 name1=value1&name2=value2 的形式。
     * @return URL 所代表远程资源的响应结果
     */
    public static String sendGet(String url) {
        String result = "";
        BufferedReader in = null;
        try {
            URL realUrl = new URL(url);
            // 打开和URL之间的连接
            URLConnection connection = realUrl.openConnection();
            // 设置通用的请求属性
            connection.setRequestProperty("accept", "*/*");
            connection.setRequestProperty("connection", "Keep-Alive");
            connection.setRequestProperty("user-agent",
                    "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1;SV1)");
            connection.setConnectTimeout(15000);
            // 建立实际的连接
            connection.connect();
            // 获取所有响应头字段
            Map<String, List<String>> map = connection.getHeaderFields();
            // 定义 BufferedReader输入流来读取URL的响应
            in = new BufferedReader(new InputStreamReader(
                    connection.getInputStream()));
            String line;
            while ((line = in.readLine()) != null) {
                result += line;
            }
        } catch (Exception e) {
            System.out.println("queryZkInfo Exception");
            logger.error("send http get error", e);
        }
        // 使用finally块来关闭输入流
        finally {
            try {
                if (in != null) {
                    in.close();
                }
            } catch (Exception e) {
                logger.error("", e);
            }
        }
        return result;
    }

    /**
     * 向指定 URL 发送POST方法的请求
     *
     * @param url   发送请求的 URL
     * @param param 请求参数，请求参数应该是 name1=value1&name2=value2 的形式。
     * @return 所代表远程资源的响应结果
     */

    public static String sendPost(String url, String param) {
        PrintWriter out = null;
        BufferedReader in = null;
        String result = "";
        try {
            URL realUrl = new URL(url);
            // 打开和URL之间的连接
            URLConnection conn = realUrl.openConnection();
            conn.setConnectTimeout(5000);
            conn.setReadTimeout(5000);
            // 设置通用的请求属性
            conn.setRequestProperty("accept", "*/*");
            conn.setRequestProperty("connection", "Keep-Alive");
            conn.setRequestProperty("user-agent",
                    "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1;SV1)");
            // 发送POST请求必须设置如下两行
            conn.setDoOutput(true);
            conn.setDoInput(true);
            // 获取URLConnection对象对应的输出流
            out = new PrintWriter(conn.getOutputStream());
            // 发送请求参数
            out.print(param);
            // flush输出流的缓冲
            out.flush();
            // 定义BufferedReader输入流来读取URL的响应
            in = new BufferedReader(
                    new InputStreamReader(conn.getInputStream()));
            String line;
            while ((line = in.readLine()) != null) {
                result += line;
            }
        } catch (Exception e) {
            logger.error("", e);
            e.printStackTrace();
        }
        //使用finally块来关闭输出流、输入流
        finally {
            try {
                if (out != null) {
                    out.close();
                }
                if (in != null) {
                    in.close();
                }
            } catch (IOException ex) {
                logger.error("", ex);
            }
        }
        return result;
    }

    /**
     * 查询zookeeper信息
     *
     * @param bootConfig
     * @return
     */
    public static CommonResult queryZkInfo(BootConfig bootConfig) {
        CommonResult commonResult = null;
        StringBuffer paramsBuffer = new StringBuffer();
        paramsBuffer.append("?").append("project=").append(bootConfig.getMeteData().getProject())
                .append("&group=").append(bootConfig.getMeteData().getGroup())
                .append("&service=").append(bootConfig.getMeteData().getService())
                .append("&version=").append(bootConfig.getMeteData().getVersion());
        String companionUrl = bootConfig.getCompanionUrl() + Constants.COMPANION_URL_PRE + paramsBuffer.toString();
        try {
            String result = Request.Get(companionUrl)
                    .socketTimeout(5000)
                    .execute().returnContent().asString();
//            如果获取不到结果重试一次
            if (null == result || result.isEmpty()) {
                result = sendGet(companionUrl);
            }
            logger.info(String.format("query_zk_info result:%s", result));
            commonResult = JacksonUtils.toObject(result, CommonResult.class);
        } catch (Exception e) {
            commonResult = new CommonResult();
            commonResult.setRet(ErrorCode.QUERY_ZK_INFO_FAIL);
            commonResult.setMsg(e.getMessage());
            logger.error("queryZkInfo error", e);
        }
        return commonResult;
    }


    public static void pushServiceFeedback(FinderManager finderManager, String pushId, String provider, String updateStatus, String loadStatus, String updateTime, String loadTime, String apiVersion,String type) {
        StringBuffer paramsBuffer = new StringBuffer();
        paramsBuffer.append("?").append("project=").append(finderManager.getBootConfig().getMeteData().getProject())
                .append("&group=").append(finderManager.getBootConfig().getMeteData().getGroup())
                .append("&consumer=").append(finderManager.getBootConfig().getMeteData().getAddress())
                .append("&provider=").append(provider)
                .append("&push_id=").append(pushId)
                .append("&consumer_version=").append(finderManager.getBootConfig().getMeteData().getVersion())
                .append("&provider_version=")//.append(apiVersion)
                .append("&api_version=").append(apiVersion)
                .append("&type=").append(type)
                .append("&addr=").append(finderManager.getBootConfig().getMeteData().getAddress())
                .append("&update_status=").append(updateStatus)
                .append("&update_time=").append(updateTime)
                .append("&load_status=").append(loadStatus)
                .append("&load_time=").append(loadTime);

        ;
        String companionUrl = finderManager.getBootConfig().getCompanionUrl() + Constants.SERVICE_FEEDBACK_URL + paramsBuffer.toString();
        try {
            String result = Request.Post(companionUrl)
                    .socketTimeout(5000)
                    .execute().returnContent().asString();
            logger.info("result:" + result);
        } catch (IOException e) {
            logger.error("", e);
        }

    }

    public static String pushConfigFeedback(FinderManager finderManager, String groupId, String configName, String pushId, String updateStatus, String loadStatus, String updateTime, String loadTime) {
        StringBuffer paramsBuffer = new StringBuffer();
        paramsBuffer.append("?").append("project=").append(finderManager.getBootConfig().getMeteData().getProject())
                .append("&group=").append(finderManager.getBootConfig().getMeteData().getGroup())
                .append("&service=").append(finderManager.getBootConfig().getMeteData().getService())
                .append("&push_id=").append(pushId)
                .append("&version=").append(finderManager.getBootConfig().getMeteData().getVersion())
                .append("&config=").append(configName)
                .append("&gray_group_id=").append(groupId)
                .append("&addr=").append(finderManager.getBootConfig().getMeteData().getAddress())
                .append("&update_status=").append(updateStatus)
                .append("&update_time=").append(updateTime)
                .append("&load_status=").append(loadStatus)
                .append("&load_time=").append(loadTime);
        String companionUrl = finderManager.getBootConfig().getCompanionUrl() + Constants.CONFIG_FEEDBACK_URL + paramsBuffer.toString();
        try {
            String result = Request.Post(companionUrl)
                    .socketTimeout(5000)
                    .execute().returnContent().asString();
            logger.info("result:" + result);
            return result;
        } catch (IOException e) {
            logger.error("", e);
        }
        return null;
    }

    public static String registerServiceInfo(BootConfig bootConfig, String apiVersion) {
        StringBuffer paramsBuffer = new StringBuffer();
        paramsBuffer.append("?").append("project=").append(bootConfig.getMeteData().getProject())
                .append("&group=").append(bootConfig.getMeteData().getGroup())
                .append("&service=").append(bootConfig.getMeteData().getService())
                .append("&api_version=").append(apiVersion);
        String companionUrl = bootConfig.getCompanionUrl() + Constants.REGISTER_SERVICE_INFO_URL + paramsBuffer.toString();
        try {
            String result = Request.Post(companionUrl)
                    .socketTimeout(5000).connectTimeout(5000)
                    .execute().returnContent().asString();
//            String result = sendPost(companionUrl, paramsBuffer.substring(1));
            logger.info(result);
            return result;
        } catch (Exception e) {
            logger.error("", e);
        }
        return null;
    }

//    public static String registerService(ConfigManager configManager) {
//        StringBuffer paramsBuffer = new StringBuffer();
//        paramsBuffer.append("?").append("project=").append(configManager.getStringConfigByKey(Constants.KEY_PROJECT))
//                .append("&group=").append(configManager.getStringConfigByKey(Constants.KEY_GROUP))
//                .append("&addr=").append(configManager.getStringConfigByKey(Constants.IP_ADDR))
//                .append("&version=").append(configManager.getStringConfigByKey(Constants.KEY_VERSION))
//                .append("&service_name=").append(configManager.getStringConfigByKey(Constants.KEY_SERVICE));
//        String companionUrl = configManager.getStringConfigByKey(ConfigManager.KEY_WEBSITE_URL) + Constants.REGISTER_SERVICE_URL + paramsBuffer.toString();
//        try {
//            String result = Request.Post(companionUrl)
//                    .socketTimeout(5000)
//                    .execute().returnContent().asString();
//            logger.info("result");
//            return result;
//        } catch (IOException e) {
//            logger.error("", e);
//        }
//        return null;
//    }

//    public static String unRegisterService(ConfigManager configManager) {
//        StringBuffer paramsBuffer = new StringBuffer();
//        paramsBuffer.append("?").append("project=").append(configManager.getStringConfigByKey(Constants.KEY_PROJECT))
//                .append("&group=").append(configManager.getStringConfigByKey(Constants.KEY_GROUP))
//                .append("&addr=").append(configManager.getStringConfigByKey(Constants.IP_ADDR))
//                .append("&version=").append(configManager.getStringConfigByKey(Constants.KEY_VERSION))
//                .append("&service_name=").append(configManager.getStringConfigByKey(Constants.KEY_SERVICE));
//        String companionUrl = configManager.getStringConfigByKey(ConfigManager.KEY_WEBSITE_URL) + Constants.UNREGISTER_SERVICE_URL + paramsBuffer.toString();
//        try {
//            String result = Request.Post(companionUrl)
//                    .socketTimeout(5000)
//                    .execute().returnContent().asString();
//            logger.info("result");
//            return result;
//        } catch (IOException e) {
//            logger.error("", e);
//        }
//        return null;
//    }

}
