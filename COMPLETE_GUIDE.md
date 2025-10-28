# Complete FluxCD Guide: Adding Applications

This guide shows you **exactly** how to add applications using **go-crud-api** as a real working example.

---

## 📁 What You Have

### Two Separate GitHub Repositories:
1. **`fluxcd-infrastructure`** (https://github.com/misua/fluxcd-infrastructure)
   - Infrastructure: namespaces, resource quotas, operators
   
2. **`fluxcd-apps`** (https://github.com/misua/fluxcd-apps)
   - Your applications and their Kubernetes manifests

### Sample Application in This Repo:
- **`go-crud-api/`** - A working Go REST API
  - `main.go` - Application code
  - `Dockerfile` - How to build the Docker image
  - `README.md` - API documentation

---

## 🎯 Real Example: How go-crud-api Was Added

### Step 1: Application Code (This Repo)

The app code lives here in `go-crud-api/`:
```
go-crud-api/
├── main.go           # Go REST API code
├── Dockerfile        # Docker build instructions
├── go.mod            # Go dependencies
├── go.sum            # Dependency checksums
├── README.md         # API docs
└── test-api.sh       # Test script
```

### Step 2: Kubernetes Manifests (fluxcd-apps Repo)

The Kubernetes files live in the `fluxcd-apps` repository:

```
fluxcd-apps/
└── apps/
    ├── base/
    │   ├── go-crud-api/                    ← BASE CONFIGURATION
    │   │   ├── namespace.yaml              # Creates namespace
    │   │   ├── deployment.yaml             # Defines how app runs
    │   │   ├── service.yaml                # Exposes the app
    │   │   └── kustomization.yaml          # Lists all files
    │   └── kustomization.yaml              ← ADD APP HERE
    └── overlays/
        ├── dev/
        │   ├── go-crud-api-patch.yaml      ← DEV OVERRIDES
        │   └── kustomization.yaml          ← REFERENCE PATCH HERE
        ├── staging/
        │   ├── go-crud-api-patch.yaml      ← STAGING OVERRIDES
        │   └── kustomization.yaml          ← REFERENCE PATCH HERE
        └── production/
            ├── go-crud-api-patch.yaml      ← PRODUCTION OVERRIDES
            └── kustomization.yaml          ← REFERENCE PATCH HERE
```

---

## 📝 Step-by-Step: What Files Were Created

### File 1: `apps/base/go-crud-api/namespace.yaml`
**Purpose:** Creates a namespace for the app

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: go-crud-api
```

### File 2: `apps/base/go-crud-api/deployment.yaml`
**Purpose:** Defines how the app runs (base configuration)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-crud-api
  namespace: go-crud-api
spec:
  replicas: 2                    # Default: 2 replicas
  selector:
    matchLabels:
      app: go-crud-api
  template:
    metadata:
      labels:
        app: go-crud-api
    spec:
      containers:
      - name: go-crud-api
        image: go-crud-api:v1.0.0
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 8080
        env:
        - name: ENVIRONMENT
          value: "production"    # Default environment
        resources:
          requests:
            cpu: 100m            # Default resources
            memory: 64Mi
          limits:
            cpu: 200m
            memory: 128Mi
```

### File 3: `apps/base/go-crud-api/service.yaml`
**Purpose:** Exposes the app on port 8080

```yaml
apiVersion: v1
kind: Service
metadata:
  name: go-crud-api
  namespace: go-crud-api
spec:
  type: ClusterIP
  selector:
    app: go-crud-api
  ports:
  - port: 8080
    targetPort: 8080
```

### File 4: `apps/base/go-crud-api/kustomization.yaml`
**Purpose:** Lists all the base files

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - namespace.yaml
  - deployment.yaml
  - service.yaml
```

### File 5: `apps/base/kustomization.yaml` (UPDATED)
**Purpose:** Tells Flux about the new app

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - podinfo/
  - go-crud-api/      # ← ADDED THIS LINE
```

### File 6: `apps/overlays/dev/go-crud-api-patch.yaml`
**Purpose:** Override settings for dev environment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-crud-api
  namespace: go-crud-api
spec:
  replicas: 1              # Dev: Only 1 replica
  template:
    spec:
      containers:
      - name: go-crud-api
        env:
        - name: ENVIRONMENT
          value: "development"    # Dev environment
        resources:
          requests:
            cpu: 50m              # Dev: Lower resources
            memory: 32Mi
          limits:
            cpu: 100m
            memory: 64Mi
```

### File 7: `apps/overlays/dev/kustomization.yaml` (UPDATED)
**Purpose:** Apply the dev patch

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../base
patches:
  - path: podinfo-patch.yaml
  - path: go-crud-api-patch.yaml    # ← ADDED THIS LINE
```

### File 8 & 9: Staging and Production Patches
Same as dev, but with different values:
- **Staging:** 2 replicas, `ENVIRONMENT=staging`
- **Production:** 3 replicas, `ENVIRONMENT=production`

---

## 🚀 How to Add YOUR Application

### Step 1: Create Your App Code
Put your application code in this repo (optional, for reference):
```bash
mkdir my-app
cd my-app
# Add your code, Dockerfile, etc.
```

### Step 2: Go to fluxcd-apps Repo
```bash
cd fluxcd-apps/apps/base/
mkdir my-app
cd my-app
```

### Step 3: Create 4 Base Files

**namespace.yaml:**
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: my-app
```

**deployment.yaml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: my-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: my-app
        image: my-app:v1.0.0
        ports:
        - containerPort: 8080
        resources:
          requests:
            cpu: 100m
            memory: 64Mi
```

**service.yaml:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-app
  namespace: my-app
spec:
  type: ClusterIP
  selector:
    app: my-app
  ports:
  - port: 8080
    targetPort: 8080
```

**kustomization.yaml:**
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - namespace.yaml
  - deployment.yaml
  - service.yaml
```

### Step 4: Add to Base Kustomization
Edit `apps/base/kustomization.yaml`:
```yaml
resources:
  - podinfo/
  - go-crud-api/
  - my-app/          # ← ADD THIS
```

### Step 5: Create Environment Patches

**Dev:** `apps/overlays/dev/my-app-patch.yaml`
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: my-app
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: my-app
        resources:
          requests:
            cpu: 50m
            memory: 32Mi
```

**Staging:** `apps/overlays/staging/my-app-patch.yaml`
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: my-app
spec:
  replicas: 2
```

**Production:** `apps/overlays/production/my-app-patch.yaml`
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: my-app
spec:
  replicas: 3
```

### Step 6: Reference Patches in Overlays

Edit each `apps/overlays/{dev,staging,production}/kustomization.yaml`:
```yaml
patches:
  - path: podinfo-patch.yaml
  - path: go-crud-api-patch.yaml
  - path: my-app-patch.yaml      # ← ADD THIS
```

### Step 7: Commit and Push
```bash
cd fluxcd-apps
git add apps/
git commit -m "Add my-app"
git push origin main
```

### Step 8: Deploy
Flux auto-deploys to dev in 1 minute, or force it:
```bash
flux reconcile kustomization apps --with-source
```

### Step 9: Verify
```bash
kubectl get pods -n my-app
kubectl get svc -n my-app
```

---

## 📊 Quick Reference

### Where Things Live:

| What | Where |
|------|-------|
| **App code** | This repo: `go-crud-api/` |
| **Base manifests** | `fluxcd-apps/apps/base/go-crud-api/` |
| **Dev config** | `fluxcd-apps/apps/overlays/dev/go-crud-api-patch.yaml` |
| **Staging config** | `fluxcd-apps/apps/overlays/staging/go-crud-api-patch.yaml` |
| **Production config** | `fluxcd-apps/apps/overlays/production/go-crud-api-patch.yaml` |

### Files to Edit When Adding App:

1. ✅ `apps/base/my-app/` - Create 4 files (namespace, deployment, service, kustomization)
2. ✅ `apps/base/kustomization.yaml` - Add app to resources list
3. ✅ `apps/overlays/dev/my-app-patch.yaml` - Create dev overrides
4. ✅ `apps/overlays/dev/kustomization.yaml` - Add patch reference
5. ✅ `apps/overlays/staging/my-app-patch.yaml` - Create staging overrides
6. ✅ `apps/overlays/staging/kustomization.yaml` - Add patch reference
7. ✅ `apps/overlays/production/my-app-patch.yaml` - Create production overrides
8. ✅ `apps/overlays/production/kustomization.yaml` - Add patch reference

---

## 🧪 Testing go-crud-api

```bash
# Port forward
kubectl port-forward -n go-crud-api svc/go-crud-api 8080:8080

# Test health
curl http://localhost:8080/health

# Create item
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","description":"Testing"}'

# Get all items
curl http://localhost:8080/items
```

---

## ✅ That's It!

Everything in one place:
- **Sample app code:** `go-crud-api/` (in this repo)
- **Kubernetes manifests:** `fluxcd-apps/apps/` (in apps repo)
- **How it works:** This guide
- **No jumping between docs!**
