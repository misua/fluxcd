# Design: FluxCD Multi-Environment GitOps

## Context

FluxCD is a CNCF graduated project that implements GitOps for Kubernetes. This design establishes a multi-environment setup for dev, staging, and production clusters, with a dedicated infrastructure folder for shared resources. The repository owner is `misua`, and authentication uses a GitHub PAT stored in the `GITHUB_TOKEN` environment variable.

**Constraints:**
- Must work with existing GitHub repository: `https://github.com/misua/fluxcd`
- GitHub PAT already available as environment variable
- Need to support three distinct environments with different promotion strategies
- Infrastructure changes must be versioned and tracked separately from application deployments

**Stakeholders:**
- DevOps engineers managing cluster infrastructure
- Development teams deploying applications
- Platform team maintaining GitOps workflows

## Goals / Non-Goals

**Goals:**
- Automated GitOps deployment for dev environment
- Manual promotion workflow for staging and production
- Centralized infrastructure management with environment-specific overlays
- Drift detection and self-healing for all environments
- Clear separation between infrastructure and application concerns
- Audit trail via Git history for all cluster changes

**Non-Goals:**
- Multi-cluster federation (each environment is independent)
- Secret management solution (use existing Kubernetes secrets or external solutions)
- CI/CD pipeline for building container images (Flux handles deployment only)
- Multi-tenancy within environments (single tenant per environment)

## Decisions

### Decision 1: Multi-Repo Structure (Infrastructure + Apps)

**What:** Use two separate repositories - one for infrastructure, one for application manifests.

**Why:**
- Clear separation of concerns (infrastructure vs applications)
- Different teams can own different repos (platform team vs app teams)
- Better access control and security boundaries
- Infrastructure changes don't clutter application history
- Scales better with multiple applications
- Aligns with production best practices

**Structure:**

**Repository 1: fluxcd-infrastructure**
```
fluxcd-infrastructure/
├── clusters/
│   ├── dev/
│   ├── staging/
│   └── production/
└── infrastructure/
    ├── base/
    └── overlays/
        ├── dev/
        ├── staging/
        └── production/
```

**Repository 2: fluxcd-apps**
```
fluxcd-apps/
├── apps/
│   ├── base/
│   └── overlays/
│       ├── dev/
│       ├── staging/
│       └── production/
└── clusters/
    ├── dev/
    ├── staging/
    └── production/
```

**Alternatives considered:**
- Monorepo: Rejected due to mixing infrastructure and app concerns, harder to scale with multiple teams
- Per-environment repos: Rejected due to excessive repo sprawl and cross-environment tracking difficulty

### Decision 2: Kustomize for Configuration Management

**What:** Use Kustomize overlays for environment-specific customization.

**Why:**
- Native Kubernetes tool, no additional dependencies
- Flux has first-class Kustomize support
- Clear base + overlay pattern for DRY configuration
- Easy to understand and maintain

**Alternatives considered:**
- Helm: Rejected for infrastructure (Kustomize is simpler for manifests); Helm still available for apps if needed
- Plain YAML: Rejected due to duplication across environments

### Decision 3: Automated Dev, Manual Staging/Prod

**What:**
- Dev: Auto-sync on every commit to main branch
- Staging: Manual approval via Git tags (e.g., `staging-v1.2.3`)
- Production: Manual approval via Git tags (e.g., `prod-v1.2.3`)

**Why:**
- Dev environment benefits from rapid iteration
- Staging/prod require human oversight for safety
- Git tags provide clear audit trail
- Aligns with common GitOps promotion patterns

**Implementation:**
- Dev: `spec.interval: 1m`, `spec.prune: true`
- Staging/Prod: `spec.ref.tag` pointing to specific release tags

**Alternatives considered:**
- Separate branches per environment: Rejected due to merge complexity and drift risk
- Fully automated: Rejected due to production safety requirements

### Decision 4: Infrastructure Folder for Shared Resources

**What:** Dedicated `infrastructure/` folder for cluster-wide resources (CRDs, operators, networking, RBAC).

**Why:**
- Clear separation from application deployments
- Infrastructure changes often require different approval workflows
- Easier to manage dependencies (infrastructure before apps)
- Supports progressive rollout (test infra changes in dev first)

**Contents:**
- Kubernetes operators (cert-manager, ingress-nginx, etc.)
- Custom Resource Definitions (CRDs)
- Cluster-wide RBAC policies
- Network policies and service meshes
- Monitoring and logging infrastructure

### Decision 5: Flux Bootstrap per Cluster with Multi-Repo GitRepository Resources

