## Sync into `group/project@main`

Job: [generate](https://gitlab.test/group/project/-/jobs/1) | Pipeline: [#2](https://gitlab.test/group/project/-/pipelines/2) | Commit: `abcdef01`

6 file(s) changed.

`README.md`

`````diff
@@ -1,5 +1,5 @@
 # App
 
 ````go
-fmt.Println("a")
+fmt.Println("b")
 ````
`````

`apps/config.yaml`

```diff
similarity index 73%
rename from apps/values.yaml
rename to apps/config.yaml
@@ -1,4 +1,4 @@
-replicas: 1
+replicas: 2
 image: app:1.0
 port: 80
 name: app
```

`apps/kustomization.yaml`

```diff
@@ -1,3 +1,3 @@
 kind: Kustomization
 resources:
-  - old.yaml
+  - new.yaml
```

`system/logo.png`

```diff
new file mode 100644
GIT binary patch
literal 10
RcmZQzWMWRr%u6h)1ON!h0*n9v

literal 0
HcmV?d00001

```

`system/namespace-legacy.yaml`

```diff
similarity index 100%
rename from system/namespace.yaml
rename to system/namespace-legacy.yaml
```

`system/removed.yaml`

```diff
deleted file mode 100644
@@ -1,4 +0,0 @@
-apiVersion: v1
-kind: ConfigMap
-metadata:
-  name: removed
```
