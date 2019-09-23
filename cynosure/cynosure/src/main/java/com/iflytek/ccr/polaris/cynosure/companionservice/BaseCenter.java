package com.iflytek.ccr.polaris.cynosure.companionservice;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;
import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.PushDetailResult;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.Result;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.response.ServiceResponse;
import com.iflytek.ccr.polaris.cynosure.domain.Region;
import com.iflytek.ccr.polaris.cynosure.exception.GlobalExceptionUtil;
import com.iflytek.ccr.polaris.cynosure.network.CustomHttpParams;
import com.iflytek.ccr.polaris.cynosure.network.CustomHttpUtil;
import com.iflytek.ccr.polaris.cynosure.network.HttpClientUtil;
import org.apache.commons.lang3.StringUtils;
import org.apache.http.HttpEntity;
import org.apache.http.HttpVersion;
import org.apache.http.client.fluent.Request;
import org.apache.http.entity.ContentType;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.io.UnsupportedEncodingException;
import java.util.*;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;

/**
 * 基础中心
 *
 * @author sctang2
 * @create 2017-12-10 20:07
 **/
@Service
public class BaseCenter {
    public static final String PUSH_CONFIG_URL               = "/cynosure/push_config";
    public static final String DELETE_DATA_URL               = "/cynosure/del_data";
    public static final String QUERY_SERVICE_PATH_URL        = "/cynosure/refresh_service";
    public static final String QUERY_SERVICE_CONFIG_PATH_URL = "/cynosure/refresh_conf";
    public static final String PUSH_CLUSTER_CONFIG_URL       = "/feedback/push_cluster_config";
    public static final String PUSH_INSTANCE_CONFIG_URL      = "/feedback/push_instance_config";
    public static final String QUERY_SERVICE_PROVIDER_URL    = "/cynosure/refresh_provider";
    public static final String QUERY_SERVICE_CONSUMER_URL    = "/cynosure/refresh_consumer";

    //2.0新增
    public static final String PUSH_GRAY_CONFIG_URL            = "/gray/push_cluster_config";
    public static final String DELETE_GRAY_GROUP_URL           = "/gray/del_gray_group";
    public static final String Query_CONFIG_CONSUMER_URL       = "/cynosure/refresh_consumer";
    public static final String PUSH_SERVICE_CLUSTER_CONFIG_URL = "/service/push_cluster_config";

    private static final EasyLogger logger = EasyLoggerFactory.getInstance(BaseCenter.class);

    @Autowired
    private CustomHttpUtil customHttpUtil;

    @Autowired
    private HttpClientUtil httpClientUtil;

    /**
     * 创建
     *
     * @param path
     * @param fileName
     * @param grayGroupId
     * @param grayServers
     * @param pushId
     * @param bt
     * @return
     */
    protected CustomHttpParams create(String path, String fileName, String grayGroupId, String grayServers, String pushId, byte[] bt) {
        CustomHttpParams params = new CustomHttpParams();
        HashMap<String, Object> map = new HashMap<>();
        try {
            map.put("path", new String(path.getBytes(), "utf-8"));
        } catch (UnsupportedEncodingException e) {
            GlobalExceptionUtil.log(e);
        }
        try {
            map.put("fileName", new String(fileName.getBytes(), "utf-8"));
        } catch (UnsupportedEncodingException e) {
            GlobalExceptionUtil.log(e);
        }
        map.put("grayGroupId", grayGroupId);

        if (StringUtils.isNotBlank(grayServers)) {
            map.put("grayServers", grayServers);
        }
        map.put("pushId", pushId);
        params.setMap(map);
        params.setBt(bt);
        return params;
    }

    /**
     * 创建
     *
     * @param pathList
     * @param pushId
     * @param btList
     * @return
     */
    protected List<CustomHttpParams> create(List<String> pathList, List<String> fileNameList, String grayGroupId, String grayServers, String pushId, List<byte[]> btList) {
        List<CustomHttpParams> customHttpParamsList = new ArrayList<>();
        int size = pathList.size();
        for (int i = 0; i < size; i++) {
            String path = pathList.get(i);
            String fileName = fileNameList.get(i);
            byte[] bt = btList.get(i);
            CustomHttpParams params = new CustomHttpParams();
            HashMap<String, Object> map = new HashMap<>();
            try {
                map.put("fileName", new String(fileName.getBytes(), "utf-8"));
                map.put("path", new String(path.getBytes(), "utf-8"));
            } catch (UnsupportedEncodingException e) {
                GlobalExceptionUtil.log(e);
            }
            map.put("grayGroupId", grayGroupId);

            if (StringUtils.isNotBlank(grayServers)) {
                map.put("grayServers", grayServers);
            }
            map.put("pushId", pushId);
            params.setMap(map);
            params.setBt(bt);
            customHttpParamsList.add(params);
        }
        return customHttpParamsList;
    }

