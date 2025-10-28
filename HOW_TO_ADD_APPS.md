# How to Add Applications to FluxCD

This is a **step-by-step guide** to add new applications to your FluxCD setup.

## What You Have

- **2 Git Repositories:**
  - `fluxcd-infrastructure` - For infrastructure stuff (namespaces, quotas)
  - `fluxcd-apps` - For your applications
  
- **3 Environments:**
  - `dev` - Auto-deploys from `main` branch
  - `staging` - Manual deploy via Git tags
  - `production` - Manual deploy via Git tags

## Adding a New Application (Step-by-Step)

### Step 1: Create Your App's Base Manifests

Go to the `fluxcd-apps` repo and create a folder for your app:

```bash
cd fluxcd-apps/apps/base/
mkdir my-new-app
cd my-new-app
```

Create 3 files:

**1. `namespace.yaml`** - Creates a namespace for your app
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: my-new-app
```

**2. `deployment.yaml`** - Defines how your app runs
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-new-app
  namespace: my-new-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: my-new-app
  template:
    metadata:
      labels:
        app: my-new-app
    spec:
      containers:
      - name: my-new-app
        image: your-image:v1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: ENVIRONMENT
          value: "production"
        resources:
          requests:
            cpu: 100m
            memory: 64Mi
          limits:
            cpu: 200m
            memory: 128Mi
```

**3. `service.yaml`** - Exposes your app
```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-new-app
  namespace: my-new-app
spec:
  type: ClusterIP
  selector:
    app: my-new-app
  ports:
  - port: 8080
    targetPort: 8080
```

**4. `kustomization.yaml`** - Lists all files
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - namespace.yaml
  - deployment.yaml
  - service.yaml
```

### Step 2: Add Your App to Base Kustomization

Edit `fluxcd-apps/apps/base/kustomization.yaml`:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - podinfo/
  - go-crud-api/
  - my-new-app/        # Add this line
```

### Step 3: Create Environment-Specific Configs

For each environment (dev, staging, production), create a patch file.

**Dev:** `fluxcd-apps/apps/overlays/dev/my-new-app-patch.yaml`
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-new-app
  namespace: my-new-app
spec:
  replicas: 1          # Less replicas for dev
  template:
    spec:
      containers:
      - name: my-new-app
        env:
        - name: ENVIRONMENT
          value: "development"
        resources:
          requests:
            cpu: 50m      # Less resources for dev
            memory: 32Mi
          limits:
            cpu: 100m
            memory: 64Mi
```

**Staging:** `fluxcd-apps/apps/overlays/staging/my-new-app-patch.yaml`
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-new-app
  namespace: my-new-app
spec:
  replicas: 2
  template:
    spec:
      containers:
      - name: my-new-app
        env:
        - name: ENVIRONMENT
          value: "staging"
```

**Production:** `fluxcd-apps/apps/overlays/production/my-new-app-patch.yaml`
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-new-app
  namespace: my-new-app
spec:
  replicas: 3          # More replicas for production
  template:
    spec:
      containers:
      - name: my-new-app
        env:
        - name: ENVIRONMENT
          value: "production"
```

### Step 4: Add Patches to Environment Kustomizations

**Dev:** Edit `fluxcd-apps/apps/overlays/dev/kustomization.yaml`
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../base
patches:
  - path: podinfo-patch.yaml
  - path: go-crud-api-patch.yaml
  - path: my-new-app-patch.yaml    # Add this
```

**Staging:** Edit `fluxcd-apps/apps/overlays/staging/kustomization.yaml`
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../base
patches:
  - path: podinfo-patch.yaml
  - path: go-crud-api-patch.yaml
  - path: my-new-app-patch.yaml    # Add this
```

**Production:** Edit `fluxcd-apps/apps/overlays/production/kustomization.yaml`
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../base
patches:
  - path: podinfo-patch.yaml
  - path: go-crud-api-patch.yaml
  - path: my-new-app-patch.yaml    # Add this
```

### Step 5: Commit and Push

```bash
cd fluxcd-apps
git add apps/
git commit -m "Add my-new-app"
git push origin main
```

### Step 6: Wait or Force Deploy

**Option 1: Wait** - Flux checks every 1 minute and will auto-deploy to dev

**Option 2: Force it now:**
```bash
flux reconcile kustomization apps --with-source
```

### Step 7: Verify

```bash
# Check if namespace was created
kubectl get namespace my-new-app

# Check if pods are running
kubectl get pods -n my-new-app

# Check deployment
kubectl get deployment -n my-new-app

# Check service
kubectl get service -n my-new-app
```

## Quick Reference: File Locations

```
fluxcd-apps/
├── apps/
│   ├── base/
│   │   ├── my-new-app/              ← Step 1: Create app files here
│   │   │   ├── namespace.yaml
│   │   │   ├── deployment.yaml
│   │   │   ├── service.yaml
│   │   │   └── kustomization.yaml
│   │   └── kustomization.yaml       ← Step 2: Add app here
│   └── overlays/
│       ├── dev/
│       │   ├── my-new-app-patch.yaml    ← Step 3: Create patch
│       │   └── kustomization.yaml       ← Step 4: Add patch here
│       ├── staging/
│       │   ├── my-new-app-patch.yaml    ← Step 3: Create patch
│       │   └── kustomization.yaml       ← Step 4: Add patch here
│       └── production/
│           ├── my-new-app-patch.yaml    ← Step 3: Create patch
│           └── kustomization.yaml       ← Step 4: Add patch here
```

## Common Mistakes

❌ **Forgot to add app to base kustomization.yaml** - App won't deploy
❌ **Forgot to add patch to overlay kustomization.yaml** - Will use base config only
❌ **Typo in app name** - Must match exactly in all files
❌ **Wrong namespace** - Make sure namespace matches in all files

## That's It!

1. Create base files in `apps/base/my-new-app/`
2. Add to `apps/base/kustomization.yaml`
3. Create patches in `apps/overlays/{dev,staging,production}/`
4. Add patches to overlay `kustomization.yaml` files
5. Git commit and push
6. Done!

Flux will automatically deploy to dev. For staging/production, you'll need to create Git tags (covered in a different guide if needed).
