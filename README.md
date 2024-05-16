# DNS over TLS Proxy (DoT)

## Background  
  
DNS is the backbone of the internet, translating human-readable domain names into IP addresses that computers understand. However, traditional DNS queries are sent in plaintext, making them vulnerable to interception and eavesdropping. To address this issue, DNS providers like Cloudflare offer a DNS-over-TLS feature, which encrypts DNS queries to protect users' privacy and data integrity.

However, many applications do not natively support DNS-over-TLS, leaving users vulnerable to potential privacy breaches. To bridge this gap and enable applications to leverage DNS-over-TLS, we need a solution that acts as an intermediary between the application and the DNS server, handling DNS queries over TLS.
  
![alt text](https://i.imgur.com/rm6cQwv.jpg "Without DoT")

## Requirement

DNS-over-TLS proxy that can:

- Handle at least one DNS query, and give a result to the client.
- Work over TCP and talk to a DNS-over-TLS server that works over TCP (e.g: Cloudflare).
- Allow multiple incoming requests at the same time
- Also handle UDP requests, while still querying tcp on the other side

So, A DNS over TLS proxy that accepts simple (conventional) DNS requests and proxy it to a DNS servers running with DNS over TLS (DoT) (eg. cloudflare). So this sidecar DNS proxy will proxify our DNS queries to a DoT server.

![alt text](https://i.imgur.com/gjoygas.jpg "Title")

## How It Works

1. Listening for Requests: The proxy listens on standard DNS ports, ready to intercept DNS queries.
2. Establishing a Secure Channel:
    - When it receives a request, it establishes a secure, encrypted connection (using TLS) with the chosen DoT server (like Cloudflare's 1.1.1.1).
    - It carefully verifies the server's identity to prevent unauthorized interception of the data.
3. Relaying and Resolving: The DNS query is securely passed to the DoT server, which resolves the domain name into an IP address.
4. The proxy receives the response and sends it, completing the lookup.

## Technical Details

UDP Enhancement: For UDP queries, the proxy handles packet size adjustments to ensure smooth communication with the DoT server.

## Getting Started

### Containerization

```bash
docker compose up -d
```

## Testing

To perform unit test:

```go
go test ./...
```

To perform integration tests we can use dig.

- Test TCP DNS resolution:

```bash
dig @127.0.0.1 -p < given port in docker-compose.yml file > google.com +tcp
```

- Test UDP DNS resolution:

```bash
dig @127.0.0.1 -p < given port in docker-compose.yml file > google.com
```

The project also includes a `script` folder with a bash scripts that sends TCP/UDP requests every second in order to validate the server's concurrency.

## Benefits And Additions

1. Enhanced Privacy: The DNS queries are encrypted, making it difficult for others to see what websites you're visiting.
2. Improved Security: The proxy validates the DoT server's identity, helping protect you from DNS spoofing attacks.
3. Ease of Use: No special configuration on devices is required.
4. Potential for Filtering: The proxy's design allows for customization, making it possible to add features like DNS filtering in the future.
5. Allows you to choose between UDP, TCP and Both in deployment.
6. Allows for env and secret configurations externally.
7. Integrated testing.

## Limitations

- Additional Overhead: The encryption and decryption processes inherent in DoT can introduce a slight overhead, leading to increased latency compared to unencrypted DNS queries. While usually negligible, this might be noticeable in some situations, especially for networks with high latency.

- Single Point of Failure: If the chosen DoT resolver goes down, the DNS resolution will be disrupted. However, this can be prevented by implementing good Redundancy technique.

## Security Concerns

- Vulnerabilities in Client-to-DOT Communication: While encryption is handled once the traffic reaches the DOT, the communication between the client and the DOT server remains vulnerable to man-in-the-middle attacks.

- Privacy Trade-off: While DoT protects the DNS queries from the ISP or local network eavesdroppers, you still need to trust the upstream DNS resolver you're using.

- DoS Attacks: Be aware that DNS proxies can be targeted by denial-of-service (DoS) attacks. Blacklisting Malicious Domains and applying rate limiting etc can help.

## Improvements

- Caching: DNS caching to store the results of recent DNS lookups. When the same query is made again, the proxy can quickly return the cached answer instead of contacting the authoritative DNS server. This improves performance and reduces network traffic. Also, vulnerabilities like cache poisoning, where an attacker could manipulate cached DNS records should be prevented.

- Secrets/Envs: Better handling of secrets and config files using secret and config managers.

- Blacklisting Malicious Domains.

- Health check to properly monitor the DOT.

## Deployment Strategies (Integration)

The application can be deployed on kubernetes or other microservice orchestrators in a variety of ways depending on requirements:

- Using DOT as an Ambassador:

My approach would be to utilize the DOT as an ambassador and expose it as a service. If a pod needs to resolve a domain name, the pod sends the DNS query to one of the DOT pods. The DOT pod checks if the DNS record is already on cached (if caching feature is added), if it is not, the DoT proxy establishes a secure, encrypted DNS-over-TLS connection with the external DNS resolver (e.g., 1.1.1.1). It forwards the query to the resolver and gets the IP address back. The DoT proxy caches the response and sends the resolved IP address back to the original requesting pod. This would reduce cost of running it as a side car on every pod and would also eliminate single point of failure. However, the network latency might not be as efficient as some other approaches.

### Other Considerations

- Sidecar Injection (Per-Pod Level)

We could Deploy the DoT proxy as a sidecar container within each pod that needs DNS resolution. This establishes a very low-latency connection between the application container and the DoT proxy, as they share the same network namespace as well as give more fine grained control. However, this causes heavy increase in resource consumption as each Pod would have their own side car. Better if latency and control is a priority.

- NodeLocal DNS Cache (Node Level)

Deploy the DoT proxy as a DaemonSet, ensuring one instance runs on each node in the cluster.
Configure the DoT proxy as a NodeLocal DNSCache in Kubernetes, effectively making it the local DNS resolver for that node.
This leverages the kube-dns or CoreDNS plugin to direct DNS requests to the node-local DoT proxy. This allows improved performance and scalability but creates a Single Point of Failure (SPOF). If the DoT proxy on a node fails, it affects all pods on that node. Can be used in less critical applications.

- Custom CoreDNS Plugin (Cluster Level)

If CoreDNS is being used as the cluster's DNS service, a custom plugin to act as a forwarder to the DoT proxy can be created. The CoreDNS would then direct all DNS traffic through the proxy, ensuring DoT encryption for the entire cluster.
This allows for centralized configurations but requires deeper knowledge of CoreDNS and plugin development.

- Ingress Controller (Edge Level)

If applications are exposed through an Ingress controller (e.g., NGINX Ingress or Traefik), some of these controllers offer the ability to configure DoT as the DNS resolver for upstream services. This allows you to enforce DoT even for external DNS resolution performed by the Ingress controller ensuring DoT is applied before traffic enters the cluster. However, This is Dependent on Ingress controller capabilities: Not all Ingress controllers support this feature.
