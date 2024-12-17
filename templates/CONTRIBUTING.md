## Contributing

Contributions to this official repository are accepted providing an issue is raised in advance and approved for work. All changes must be tested. Please read the rest of this guide and the main OpenFaaS guide before opening a PR or Issue.

### License

This project is licensed under the MIT License.

## Submitting a new template

> Note: new templates are not being accepted into the official project repository, but you can still create your own templates and share them with the community.

* We have extended the CLI so that you can create templates within your own public or private Git repositories and then pull them in to use in the CLI.

Example:

```
$ ermes-cli template pull https://github.com/owner/repo
```

Multiple templates can be stored within a single Git repository. View [the CLI reference guide](https://github.com/openfaas/faas-cli) for how to build a template.

## Modifying templates

 ```./verify.sh``` can be utilized to verify your changes to templates in this repo.
