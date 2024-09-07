# Haze
An easy to use point-and-shoot web fuzzer to quickly identify crashes and weird behaviours.

It's purpose is to point your attention to parameters, headers etc. which may prove to be vulnerable upon further analysis.

## Usage
Save requests of interest in files (for example by using burp's 'save to file' feature) and point haze to fuzz all the parameters and headers with predefined payloads. All responses identified as crashes will be reported.

```bash
haze -t https://targetapp.local burp_reqs/*.txt
```

### Full list of options:
```
USAGE:
  haze [OPTION]... [REQUEST_FILE]...

ARGS:
  REQUEST_FILE    File(s) containing the raw http request(s)

GENERAL:
  -host, -t       Target host (protocol://hostname:port)
  -probe, -p      Send the probe request only. (Default: false)
  -output, -o     Directory where the report will be created. (Default: cwd)
  -threads, -th   Number of threads to use for fuzzing. (Default: 10)

MATCHERS:
  -mc             Comma-separated list of response codes to report. (Default: 500-599)
  -ml             Comma-separated list of response lengths to report
  -ms             A string to match in response

FILTERS:
  -fc             Comma-separated list of response codes to not report
  -fl             Comma-separated list of response lengths to not report
  -fs             A string to filter in response
```
