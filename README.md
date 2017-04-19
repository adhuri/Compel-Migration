# CoMPeL-LiveMigration
CoMPeL is a framework which  Monitors resource utilization of a container, predicts resource utilization and does live migration of containers to achieve efficient resource utilization.This repository consists of live migration modules.

## It has following major requirements :- 

1. Configuration to switch migration on and off
2. Verify if migration was a correct decision by using timestamp to check if migration would have been done if actual data was given to the migration module.
3. Configuration to use number of minimum decisions before migration happens for all metrics [cpu,memory]
4. Choose max,avg value of the datapoints to consider to decide migration


## Message structure
### 1. Prediction-> Migration 
```javascript
{
  Timestamp
  {
    AgentIP{
      ContainerID{
        PredictedData{
        CPU[]           
        Memory[]
        }
      }
    }
  }
}
```
### 2. Migration -> Migration Scripts
```javascript
{
  SourceAgentIP
  ContainerID
  DestinationAgentIP
}
```

### Additional Details
1. Use counter to track decisions of each Metric(Memory or CPU ) before minimum decisions for CPU [Default 3 ] or Memory [Default 1]
2. Use Timestamp similar to datafetcher for that timestamp would be required to fetch actual data for accuracy test (Whether it was correct decision to Migrate) - On demand check validity
 a. ValidChecker() -> returns ActualMetricDecision for both CPU,Memory and PredictedMetricDecision for both CPU,Memory
3. Store log of Migration details - Timestamp + PredictedData saved for later retrieval and debugging

