# Package `springcloud`

* Implements general `config.KVStore` interface. Will load config from remote spring config-server or local file.
* Setup for remote spring config
    * Set environment `SPRING_CLOUD_CONFIG_URL` with remote config host:port
    * Set environment `SPRING_CLOUD_CONFIG_PATHS` with spring config paths,
        comma separate. E.g. `/myapp/dev,/myapp/local`
* Setup for local file spring config
    * Set environment `SPRING_CLOUD_CONFIG_URL` with
        `file://<local_machine_path>`. E.g. `file://var/local/myapp.yml`
    * `SPRING_CLOUD_CONFIG_PATHS` env-var will be ignored
* To configure autorefresh, set `springcloud-config.refresh-interval` (duration) in your config. E.g. `springcloud-config.refresh-interval: 10s`

## Secret values

* Use format `{cipher}<encryptedvalues>` to use encrypted values. E.g. `{cipher}jalsdkjflasd43709uahlALJLdgajsdl`. To create encrypted string value, use `<spring_remote_host>/encrypt` endpoint.
* GCP SecretManager, use `{sm}<your_sm_key>` to use value from GCP secret
    manager. E.g. `{sm}mysecret_database_password` will load value from GCP
    secret manager for key `secret_database_password`. Probably need to setup
    `GOOGLE_APPLICATION_CREDENTIALS` environment variable too.

## References
* [GCP local dev credentials setup](https://cloud.google.com/docs/authentication/production#providing_credentials_to_your_application)
* [spring cloud config](https://cloud.spring.io/spring-cloud-config/multi/multi__spring_cloud_config_server.html)
