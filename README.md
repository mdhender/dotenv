# dotenv
`dotenv` is a wrapper around John Barton's package for loading dot files.
It updates the global list of environment variables by loading values from files,
commonly called ".env" files.

It really exists just to remind me how to use godotenv.

Note:
"Environment" is heavily overloaded here.
It can mean the environment that the application is running in (development, testing, or production).
It can also mean the list of environment variables set for the application.
Finally, it can mean the file containing environment variables that should be loaded.

## Environments
`dotenv` support the following environment names:

1. development
2. test
3. production


## Environment Files
Environment files are commonly called ".env" files.
They contain lines that set global environment values for the application.

Each environment should have its own set of environment files.

## Loading Priority
Load tries to emulate the priority list from the dotenv page at
https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use.

This is essentially .env.{environment}.local, .env.local, .env.{environment}, then .env.
Local files take priority over environment files, which take priority over the global .env file.

The environment files are loaded in the following order:

	Pri  Filename________________  .gitignore?
	1st  .env.{environment}.local  yes
	2nd  .env.local                yes
	3rd  .env.{environment}        no, but be wary of secrets
	4th  .env                      no, but be wary of secrets

The `godotenv` package (which `Load` imports) loads the environment variable file and updates the application's global environment table.
But, before updating that table, it checks to see if the variable has already been set (either by explicitly exporting it or via a file with higher priority).
If the variable has been set, then its value isn't updated.
Note that this is the default behavior for this package.
It's possible that a user can override this behavior; we assume that they don't.

Take note of that `.gitignore` column.
It's a reminder of which files are expected to contain
If `.gitignore` is `no`, then the file should never contain  sensitive information like credentials and tokens because it may be checked into Git.
If the value is `yes`, then the file might contain that information, so it should never, ever be checked in to Git.

It's important to note that `env.local` is loaded in all in all environments except for test.
The `.env` is loaded in all environments, including test.

## Load
Dot exports a single function, `Load`.

Parameters are

* prefix - string - the name, sort of, of the environment variable that contains the name of the current environme to load
* show - boolean - if true, the environment values are logged (mostly sort of)
* verbose - boolean - if true, information such as the current prefix value and the environment files loaded are written to the log

### The Environment Variable
`Load` requires that the caller set an environment variable to hold the name of the current environment.

This environment variable name defaults to "ENV" if an empty prefix is passed in.
Otherwise, it is set to the prefix with "_ENV" appended.

For example, if you call `Load` with prefix set to the empty string (""),
it will search for an environment variable named `ENV`.
If you call `Load` with prefix set to `SKIPPY`,
it will search for an environment variable named `SKIPPY_ENV`.

The `Load` function will fail and return an error if it can't find the environment variable.
It will also fail if the value of the environment variable isn't in this list

1. development
2. test
3. production

### Load Local Environment Files
`Load` searches for a ".local" file that matches the current environment.
For example, if the value is `test`, then it will look for the `.env.test.local` file.
If the file is found, the contents are loaded and used to update the global list of


## Notes

* The .env.{environment}.local files are for local overrides of environment-specific settings.
We assume that they're created by the deployment process.
They can contain sensitive information like credentials or tokens.
They are loaded first, so they will override settings in the shared files.
They should never be checked into the repository.
 
* The .env.local file is loaded in development and production; never in test.
It should never be checked into the repository.
 
* The .env.* files are shared environment-specific settings.
They should not contain sensitive information like credentials or tokens.
They should always be checked into the repository.
 
* The .env file is loaded in all environments.
It should not contain sensitive information like credentials or tokens.
It is loaded last, so all other files will override any settings.
It should always be checked into the repository.
