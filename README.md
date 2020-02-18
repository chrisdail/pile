# Pile

A simple, opinionated build tool for creating piles of containers (building them).

Pile provides a few simple capabilities:

- Simple abstraction for building containers - Configuration files contain information about the image name, building and
  registry information.
- Git-based Versioning - Consistent version numbering based on git. Rather than using `container:latest` which can be unpredictable.
- Simple Caching - Skips building of images that have already been built (by version).
- Container Testing - Runs test suite prior to build and push of an image. Can copy test results for easy integration with CI tools.

## Dependencies

Pile requires a few things to be present to work.

- Git - Pile uses an opinionated versioning scheme based on the information in git. You must have the git CLI available and
  the repo must be using git.
- Docker - Today, only builds using the `docker build` command are supported. Other tools may be supported in the future
  (see [Roadmap](#roadmap) for more details).
- Registry - For the image caching checks to work, Pile queries a registry to see if the image is present. This is implemented for:
    - Amazon ECR - Fully supported
    - Docker Registry v2 - Limited testing done. Some additional work around authentication required (see [Roadmap](#roadmap)
      for more details).

## Getting Started

Download and install the `pile` executable and put it on your path.

### Creating a new project

Choose an existing project (under version control) with a `Dockerfile` to build the container image. Put an empty `pile.yml` file
in the same directory.

Running `pile build [directory]` will build the image from the `Dockerfile` with the same name as the directory.

See [Usage](#usage) for full details on what you can do in `pile.yml`.

## Usage

`pile version [directory ...]` - Generate version information for a project subtree. See (see [Versioning](#versioning))

`pile info [directory ...]` - Lists build information for a project or set of projects.

`pile build [directory ...]` - Builds container images for projects in the directories specified.

`pile build` - Builds all container images.

The `pile.yml` file is used to mark a directory as one that can be built using pile. For a directory to be eligable to be built,
it must have both a `pile.yml` and a `Dockerfile`. The root directory of the repository is special and the `pile.yml` file there
will be treated as defaults for the rest of the projects. Typically this is where you would be information about the
registry or other common settings. If you want the root directory to be somewhere other than the git root, you can use the `-r`
command switch when running any command.

### Image Testing

Pile also supports built-in image testing. Often times you will have a set of unit tests that need to be run against source code
prior to compiling the final image. The image should only be built if the tests pass. Pile supports this by running a test
container to execute tests. Also the test results can be easily copied out of the container to integrate with CI systems.
If the tests fails, pile will exist with code `127` which is the unit test failure code for easy integration with Jenkins or other
CI systems.

### `pile.yml` Syntax

```
# Alternative name for this image. If none specified, defaults to the directory of the project
name: ""

# Alternate context directory (Directory to "build" the image from)
context_dir: ""

# Prefix for the container image name
image_prefix: ""

# Prefix to add in front of the calculated version. Useful for SemVer/CalVer or for variations of an image in the same registry
version_prefix: ""

# Template for computing the version strong
version_template: "{{if .Dirty}}dirty-{{.User}}-{{end}}{{.Commits}}.{{.Hash}}"

# Relative paths to other projects that this project depends on. These are incorporated into the version string
depends_on:
  - ../other_directory

# Arguments passed to the build command via `--build-arg`
build_args:
    Key: Value

# Optional testing configuration
test:
    # Alternate target in a multi-stage build to use for tests. Build is only successful if the tests succeed
    target:

    # Copies test results from the container to the local filesystem (via docker cp)
    copy_results:
        # Location to copy files from in the container.
        src_path: "/app/build/."

        # Location to copy files to relative to the project directory.
        dst_path: "build"

# Optional registry configuration. Required for caching images
registry:
	# Standard Docker Registry v2
    registry_v2:
        url: "http://localhost:5000"
        insecure: true

	# Amazon ECR
	ecr:
        account_id: "12345"
        region: "us-east-1"
```

### Command Chaining

Most of the output of commands run by `pile` are written to standard error. When a build command completes, the image built
is written to standard out. This allows for chaining commands like this:

```
docker run $(pile -r examples/passing-tests build)
```

### Image Descriptor

After an image is built, a descriptor file is written to `build/pile.image.json`. Example:

```
{
    "name": "passing-tests",
    "repository": "passing-tests",
    "tag": "1.de6a95f",
    "fully_qualified_image": "passing-tests:1.de6a95f"
}
```

This can be useful as an input to other automated tools or deployment pipelines.

## Examples

There are multiple example projects located in [examples/](examples/). For each the root of the project is the top level directory
under examples.

### failing-tests

Running this test suite will always fail. Example:

```
pile -r examples/failing-tests build
```

This shows an example of the tests results being written to the `build/test-results` directory.


### passing-tests

An example of running tests that will pass. Example:

```
pile -r examples/passing-tests build
```

This shows an example of the tests results being written to the `build/test-results` directory.


### registry

Some of the examples are configured to push to a local docker registry (for example [examples/registry](examples/registry)). To run
these examples, you need a docker registry. You can run one by running:

```
docker run -d -p 5000:5000 --restart=always --name registry registry:2
```

The registry examples show a more complex project building multiple containers and pushing/searching in the local registry. To run:

```
pile -r examples/registry build
```

## Versioning

Micro-services are often deployed as containers. Software running as a service is often continually deployed from a
master branch. API version becomes much more important that container version. There may not be value in adopting [Semantic Versioning](https://en.wikipedia.org/wiki/Software_versioning). The main consideration is to keep the version number unique.

Many samples show using the `latest` version (example: `container:latest`) to track the most recent version. There are many
issues with using `latest`:

- The `latest` tag may not be the newest version. It is just a tag of who pushed last.
- It is possible for older code to be pushed as latest. It is easy to overwrite the authoritative `latest` version.
- In a distributed team with multiple members, each user may want to test `latest` from their own branches.

The approach used by Pile proposes using a version string that is generated from git metadata and combines two key pieces of data:

- Number of commits on the branch - This is useful for a human readable, incrementing (usually) number. Within a branch, this
  will be increasing. It is not enough though to uniquely identify a branch.
- Git hash - This is sufficient to identify a unique build but lacks the human readability.

Combining these into one gives the human readability/incrementing nature but also uniqueness. These are also unique across branches
due to the hash (though the incrementing number is only a valid reference within a branch). The only snag to this is for a
branch with a dirty working directory (uncommitted changes). The strategy uses the username (for distributed teams) plus a dirty flag
to indicate the nature of this version. A dirty working directory cannot easily be cached.

The default computed version for Pile is: `{COMMITS}.{HASH}` or if dirty `dirty-{USER}.{COMMITS}.{HASH}`. This version can be computed
on any subtree (or multiple subtrees) in the repository. The `pile version` command is provided to use this feature standalone
outside of any container building.

Example:

```
pile version examples
```

Returns:

```
5.023aacf
```

The template used to construct the version string is also configurable from the CLI and for build images.

## Caching

Docker build is very good at caching previously run build steps locally. Some companies use ephemeral build servers which leads
to a situation where this cache is often empty. A simpler caching model is if you know the version number is unique (see [Versioning](#versioning) for why it is), then you can simply check if your registry already has that version (unless you want to pull the latest
base image for security fixes etc). This is the strategy used by pile.

Pile will build an image if:

- The working subtree is dirty. If there are local changes, the container image will always be rebuilt and pushed.
- The image does not exist in the remote registry. For supported registries, it will query the registry to see if the version already
  exists. If it does, it can skip the build.
- The `--force` flag is used. If this flag is specified, it will always perform a build and push.

## Comparison with Other Tools

Other container build tools:

- [BuildKit](https://github.com/moby/buildkit)
- [img](https://github.com/genuinetools/img)
- [orca-build](https://github.com/cyphar/orca-build)
- [buildah](https://github.com/containers/buildah)
- [kaniko](https://github.com/GoogleContainerTools/kaniko)
- [makisu](https://github.com/uber/makisu)
- [bazel](https://github.com/bazelbuild/rules_docker)

Many of these tools (`BuildKit`, `img`, `orca`, `buildah`, `kaniko`, `makisu`) focus on being a better `docker build`
by supporting native environments like Kubernetes, allowing non-root/underprivileged builds, distributed caching
or other enhancements. Others like `Buildah` focuses a lot on the interfacing building of containers. `Bazel` is more
of an overall code build tool that can also create containers.

`Pile` is trying to solve a different set of problems around the overall usability and tooling around building many
containers, integration with registries and CI systems. `Pile` can (see [Roadmap](#roadmap)) leverage some of these other tools
to perform the actual container build part. It can be thought of a orchestrator of other build tools. The caching approach
for some of these tools is very sophisticated, often requires a separate cache server to be deployed and is suitable for large teams.
Pile intends to solve the simple use case and only build images if they were not already built and pushed.

## Roadmap

- [ ] Authentication for registry v2
- [ ] Integration with alternate builders (BuildKit, img, Buildah, etc)
- [ ] CLI enhancements (colors, messages)
- [ ] Complex dependencies (containers based on other containers, aka base images)
- [ ] Parallel Container Builds
- [ ] Integration with additional registries (Docker Hub, GCR, etc)
