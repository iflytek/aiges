package com.iflytek.ccr.polaris.cynosure.redis;

import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.session.data.redis.config.annotation.web.http.EnableRedisHttpSession;
import org.springframework.stereotype.Component;

/**
 * redis session共享
 *
 * @author sctang2
 * @create 2018-01-02 10:52
 **/
@Component
@EnableRedisHttpSession(maxInactiveIntervalInSeconds = 7200)
@ConditionalOnProperty(name = "sessionShare", havingValue = "1")
public class RedisSessionConfig {
}
