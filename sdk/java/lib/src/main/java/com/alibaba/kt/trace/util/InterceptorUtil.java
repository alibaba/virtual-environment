package com.alibaba.kt.trace.util;

import com.alibaba.kt.trace.interceptor.OutgoingInterceptor;
import org.springframework.http.client.ClientHttpRequestInterceptor;
import org.springframework.web.client.RestTemplate;

import java.util.List;

public class InterceptorUtil {

    /**
     * Inject environment tag passing behavior to RestTemplate
     * 注入拦截器到RestTemplate对象
     * @param rt 需要注入拦截器的HTTP客户端
     */
    public static void enableRouteLabel(RestTemplate rt) {
        List<ClientHttpRequestInterceptor> interceptors = rt.getInterceptors();
        interceptors.add(new OutgoingInterceptor());
        rt.setInterceptors(interceptors);
    }

}
