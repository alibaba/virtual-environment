package com.alibaba.kt.trace.util;

import com.alibaba.kt.trace.interceptor.OutgoingInterceptor;
import org.springframework.http.client.ClientHttpRequestInterceptor;
import org.springframework.web.client.RestTemplate;

import java.util.List;

public class InterceptorUtil {

    /**
     * Inject a RestTemplate object with default parameter
     * 使用默认参数为RestTemplate对象注入传递环境标签的拦截器
     */
    public static void enableRouteLabel(RestTemplate rt) {
        enableRouteLabel(rt, "X-Virtual-Env", "VIRTUAL_ENV_TAG");
    }

    /**
     * Inject environment tag passing behavior to RestTemplate
     * 注入拦截器，同时指定路由标的Header名称，以及是否自动根据环境变量注入
     * @param rt 需要注入拦截器的HTTP客户端
     * @param tagHeader 作为路由标的Header
     * @param tagEnvVar 读取路由标的环境变量（为null则禁用自动注入路由标）
     */
    public static void enableRouteLabel(RestTemplate rt, String tagHeader, String tagEnvVar) {
        InterceptorGlobal.tagHeader = tagHeader;
        InterceptorGlobal.tagEnvVar = tagEnvVar;
        List<ClientHttpRequestInterceptor> interceptors = rt.getInterceptors();
        interceptors.add(new OutgoingInterceptor());
        rt.setInterceptors(interceptors);
    }

}