**What:** Run `flux bootstrap github` on infrastructure repo, then configure additional GitRepository resources for apps repo.

**Why:**
- Each cluster gets its own Flux installation
- Independent failure domains
- Infrastructure repo is bootstrapped first (contains Flux itself)
- Apps repo is added as additional GitRepository resource
- Aligns with Flux recommended practices for multi-repo setups

**Bootstrap command template:**
```bash
# Bootstrap infrastructure repo
flux bootstrap github \
  --owner=misua \
  --repository=fluxcd-infrastructure \
  --branch=main \
  --path=clusters/<environment> \
  --personal \
  --token-auth

# Apps repo is configured via GitRepository resource in infrastructure repo
```

## Risks / Trade-offs

### Risk 1: Monorepo Scaling
**Risk:** Single repository may become unwieldy with many applications.
**Mitigation:** 
- Start with monorepo for simplicity
- Document migration path to multi-repo if needed
- Use clear folder structure to maintain organization
- Consider splitting when >20 applications or >5 teams

### Risk 2: Manual Promotion Errors
**Risk:** Human error in tagging releases for staging/prod.
**Mitigation:**
- Document promotion workflow clearly
- Implement tag naming conventions
- Add validation scripts for tag format
- Consider automation tooling (e.g., GitHub Actions) for tag creation

### Risk 3: Drift Between Environments
**Risk:** Dev, staging, and prod configurations diverge over time.
**Mitigation:**
- Use Kustomize bases to enforce consistency
- Regular audits of overlay differences
- Automated testing of Kustomize builds
- Document intentional differences

### Risk 4: GitHub PAT Expiration
**Risk:** PAT expires, breaking Flux sync.
**Mitigation:**
- Set PAT expiration reminder
- Document PAT rotation procedure
- Consider GitHub App authentication (future enhancement)
- Monitor Flux sync status

## Migration Plan

**Phase 1: Repository Setup**
1. Create directory structure in `https://github.com/misua/fluxcd`
2. Add base Kustomize configurations
3. Commit and push to main branch

**Phase 2: Dev Environment Bootstrap**
1. Ensure kubectl context points to dev cluster
2. Run Flux bootstrap for dev
3. Verify Flux controllers are running
4. Test auto-sync with sample deployment

**Phase 3: Staging Environment Bootstrap**
1. Switch kubectl context to staging cluster
2. Run Flux bootstrap for staging
3. Configure tag-based sync
4. Test manual promotion workflow

**Phase 4: Production Environment Bootstrap**
1. Switch kubectl context to production cluster
2. Run Flux bootstrap for production
3. Configure tag-based sync
4. Document production deployment procedures

**Phase 5: Infrastructure Setup**
1. Add infrastructure base configurations
2. Create environment-specific overlays
3. Test infrastructure deployment in dev
4. Promote to staging, then production

**Rollback:**
- Flux can be uninstalled with `flux uninstall`
- Git history provides rollback points
- Kubernetes resources remain until explicitly pruned

### Decision 6: Multi-Repo Configuration Strategy

**What:** Infrastructure repo contains GitRepository resource pointing to apps repo.

**Why:**
- Infrastructure repo is the source of truth for cluster configuration
- Apps repo is treated as a dependency
- Single PAT can access both repos (same owner)
- Clear dependency order: infrastructure → apps

**Implementation:**
```yaml
# In fluxcd-infrastructure/clusters/<env>/apps-repo.yaml
apiVersion: source.toolkit.fluxcd.io/v1
kind: GitRepository
metadata:
  name: fluxcd-apps
  namespace: flux-system
spec:
  interval: 1m
  url: https://github.com/misua/fluxcd-apps
  ref:
    branch: main  # or tag for staging/prod
  secretRef:
    name: flux-system  # Reuse bootstrap secret
---
apiVersion: kustomize.toolkit.fluxcd.io/v1
kind: Kustomization
metadata:
  name: apps
  namespace: flux-system
spec:
  interval: 5m
  path: ./apps/overlays/<environment>
  prune: true
  sourceRef:
    kind: GitRepository
    name: fluxcd-apps
  dependsOn:
    - name: infrastructure  # Wait for infra first
```

## Open Questions

1. **Cluster provisioning:** Are Kubernetes clusters already provisioned, or should we include cluster setup in this proposal?
2. **Notification channels:** What notification channels should be configured (Slack, Discord, email, webhook)?
3. **RBAC model:** Should Flux have cluster-admin or limited permissions per namespace?
4. **Image automation:** Should we include Flux Image Automation for automatic image updates?
5. **Secrets management:** What secret management solution is preferred (SealedSecrets, SOPS, external secrets operator)?
