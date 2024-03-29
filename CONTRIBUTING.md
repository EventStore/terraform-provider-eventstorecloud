# Contributing to `terraform-provider-eventstorecloud`

Welcome to the Event Store community, and thank you for contributing! The following guide outlines the basics of how to get involved. Pull requests to update and expand this guide are very welcome!

## Getting Started

### Community Guidelines

We want the Event Store community to be as welcoming and inclusive as possible, and have adopted a [Code of Conduct][conduct] that we ask all community members to read and observe.

### Licensing

`terraform-provider-eventstorecloud` is licensed under the [MPL-v2][mpl2] license. By submitting a pull request, you represent that you have right to license your contribution to Event Store Ltd and the community, and that by submitting a patch your contributions are licensed under the MPL-v2.0 license.

## Contributing

### Security Issues

Please disclose issues which you believe to be a security threat _by e-mail_ to [security@eventstore.com](mailto:security@eventstore.com) rather than opening a public issue.

### Reporting Issues

Issues may be reported via the [GitHub Repository][github-repo]. When reporting issues, please ensure that you include relevant information which can help diagnose the problem. This includes:

- The version of `terraform-provider-eventstorecloud` - including a commit SHA if using a local build

- The version of `terraform`

- The operating system on which the issue is exhibited (on Unix, the output of `uname -a`)

- Any steps which required to reproduce the issue, including non-default Event Store server configuration.

### Testing local changes
Check the [Developer documentation](DEVELOPER.md)

### Writing a Pull Request

A good client pull request makes a single set of coherent changes, with testing, documentation and a commit message which follows our [Commit Message Template][commit-template].

All contributions to the project are submitted, reviewed and merged via pull requests to the [GitHub repository][github-repo]. We strongly encourage use of [draft pull request][drafts] to open up early discussion about contributions which are not yet ready for review - this also helps reduce work duplicated between community members.

### Checklist

Please ensure that the following steps have been undertaken before marking pull requests as "Ready for Review": 

- Rebase the pull request against the target branch, and ensure there are no merge commits.

- Ensure the build works and tests pass locally using `make install`).

- Ensure that commits are atomic units of work, and that the message follows our [Commit Message Template][commit-template]. A good way to do this is to configure your local clone of the repository to use our included template, by running: `git config commit.template .git.commit.template` in the root directory of the clone.


[conduct]:  https://github.com/EventStore/terraform-provider-eventstorecloud/tree/master/CODE-OF-CONDUCT.md
[mpl2]: https://www.mozilla.org/en-US/MPL/2.0/
[commit-template]: https://github.com/EventStore/terraform-provider-eventstorecloud/tree/master/.git.commit.template
[github-repo]: https://github.com/EventStore/terraform-provider-eventstorecloud
[drafts]: https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/about-pull-requests#draft-pull-requests
