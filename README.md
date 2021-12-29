# engula-operator demo

deploy demo Engula Engine v0.2 to k8s.

## Simple Usage

1. install [kubebuilder](https://book.kubebuilder.io/quick-start.html#installation)
2. clone & cd project folder
```bash
> git clone git@github.com:zojw/engula-operator.git
> cd engula-operator
```
3. deploy CRDs and operator to k8s
```bash
> make deploy IMG=uhub.service.ucloud.cn/zojw/engula-operator:0.9
```

> The operator image is build in my local Apple M1, so x86 cannot works, you can build image by yourself via
>  `make docker-build docker-push IMG=your-image-repo`

4. create cluster namespace(e.g. `t2`)
```bash
> kubectl create ns t2
```
5. deploy Engula demo cluster in namespace
```
> kubectl -n t2 apply -f config/samples/cluster_v1alpha1_cluster.yaml
```
it will deploy 1 journal, 1 storage, 1 kernel and 1 hashengine

> The operator image is build from [https://github.com/zojw/engula/tree/kube-dns](https://github.com/zojw/engula/tree/kube-dns) and upload to UCloud Docker Hub which is fast in China mainland.
> it's also build in Apple M1 now, x86 user can build by yourself via its `amd64.dockerfile` and change image property in `cluster_v1alpha1_cluster.yaml`

6. list pod to check ready

```bash
> kubectl -n t2 get po
NAME                       READY   STATUS    RESTARTS   AGE
cluster-sample-engine-0    1/1     Running   0          5s
cluster-sample-journal-0   1/1     Running   0          5s
cluster-sample-kernel-0    1/1     Running   0          5s
cluster-sample-storage-0   1/1     Running   0          5s

> kubectl -n t2 get svc
NAME                     TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)     AGE
cluster-sample-engine    ClusterIP   None         <none>        24568/TCP   2m25s
cluster-sample-journal   ClusterIP   None         <none>        24570/TCP   2m24s
cluster-sample-kernel    ClusterIP   None         <none>        24567/TCP   2m25s
cluster-sample-storage   ClusterIP   None         <none>        24569/TCP   2m25s
```

7. test engine with `curl`

```bash
> kubectl -n t2 run -it --rm dns-test --image=uhub.service.ucloud.cn/zojw/tiny-tools:3.12
> # curl http://cluster-sample-engine-0.cluster-sample-engine:24568/put\?key\=1\&value\=111
> # put success
> # curl http://cluster-sample-engine-0.cluster-sample-engine:24568/get\?key\=1
> # 111
```

## Some implement details

### CRDs and Controller

implement 6 CRDs's reconcile logic, maybe can view it in one picture

![crds](https://raw.githubusercontent.com/zojw/misc/main/ctrl2.drawio.svg)

ClusterController will reconcile ClusterCR and apply other resource's CRs, for each sub-CR will have a controller to reconcile it.

Finally, we can get statefulet that paired with headless-service(also configured with right ServiceAccount)


### Service Discovery for components

In two places, we (maybe or not) need k8s help to discovery nodes.

1. Kernel need find all Journals and Storages
2. Engine need find all Kernel nodes(then find real leader)

in this version, it use headless-service...kernel can get all storage node in cluster:

```
nslookup <service_name>
```

for example

```
> nslookup cluster-sample-storage
Server:		10.96.0.10
Address:	10.96.0.10#53

Name:	cluster-sample-storage.t2.svc.cluster.local
Address: 10.1.0.226
Name:	cluster-sample-storage.t2.svc.cluster.local
Address: 10.1.0.227
```

so in this implements, it just simple inject Storage/Journal's service name and Port info into kernel by EnvVars

[https://github.com/zojw/engula/tree/kube-dns](https://github.com/zojw/engula/tree/kube-dns) demonstrates how to use them to start a engula without any addition local configure.


### The Questions

1. Not sure where is suitable to place Scaling API

now, user can modify & apply CR's `replca` property to scale up cluster(scale down isn't implement, due to it's highly relevant to engula implementation in next steps), but "where is suitable to place programming callable API" or "how to security expose Scaling API?" is the question.

- option1: use operator-process, it seems already expose port to webhook and metric and all of them is https-only, but is security to expose current certificate info to engula-kernel pod?
- option2: use kernel-process, kernel can do that, but it need kubeapi's privilege
- option3: another-process?????

2. Not sure multiple nodes running mode for each component

it's early stage for Engula, we can setup multiple nodes(e.g. multiple journal node), but it's still unclear for how those nodes works together.

specially, how multiple nodes boostrap as cluster is a big question to implement deploy tools like engula-operator.

some components maybe need sth like "start one node first, then others join previous" or "discovery each other via addition service without join", it's not clear for now

3. Lack retry mechanism

in [https://github.com/zojw/engula/tree/kube-dns](https://github.com/zojw/engula/tree/kube-dns), it's tried to use some dirty way to retry.

tonic seems lack of built-in rediscovery-reconnect-retry mechanism like grpc-go, it makes v0.2 hard to deploy(for example: kernel start first and cannot find any usable storage/journal will exit).

