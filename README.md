## Team CLI

A Command Line Interface for retrieving player information.

### Configuration

The configuration for the CLI is read from a file. The configuration includes the URL for retrieving data, and the list of target teams.

The config file should be named config.json and have the following structure:
```
{
  "GetTeamURL": "https://api.example.com/score-one-proxy/api/teams/en/%d.json",
  "TargetTeams": [
    "Germany",
    "Bayern Munich"
  ]
}
```

### Usage

1. Build the cli using `go build`

```
./team-cli -max_routines [num_of_goroutines]
```

The -max_routines flag is optional and allows you to specify the number of goroutines to use for retrieving data.