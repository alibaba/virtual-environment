package com.example.opentracingspringbootdemo;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.Map;

import javax.servlet.http.HttpServletRequest;

import io.opentracing.Scope;
import io.opentracing.SpanContext;
import io.opentracing.propagation.Format.Builtin;
import io.opentracing.propagation.TextMapExtractAdapter;
import io.opentracing.propagation.TextMapInjectAdapter;
import io.opentracing.util.GlobalTracer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.ApplicationArguments;
import org.springframework.util.StringUtils;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class DemoController {

    @Autowired
    private ApplicationArguments applicationArguments;

    private final static String HEADER_ENV_MARK_NAME = "ali-env-mark";
    private final static String LINE_BREAK_TEXT = "\n";
    private final static String ENV_MARK_KEY = "envMark";
    private final static String URL_KEY = "url";

    @RequestMapping("/demo")
    public String demo(HttpServletRequest request) {
        Map<String, String> recieveHeaders = getHeaderMap(request);
        SpanContext parentContext = GlobalTracer.get().extract(Builtin.HTTP_HEADERS,
            new TextMapExtractAdapter(recieveHeaders));

        Scope orderSpanScope = GlobalTracer.get().buildSpan("demo").asChildOf(parentContext).startActive(true);

        String url = System.getenv(URL_KEY);
        String envMark = System.getenv(ENV_MARK_KEY);
        System.out.println(String.format("url is: %s, envMark is: %s", url, envMark));

        String requestText = "";
        if (!StringUtils.isEmpty(url)) {
            Map<String, String> sendHeaders = new HashMap<>();
            GlobalTracer.get().inject(orderSpanScope.span().context(), Builtin.HTTP_HEADERS,
                new TextMapInjectAdapter(sendHeaders));
            try {
                requestText = httpGetCall(url, sendHeaders);
            } catch (IOException e) {
                requestText = String.format("call %s failed", url);
            }
        }

        return (StringUtils.isEmpty(requestText) ? "" : requestText + LINE_BREAK_TEXT) + String.format(
            "[springboot @ %s] <-%s", StringUtils.isEmpty(envMark) ? "dev" : envMark,
            StringUtils.isEmpty(orderSpanScope.span().getBaggageItem(HEADER_ENV_MARK_NAME)) ? "empty"
                : orderSpanScope.span().getBaggageItem(HEADER_ENV_MARK_NAME) );
    }

    private static String httpGetCall(String url, Map<String, String> headers) throws IOException {
        URL getUrl = new URL(url);
        HttpURLConnection connection = (HttpURLConnection)getUrl.openConnection();
        for (String oneKey : headers.keySet()) {
            connection.setRequestProperty(oneKey, headers.get(oneKey));
        }
        connection.connect();
        BufferedReader reader = new BufferedReader(new InputStreamReader(
            connection.getInputStream()));
        String lines;
        StringBuilder sb = new StringBuilder();
        while ((lines = reader.readLine()) != null) {
            sb.append(lines);
        }
        reader.close();
        connection.disconnect();
        return sb.toString();
    }

    private Map<String, String> getHeaderMap(HttpServletRequest request) {
        Map<String, String> headers = new HashMap<>();
        Enumeration<String> headerNames = request.getHeaderNames();
        while (headerNames.hasMoreElements()) {
            String key = headerNames.nextElement();
            headers.put(key, request.getHeader(key));
        }
        return headers;
    }
}