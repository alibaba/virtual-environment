package com.example.opentracingspringbootdemo;

import java.net.MalformedURLException;

import io.jaegertracing.Configuration;
import io.jaegertracing.Configuration.ReporterConfiguration;
import io.jaegertracing.Configuration.SamplerConfiguration;
import io.jaegertracing.Configuration.SenderConfiguration;
import io.jaegertracing.samplers.ConstSampler;
import io.opentracing.Tracer;
import io.opentracing.util.GlobalTracer;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class OpentracingSpringbootDemoApplication {

    public static void main(String[] args) throws Exception {
        if (!configureGlobalTracer("springbootdemo")) {
            throw new Exception("Could not configure the global tracer");
        }

        SpringApplication.run(OpentracingSpringbootDemoApplication.class, args);
    }

    static boolean configureGlobalTracer(String componentName)
        throws MalformedURLException {
        Tracer tracer = null;
        SamplerConfiguration samplerConfig = new SamplerConfiguration()
            .withType(ConstSampler.TYPE)
            .withParam(1);
        SenderConfiguration senderConfig = new SenderConfiguration();
        ReporterConfiguration reporterConfig = new ReporterConfiguration()
            .withLogSpans(true)
            .withFlushInterval(1000)
            .withMaxQueueSize(10000)
            .withSender(senderConfig);
        tracer = new Configuration(componentName).withSampler(samplerConfig).withReporter(reporterConfig).getTracer();
        GlobalTracer.register(tracer);
        return true;
    }

}
