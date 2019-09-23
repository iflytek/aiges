package com.iflytek.ccr.polaris.cynosure.druid;

import com.alibaba.druid.support.http.StatViewServlet;
import com.alibaba.druid.support.http.WebStatFilter;
import com.iflytek.ccr.polaris.cynosure.util.PropUtil;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.boot.web.servlet.ServletRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * druid配置
 *
 * @author sctang2
 * @create 2018-03-12 20:56
 **/
@Configuration
public class DruidConfig {
    @Autowired
    private PropUtil propUtil;

    /**
     * 注册ServletRegistrationBean
     *
     * @return
     */
    @Bean
    public ServletRegistrationBean druidServlet() {
        ServletRegistrationBean reg = new ServletRegistrationBean();
        reg.setServlet(new StatViewServlet());
        reg.addUrlMappings("/druid/*");
        reg.addInitParameter("loginUsername", "sctang2");
        reg.addInitParameter("loginPassword", "2017007476");
        reg.addInitParameter("allow", propUtil.IP);
        reg.addInitParameter("deny", "");
        reg.addInitParameter("resetEnable", "false");
        reg.addInitParameter("mergeSql", "true");
        reg.addInitParameter("slowSqlMillis", "10");
        reg.addInitParameter("logSlowSql", "true");
        return reg;
    }

    /**
     * 注册FilterRegistrationBean
     *
     * @return
     */
    @Bean
    public FilterRegistrationBean filterRegistrationBean() {
        FilterRegistrationBean filterRegistrationBean = new FilterRegistrationBean();
        filterRegistrationBean.setFilter(new WebStatFilter());
        filterRegistrationBean.addUrlPatterns("/*");
        filterRegistrationBean.addInitParameter("exclusions", "*.js,*.gif,*.jpg,*.png,*.css,*.ico,/druid/*");
        filterRegistrationBean.addInitParameter("profileEnable", "true");
        return filterRegistrationBean;
    }
}
