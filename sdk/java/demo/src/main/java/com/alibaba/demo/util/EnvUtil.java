package com.alibaba.demo.util;

import com.alibaba.kt.trace.util.InterceptorGlobal;
import javax.servlet.http.HttpServletRequest;

public class EnvUtil {

    static public String getCurrentEnv() {
        String env = System.getenv(InterceptorGlobal.tagEnvVar);
        return (env != null) ? env : "UNKNOWN";
    }

    static public String getRequestEnv(HttpServletRequest request) {
        String env = request.getHeader(InterceptorGlobal.tagHeader);
        return (env != null) ? env : "UNKNOWN";
    }

}
