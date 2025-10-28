# Environment Management Specification

## ADDED Requirements

### Requirement: Multi-Environment Directory Structure

The system SHALL organize configurations using environment-specific directories with Kustomize overlays.

#### Scenario: Repository structure for three environments

- **GIVEN** a Git repository at https://github.com/misua/fluxcd
- **WHEN** repository structure is created
- **THEN** clusters/ directory contains dev, staging, and production subdirectories
- **AND** apps/ directory contains base and overlays for each environment
- **AND** infrastructure/ directory contains base and overlays for each environment
- **AND** each environment has isolated configuration

#### Scenario: Environment-specific Kustomize overlays

- **GIVEN** base application configuration exists
- **WHEN** environment-specific overlay is applied
- **THEN** environment-specific patches are merged with base
- **AND** resource limits differ per environment (dev: low, prod: high)
- **AND** replica counts differ per environment (dev: 1, prod: 3)
- **AND** environment variables are injected per environment

### Requirement: Automated Dev Environment Sync

The system SHALL automatically sync dev environment on every commit to main branch.

#### Scenario: Auto-sync dev on commit

- **GIVEN** Flux is bootstrapped on dev cluster
- **AND** GitRepository interval is set to 1 minute
- **WHEN** a commit is pushed to main branch
- **THEN** Flux detects change within 1 minute
- **AND** Kustomization is reconciled automatically
- **AND** Kubernetes resources are updated to match Git state
- **AND** old resources are pruned if removed from Git

#### Scenario: Dev drift detection and correction

- **GIVEN** Flux is managing dev cluster
- **AND** a manual change is made to a managed resource
- **WHEN** Flux reconciliation runs
- **THEN** drift is detected
- **AND** resource is reverted to Git state
- **AND** drift event is logged

### Requirement: Manual Staging Promotion

The system SHALL require manual approval for staging deployments via Git tags.

#### Scenario: Promote to staging via Git tag

- **GIVEN** changes are tested in dev environment
- **AND** staging cluster is configured with tag-based sync
- **WHEN** Git tag `staging-v1.2.3` is created and pushed
- **THEN** staging GitRepository detects new tag
- **AND** Kustomization reconciles to tagged commit
- **AND** staging cluster state matches tagged version
- **AND** deployment event is logged

#### Scenario: Reject invalid staging tag format

- **GIVEN** staging cluster expects tags matching `staging-v*` pattern
- **WHEN** Git tag `prod-v1.0.0` is created
- **THEN** staging GitRepository ignores the tag
- **AND** staging cluster state remains unchanged

### Requirement: Manual Production Promotion

The system SHALL require manual approval for production deployments via Git tags.

#### Scenario: Promote to production via Git tag

- **GIVEN** changes are validated in staging environment
- **AND** production cluster is configured with tag-based sync
- **WHEN** Git tag `prod-v1.2.3` is created and pushed
- **THEN** production GitRepository detects new tag
- **AND** Kustomization reconciles to tagged commit
- **AND** production cluster state matches tagged version
- **AND** deployment event is logged with audit trail

#### Scenario: Production rollback via tag

- **GIVEN** production is running version `prod-v1.2.3`
- **AND** a critical bug is discovered
- **WHEN** GitRepository tag reference is updated to `prod-v1.2.2`
- **THEN** Flux reconciles to previous version
- **AND** production cluster rolls back to stable state
- **AND** rollback event is logged

### Requirement: Environment Isolation

The system SHALL ensure complete isolation between dev, staging, and production environments.

#### Scenario: Dev changes do not affect staging

- **GIVEN** dev environment has auto-sync enabled
- **AND** staging environment uses tag-based sync
- **WHEN** breaking change is committed to main branch
- **THEN** dev environment is updated immediately
- **AND** staging environment remains unchanged
- **AND** production environment remains unchanged

#### Scenario: Independent environment failures

- **GIVEN** all three environments are running
- **AND** dev cluster experiences network issues
- **WHEN** dev Flux controllers fail to sync
- **THEN** staging and production continue syncing normally
- **AND** dev failure is isolated

### Requirement: Promotion Workflow Documentation

The system SHALL provide clear documentation for environment promotion workflows.

#### Scenario: Document dev to staging promotion

- **GIVEN** a developer wants to promote changes to staging
- **WHEN** promotion documentation is consulted
- **THEN** step-by-step instructions are provided
- **AND** tag naming conventions are documented
- **AND** validation steps are included
- **AND** rollback procedures are explained

#### Scenario: Document staging to production promotion

- **GIVEN** a release manager wants to promote to production
- **WHEN** promotion documentation is consulted
- **THEN** approval process is documented
- **AND** production tag creation steps are provided
- **AND** monitoring and validation steps are included
- **AND** emergency rollback procedures are documented

### Requirement: Environment-Specific Resource Quotas

The system SHALL apply environment-appropriate resource quotas and limits.

#### Scenario: Dev environment resource constraints

- **GIVEN** dev environment Kustomize overlay
- **WHEN** resources are deployed to dev cluster
- **THEN** CPU limits are set to development-appropriate values
- **AND** memory limits are set to development-appropriate values
- **AND** replica count is set to 1 for cost efficiency

#### Scenario: Production environment resource allocation

- **GIVEN** production environment Kustomize overlay
- **WHEN** resources are deployed to production cluster
- **THEN** CPU limits are set to production-appropriate values
- **AND** memory limits are set to production-appropriate values
- **AND** replica count is set to 3+ for high availability
- **AND** resource quotas prevent over-allocation
