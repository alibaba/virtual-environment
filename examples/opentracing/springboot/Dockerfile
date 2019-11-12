FROM openjdk:8
MAINTAINER wangte.wt@alibaba-inc.com
LABEL app="opentracingSpringbootDemo" version="0.0.1" by="wangte.wt@alibaba-inc.com"
COPY ./target/opentracing-springboot-demo-0.0.1-SNAPSHOT.jar opentracingSpringbootDemo.jar
ENV url "none"
ENV envMark "dev"
ENTRYPOINT java -jar opentracingSpringbootDemo.jar --envMark=$envMark --url=$url