    /**
     * 创建
     *
     * @param pathList
     * @param pushId
     * @param btList
     * @return
     */
    protected List<CustomHttpParams> create(List<String> pathList, String pushId, List<byte[]> btList) {
        List<CustomHttpParams> customHttpParamsList = new ArrayList<>();
        int size = pathList.size();
        for (int i = 0; i < size; i++) {
            String path = pathList.get(i);
            byte[] bt = btList.get(i);
            CustomHttpParams params = new CustomHttpParams();
            HashMap<String, Object> map = new HashMap<>();
            try {
                map.put("path", new String(path.getBytes(), "utf-8"));
            } catch (UnsupportedEncodingException e) {
                GlobalExceptionUtil.log(e);
            }
            map.put("pushId", pushId);
            params.setMap(map);
            params.setBt(bt);
            customHttpParamsList.add(params);
        }
        return customHttpParamsList;
    }

    /**
     * 创建
     *
     * @param path
     * @param pushId
     * @param bt
     * @return
     */
    protected CustomHttpParams create(String path, String pushId, byte[] bt) {
        CustomHttpParams params = new CustomHttpParams();
        HashMap<String, Object> map = new HashMap<>();
        try {
            map.put("path", new String(path.getBytes(), "UTF-8"));
        } catch (UnsupportedEncodingException e) {
            GlobalExceptionUtil.log(e);
        }
        map.put("pushId", pushId);
        params.setMap(map);
        params.setBt(bt);
        return params;
    }

    /**
     * 创建
     *
     * @param path
     * @param pushId
     * @param bt
     * @return
     */
    protected CustomHttpParams create(String path, String grayGroupId, String pushId, byte[] bt) {
        CustomHttpParams params = new CustomHttpParams();
        HashMap<String, Object> map = new HashMap<>();
        try {
            map.put("path", new String(path.getBytes(), "utf-8"));
        } catch (UnsupportedEncodingException e) {
            GlobalExceptionUtil.log(e);
        }
        map.put("pushId", pushId);
        map.put("grayGroupId", grayGroupId);
        params.setMap(map);
        params.setBt(bt);

        return params;
    }

    /**
     * 创建
     *
     * @param path
     * @param pushId
     * @param bt
     * @return
     */
    protected CustomHttpParams create(String path, String pushId, JSONArray sdk, JSONObject user, byte[] bt) {
        CustomHttpParams params = new CustomHttpParams();
        HashMap<String, Object> map = new HashMap<>();
        try {
            map.put("path", new String(path.getBytes(), "UTF-8"));
        } catch (UnsupportedEncodingException e) {
            GlobalExceptionUtil.log(e);
        }

        if (null != sdk) {
            map.put("sdk",sdk);
            logger.info("nameFormatBytes"+"---------"+ sdk);
        }
        map.put("user", user);
        map.put("pushId", pushId);
        params.setMap(map);
        params.setBt(bt);
        return params;
    }

    /**
     * 解析推送参数
     *
     * @param cacheCenterResults
     * @return
     */
    protected List<PushDetailResult> parsePush(List<Result> cacheCenterResults) {
        List<PushDetailResult> results = new ArrayList<>();
        cacheCenterResults.forEach(x -> {
            PushDetailResult result = new PushDetailResult();
            String name = x.getName();
            result.setName(name);
            try {
                ServiceResponse cacheServiceResponse = JSON.parseObject(x.getResult(), ServiceResponse.class);
                if (null == cacheServiceResponse) {
                    result.setSuccessed(-1);
                } else {
                    int ret = cacheServiceResponse.getRet();
                    if (0 == ret) {
                        result.setSuccessed(1);
                    } else {
                        result.setSuccessed(-1);
                    }
                }
            } catch (Exception ex) {
                GlobalExceptionUtil.log(ex);
                result.setSuccessed(-1);
            }
            results.add(result);
        });
        return results;
    }

