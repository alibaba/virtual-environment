package com.alibaba.kt.trace.interceptor;

import com.alibaba.ttl.TransmittableThreadLocal;

/**
 * Thread context storage for environment tag
 * 存储在线程上下文的环境标签
 */
public class ThreadLocalStoragedVar {

    private static TransmittableThreadLocal<String> envTag = new TransmittableThreadLocal<>();

    public static String get() {
        return envTag.get();
    }

    public static void set(String id) {
        envTag.set(id);
    }

    private ThreadLocalStoragedVar() {
    }
}
