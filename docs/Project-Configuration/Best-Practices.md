# Best Practices

## Secrets

Note that you should not pass secrets as command-line arguments but instead as
environment variables or from a file, as Evergreen runs `ps` periodically, which
will log command-line arguments. You can use the `silent` parameter in
`shell.exec` or `subprocess.exec` to avoid logging output.

## Task Directory

Evergreen creates a temporary task directory for each task (aka the [expansion](./Project-Configuration-Files.md#expansions) `${workdir}`). Commands by default execute in that directory. Evergreen will clean up that directory in between tasks unless explicit configured not to. Please don't write outside this directory, as Evergreen won't be able to delete the data your task has written.

## subprocess.exec

In general, use [subprocess.exec](Project-Commands#subprocessexec) instead of shell.exec.

The reasons to prefer subprocess.exec include:

1. Evergreen uses expansions with the same syntax as shell expansions.
2. The shell isn't always bash. Debian/Ubuntu sh is dash.
3. Debugging inline shell scripts is hard.
4. When they're inline, it becomes hard for people to execute them themselves, which makes it difficult to reproduce what Evergreen has done locally.

You can pass environment variables to subprocess.exec if you'd like to pass expansions. It's a good idea to avoid shell.exec as much as possible.

## Task Tags

Use [task tags](Project-Configuration-Files#task-and-variant-tags) to reduce repetition in your Evergreen configuration file.

## Expansions

Be cautious about Evergreen's [expansion syntax](Project-Configuration-Files#expansions).

Evergreen chose an expansion syntax that unfortunately conflicts with bash. This means that you cannot do something like this.

```bash
# WRONG
export foo=bar
echo ${foo}
```

Evergreen will first look for an expansion called `foo` and will substitute that expansion, or, if it doesn't exist, the empty string. You must drop the curly braces if you would like to use a bash variable.

```bash
export foo=bar
echo $foo
```

## Distro Choice

Tasks on more popular distros are often run quicker than tasks on less popular ones. Prefer more popular distros where possible. For more information about available distro choices see [Guidelines around Evergreen distros](https://wiki.corp.mongodb.com/x/CZ7yBg)

## Secrets in Project Variables

If possible, prefer to use private variables instead of public variables to store secrets. Using private variables
ensures that the plaintext secret is not visible in user-facing endpoints and the UI, which reduces the possibility of
someone accidentally being able to view the plaintext value that shouldn't have access. Also, private variables are
always redacted from task logs, making it less likely that the secret will be leaked accidentally (e.g. by echoing the
secret in a script). If you use a public variable for a secret, Evergreen may still attempt to redact it if it contains
a suspicious pattern (e.g. if it looks like a GitHub API key), but that automatic redaction logic is more of a last
resort, best effort mechanism and by nature will never be perfectly reliable.

If you wish to securely provide access to the plaintext secret value to other people, it's recommended to use an
external secrets management service such as 1Password rather than relying on setting an Evergreen public variable with a
secret.
