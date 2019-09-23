package com.iflytek.ccr.polaris.cynosure.interceptor;

import org.springframework.beans.BeansException;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ApplicationContextAware;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.EnableWebMvc;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.ResourceHandlerRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurerAdapter;

/**
 * 配置权限url
 *
 * @author sctang2
 * @create 2017-11-10 14:46
 **/
@Configuration
@EnableWebMvc
@ComponentScan
public class WebConfig extends WebMvcConfigurerAdapter implements ApplicationContextAware {
	private ApplicationContext applicationContext;

	public WebConfig() {
		super();
	}

	@Bean
	SecurityInterceptor securityInterceptor() {
		return new SecurityInterceptor();
	}

	@Override
	public void addResourceHandlers(ResourceHandlerRegistry registry) {
		registry.addResourceHandler("/**/favicon.ico", "/**/index.html", "/**/static/**").addResourceLocations("classpath:/public/");
		registry.addResourceHandler("swagger-ui.html").addResourceLocations("classpath:/META-INF/resources/");
		registry.addResourceHandler("/webjars/**").addResourceLocations("classpath:/META-INF/resources/webjars/");
		super.addResourceHandlers(registry);
	}

	@Override
	public void setApplicationContext(ApplicationContext applicationContext) throws BeansException {
		this.applicationContext = applicationContext;
	}

	@Override
	public void addInterceptors(InterceptorRegistry registry) {
		// 拦截规则
		registry.addInterceptor(securityInterceptor()).addPathPatterns("/**").excludePathPatterns("/", "/api/**/user/login", //
				"/api/**/service/config/feedback", "/api/**/service/discovery/add", "/api/**/service/discovery/feedback");

		super.addInterceptors(registry);
	}
}
