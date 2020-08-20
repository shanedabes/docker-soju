#!/usr/bin/env sh

WORKDIR="$(dirname $0)"

if [ ! -f "${WORKDIR}/soju.db" ]; then
    "${WORKDIR}/dbgen"
fi

# Add certs to trust cos synirc is dogshit
if [ ! -z "${SOJU_TRUST_ADD}" ]; then
    # mkdir /usr/local/share/ca-certificates/extra
    for s in "${SOJU_TRUST_ADD}"; do
        openssl s_client -connect "${s}" 2>/dev/null </dev/null | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > "/usr/local/share/ca-certificates/${s}.pem"
    done
    update-ca-certificates
fi

"${WORKDIR}/soju" -listen irc+insecure://0.0.0.0:6667
