jira-file-sync
==============

This is a simple tool for synchronizing attachments from Jira tickets to a 
local directory and supports a single file upload mode too.


Config
------

In order to keep the CLI tool simple, it requires the user configure a config
file that specifies their credentials and site information.

This should be placed in the users home directory and named ```.jira-file-sync.json```.

And example of such a config file is:

```
{
    username: "myusename",
    password: "mypassword",
    sites: [
        {
            "name": "mysite",
            "url" : "https://myurl.domain.com"
        }
    ]
}
```

Note: you will want to be careful about the permissions on this file given it
contains a set of credentials in plaintext.


Usage
---

```
$ jira-file-sync --help

Usage of jira-file-sync:
   -issue string
        Issue references
   -loc string
        Site to use as defined in config
   -upload string
        File to upload
```
