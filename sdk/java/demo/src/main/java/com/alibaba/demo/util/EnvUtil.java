package com.alibaba.demo.util;

import javax.servlet.http.HttpServletRequest;

public class EnvUtil {

    static public String getCurrentEnv() {
        String env = System.getenv("TB_APP_ENV");
        return (env != null) ? env : "UNKNOWN";
    }

    static public String getRequestEnv(HttpServletRequest request) {
        String env = request.getHeader("X-Virtual-Env");
        return (env != null) ? env : "UNKNOWN";
    }

}
