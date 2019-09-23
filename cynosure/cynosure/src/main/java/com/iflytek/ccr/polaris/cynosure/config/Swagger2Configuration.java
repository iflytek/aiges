package com.iflytek.ccr.polaris.cynosure.config;

import org.springframework.boot.SpringBootConfiguration;
import org.springframework.context.annotation.Bean;

import com.google.common.base.Predicate;

import io.swagger.annotations.ApiOperation;
import springfox.documentation.RequestHandler;
import springfox.documentation.builders.ApiInfoBuilder;
import springfox.documentation.builders.PathSelectors;
import springfox.documentation.service.ApiInfo;
import springfox.documentation.spi.DocumentationType;
import springfox.documentation.spring.web.plugins.Docket;

/**
 * Swagger2配置中心
 * 
 * @author jianchen15
 *
 */
@SpringBootConfiguration
public class Swagger2Configuration {

	@Bean
	public Docket createRestApi() {
		Predicate<RequestHandler> predicate = new Predicate<RequestHandler>() {
			@Override
			public boolean apply(RequestHandler input) {
				if (input.isAnnotatedWith(ApiOperation.class))// 只有添加了ApiOperation注解的method才在API中显示
					return true;
				return false;
			}
		};
		return new Docket(DocumentationType.SWAGGER_2).//
				apiInfo(apiInfo()).select().//
				apis(predicate).//
				paths(PathSelectors.any()).build();
	}

	private ApiInfo apiInfo() {
		return new ApiInfoBuilder().title("配置中心接口").//
				description("配置中心接口文档").//
				version("1.0.0").build();
	}
}
