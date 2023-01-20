#!/usr/bin/env bash

git add bin
git diff-index --quiet HEAD bin

if [[ $? -eq 0 ]]; then
    echo "👍 No changes to push!"
    exit 0
fi

echo "Pushing changes for updated binary files..."

commit_message=$(cat<<EOF
Continuous Integration Build Artifacts

$(git status --porcelain --untracked-files=no bin)
EOF
)

echo "${commit_message}"

# CI variable is always set to "true" for GitHub actions
# This will if we're on GitHub and if it's safe to push
if [[ -n "${CI}" ]]; then
    git config user.name "Continuous Integration"
    git config user.email 'heroku-production-services@salesforce.com'
    git commit -a -m "${commit_message}"
    git push
fi
