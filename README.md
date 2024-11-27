## timelimitx

timelimitx is alternative advanced version of timeout command.

## Installation
```
go install github.com/rix4uni/timelimitx@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/timelimitx/releases/download/v0.0.1/timelimitx-linux-amd64-0.0.1.tgz
tar -xvzf timelimitx-linux-amd64-0.0.1.tgz
rm -rf timelimitx-linux-amd64-0.0.1.tgz
mv timelimitx ~/go/bin/timelimitx
```
Or download [binary release](https://github.com/rix4uni/timelimitx/releases) for your platform.

## Compile from source
```
git clone --depth 1 github.com/rix4uni/timelimitx.git
cd timelimitx; go install
```

## Usage
```
Usage of timelimitx:
  -s, --signal string   Signal to send on timeout (e.g., SIGTERM, SIGINT, SIGKILL) (default "SIGTERM")
  -t, --time string     Time limit (e.g., 1s, 1m, 1h)
      --verbose         Enable verbose output
      --version         Print the version of the tool and exit.
```

## Examples
```
▶ timelimitx -t 1s ping google.com

OR
▶ timelimitx -t 1s "ping google.com"

OR, If you debugging
▶ time timelimitx -t 1s ping google.com
```

#### If you using `shell pipelines (|) and redirection (2>/dev/null)` then use your command in single or double quotes.
```
▶ timelimitx -t 1m "echo target.com | uforall -silent 2>/dev/null | unew -el -i -t -q target.com.txt"

OR, 1 minute for complete command
timelimitx -t 1m "for target in $(cat subs.txt);do echo $target | uforall -silent 2>/dev/null | unew -el -i -t -q $target.txt;done"

OR, 1 minute for every uforall command
for target in $(cat subs.txt);do timelimitx -t 1m "echo $target | uforall -silent 2>/dev/null | unew -el -i -t -q $target.txt";done
```