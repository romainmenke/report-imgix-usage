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
- login to IMGIX in the browser and use the inspector to grab your Auth token
- enter the Auth token when prompted

![Screenshot 2019-06-02 at 16 09 30](https://user-images.githubusercontent.com/11521496/58762738-11dd3400-8553-11e9-92ee-6de0285e3238.png)
