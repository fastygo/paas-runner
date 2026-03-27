# Deploy Checklist

Use this checklist when setting up a new project with `paas.exe`.

## One-time machine setup

1. Create or reuse an SSH key in `C:/Users/<you>/.ssh/`.
2. Add the public key to the target server.
3. Verify SSH login works:

```bash
ssh root@your-server
```

4. Create `C:/Users/<you>/.config/paas/servers.yml` from `servers.example.yml`.

## Project setup

1. Copy `paas.exe` into the project root, or call it from a shared tools directory.
2. Copy `config.example.yml` to `.paas/config.yml`.
3. Adjust `.paas/extensions/deploy.yml` for the project.
4. Keep `.paas/config.yml` out of git if it contains local values.
5. Do not store passwords in the project config.

## Before each deploy

1. Load the SSH key into the agent if needed:

```bash
eval "$(ssh-agent -s)"
ssh-add ~/.ssh/id_ed25519
```

2. Export runtime secrets:

```bash
export INPUT_REGISTRY_USERNAME="..."
export INPUT_REGISTRY_PASSWORD="..."
export INPUT_DASHBOARD_PASS="..."
```

3. Validate the extension:

```bash
./paas.exe validate deploy
```

4. Run deploy with the minimal environment wrapper:

```bash
env -i \
  HOME="$HOME" \
  USERPROFILE="$USERPROFILE" \
  HOMEDRIVE="$HOMEDRIVE" \
  HOMEPATH="$HOMEPATH" \
  PATH="$PATH" \
  TERM="${TERM:-xterm-256color}" \
  LANG="${LANG:-en_US.UTF-8}" \
  SSH_AUTH_SOCK="$SSH_AUTH_SOCK" \
  SSH_AGENT_PID="$SSH_AGENT_PID" \
  INPUT_REGISTRY_USERNAME="$INPUT_REGISTRY_USERNAME" \
  INPUT_REGISTRY_PASSWORD="$INPUT_REGISTRY_PASSWORD" \
  INPUT_DASHBOARD_PASS="$INPUT_DASHBOARD_PASS" \
  ./paas.exe run deploy
```

## After deploy

1. Check the container status on the server.
2. Check application logs if the site is unavailable.
3. Verify the public health endpoint.
4. If the container fails immediately, inspect entrypoint and line endings first.
