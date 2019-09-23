package com.iflytek.ccr.polaris.cynosure.guava;

import com.google.common.cache.CacheBuilder;
import org.springframework.cache.CacheManager;
import org.springframework.cache.annotation.EnableCaching;
import org.springframework.cache.guava.GuavaCacheManager;
import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Component;

import java.util.concurrent.TimeUnit;

/**
 * guava缓存配置
 *
 * @author sctang2
 * @create 2018-02-05 19:54
 **/
@Component
@EnableCaching
public class GuavaConfig {
    /**
     * 配置全局缓存参数，10秒过期，最大个数1000
     */
    @Bean
    public CacheManager cacheManager() {
        GuavaCacheManager cacheManager = new GuavaCacheManager();
        cacheManager.setCacheBuilder(CacheBuilder.newBuilder().expireAfterWrite(10, TimeUnit.SECONDS).maximumSize(1000));
        return cacheManager;
    }
}
