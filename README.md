# helm-apigroup-fixer
A utility to read config maps and update the api group so they work with kubernetes 1.16.x.  There are two modes to run this utility, ingrestion, and non-ingestion.  Ingestion means that the utility will read a manifest.yaml file from the current directory and apply that manifest to the helm chart.  Non-ingestion means the utility will look up the release from the kubernetes API, and write the manifest to the manifest.yaml file in the local directory.  

There is a need for this utility to reuse the protocol buffer files from helm in order to ingest or egest the manifest section of the release.  The releases are represented as protocol buffer objects, that are zipped and base64 encoded inside a configmap.

The problem this utility fixes is that the last release of a helm chart will have references to removed apigroups while running on kubernetes 1.16.x.  This means we need to update that last release to edit the resources to match what is actually deployed into kubernetes.  For instance, deployments at extensions/v1beta1 are no longer supported.  If you did not update you helm charts before kubernetes version 1.16.x, your chart updates will now fail because kubernetes does not recognize extensions/v1beta1 anymore for deployments.

The normal flow to update an existing helm chart would be to run the utilty in non-ingestion mode first to get the manifest from the **last successful** release.  Then, the user would make the changes necessary in the manifest.yaml file by changing extensions/v1beta1 deployments to apps/v1.  Or, updating any other resource.  After that, the user would run the utility in ingestion mode (the --ingest flag) to read the manifest.yaml file and apply it to the release.

This utility uses glide, since it was branched off of helm at version 2.14.1.  Go modules were not used in helm until version 3.  As such, you need to run glide to get the vendored dependencies.


### NOTE
This client application does not currently contain additional authorization libraries like oidc.  


## Using

To use the utility to create a manifest.yaml file for a release called dev-wildcard-cert at version 25 in the kube-system namespace, run: 
```
go run cmd/fixer.go -release dev-wildcard-cert.v25 -tiller-namespace kube-system
```

I know that version 25 is the last version of my chart since I see that in **helm history dev-wildcard-cert** command.

I make my changes locally to manifest.yaml file, and then run the utility to update the configmap in place for the release, passing the --ingest flag now:


```
go run cmd/fixer.go -release dev-wildcard-cert.v25 -tiller-namespace kube-system --ingest
```

