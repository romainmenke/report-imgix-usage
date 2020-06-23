![Code scanning - action](https://github.com/romainmenke/report-imgix-usage/workflows/Code%20scanning%20-%20action/badge.svg)

# Report IMGIX usage

__⚠️ uses undocumented API's to gather the reports, might break at any time ⚠️__

### How it works :

- fetches all your IMGIX sources
- fetches all report data for all sources since source creation
- interactive CLI to check some data points
- csv exports of all data

### Notes :

- IMGIX might not have reports for the entire time your source has existed
- creates a locale copy of your report data in a `imgix-report-cache.db` file
- cost data is an estimate based on the formula published by IMGIX, this might not be exactly the same.

### Install :

`$ go get -u github.com/romainmenke/report-imgix-usage/...`

`$ report-imgix-usage`

Follow the interactive CLI steps

### Authentication :

- these are undocumented API's without a proper login/auth flow
- uses the same API calls as the web interface to login with your email and password
- **as always inspect the source code of any tool before trusting it with sensitive information**
