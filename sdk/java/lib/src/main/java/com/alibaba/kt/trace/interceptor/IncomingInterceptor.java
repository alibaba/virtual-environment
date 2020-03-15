package com.alibaba.kt.trace.interceptor;

import com.alibaba.kt.trace.util.InterceptorGlobal;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.ModelAndView;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

/**
 * Entry interceptor, save environment tag from request header to thread context
 * 入口拦截器，将收到请求中的环境标签保存到线程上下文
 */
public class IncomingInterceptor implements HandlerInterceptor {

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) {
        ThreadLocalStoragedVar.set(request.getHeader(InterceptorGlobal.tagHeader));
        return true;
    }

    @Override
    public void postHandle(HttpServletRequest request, HttpServletResponse response, Object handler, ModelAndView modelAndView) {
    }

    @Override
    public void afterCompletion(HttpServletRequest request, HttpServletResponse response, Object handler, Exception ex) {
    }
}