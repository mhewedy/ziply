# ziply

zip some file and offer a zip download link

## Why?

Sometimes you find yourself in a restricted environement where you can not download certain type of files, for example windows binary files, using this service (on heroku as below, or by deploying it in any cloud environemnt of your choice) you can bypass such constraint and got your file compressed for free.

#### Why not other online alternatives?
This is safe. You are seeing the code and you deploying it yourself, so no chance for your binary to be trojaned by the service provider.


## Usage:

Navigate to the home page: `http://<host:port>` (e.g. https://ziply.herokuapp.com/)

Or use the API directly:
```
http://<host:port>/dl?url=<url contains file to be downloaded>
```
e.g.: (ziply is already deployed to herokuapp)
```
https://ziply.herokuapp.com/dl?url=https://download-office.grammarly.com/latest/GrammarlyAddInSetup.exe
```

