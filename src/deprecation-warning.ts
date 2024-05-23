import { warning, summary } from "@actions/core";

function dedent(str: string): string {
    const strings = str.split("\n");
    const match = str.match(/\n\s+/);

    if (!match) {
        return str;
    }

    const indentRegex = new RegExp(match[0], "g");
    const fixedStr = str.replace(indentRegex, "\n").replace(/^\n/, "");

    return fixedStr;
}

export default function deprecationWarning() {
    const oldActionName = "heroku/use-app-token-action";
    const oldActionUrl = `https://github.com/${oldActionName}`;
    const newActionName = "actions/create-github-app-token";
    const newActionUrl = `https://github.com/${newActionName}`;
    const summaryMarkdown = dedent(`
    ## ⚠️ Deprecation warning: [${oldActionName}](${oldActionUrl})

    > [!WARNING]
    > Please note that [this action](${oldActionUrl}) is deprecated and will be removed in the future.
    > We recommend using the [${newActionName}](${newActionUrl}) action instead.
    `);

    warning(
        `This action is deprecated. Please use the '${newActionName}' action from the GitHub Marketplace instead.`,
        { title: `Deprecation warning: ${oldActionName}` }
    );
    summary.addRaw(summaryMarkdown);
}
