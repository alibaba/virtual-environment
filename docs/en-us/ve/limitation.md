# Limitations

The isolation capability provided by Virtual Environment is based on Service Mesh rules. While currently it only support Istio as Service Mesh controller, many limitations are caused by the capability of Istio.

Including:

- Dose not support other service discover methods beyond Kubernetes, currently Istio route rules could not work with frameworks using its own service discovery mechanism such as Dubbo and SpringCloud.
- Dose not support other protocol except HTTP, currently non-HTTP protocols cannot perform fine routing control in Istio.

Besides, because the Sidecar mechanism does not invade the internal logic of the application, it is necessary for users to implement the transfer of tag Header between requests in the application by their own. If the project is already using the OpenTracing SDK, its baggage mechanism can be reused to achieve the tag transparent transmission. You can also [use SDK](en-us/ve/use-sdk.md) or do it manually in the code.
