oragono-ldap
============

This is an authentication plugin for Oragono that defers password checking to an LDAP server.

See `example.yaml` for LDAP authentication options.

To configure oragono to use this plugin, add a section like this to your `accounts` block:

```yaml
    auth-script:
        enabled: true
        command: "/path/to/oragono-ldap"
        # constant list of args to pass to the command; the actual authentication
        # data is transmitted over stdin/stdout:
        args: ["/path/to/ldap-config.yaml"]
        # should we automatically create users if the plugin returns success?
        autocreate: true
        # timeout for process execution, after which we send a SIGTERM:
        timeout: 9s
        # how long after the SIGTERM before we follow up with a SIGKILL:
        kill-timeout: 1s
```
