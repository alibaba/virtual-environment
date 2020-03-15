package com.alibaba.kt.trace.interceptor;

import com.alibaba.kt.trace.util.InterceptorGlobal;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpRequest;
import org.springframework.http.client.ClientHttpRequestExecution;
import org.springframework.http.client.ClientHttpRequestInterceptor;
import org.springframework.http.client.ClientHttpResponse;

import java.io.IOException;

/**
 * Exit interceptor, restore environment from thread context to next request header
 * 出口拦截器，在消息发出前透传环境标签
 */
public class OutgoingInterceptor implements ClientHttpRequestInterceptor {

    public OutgoingInterceptor() {
    }

    @Override
    public ClientHttpResponse intercept(HttpRequest request, byte[] body,
                                        ClientHttpRequestExecution execution) throws IOException {
        HttpHeaders headers = request.getHeaders();
        // Try fetch environment tag from thread context
        // 从线程上下文中取出环境标签
        String envTag = ThreadLocalStoragedVar.get();
        if (envTag != null && envTag.length() > 0) {
            // When environment tag exist, pass it to downstream request
            // 如果当前上下文有环境标签，透传此标签到后续请求
            headers.add(InterceptorGlobal.tagHeader, envTag);
        } else if (InterceptorGlobal.tagEnvVar != null && InterceptorGlobal.tagEnvVar.length() > 0) {
            // When allow fetching environment tag from environment variable, try it
            // 如果启用了从环境变量读取，则尝试将环境变量中的标签设置到请求中
            headers.add(InterceptorGlobal.tagHeader, System.getenv(InterceptorGlobal.tagEnvVar));
        }
        return execution.execute(request, body);
    }
}
