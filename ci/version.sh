#!/usr/bin/env bash

branch="git rev-parse --abbrev-ref HEAD"
commit_hash="$(git rev-parse --short HEAD)"
tag="v${VERSION}"
message="CI production version ${tag}"

if [[ "${branch}" != "main" ]] || [[ "${branch}" == "master" ]]; then
    if [[ "${branch}" =~ ^release ]]; then
    message="CI pre-release $"
        tag="${tag}-release_${commit_hash}"
        message="CI pre-production release ${tag}"
    else
        tag="${tag}-beta_${commit_hash}"
        message="CI deveopment release ${tag}"
    fi
fi

# CI variable is always set to "true" for GitHub actions
# This will if we're on GitHub and if it's safe to tag and push
if [[ -n "${CI}" ]]; then
    git tag -a ${tag} -m "${message}"
    git push origin ${tag}
fi

echo "${tag}"
