package com.iflytek.ccr.polaris.cynosure.util;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

/**
 * 配置属性读取
 *
 * @author sctang2
 * @create 2017-12-10 19:02
 **/
@Component
public class PropUtil {
    @Value("${configPath}")
    public String CONFIG_PATH;

    @Value("${servicePath}")
    public String SERVICE_PATH;

    @Value("${ip}")
    public String IP;

    @Value("${maxInterval}")
    public int MAXINTERVAL;

    @Value("${httpport}")
    public int HTTPPORT;

    @Value("${server.port}")
    public int HTTPSPORT;

}
