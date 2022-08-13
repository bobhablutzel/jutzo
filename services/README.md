# Jutzo Services

The backend for the Jutzo content management system.

Jutzo provides a set of APIs that allow the user to access the content, manage their
permissions, schedule events, and other activities. 



Environment variables:
- **JUTZO_DB_URL** [required]: A Postgresql connection URL in the format postgresql://(user(:pass)?@)?host(:port)?
- **JUTZO_JWT_SECRET** [required]: A hex string representing the secret value used for signing 
  JWT tokens. Should be unique for each environment but shared across all servers in a 
  given environment
- **JUTZO_REDIS_URL** [required]: a Redis connection url in the format redis://

- **JUTZO_ADMIN_USER** [required for first run]: The administrative username
- **JUTZO_ADMIN_PASS** [required for first run]: The administrative password
- **JUTZO_ADMIN_EMAIL** [required for first run]: The administrative email

- **JUTZO_SERVER_PORT** [optional, default 8080]: The port number to listen on
- **JUTZO_HASH_COST** [optional, default 15]: The bcrypt password hashing cost. Larger values will impact login performance.
- **GIN_MODE** [optional]: Set to "release" in production environment

