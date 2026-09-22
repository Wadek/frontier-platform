# Self-hosted GitHub Actions runner (Linux) with Stage 1 scanners on PATH.
# No sudo on the host: run via Docker Desktop Compose (see compose.yml).
FROM myoung34/github-runner:ubuntu-noble

# Install tooling as root; the upstream entrypoint expects root so it can
# manage the runner service and optional Docker socket.
# checkov:skip=CKV_DOCKER_8: GitHub Actions runner image must stay root for registration and Docker socket access
USER root
RUN apt-get update \
  && apt-get install -y --no-install-recommends curl ca-certificates python3 python3-pip python3-venv jq procps \
  && rm -rf /var/lib/apt/lists/*

# Pin all Stage 1 scanners into /usr/local/bin (no per-job pip).
ARG GITLEAKS_VERSION=8.30.1
ARG TRIVY_VERSION=0.74.0
ARG CHECKOV_VERSION=3.2.334
ARG SEMGREP_VERSION=1.97.0
RUN curl -fsSL -o /tmp/gitleaks.tgz \
      "https://github.com/gitleaks/gitleaks/releases/download/v${GITLEAKS_VERSION}/gitleaks_${GITLEAKS_VERSION}_linux_x64.tar.gz" \
  && tar -xzf /tmp/gitleaks.tgz -C /usr/local/bin gitleaks \
  && chmod +x /usr/local/bin/gitleaks \
  && curl -fsSL -o /tmp/trivy.tgz \
      "https://github.com/aquasecurity/trivy/releases/download/v${TRIVY_VERSION}/trivy_${TRIVY_VERSION}_Linux-64bit.tar.gz" \
  && tar -xzf /tmp/trivy.tgz -C /usr/local/bin trivy \
  && chmod +x /usr/local/bin/trivy \
  && python3 -m venv /opt/checkov \
  && /opt/checkov/bin/pip install --no-cache-dir "checkov==${CHECKOV_VERSION}" \
  && python3 -m venv /opt/semgrep \
  && /opt/semgrep/bin/pip install --no-cache-dir "semgrep==${SEMGREP_VERSION}" "setuptools==75.8.2" \
  && ln -sf /opt/checkov/bin/checkov /usr/local/bin/checkov \
  && ln -sf /opt/semgrep/bin/semgrep /usr/local/bin/semgrep \
  && rm -f /tmp/gitleaks.tgz /tmp/trivy.tgz \
  && gitleaks version \
  && trivy --version | head -1 \
  && checkov --version \
  && semgrep --version

# Runner process liveness (Runner.Listener under /actions-runner).
HEALTHCHECK --interval=60s --timeout=10s --start-period=120s --retries=3 \
  CMD pgrep -f Runner.Listener >/dev/null || exit 1

# Keep the image's default entrypoint/cmd for the runner.
