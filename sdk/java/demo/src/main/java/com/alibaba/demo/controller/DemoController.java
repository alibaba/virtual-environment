package com.alibaba.demo.controller;

import com.alibaba.demo.util.EnvUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

import javax.servlet.http.HttpServletRequest;

@RestController
public class DemoController {

    Logger logger = LoggerFactory.getLogger(DemoController.class);

    @Autowired
    private RestTemplate restTemplate;

    @GetMapping("/a")
    String callA(HttpServletRequest request) {
        logger.info("accessing service a at " + EnvUtil.getCurrentEnv());
        return "a-[" + EnvUtil.getCurrentEnv() + "] received " + EnvUtil.getRequestEnv(request) + "\n" +
            restTemplate.getForEntity("http://service-b:9000/b", String.class);
    }

    @GetMapping("/b")
    String callB(HttpServletRequest request) {
        logger.info("accessing service b at " + EnvUtil.getCurrentEnv());
        return "b-[" + EnvUtil.getCurrentEnv() + "] received " + EnvUtil.getRequestEnv(request) + "\n" +
            restTemplate.getForEntity("http://service-c:9000/c", String.class);
    }

    @GetMapping("/c")
    String callC(HttpServletRequest request) {
        logger.info("accessing service c at " + EnvUtil.getCurrentEnv());
        return "c-[" + EnvUtil.getCurrentEnv() + "] received " + EnvUtil.getRequestEnv(request) + "\n";
    }

}
