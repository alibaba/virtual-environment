package com.alibaba.kt.trace.processor;

import com.alibaba.kt.trace.util.InterceptorUtil;
import org.springframework.beans.BeansException;
import org.springframework.beans.factory.config.BeanPostProcessor;
import org.springframework.web.client.RestTemplate;

/**
 * Watch for bean creation, and auto add tag-passing-down behavior to any RestTemplate bean
 * 监听Bean创建，为Spring容器中的RestTemplate对象添加出口拦截器
 */
public class TraceRestTemplateProcessor implements BeanPostProcessor {

    @Override
    public Object postProcessBeforeInitialization(Object bean, String beanName) throws BeansException {
        return bean;
    }

    @Override
    public Object postProcessAfterInitialization(Object bean, String beanName) throws BeansException {
        // When find RestTemplate bean, add injection to it
        // 若存在RestTemplate类型的Bean，自动为其增加透传环境标签Header能力
        if (bean instanceof RestTemplate) {
            RestTemplate rt = (RestTemplate) bean;
            InterceptorUtil.enableRouteLabel(rt);
        }
        return bean;
    }

}
