# Fabric Docker Image

This directory provides a simple Docker setup for running the [Fabric](https://github.com/creativerezz/fabric) CLI.

## Build

Build the image from the repository root:

```bash
docker build -t aio -f scripts/docker/Dockerfile .
```

## Persisting configuration

Fabric stores its configuration in `~/.config/aio/.env`. Mount this path to keep your settings on the host.

### Using a host directory

```bash
mkdir -p $HOME/.aio-config
# Run setup to create the .env and download patterns
 docker run --rm -it -v $HOME/.aio-config:/root/.config/aio aio --setup
```

Subsequent runs can reuse the same directory:

```bash
docker run --rm -it -v $HOME/.aio-config:/root/.config/aio aio -p your-pattern
```

### Mounting a single .env file

If you only want to persist the `.env` file:

```bash
# assuming .env exists in the current directory
docker run --rm -it -v $PWD/.env:/root/.config/aio/.env aio -p your-pattern
```

## Running the server

Expose port 8080 to use Fabric's REST API:

```bash
docker run --rm -it -p 8080:8080 -v $HOME/.aio-config:/root/.config/aio aio --serve
```

The API will be available at `http://localhost:8080`.

## Multi-arch builds and GHCR packages

For multi-arch Docker builds (such as those used for GitHub Container Registry packages), the description should be set via annotations in the manifest instead of the Dockerfile LABEL. When building multi-arch images, ensure the build configuration includes:

```json
"annotations": {
  "org.opencontainers.image.description": "A Docker image for running the Fabric CLI. See https://github.com/creativerezz/Fabric/tree/main/scripts/docker for details."
}
```

This ensures that GHCR packages display the proper description.
