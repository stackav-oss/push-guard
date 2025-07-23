# Overview
PushGuard is a software security tool designed to prevent accidental code pushes to public repositories such as GitHub, ensuring the integrity and security of a company's codebase. This tool operates seamlessly within a development environment, offering a high-level defense against unintentional commits and pushes that could compromise your project, source code, or secrets.

PushGuard is a proactive solution that empowers development teams to interact with the Git client as they normally would, while offering a comprehensive defense against accidental pushes. The tool promotes a smooth and efficient development process and not impact development productivity and speed. Developers are given a message for any git push command run that targets a public Git repository. Developers can choose to allow the push to continue in the event that the action is an approved use case or they can chose to abort the operation.

PushGuard is not designed or intended to be a tool that would prevent bad actors from stealing source code. Attempts to work around this tool may be viewed as a security violation - The tool’s main purpose is to try to protect ourselves from making mistakes.

Please reach out and open an issue with any questions, concerns, or bugs surrounding the tooling.

# Build Binary
Use the local os and arch by default
```
make build
```
Or you could specify the os and arch
```
make build os=<OS> arch=<ARCH>
```
# Clean Workspace
```
make clean
```
