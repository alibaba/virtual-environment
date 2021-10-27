package com.alibaba.kt.trace.config;

import com.alibaba.kt.trace.interceptor.IncomingInterceptor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.EnableWebMvc;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurerAdapter;

/**
 * Use a interceptor to save environment tag from request
 * 拦截器注入，为所有收到的请求添加入口拦截器
 */
@EnableWebMvc
@Configuration
public class WebMvcConfig extends WebMvcConfigurerAdapter {

    @Bean
    IncomingInterceptor getTraceInterceptor() {
        return new IncomingInterceptor();
    }

    @Override
    public void addInterceptors(InterceptorRegistry registry) {
        // For all request
        // 对所有HTTP请求生效
        registry.addInterceptor(getTraceInterceptor()).addPathPatterns("/**");
    }

}