# Compare with other methods

- **Mehod 1. Have multiple test environments permanently**

**Advantage**: Test environment is always available

**Shortcoming**: Dose not fix the root problem, test environments will become insufficient when more works run in parallel; Multiple times waste of resources when no one is using them; Additional resource allocation and management problems

- **Method 2. Always recreate whole test environment use tools such as helm**

**Advantage**: Minimal waste of resources when idle

**Shortcoming**: Pull up all services in batches may take as long as several minutes to tens of minutes; Sometimes have to pull the whole set of services in order to test just one of them, cause much more waste; A complete environment not only needs to run service instances, but also imports service data, maintaining data and related script require extra works
