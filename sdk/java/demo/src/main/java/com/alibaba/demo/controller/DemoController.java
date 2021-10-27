package com.alibaba.demo.controller;

import static com.alibaba.demo.util.EnvUtil.*;
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
        logger.info("accessing service a at " + getCurrentEnv());
        return "a-[" + getCurrentEnv() + "] received " + getRequestEnv(request) + "\n" +
            restTemplate.getForObject("http://127.0.0.1:9002/b", String.class);
    }

    @GetMapping("/b")
    String callB(HttpServletRequest request) {
        logger.info("accessing service b at " + getCurrentEnv());
        return "b-[" + getCurrentEnv() + "] received " + getRequestEnv(request) + "\n" +
            restTemplate.getForObject("http://127.0.0.1:9003/c", String.class);
    }

    @GetMapping("/c")
    String callC(HttpServletRequest request) {
        logger.info("accessing service c at " + getCurrentEnv());
        return "c-[" + getCurrentEnv() + "] received " + getRequestEnv(request) + "\n";
    }

}
