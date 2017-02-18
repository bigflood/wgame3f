## game data server (Google App Engine)

## dev env setup

 . dev_setup.sh
    

## deploy

 1. dev env setup with
   ```
   . dev_setup.sh
   ```

 1. build 
   ```
   bash build_gae_wgame.sh
   ```
   
 1. Use the [Admin Console](https://appengine.google.com) to create a
   project/app id. (App id and project id are identical)

   * update app.yaml and service_account.json (your-app-id)
   * update src/gae_wgame/db/config.go (session secret, client id and client secret)

 1. Deploy the application with
   ```
   bash deploy_gae_wgame.sh
   ```

 1. deploy datastore index (if needed)
   ```
   gcloud app deploy src/gae_wgame_run/index.yaml
   ```

 1. Congratulations!  Your application is now live at your-app-id.appspot.com


## references
 * https://cloud.google.com/appengine/docs
 
