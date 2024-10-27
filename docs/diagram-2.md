This text is optional and is not rendered.
Only a single `mermaid` block is allowed in each `.md` file.
This diagram is created by https://claude.ai for example purposes.

```mermaid
stateDiagram-v2
    [*] --> CodeCommitted

    CodeCommitted --> StaticAnalysis: Start Analysis
    StaticAnalysis --> UnitTesting: Analysis Passed
    StaticAnalysis --> CodeReview: Analysis Failed
    
    CodeReview --> CodeCommitted: Changes Required
    CodeReview --> UnitTesting: Approved
    
    UnitTesting --> IntegrationQueue: Tests Passed
    UnitTesting --> CodeCommitted: Tests Failed
    
    IntegrationQueue --> IntegrationTesting: Resources Available
    IntegrationQueue --> WaitingQueue: Resources Busy
    WaitingQueue --> IntegrationQueue: Resources Freed
    
    IntegrationTesting --> StagingDeployment: Tests Passed
    IntegrationTesting --> CodeCommitted: Tests Failed
    
    StagingDeployment --> SecurityScan: Deployment Successful
    StagingDeployment --> DeploymentFailed: Deploy Error
    DeploymentFailed --> CodeCommitted: Critical Issues
    DeploymentFailed --> StagingDeployment: Retry Deploy
    
    SecurityScan --> PerformanceTesting: No Vulnerabilities
    SecurityScan --> SecurityReview: Vulnerabilities Found
    SecurityReview --> CodeCommitted: Major Issues
    SecurityReview --> PerformanceTesting: Approved With Risks
    
    PerformanceTesting --> UserAcceptance: Performance OK
    PerformanceTesting --> PerformanceReview: Performance Issues
    PerformanceReview --> CodeCommitted: Major Issues
    PerformanceReview --> UserAcceptance: Minor Issues Accepted
    
    UserAcceptance --> ProductionDeployment: Approved
    UserAcceptance --> StagingDeployment: Changes Required
    
    ProductionDeployment --> Monitoring: Deploy Success
    ProductionDeployment --> RollbackInitiated: Deploy Failed
    
    Monitoring --> [*]: System Stable
    Monitoring --> RollbackInitiated: Issues Detected
    
    RollbackInitiated --> StagingDeployment: Rollback Complete
    
    note right of SecurityScan
        Automated security checks
        and vulnerability scanning
    end note
    
    note right of PerformanceTesting
        Load testing and
        resource utilization checks
    end note
    
    note right of Monitoring
        24/7 system monitoring
        with automated alerts
    end note
```
