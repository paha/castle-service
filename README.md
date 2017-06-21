# Castle trivial project

[![Build Status](https://travis-ci.org/paha/castle-service.svg?branch=master)](https://travis-ci.org/paha/castle-service)

Deployment automation for a 2 tier web service.

_The resulting service can be reached at [http://www.castle.snagovsky.com][2]_

**tl;dr:** _Cloud native style_. Deployment and service lyfecycle are managed via CI/CD pipeline in GitHub/Travis CI, with rolling updates and existing tests validation. The service components dockerized, artifacts are stored on Docker Hub. The service is deployed to kubernetes cluster in AWS with Prometheus monitoring. All components are resilient to failure and setup for reactive auto-scaling based on demand.

_Reach out to me with any additional questions or to get more details on the project._

-----

**Contents:**
* [Initial setup][4]
* [Design][5]
    + [Service and kube cluster][16]
    + [CI/CD][17]
* [Improvements for the next iteration][6]
* [Additional details][15]

-----

## Initial setup

To meet set criteria a [Kubernetes][3] cluster was setup in AWS ([kops deployment][7]) with autoscalable worker nodes and [Prometheus][8]. Applications have been dockerized.

* Route53 managed delegated subdomain for the project - castle.snagovsky.com
* Kubernetes cluster in AWS, kube API - [https://api.castle.snagovsky.com][9], kube dashboard - [https://kube.castle.snagovsky.com/ui][10]
* [Docker Hub][11] repositories setup to store docker image artifacts.
* [GitHub][12] + [Travis-Ci][13] projects setup

:info: The cluster will be _live_ for a week from Wed, Jun 21 2017. Feel free to ask for kube cluster creds if you would like to get your hands on the setup.

## Design

The service design follows traditional Microservice (contemporary SOA) design pattern. Each tier is part of separate deployment allowing to apply separate resilience and scalability models. Components interact over HTTP.

Each tier is fronted by a loadbalancer, ELB + kube-proxy in the case of the _www app_ and kube-proxy for the _backend app_ loadbalancing between kube replica-sets for corresponding app with simple healthchecks.

Consistency and repeatability of the service is achieved by describing service via Dockerfiles, Kubernetes and CI/CD config files.

### Service and kube cluster.

1. The service endpoint authoritative A record is managed by Route53 aliasing it to ELB resource used for the first tier app ingress.
1. The service is deployed to a **single** Kubernetes cluster running in the _us-west-2_ AWS region, worker nodes are in a single AZ, autoscaled with ASG.
1. The only supported protocol for the service is HTTP.
1. There are two endpoints exposed: a) The service itself - http://www.castle.snagovsky.com; b) Kubernetes API and dashboard https://api.castle.snagovsky.com.
1. The first tier (www|frontend) has limited throughput and autoscaled on demand by [_HorizontalPodAutoscaler_][14] with a minimum replicas of 5 to allow at least 50 rps.
1. The second tier isn't exposed publicly, it has high throughput and deployed with 2 replica-sets for resiliency. With a max of 10 replicas on the first tier, the backend can't become a bottle neck. In case of a failure of a single replica-set it can handle max workloads on a single replica while kube brings up a second healthy one.
1. Basic metrics are colleted by Prometheus node-exporter and can be presented via Grafana deployed in _monitoring namespace_, or via default Prometheus UI. No alerting or custom metrics are setup. Prometheus and its components are not exposed publicly.
1. Logs are written locally on a none-persistent volumes. No aggregation or filtering/analytics setup.

![alt text][1]

### CI CD:

1. Merges and pushes to the project git repository trigger a test job at Travis-CI. The job runs basic tests and docker builds. Initially only `go fmt && go build` and `docker build`. NOTE: resulting images are not uploaded to Docker Hub, this test validates that the code complies and images build. See bellow for deployment job details.
1. On merges (or pushes if you dare) to deploy branch (`master`) besides all tests from the job above, docker images will be uploaded to Docker Hub and the service will be deployed using `kubectl` without service interruption using rolling update.
1. Sensitive data is stored using Travis-CI [Encryption keys][19]


## Improvements for the next iteration

1. AWS multiple regions for resiliency and perfomance (geographical approximaty) with Route53 latency based routing rules cross region loadbalancing.
1. Consider serverless architecture for this service, Kube cluster is overkill, plus based on use patens it might save significant amounts.
1. Design and implement more tests for CI/CD: a) go unit tests for each app; b) docker image tests; c) functional testing etc.
1. Improve monitoring tracking application specific metrics, add alerting and setup external monitoring for publicly facing service.
1. HorizontalPodAutoscaler for the www app should be using a custom metrics to trigger up/down scaling, current use of CPU utilization is to demonstrate the concept only. A custom kube controller can be written to base scaling decisions on ELB CloudWatch metrics as an alternative.
1. CI/CD versioned deployments based on [Semantic versioning][18] git tags.
1. Implement more service lifeCycle actions for CI/CD: destroy service, deploy to a different cluster etc.
1. Log aggregation and analysis.

## Additional details

> _Reach out to me with any additional questions or to get more details on the project._

Repository contents:

```bash
├── .travis.yml                 # Travis-CI configuration
├── LICENSE                     # Licence
├── README.md                   # This README
├── back-deployment.yaml        # Backend kube deployment, service and autoscaling
├── backend
│   ├── Dockerfile              # Dockerfile for backend container
│   └── backend.go              # backend source
├── castle-namespace.yaml       # Castle namespace kube config
├── docs
│   └── arch_diagram.png        # Architectural diagram
├── www
│   ├── Dockerfile              # Dockerfile for www container
│   └── www.go                  # www source
└── www-deployment.yaml         # www kube deployment and service
```

-----

[1]: docs/arch_diagram.png
[2]: http://www.castle.snagovsky.com
[3]: https://kubernetes.io
[4]: README.md#initial-setup
[5]: README.md#design
[6]: README.md#improvements-for-the-next-iteration
[7]: https://github.com/kubernetes/kops
[8]: https://prometheus.io
[9]: https://api.castle.snagovsky.com
[10]: https://kube.castle.snagovsky.com/ui
[11]: https://hub.docker.com/r/pahatmp/
[12]: https://githum.com/paha/castle-service
[13]: https://travis-ci.org/paha
[14]: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale
[15]: README.md#additional-details
[16]: README.md#service-and-kube-cluster
[17]: README.md#ci-cd
[18]: http://semver.org
[19]: https://docs.travis-ci.com/user/encryption-keys/
