# Run Once Duration Override

The Run Once Duration Override mutating admission webhook allows to override `activeDeadlineSeconds` field
for a pod with `RestartPolicy` set to `Never` or `OnFailure`. The so called run-once pods.

## Deploy the Operator

### Quick Development

1. Build and push the operand image to a registry:
   ```sh
   export QUAY_USER=${your_quay_user_id}
   export IMAGE_TAG=${your_image_tag}
   podman build -t quay.io/${QUAY_USER}/run-once-duration-override:${IMAGE_TAG} .
   podman login quay.io -u ${QUAY_USER}
   podman push quay.io/${QUAY_USER}/run-once-duration-override:${IMAGE_TAG}
   ```

1. Generate manifests deploying the admission webhook:
   ```sh
   make manifests
   ```

1. Update the image spec under `.spec.template.spec.containers[0].image` field in the `_output/manifests/500_deployment.yaml` Deployment to point to the newly built image.

1. Deploy the admission webhook:
   ```sh
   oc apply -f _output/manifests
   ```

1. Check all DaemonSet pods are running:
   ```sh
   oc get pods -n run-once-duration-override
   ```

## Example

1. Create or choose a namespace. E.g. `test`
   ```
   $ oc create ns test
   ```

1. Label the namespace with `runoncedurationoverrides.admission.runoncedurationoverride.openshift.io/enabled: "true"`
   ```
   $ oc label ns test runoncedurationoverrides.admission.runoncedurationoverride.openshift.io/enabled=true
   ```

1. Create a testing pod in the namespace with RestartPolicy set to Never. E.g.
   ```
   $ cat pod.yaml
   apiVersion: v1
   kind: Pod
   metadata:
     name: example
     namespace: test
   spec:
     restartPolicy: Never
     containers:
       - name: busybox
         securityContext:
           allowPrivilegeEscalation: false
           capabilities:
             drop: ["ALL"]
           runAsNonRoot:
             true
           seccompProfile:
             type: "RuntimeDefault"
         image: busybox:1.25
         command:
           - /bin/sh
           - -ec
           - |
             while sleep 5; do date; done
   ```
   The manifest is also located under `examples/pod.yaml`.

   ```sh
   $ oc apply -f pod.yaml
   ```

1. Checking the `.spec.activeDeadlineSeconds` field was set to 3600:
   ```sh
   $ oc get pods -n test -o json | jq '.spec.activeDeadlineSeconds'
   3600
   ```