    /**
     * 批量post请求
     *
     * @param path
     * @param param
     * @param regions
     * @return
     */
    protected List<Result> batchPost(String path, CustomHttpParams param, List<Region> regions) {
        List<CustomHttpParams> networkParams = new ArrayList<>();
        networkParams.add(param);
        return this.batchPost(path, networkParams, regions);
    }

    /**
     * 批量post请求
     *
     * @param path
     * @param customHttpParams
     * @param regions
     * @return
     */
    protected List<Result> batchPost(String path, List<CustomHttpParams> customHttpParams, List<Region> regions) {
        int size = regions.size();
        List<Future<Result>> httpFutureList = new ArrayList<>();
        ExecutorService pool = Executors.newFixedThreadPool(size);
        regions.forEach(x -> {
            Future<Result> httpFuture = pool.submit(() -> {
                Result result = new Result();
                String host = x.getPushUrl();
                String name = x.getName();
                result.setName(name);
                try {
                    String response = this.customHttpUtil.doPostByMultipartMixed(host + path, customHttpParams);
                    logger.info("cache server path " + host + path + ",get result:" + response);
                    result.setResult(response);
                } catch (Exception e) {
                    GlobalExceptionUtil.log(e);
                    result.setResult(null);
                }
                return result;
            });
            httpFutureList.add(httpFuture);
        });
        pool.shutdown();
        List<Result> result = this.getResult(httpFutureList);
        return result;
    }

    /**
     * 批量post请求
     *
     * @param path
     * @param maps
     * @param regions
     * @return
     */
    protected List<Result> batchPost(String path, Map<String, String> maps, List<Region> regions) {
        int size = regions.size();
        List<Future<Result>> httpFutureList = new ArrayList<>();

        //创建一份可重用固定线程数的线程池实例，以共享的无界队列方式来运行这些线程
        ExecutorService pool = Executors.newFixedThreadPool(size);
        regions.forEach(x -> {
            //使用ExecutorService执行Callable类型的任务，并将结果保存在future变量中
            Future<Result> httpFuture = pool.submit(() -> {
                Result result = new Result();
                String host = x.getPushUrl();//集群地址
                String name = x.getName();//集群名
                result.setName(name);
                try {
                    String response = this.httpClientUtil.sendHttpPost(host + path, maps);
                    logger.info("cache server path " + host + path + ",get result:" + response);
                    result.setResult(response);
                } catch (Exception e) {
                    GlobalExceptionUtil.log(e);
                    result.setResult(null);
                }
                return result;
            });
            httpFutureList.add(httpFuture);
        });
        pool.shutdown();
        List<Result> result = this.getResult(httpFutureList);
        return result;
    }

    /**
     * 批量get请求
     *
     * @param path
     * @param querys
     * @param regions
     * @return
     */
    protected List<Result> batchGet(String path, String querys, List<Region> regions) {
        int size = regions.size();
        List<Future<Result>> httpFutureList = new ArrayList<>();
        ExecutorService pool = Executors.newFixedThreadPool(size);
        regions.forEach(x -> {
            Future<Result> httpFuture = pool.submit(() -> {
                Result result = new Result();
                String host = x.getPushUrl();
                String name = x.getName();
                result.setName(name);
                try {//从compansion端获得服务下的版本数据
                    String response = this.httpClientUtil.sendHttpGet(host + path + querys);
                    logger.info("cache server path " + host + path + querys + ",get result:" + response);
                    result.setResult(response);
                } catch (Exception e) {
                    GlobalExceptionUtil.log(e);
                    result.setResult(null);
                }
                return result;
            });
            httpFutureList.add(httpFuture);
        });
        pool.shutdown();
        List<Result> result = this.getResult(httpFutureList);
        return result;
    }

    /**
     * get请求
     *
     * @param path
     * @param querys
     * @param region
     * @return
     */
    protected Result get(String path, String querys, Region region) {
        List<Region> regions = new ArrayList<>();
        regions.add(region);
        List<Result> result = this.batchGet(path, querys, regions);
        return result.get(0);
    }

    /**
     * 获取结果
     *
     * @param httpFutureList
     * @return
     */
    private List<Result> getResult(List<Future<Result>> httpFutureList) {
        List<Result> results = new ArrayList<>();
        httpFutureList.forEach(x -> {
            Result result = null;
            try {
                result = x.get();
            } catch (Exception e) {
                GlobalExceptionUtil.log(e);
            }
            results.add(result);
        });
        return results;
    }
}
