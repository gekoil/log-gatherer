FROM bash:latest


ENTRYPOINT ["/usr/local/bin/bash", "-c", "i=0; while true; do echo \"[$(uname -n)] $(date)\"; i=$((i+1)); sleep 1; done"]