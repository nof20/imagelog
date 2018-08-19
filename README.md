# imagelog

A very simple Go App Engine app to record when a webpage is visited using a tracking pixel.

## Setup

* Clone this repo.
* Install and set up the [Cloud SDK](https://cloud.google.com/sdk/install).
* Enable the [Google Sheets API](https://console.developers.google.com/apis/api/sheets.googleapis.com) for your cloud project.
* Create a new Google Sheet and copy the sheet ID (the long string in the URL).  Create a tab called "Log".
* Paste the sheet ID into the environment variable section in `app.yaml`.
* From Google Sheets, Choose File > Share and add write access for the account `PROJECTID@appspot.gserviceaccount.com`, where `PROJECTID` is your cloud project ID.
* Deploy the app using `gcloud app deploy` on the command line.

Visitors to `https://PROJECT_ID.appspot.com/imagelog.png` will now be logged in your sheet.

![imagelog](https://franklin-labs.appspot.com/imagelog.png)