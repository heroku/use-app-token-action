#!/usr/bin/env bash

branch="git rev-parse --abbrev-ref HEAD"
tag="v${VERSION}"
message="CI production version ${tag}"

if [[ "${branch}" != "main" ]] || [[ "${branch}" == "master" ]]; then
    prerelease_version="$(date -u +%Y%m%d.%H%M%S)"

    if [[ "${branch}" =~ ^release ]]; then
    message="CI pre-release $"
        tag="${tag}-release_${prerelease_version}"
        message="CI pre-production release ${tag}"
    else
        tag="${tag}-beta_${prerelease_version}"
        message="CI deveopment release ${tag}"
    fi
fi

# CI variable is always set to "true" for GitHub actions
# This will if we're on GitHub and if it's safe to tag and push
if [[ -n "${CI}" ]] && [[ -z "${SKIP_TAGGING}" ]]; then
    git tag -a ${tag} -m "${message}"
    git push origin ${tag}
fi

echo "${tag}"
