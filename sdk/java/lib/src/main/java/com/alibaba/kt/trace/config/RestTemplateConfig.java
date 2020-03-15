package com.alibaba.kt.trace.config;

import com.alibaba.kt.trace.processor.TraceRestTemplateProcessor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * Create BeanPostProcessor as an auto-injector
 * 创建请求上下文自动注入器的BeanPostProcessor
 */
@Configuration
public class RestTemplateConfig {

    @Bean
    public TraceRestTemplateProcessor traceRestTemplateBPP() {
        return new TraceRestTemplateProcessor();
    }

}
