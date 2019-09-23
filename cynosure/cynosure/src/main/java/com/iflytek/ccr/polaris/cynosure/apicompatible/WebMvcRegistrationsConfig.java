package com.iflytek.ccr.polaris.cynosure.apicompatible;

import org.springframework.boot.autoconfigure.web.WebMvcRegistrationsAdapter;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerMapping;

/**
 * 继承WebMvcRegistrationsAdapter，重写getRequestMappingHandlerMapping，返回自定义RequestMappingHandlerMapping
 *
 * @author sctang2
 * @create 2017-11-09 15:48
 **/
@Configuration
public class WebMvcRegistrationsConfig extends WebMvcRegistrationsAdapter {
    @Override
    public RequestMappingHandlerMapping getRequestMappingHandlerMapping() {
        return new ApiRequestMappingHandlerMapping();
    }
}
