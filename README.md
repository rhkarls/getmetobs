# getmetobs

getmetobs is a CLI tool that downloads meteorological observation data provided by SMHI.
Data is downloaded and saved in the specified directory (or current directory if not specified) with a standardized filename.

```shell  
Usage:
getmetobs <parameter> <station> <period> [flags]

Examples:
getmetobs 1 159880 latest-day --output /path/to/directory

Flags:
-e, --ext string       File extension for the data (e.g., csv, json) (default "csv")
-h, --help             help for getmetobs
-o, --output string    Directory to save the downloaded file (default is current directory) (default ".")
-v, --version string   API version to use (default "1.0")

Arguments:
parameter    The meteorological parameter to retrieve provided as integer ID, 
             see https://opendata.smhi.se/apidocs/metobs/parameter.html
station      The weather station identifier as integer ID, 
             see https://www.smhi.se/data/meteorologi/ladda-ner-meteorologiska-observationer
period       One of four Periods. Valid values are latest-hour, latest-day, latest-months 
             or corrected-archive. Notice that all Stations do not have all four Periods 
             so make sure to check which ones are available in the Period level.
```