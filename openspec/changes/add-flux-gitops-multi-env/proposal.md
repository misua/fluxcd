# Proposal: Add FluxCD GitOps Multi-Environment Setup

## Why

Organizations need a reliable, automated way to manage Kubernetes deployments across multiple environments (dev, staging, production) with infrastructure-as-code practices. FluxCD provides GitOps-based continuous delivery, ensuring that cluster state matches the desired state defined in Git, with built-in drift detection and self-healing capabilities.

This proposal establishes a production-ready FluxCD setup using a **multi-repository approach** that separates infrastructure from application manifests, monitoring three environments (dev, staging, production) through declarative, Git-driven workflows.

## What Changes

- Bootstrap FluxCD on target Kubernetes clusters using GitHub PAT authentication
- Create **two separate repositories**: infrastructure repo and application manifests repo
- Establish infrastructure repository (`fluxcd-infrastructure`) for cluster-wide resources, CRDs, and operators
- Establish application manifests repository (`fluxcd-apps`) for Kubernetes application deployments
- Create environment-specific folder structure (dev, staging, prod) in both repositories
- Configure automated reconciliation for dev environment and manual promotion workflow for staging/prod
- Set up Flux notification system for deployment events and alerts
- Implement Kustomize-based configuration management with environment-specific overlays
- Configure RBAC and security policies for Flux components
- Add health checks and monitoring for Flux controllers

**BREAKING**: This is a new capability with no breaking changes to existing systems.

## Impact

**Affected specs:**
- `flux-bootstrap` (new) - FluxCD installation and GitHub integration
- `environment-management` (new) - Multi-environment configuration and promotion
- `infrastructure-management` (new) - Infrastructure-as-code deployment patterns

**Affected code:**
- **Repository 1** (`github.com/misua/fluxcd-infrastructure`):
  - `clusters/`: Per-environment Flux bootstrap configs
  - `infrastructure/`: Shared infrastructure manifests with environment overlays
- **Repository 2** (`github.com/misua/fluxcd-apps`):
  - `apps/`: Application manifests with environment overlays
  - `clusters/`: GitRepository and Kustomization resources per environment

**Dependencies:**
- Kubernetes clusters (1.21+) for dev, staging, and prod environments
- GitHub Personal Access Token (PAT) with repo permissions (already available as `GITHUB_TOKEN`)
- FluxCD CLI (`flux` binary) for bootstrapping
- kubectl configured for target clusters

**Migration:**
- Fresh installation, no migration required
- Existing Kubernetes resources remain unaffected unless explicitly managed by Flux
