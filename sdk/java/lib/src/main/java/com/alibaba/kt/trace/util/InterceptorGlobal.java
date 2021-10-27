package com.alibaba.kt.trace.util;

public class InterceptorGlobal {

    public static String tagHeader;
    public static String tagEnvVar;

    /**
     * Initialize virtual environment interceptor with default parameters
     * 使用默认参数初始化透传环境标签的拦截器
     */
    public static void setupInterceptors() {
        setupInterceptors("X-Virtual-Env", "APP_VIRTUAL_ENV");
    }

    /**
     * Initialize virtual environment interceptor with specified parameters
     * 使用指定路由标的Header名称和环境变量名初始化透传环境标签的拦截器
     * @param header 作为路由标的Header
     * @param envVar 读取路由标的环境变量（为null则禁用自动注入路由标）
     */
    public static void setupInterceptors(String header, String envVar) {
        InterceptorGlobal.tagHeader = header;
        InterceptorGlobal.tagEnvVar = envVar;
    }

}
